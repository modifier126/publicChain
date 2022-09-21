package blc

import (
	"blockDemo/blc/utils"
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
	"os"
)

// 目标难度值
const targetBit = 16

// 工作量证明结构
type ProofOfWork struct {
	// 需要共识验证的区块
	Block *Block
	// 目标难度的哈希(大数据存储)
	target *big.Int
}

// 创建pow对象
func NewProofOfWork(block *Block) *ProofOfWork {
	target := big.NewInt(1)
	target = target.Lsh(target, 256-targetBit)
	return &ProofOfWork{Block: block, target: target}
}

// 执行pow,比较哈希
// 返回哈希值，以及碰撞次数
func (pow *ProofOfWork) Run() ([]byte, int64) {
	// 碰撞次数
	var nonce = 0

	var hashInt big.Int

	var hash [32]byte

	// 无限循环生成符合条件的哈希
	for {
		// 生成准备数据
		dataBytes := pow.prepareData(int64(nonce))
		hash = sha256.Sum256(dataBytes)

		// 字节数值转big.Int
		hashInt.SetBytes(hash[:])

		// 检测生成的哈希值
		if pow.target.Cmp(&hashInt) == 1 {
			// 如果满足条件,终止循环
			break
		}
		// 修改值继续循环
		nonce++

		fmt.Fprintf(os.Stdout, "%x\r", hash)
		//time.Sleep(1 * time.Millisecond)
	}
	//fmt.Printf("碰撞次数:%d\n", nonce)
	return hash[:], int64(nonce)
}

// 生成准备数据
func (pow *ProofOfWork) prepareData(nonce int64) []byte {
	var data []byte
	// 拼接区块属性，进行哈希计算
	timeStamp := utils.IntToHex(pow.Block.TimeStamp)
	height := utils.IntToHex(pow.Block.Height)

	data = bytes.Join([][]byte{
		timeStamp,
		height,
		pow.Block.PreHash,
		pow.Block.HasTransaction(),
		utils.IntToHex(nonce),
		utils.IntToHex(targetBit),
	}, []byte{})

	return data
}

// 校验区块是否有效
func (pow *ProofOfWork) IsValid() bool {
	hash := pow.Block.Hash
	var hashInt big.Int
	hashInt.SetBytes(hash)
	return pow.target.Cmp(&hashInt) == 1
}
