package blc

import (
	"blockDemo/blc/utils"
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"log"
	"math/big"
	"time"
)

type Transaction struct {
	// 交易哈希
	TxHash []byte
	// 输入
	Vins []*TxInput
	// 输出
	Vouts []*TxOutput
}

func NewCoinBaseTrans(address string) *Transaction {
	txIn := &TxInput{[]byte{}, -1, nil, []byte{}}
	txOut := NewTxOutput(10, address)
	tx := &Transaction{[]byte{}, []*TxInput{txIn}, []*TxOutput{txOut}}
	tx.SetTxHash()
	return tx
}

func (t *Transaction) SetTxHash() {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(t)
	if err != nil {
		log.Panic(err)
	}
	bytesJoin := bytes.Join([][]byte{buf.Bytes(),
		utils.IntToHex(time.Now().Unix())}, []byte{})
	hash := sha256.Sum256(bytesJoin)
	t.TxHash = hash[:]
}

func NewSimpleTransaction(from string, to string, amount int64, utoxSet *UTXOSet, txs []*Transaction, nodeId string) *Transaction {
	var txIns []*TxInput
	var txOuts []*TxOutput

	wallets, _ := NewWalletList(nodeId)
	wallet := wallets.Wallets

	money, spendableUTXODic := utoxSet.FindSpendableUTXOs(from, amount, txs)

	for txHash, indexArr := range spendableUTXODic {
		bytes, err := hex.DecodeString(txHash)
		if err != nil {
			log.Panic(err)
		}
		for _, index := range indexArr {
			txIn := &TxInput{bytes, index, nil, wallet[from].PublicKey}
			txIns = append(txIns, txIn)
		}
	}

	txOut := NewTxOutput(int64(amount), to)
	txOuts = append(txOuts, txOut)
	txOut = NewTxOutput(int64(money)-int64(amount), from)
	txOuts = append(txOuts, txOut)

	tx := &Transaction{
		[]byte{},
		txIns,
		txOuts,
	}
	tx.SetTxHash()

	// 签名
	utoxSet.BlockChain.SignTransaction(tx, wallet[from].PrivateKey, txs)
	return tx
}

func (tx *Transaction) IsCoinBaseTrans() bool {
	return len(tx.Vins[0].TxHash) == 0 && tx.Vins[0].Vout == -1
}

func (tx *Transaction) Sign(privKey ecdsa.PrivateKey, prevTxs map[string]Transaction) {
	if tx.IsCoinBaseTrans() {
		return
	}

	for _, vin := range tx.Vins {
		if prevTxs[hex.EncodeToString(vin.TxHash)].TxHash == nil {
			log.Panic("Error: sign previous Transaction not correct")
		}
	}

	copyTx := tx.CopyTx()

	for idx, vin := range copyTx.Vins {
		prevTx := prevTxs[hex.EncodeToString(vin.TxHash)]
		copyTx.Vins[idx].Signature = nil
		copyTx.Vins[idx].PublicKey = prevTx.Vouts[vin.Vout].Ripemd160Hash
		copyTx.TxHash = copyTx.Hash()
		// 签名代码
		r, s, err := ecdsa.Sign(rand.Reader, &privKey, copyTx.TxHash)
		if err != nil {
			log.Panic(err)
		}
		signature := append(r.Bytes(), s.Bytes()...)

		tx.Vins[idx].Signature = signature
		copyTx.Vins[idx].PublicKey = nil
	}

}

func (tx *Transaction) CopyTx() Transaction {
	var txInput []*TxInput
	var txOutput []*TxOutput

	for _, vin := range tx.Vins {
		txInput = append(txInput, &TxInput{vin.TxHash, vin.Vout, nil, nil})
	}

	for _, vout := range tx.Vouts {
		txOutput = append(txOutput, &TxOutput{vout.Money, vout.Ripemd160Hash})
	}

	txCopy := Transaction{tx.TxHash, txInput, txOutput}
	return txCopy
}

func (tx *Transaction) Hash() []byte {
	var hash [32]byte
	txCopy := *tx
	txCopy.TxHash = []byte{}
	hash = sha256.Sum256(txCopy.Serialize())
	return hash[:]
}

func (tx *Transaction) Serialize() []byte {
	var encoded bytes.Buffer
	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	return encoded.Bytes()
}

func (tx *Transaction) Verify(prevTxs map[string]Transaction) bool {
	if tx.IsCoinBaseTrans() {
		return true
	}
	for _, vin := range tx.Vins {
		if prevTxs[hex.EncodeToString(vin.TxHash)].TxHash == nil {
			log.Panic("Error: verify previous Transaction not correct")
		}
	}

	copyTx := tx.CopyTx()
	curve := elliptic.P256()

	for idx, vin := range tx.Vins {
		prevTx := prevTxs[hex.EncodeToString(vin.TxHash)]
		copyTx.Vins[idx].Signature = nil
		copyTx.Vins[idx].PublicKey = prevTx.Vouts[vin.Vout].Ripemd160Hash
		copyTx.TxHash = copyTx.Hash()

		r := big.Int{}
		s := big.Int{}
		sigLen := len(vin.Signature)
		r.SetBytes(vin.Signature[:(sigLen / 2)])
		s.SetBytes(vin.Signature[(sigLen / 2):])

		x := big.Int{}
		y := big.Int{}
		pkLen := len(vin.PublicKey)
		x.SetBytes(vin.PublicKey[:(pkLen / 2)])
		y.SetBytes(vin.PublicKey[(pkLen / 2):])
		rawPubKey := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}
		if !ecdsa.Verify(&rawPubKey, copyTx.TxHash, &r, &s) {
			return false
		}

		copyTx.Vins[idx].PublicKey = nil

	}
	return true
}

func DeSerializeTransaction(data []byte) *Transaction {
	var tx Transaction
	dec := gob.NewDecoder(bytes.NewReader(data))
	err := dec.Decode(&tx)
	if err != nil {
		log.Panic(err)
	}
	return &tx
}
