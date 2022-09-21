package blc

import "fmt"

func (cli *CLI) CreateWallet(nodeId string) {
	ws, _ := NewWalletList(nodeId)
	ws.CreateNewWallet(nodeId)
	fmt.Println(len(ws.Wallets))
}
