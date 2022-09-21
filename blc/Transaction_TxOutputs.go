package blc

import (
	"bytes"
	"encoding/gob"
	"log"
)

type TxOutputs struct {
	Utxos []*UTXO
}

func (txOps *TxOutputs) Serialize() []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(txOps)
	if err != nil {
		log.Panic(err)
	}
	return buf.Bytes()
}

func DeSerializeTxOutputs(data []byte) *TxOutputs {
	var txOutputs TxOutputs
	dec := gob.NewDecoder(bytes.NewReader(data))
	err := dec.Decode(&txOutputs)
	if err != nil {
		log.Panic(err)
	}
	return &txOutputs
}
