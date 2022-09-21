package blc

func (cli *CLI) createBlockChainWithGenesis(address string, nodeId string) {
	bc := CreateBlockChainWithGenesis(address, nodeId)
	defer bc.DB.Close()
	utxoSet := &UTXOSet{bc}
	utxoSet.ResetUTXOSet()
}
