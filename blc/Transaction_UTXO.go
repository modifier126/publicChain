package blc

type UTXO struct {
	TxHash []byte
	Index  int
	Out    *TxOutput
}
