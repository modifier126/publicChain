package blc

import (
	"blockDemo/blc/utils"
	"bytes"
)

type TxOutput struct {
	Money         int64
	Ripemd160Hash []byte
}

// 解锁
func (out *TxOutput) UnlockPublicScriptKeyWithAddr(address string) bool {
	pubKeyBytes := utils.Base58Decode([]byte(address))
	ripemd160Hash := pubKeyBytes[1 : len(pubKeyBytes)-4]
	return bytes.Equal(out.Ripemd160Hash, ripemd160Hash)
}

func NewTxOutput(value int64, address string) *TxOutput {
	txOutput := TxOutput{value, nil}
	txOutput.Lock(address)
	return &txOutput
}

func (tx *TxOutput) Lock(address string) {
	pubKeyBytes := utils.Base58Decode([]byte(address))
	ripemd160Hash := pubKeyBytes[1 : len(pubKeyBytes)-4]
	tx.Ripemd160Hash = ripemd160Hash
}
