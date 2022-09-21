package blc

import (
	"blockDemo/blc/utils"
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
	"time"

	"github.com/boltdb/bolt"
)

const dbName = "blockChain_%s.db"
const tableName = "block.bucket"

// 定义区块链结构
type BlockChain struct {
	//Blocks []*Block

	Tip []byte
	DB  *bolt.DB
}

// 初始化区块链
func CreateBlockChainWithGenesis(address string, nodeId string) *BlockChain {
	// 数据库是否存在
	if DBExists(nodeId) {
		fmt.Println("创世区块已存在...")
		os.Exit(1)
	}
	dbName := fmt.Sprintf(dbName, nodeId)
	var hash []byte
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte(tableName))
		if err != nil {
			log.Panic(err)
		}
		cbt := NewCoinBaseTrans(address)
		block := CreateGenesisBlock([]*Transaction{cbt})
		hash = block.Hash
		dataByte := block.Serialize()
		// 哈希与区块关系
		b.Put(hash, dataByte)
		// 最后一个区块的哈希
		b.Put([]byte("l"), hash)
		return nil
	})

	if err != nil {
		log.Panic(err)
	}
	return &BlockChain{hash, db}
}

// 添加区块到区块链中
func (bc *BlockChain) AddBlock(txs []*Transaction) {
	var height int64
	var preBlockHash []byte

	bc.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(tableName))
		if b != nil {
			dataByte := b.Get(bc.Tip)
			block := DeSerializeBlock(dataByte)
			height = block.Height + 1
			preBlockHash = block.Hash
		}
		return nil
	})

	err := bc.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(tableName))
		if b == nil {
			var err error
			b, err = tx.CreateBucket([]byte(tableName))
			if err != nil {
				log.Panic(err)
			}
		}

		// 生成新的区块
		block := NewBlock(height, preBlockHash, txs)
		b.Put(block.Hash, block.Serialize())
		b.Put([]byte("l"), block.Hash)
		// 重新设置
		bc.Tip = block.Hash
		return nil
	})

	if err != nil {
		log.Panic(err)
	}
}

func (bc *BlockChain) PrintChain() {

	bci := bc.NewBlockChainIter()
	for {
		block := bci.Next()
		// 输出区块
		fmt.Printf("block.Height:%d\n", block.Height)
		fmt.Printf("block.TimeStamp:%s\n", time.Unix(block.TimeStamp, 0).Format("2006-01-02 3:4:5 pm"))
		fmt.Printf("block.Hash:%x\n", block.Hash)
		fmt.Printf("block.PreHash:%x\n", block.PreHash)
		fmt.Printf("block.Nonce:%d\n", block.Nonce)
		//fmt.Printf("block.txs:%v\n", block.Txs)
		fmt.Println()
		fmt.Println("tx:")
		for _, tx := range block.Txs {
			fmt.Printf("txHash:%x\n", tx.TxHash)
			fmt.Println("Vins:")
			for _, vin := range tx.Vins {
				fmt.Printf("vin.TxHash:%x\n", vin.TxHash)
				fmt.Printf("vin.Vout:%d\n", vin.Vout)
				fmt.Printf("vin.PublicKey:%x\n", vin.PublicKey)
			}
			fmt.Println("Vouts:")
			for _, vout := range tx.Vouts {
				fmt.Printf("vout.Money:%d\n", vout.Money)
				fmt.Printf("vout.Ripemd160Hash:%x\n", vout.Ripemd160Hash)
			}
		}

		fmt.Printf("-------------------------------------------------------\n")

		var bigInt big.Int
		bigInt.SetBytes(block.PreHash)
		if big.NewInt(0).Cmp(&bigInt) == 0 {
			break
		}
	}

}

// 新建区块链迭代器
func (bc *BlockChain) NewBlockChainIter() *BlockChainIter {
	return &BlockChainIter{
		CurHash: bc.Tip,
		DB:      bc.DB,
	}
}

func DBExists(nodeId string) bool {
	dbName := fmt.Sprintf(dbName, nodeId)
	if _, err := os.Stat(dbName); os.IsNotExist(err) {
		return false
	}
	return true
}

func BlockChainObj(nodeId string) *BlockChain {
	dbName := fmt.Sprintf(dbName, nodeId)
	var hash []byte
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(tableName))
		if b != nil {
			hash = b.Get([]byte("l"))
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	return &BlockChain{Tip: hash, DB: db}
}

func (bc *BlockChain) MineNewBlock(txs []*Transaction) *Block {
	var block *Block
	err := bc.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(tableName))
		if b != nil {
			hash := b.Get([]byte("l"))
			dataBytes := b.Get(hash)
			block = DeSerializeBlock(dataBytes)
		}
		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	//2、建立新区块
	block = NewBlock(block.Height+1, block.Hash, txs)
	err = bc.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(tableName))
		if b != nil {
			b.Put(block.Hash, []byte(block.Serialize()))
			// 最后区块哈希
			b.Put([]byte("l"), block.Hash)
			bc.Tip = block.Hash
		}
		return nil
	})

	if err != nil {
		log.Panic(err)
	}
	return block
}

func (bc *BlockChain) UnUTXOs(address string, txs []*Transaction) []*UTXO {

	var unUTXOs []*UTXO
	var spentOutputMap = make(map[string][]int)

	for _, tx := range txs {

	work1:
		for index, out := range tx.Vouts {
			if out.UnlockPublicScriptKeyWithAddr(address) {
				if len(spentOutputMap) == 0 {
					utxo := &UTXO{
						TxHash: tx.TxHash,
						Index:  index,
						Out:    out,
					}
					unUTXOs = append(unUTXOs, utxo)
				} else {
					for hash, indexArr := range spentOutputMap {
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

	bcIter := bc.NewBlockChainIter()
	for {
		block := bcIter.Next()
		for i := len(block.Txs) - 1; i >= 0; i-- {
			tx := block.Txs[i]
			if !tx.IsCoinBaseTrans() {
				for _, in := range tx.Vins {
					pubKeyBytes := utils.Base58Decode([]byte(address))
					ripemd160Hash := pubKeyBytes[1 : len(pubKeyBytes)-4]
					if in.UnlockRipemd160Hash(ripemd160Hash) {
						// 映射的是哈希区块
						key := hex.EncodeToString(in.TxHash)
						spentOutputMap[key] = append(spentOutputMap[key], in.Vout)
					}
				}
			}
		work2:
			for index, out := range tx.Vouts {
				if out.UnlockPublicScriptKeyWithAddr(address) {
					if len(spentOutputMap) != 0 {
						var isSpentUTXO bool

						for txHash, idxArr := range spentOutputMap {
							for _, idx := range idxArr {
								if index == idx && txHash == hex.EncodeToString(tx.TxHash) {
									isSpentUTXO = true
									continue work2
								}
							}
						}

						if !isSpentUTXO {
							utxo := &UTXO{tx.TxHash, index, out}
							unUTXOs = append(unUTXOs, utxo)
						}
					} else {
						utxo := &UTXO{tx.TxHash, index, out}
						unUTXOs = append(unUTXOs, utxo)
					}
				}
			}
		}

		var bigInt big.Int
		bigInt.SetBytes(block.PreHash)
		if bigInt.Cmp(big.NewInt(0)) == 0 {
			break
		}
	}
	return unUTXOs
}

func (bc *BlockChain) GetBalance(address string) int64 {
	var amount int64
	utxos := bc.UnUTXOs(address, []*Transaction{})
	for _, utxo := range utxos {
		amount = amount + utxo.Out.Money
	}
	return amount
}

func (bc *BlockChain) FindSpendableUTXOs(address string, amount int, txs []*Transaction) (int64, map[string][]int) {

	var value int64
	var spendableUtxoDic = make(map[string][]int)

	// 先获取所有utxo
	utxos := bc.UnUTXOs(address, txs)
	for _, utxo := range utxos {
		hash := hex.EncodeToString(utxo.TxHash)
		spendableUtxoDic[hash] = append(spendableUtxoDic[hash], utxo.Index)
		value += utxo.Out.Money
		if value >= int64(amount) {
			break
		}
	}

	if value < int64(amount) {
		fmt.Printf("%s的余额不足", address)
		os.Exit(1)
	}
	return value, spendableUtxoDic
}

func (bc *BlockChain) SignTransaction(tx *Transaction, privKey ecdsa.PrivateKey, txs []*Transaction) {
	prevTxs := make(map[string]Transaction)

	for _, vin := range tx.Vins {
		prevTx, err := bc.FindTransaction(vin.TxHash, txs)
		if err != nil {
			log.Panic(err)
		}

		prevTxs[hex.EncodeToString(prevTx.TxHash)] = prevTx
	}
	// 签名
	tx.Sign(privKey, prevTxs)

}

func (bc *BlockChain) FindTransaction(id []byte, txs []*Transaction) (Transaction, error) {
	for _, tx := range txs {
		if bytes.Equal(tx.TxHash, id) {
			return *tx, nil
		}
	}

	bci := bc.NewBlockChainIter()
	for {
		b := bci.Next()
		for _, tx := range b.Txs {
			if bytes.Equal(tx.TxHash, id) {
				return *tx, nil
			}
		}

		var bigInt big.Int
		bigInt.SetBytes(b.PreHash)
		if big.NewInt(0).Cmp(&bigInt) == 0 {
			break
		}
	}
	return Transaction{}, nil
}

func (bc *BlockChain) VerifyTransaction(tx *Transaction, txs []*Transaction) bool {
	prevTxs := make(map[string]Transaction)

	for _, vin := range tx.Vins {
		prevTx, err := bc.FindTransaction(vin.TxHash, txs)
		if err != nil {
			log.Panic(err)
		}
		prevTxs[hex.EncodeToString(prevTx.TxHash)] = prevTx
	}
	return tx.Verify(prevTxs)
}

func (bc *BlockChain) FindUtoxMap() map[string]*TxOutputs {
	bcIter := bc.NewBlockChainIter()
	var spendableUTXOsMap = make(map[string][]*TxInput)
	utoxMaps := make(map[string]*TxOutputs)

	for {
		block := bcIter.Next()
		for i := len(block.Txs) - 1; i >= 0; i-- {
			txOutputs := &TxOutputs{[]*UTXO{}}
			tx := block.Txs[i]
			if !tx.IsCoinBaseTrans() {
				for _, vin := range tx.Vins {
					inputTxHash := hex.EncodeToString(vin.TxHash)
					spendableUTXOsMap[inputTxHash] = append(spendableUTXOsMap[inputTxHash], vin)
				}
			}

			txHash := hex.EncodeToString(tx.TxHash)

		work1:
			for idx, vout := range tx.Vouts {
				txInputs := spendableUTXOsMap[txHash]
				if len(txInputs) > 0 {

					isSpent := false

					for _, txInput := range txInputs {
						outPubKey := vout.Ripemd160Hash
						inPubKey := txInput.PublicKey
						if bytes.Equal(outPubKey, Ripemd160Hash(inPubKey)) {
							if idx == txInput.Vout {
								isSpent = true
								continue work1
							}
						}
					}
					if !isSpent {
						utox := &UTXO{tx.TxHash, idx, vout}
						txOutputs.Utxos = append(txOutputs.Utxos, utox)
					}
				} else {
					utox := &UTXO{tx.TxHash, idx, vout}
					txOutputs.Utxos = append(txOutputs.Utxos, utox)
				}
			}
			utoxMaps[txHash] = txOutputs
		}

		var bigInt big.Int
		bigInt.SetBytes(block.PreHash)
		if big.NewInt(0).Cmp(&bigInt) == 0 {
			break
		}
	}
	return utoxMaps
}

func (bc *BlockChain) GetBestHeight() int64 {
	bcIter := bc.NewBlockChainIter()
	b := bcIter.Next()
	return b.Height
}

func (bc *BlockChain) GetBlockHashes() [][]byte {
	bcIter := bc.NewBlockChainIter()
	var hashes [][]byte
	for {
		block := bcIter.Next()
		hashes = append(hashes, block.Hash)
		var bigInt big.Int
		bigInt.SetBytes(block.PreHash)
		if big.NewInt(0).Cmp(&bigInt) == 0 {
			break
		}
	}
	return hashes

}

func (bc *BlockChain) GetBlock(blockHash []byte) (*Block, error) {

	var block *Block

	err := bc.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(tableName))
		if b != nil {
			blockData := b.Get(blockHash)
			block = DeSerializeBlock(blockData)
		}
		return nil
	})

	if err != nil {
		return block, err
	}
	return block, nil
}

func (bc *BlockChain) SynBlock(block *Block) error {

	err := bc.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(tableName))
		if b != nil {
			blockExist := b.Get(block.Hash)

			if blockExist == nil {
				err := b.Put(block.Hash, block.Serialize())
				if err != nil {
					log.Panic(err)
				}

				blockBytes := b.Get(b.Get([]byte("l")))
				blockInDB := DeSerializeBlock(blockBytes)

				if blockInDB.Height < block.Height {
					b.Put([]byte("l"), block.Hash)
					bc.Tip = block.Hash
				}
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (bc *BlockChain) NewTrans(from []string, to []string, amount []string, nodeId string) []*Transaction {
	utxoSet := &UTXOSet{bc}
	//1、通过算法建立Transaction
	var txs []*Transaction

	for i, addr := range from {
		val, _ := strconv.Atoi(amount[i])
		tx := NewSimpleTransaction(addr, to[i], int64(val), utxoSet, txs, nodeId)
		txs = append(txs, tx)
	}
	// 奖励
	tx := NewCoinBaseTrans(from[0])
	txs = append(txs, tx)

	// 验证交易
	var _txs []*Transaction
	for _, tx := range txs {
		if !bc.VerifyTransaction(tx, _txs) {
			log.Panic("签名失败")
		}
		_txs = append(_txs, tx)
	}
	return txs
}
