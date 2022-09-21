package blc

import (
	"blockDemo/blc/utils"
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
	"time"
)

type Block struct {
	TimeStamp int64          // 区块时间戳,代码区块时间
	Hash      []byte         // 当前区块哈希
	PreHash   []byte         //  前区块哈希
	Height    int64          // 区块高度
	Txs       []*Transaction // 交易数据
	Nonce     int64          // 碰撞次数
}

func NewBlock(height int64, preHash []byte, txs []*Transaction) *Block {
	block := Block{
		TimeStamp: time.Now().Unix(),
		Hash:      nil,
		PreHash:   preHash,
		Height:    height,
		Txs:       txs,
		Nonce:     0,
	}
	// 通过pow生成新的哈希
	pow := NewProofOfWork(&block)
	hash, nonce := pow.Run()
	block.Hash = hash
	block.Nonce = nonce
	return &block
}

// 计算区块哈希
func (b *Block) SetHash() {
	// 调用sha256实现哈希生成
	timeStamp := utils.IntToHex(b.TimeStamp)
	height := utils.IntToHex(b.Height)

	blockBytes := bytes.Join([][]byte{
		timeStamp,
		height,
		b.Hash,
		b.PreHash,
		b.HasTransaction(),
	}, []byte{})
	hash := sha256.Sum256(blockBytes)
	b.Hash = hash[:]
}

// 生成创世区块
func CreateGenesisBlock(txs []*Transaction) *Block {
	fmt.Println("正在创建创世区块...")
	b := NewBlock(1, []byte{0}, txs)
	fmt.Printf("%x\n", b.Hash)
	return b
}

// 区块反序列号
func (b *Block) Serialize() []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(b)
	if err != nil {
		log.Panic(err)
	}
	return buf.Bytes()
}

// 反序列化字节内容
func DeSerializeBlock(data []byte) *Block {
	var block Block
	dec := gob.NewDecoder(bytes.NewReader(data))
	err := dec.Decode(&block)
	if err != nil {
		log.Panic(err)
	}
	return &block
}

func (b *Block) HasTransaction() []byte {
	var txHash []byte
	var txHashes [][]byte

	for _, tx := range b.Txs {
		txHashes = append(txHashes, tx.TxHash)
	}
	//txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	mt := NewMerkleTree(txHashes)
	txHash = mt.RootNode.Data
	return txHash[:]
}
