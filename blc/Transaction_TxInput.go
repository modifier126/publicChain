package blc

import "bytes"

type TxInput struct {
	TxHash    []byte // 交易哈希
	Vout      int    // 存储在Txoutput的Vouts的索引
	Signature []byte // 签名
	PublicKey []byte // 公钥,钱包里面的
}

// 解锁地址
func (in *TxInput) UnlockRipemd160Hash(ripemd160Hash []byte) bool {
	publicKey := Ripemd160Hash(in.PublicKey)
	return bytes.Equal(publicKey, ripemd160Hash)
}
