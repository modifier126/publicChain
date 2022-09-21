package blc

func (cli *CLI) printChain(nodeId string) {
	if DBExists(nodeId) {
		blockChain := BlockChainObj(nodeId)
		blockChain.PrintChain()
		defer blockChain.DB.Close()
	}
}
