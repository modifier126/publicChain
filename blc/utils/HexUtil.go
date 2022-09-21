package utils

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"encoding/json"
	"log"
)

// 实现int64->hash
func IntToHex(data int64) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, data)
	if err != nil {
		log.Panicf("int transact to []byte failure!%v\n", err)
	}
	return buf.Bytes()
}

// json字符串转字符串数组
func JsonToArray(jsonStr string) []string {
	var sArr []string
	if err := json.Unmarshal([]byte(jsonStr), &sArr); err != nil {
		log.Panic(err)
	}
	return sArr
}

func GodEncode(data interface{}) []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(data)
	if err != nil {
		log.Panic(err)
	}
	return buf.Bytes()
}
