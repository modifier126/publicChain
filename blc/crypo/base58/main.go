package main

import (
	"bytes"
	"fmt"
	"math/big"
	"math/bits"
)

const (
	_S = _W / 8 // word size in bytes

	_W = bits.UintSize // word size in bits
	_B = 1 << _W       // digit base
	_M = _B - 1        // digit mask
)

// 种子码
var b58Alphabet = []byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")

func ReverseBytes(data []byte) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}

func Base58Encoding(input []byte) []byte {
	var result []byte
	x := big.NewInt(0).SetBytes(input)

	base := big.NewInt(int64(len(b58Alphabet)))
	zero := big.NewInt(0)
	mod := &big.Int{}

	for x.Cmp(zero) != 0 {
		x.DivMod(x, base, mod)
		result = append(result, b58Alphabet[mod.Int64()])
	}

	ReverseBytes(result)

	for b := range input {
		if b == 0x00 {
			result = append([]byte{b58Alphabet[0]}, result...)
		} else {
			break
		}
	}

	return result
}

func Base58Decode(input []byte) []byte {
	result := big.NewInt(0)
	zeroBytes := 0

	for b := range input {
		if b == 0x00 {
			zeroBytes++
		}
	}

	payload := input[zeroBytes:]

	for _, b := range payload {
		charIndex := bytes.IndexByte(b58Alphabet, b)
		result.Mul(result, big.NewInt(58))
		result.Add(result, big.NewInt(int64(charIndex)))
	}

	decoded := result.Bytes()
	decoded = append(bytes.Repeat([]byte{byte(0x00)}, zeroBytes), decoded...)
	return decoded
}

func main() {
	res := Base58Encoding([]byte("modifier126"))
	fmt.Printf("%s\n", res)

	res = Base58Decode(res)
	fmt.Printf("%s\n", res)

	fmt.Printf("_W=%d\n", _W)

	fmt.Println(uint(0))

	fmt.Println(^uint(0))
	fmt.Println(^uint(0) >> 63)

	fmt.Println(32 << (^uint(0) >> 63))

	//fmt.Println(8 << 2)
	// fmt.Println(8 >> 1)
	// fmt.Println(8 >> 2)

	fmt.Println(_S)

	//fmt.Println(len("1000000 00000000 00000000 00000000 0000000000 00000000 00000000 00000000"))

	// var a [2]int
	// fmt.Println(cap(a))
}
