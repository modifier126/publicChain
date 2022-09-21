package blc

import "fmt"

func (cli *CLI) TestMethod(nodeId string) {
	bc := BlockChainObj(nodeId)
	utoxMap := bc.FindUtoxMap()
	fmt.Println(utoxMap)
	utxoSet := &UTXOSet{bc}
	utxoSet.ResetUTXOSet()
}
