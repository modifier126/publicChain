package blc

import (
	"blockDemo/blc/utils"
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"log"

	"golang.org/x/crypto/ripemd160"
)

const version = byte(0x00)
const addressChecksumLen = 4

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

func NewWallet() *Wallet {
	privateKey, publicKey := newPairKey()
	return &Wallet{privateKey, publicKey}
}

func newPairKey() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
	return *private, pubKey
}

func (w *Wallet) GetAddress() []byte {
	ripemd160Hash := Ripemd160Hash(w.PublicKey)
	versionPayload := append([]byte{version}, ripemd160Hash...)
	checksum := checkSum(versionPayload)
	fullPayload := append(versionPayload, checksum...)
	address := utils.Base58Encoding(fullPayload)
	return address
}

func Ripemd160Hash(publicKey []byte) []byte {
	hash := sha256.New()
	_, err := hash.Write(publicKey)
	if err != nil {
		log.Panic(err)
	}

	ripemd160 := ripemd160.New()
	ripemd160.Write(hash.Sum(nil))
	return ripemd160.Sum(nil)
}

func IsValidForAddress(address []byte) bool {
	version_public_checkSumBytes := utils.Base58Decode(address)
	checkSumBytes := version_public_checkSumBytes[len(version_public_checkSumBytes)-addressChecksumLen:]
	version_ripemd160 := version_public_checkSumBytes[:len(version_public_checkSumBytes)-addressChecksumLen]
	checkBytes := checkSum(version_ripemd160)
	return bytes.Equal(checkSumBytes, checkBytes)
}

func checkSum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])
	return secondSHA[:addressChecksumLen]
}
