package blc

import "fmt"

func (cli *CLI) GetWalletlist(nodeId string) {
	fmt.Println("打印所有钱包地址:")
	ws, _ := NewWalletList(nodeId)
	for address := range ws.Wallets {
		fmt.Println(address)
	}
}
