package blc

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

const utxoTableName = "utxoTableName"

type UTXOSet struct {
	BlockChain *BlockChain
}

func (us *UTXOSet) ResetUTXOSet() {
	bc := us.BlockChain
	err := bc.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoTableName))
		if b != nil {
			err := tx.DeleteBucket([]byte(utxoTableName))
			if err != nil {
				log.Panic(err)
			}
		}
		b, _ = tx.CreateBucket([]byte(utxoTableName))
		if b != nil {
			txOutputsMap := us.BlockChain.FindUtoxMap()
			for keyHash, txOps := range txOutputsMap {
				txHash, _ := hex.DecodeString(keyHash)
				b.Put(txHash, txOps.Serialize())
			}
		}
		return nil
	})

	if err != nil {
		log.Panic(err)
	}
}

func (us *UTXOSet) GetBalance(address string) int64 {
	var amount int64
	utoxs := us.FindUtoxForAddress(address)
	for _, utox := range utoxs {
		amount = amount + utox.Out.Money
	}
	return amount
}

func (us *UTXOSet) FindUtoxForAddress(address string) []*UTXO {
	var utoxs []*UTXO
	bc := us.BlockChain
	err := bc.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoTableName))
		if b != nil {
			c := b.Cursor()
			for k, v := c.First(); k != nil; k, v = c.Next() {
				txOutputs := DeSerializeTxOutputs(v)
				for _, utxo := range txOutputs.Utxos {
					if utxo.Out.UnlockPublicScriptKeyWithAddr(address) {
						utoxs = append(utoxs, utxo)
					}
				}
			}
		}
		return nil
	})

	if err != nil {
		log.Panic(err)
	}
	return utoxs
}

func (us *UTXOSet) FindSpendableUTXOs(address string, amount int64, txs []*Transaction) (int64, map[string][]int) {
	var money int64
	var spentOutputMap = make(map[string][]int)

	utxos := us.FindUnPackageSpendableUTXOs(address, txs)

	var fullMoney bool
	for _, utxo := range utxos {
		money = money + utxo.Out.Money
		txHash := hex.EncodeToString(utxo.TxHash)
		spentOutputMap[txHash] = append(spentOutputMap[txHash], utxo.Index)
		if money >= amount {
			fullMoney = true
			break
		}
	}

	if !fullMoney {
		db := us.BlockChain.DB
		err := db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(utxoTableName))
			if b != nil {
				c := b.Cursor()
			work1:
				for k, v := c.First(); k != nil; k, v = c.Next() {
					txOutputs := DeSerializeTxOutputs(v)
					utxos := txOutputs.Utxos

					for _, utxo := range utxos {
						money = money + utxo.Out.Money
						txHash := hex.EncodeToString(utxo.TxHash)
						spentOutputMap[txHash] = append(spentOutputMap[txHash], utxo.Index)
						if money >= amount {
							break work1
						}
					}
				}
			}
			return nil
		})
		if err != nil {
			log.Panic(err)
		}
	}

	if money < amount {
		log.Panic("余额不足")
	}

	return money, spentOutputMap
}

func (us *UTXOSet) FindUnPackageSpendableUTXOs(address string, txs []*Transaction) []*UTXO {
	var unUTXOs []*UTXO
	var unPackageSpentOutputMap = make(map[string][]int)

	for _, tx := range txs {
	work1:
		for index, out := range tx.Vouts {
			if out.UnlockPublicScriptKeyWithAddr(address) {
				if len(unPackageSpentOutputMap) == 0 {
					utxo := &UTXO{
						TxHash: tx.TxHash,
						Index:  index,
						Out:    out,
					}
					unUTXOs = append(unUTXOs, utxo)
				} else {
					for hash, indexArr := range unPackageSpentOutputMap {
						txHashStr := hex.EncodeToString(tx.TxHash)
						if hash == txHashStr {
							var isSpentUTXO bool

							for _, idx := range indexArr {
								if idx == index {
									isSpentUTXO = true
									continue work1
								}
							}

							if !isSpentUTXO {
								utxo := &UTXO{TxHash: tx.TxHash, Index: index, Out: out}
								unUTXOs = append(unUTXOs, utxo)
							}
						} else {
							utxo := &UTXO{TxHash: tx.TxHash, Index: index, Out: out}
							unUTXOs = append(unUTXOs, utxo)
						}
					}
				}
			}
		}
	}
	return unUTXOs
}

func (us *UTXOSet) Update() {
	bcIter := us.BlockChain.NewBlockChainIter()
	block := bcIter.Next()
	var ins []*TxInput

	outsMap := make(map[string]*TxOutputs)

	for _, tx := range block.Txs {
		ins = append(ins, tx.Vins...)
	}

	for _, tx := range block.Txs {
		utxos := []*UTXO{}

		for idx, vout := range tx.Vouts {
			isSpent := false
			for _, in := range ins {
				if in.Vout == idx && bytes.Equal(in.TxHash, tx.TxHash) && bytes.Equal(Ripemd160Hash(in.PublicKey), vout.Ripemd160Hash) {
					isSpent = true
					break
				}
			}

			if !isSpent {
				utxo := &UTXO{TxHash: tx.TxHash, Index: idx, Out: vout}
				utxos = append(utxos, utxo)
			}
		}

		if len(utxos) > 0 {
			txHash := hex.EncodeToString(tx.TxHash)
			outsMap[txHash] = &TxOutputs{utxos}
		}
	}

	db := us.BlockChain.DB

	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoTableName))
		if b != nil {
			for _, in := range ins {
				txOutputsBytes := b.Get(in.TxHash)
				if len(txOutputsBytes) == 0 {
					continue
				}
				txOutputs := DeSerializeTxOutputs(txOutputsBytes)
				utxos := []*UTXO{}
				isNeedDel := false

				for _, utox := range txOutputs.Utxos {
					if utox.Index == in.Vout {
						if bytes.Equal(utox.Out.Ripemd160Hash, Ripemd160Hash(in.PublicKey)) {
							isNeedDel = true
						} else {
							utxos = append(utxos, utox)
						}
					}
				}

				if isNeedDel {
					err := b.Delete(in.TxHash)
					if err != nil {
						log.Panic(err)
					}

					if len(utxos) > 0 {
						fmt.Printf("into here...")
						txHash := hex.EncodeToString(in.TxHash)
						preTxOutputs := outsMap[txHash]
						preTxOutputs.Utxos = append(preTxOutputs.Utxos, utxos...)
						outsMap[txHash] = preTxOutputs
					}
				}

				if len(outsMap) > 0 {
					for keyHash, output := range outsMap {
						keyHashBytes, _ := hex.DecodeString(keyHash)
						b.Put(keyHashBytes, output.Serialize())
					}
				}
			}
		}
		return nil
	})

}
