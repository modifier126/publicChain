package blc

import (
	"fmt"
	"log"
)

func (cli *CLI) StartNode(nodeId, minerAdd string) {

	if len(minerAdd) > 0 {
		if IsValidForAddress([]byte(minerAdd)) {
			fmt.Println("Mining is on. Address to receive rewards: ", minerAdd)
		} else {
			log.Panic("Wrong miner address!")
		}
	}

	StartServer(nodeId, minerAdd)
}
