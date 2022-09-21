package blc

import "fmt"

func (cli *CLI) getBalance(address string, nodeId string) {
	bc := BlockChainObj(nodeId)
	defer bc.DB.Close()
	utxoSet := &UTXOSet{bc}
	amount := utxoSet.GetBalance(address)
	fmt.Printf("%s一共有%d个token\n", address, amount)
}
