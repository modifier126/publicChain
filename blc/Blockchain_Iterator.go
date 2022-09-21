package blc

import (
	"log"

	"github.com/boltdb/bolt"
)

type BlockChainIter struct {
	CurHash []byte
	DB      *bolt.DB
}

// 下一个元素
func (bcIter *BlockChainIter) Next() *Block {
	// 定义变量
	var block *Block

	err := bcIter.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(tableName))
		if b != nil {
			dataByte := b.Get(bcIter.CurHash)
			block = DeSerializeBlock(dataByte)
			bcIter.CurHash = block.PreHash
		}
		return nil
	})

	if err != nil {
		log.Panic(err)
	}
	return block
}
