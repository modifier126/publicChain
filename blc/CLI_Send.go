package blc

import (
	"fmt"
	"os"
)

func (cli *CLI) send(from []string, to []string, amount []string, nodeId string, mineNow bool) {
	if !DBExists(nodeId) {
		fmt.Println("数据对象不存在")
		os.Exit(1)
	}

	bc := BlockChainObj(nodeId)
	defer bc.DB.Close()
	txs := bc.NewTrans(from, to, amount, nodeId)
	if mineNow {
		bc.MineNewBlock(txs)
		utxoSet := UTXOSet{bc}
		utxoSet.Update()
	} else {
		SendTx(KnownNodes[0], Transaction{})
		fmt.Println("send tx")
	}
	fmt.Println("Success!")
}
