package blc

import (
	"blockDemo/blc/utils"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
)

func SendVersion(toAddress string, bc *BlockChain) {
	bestHeight := bc.GetBestHeight()
	payload := utils.GodEncode(Version{NODE_VERSION, bestHeight, nodeAddress})
	request := append(CmdToBytes(VERSION), payload...)
	SendData(toAddress, request)
}

func SendData(addr string, data []byte) {
	fmt.Printf("发送数据到:%s\n", addr)
	conn, err := net.Dial(PROTOCOL, addr)
	if err != nil {
		fmt.Printf("%s addr not available", addr)
		var updateNodes []string
		for _, node := range KnownNodes {
			if addr != node {
				updateNodes = append(updateNodes, node)
			}
		}
		KnownNodes = updateNodes
		return
	}
	defer conn.Close()

	_, err = io.Copy(conn, bytes.NewReader(data))
	if err != nil {
		log.Panic(err)
	}
}

func SendGetBlocks(addr string) {
	payload := GobEncode(GetBlocks{nodeAddress})
	request := append(CmdToBytes(GETBLOCKS), payload...)
	SendData(addr, request)
}

func SendInv(addr, kind string, items [][]byte) {
	inventory := Inv{nodeAddress, kind, items}
	payload := GobEncode(inventory)
	request := append(CmdToBytes(INV), payload...)
	SendData(addr, request)
}

func SendGetData(addr, kind string, id []byte) {
	payload := GobEncode(GetData{nodeAddress, kind, id})
	request := append(CmdToBytes(GETDATA), payload...)
	SendData(addr, request)
}

func SendBlockData(addr string, b *Block) {
	payload := GobEncode(BlockData{nodeAddress, b.Serialize()})
	request := append(CmdToBytes(BLOCK), payload...)
	SendData(addr, request)
}

func SendTx(addr string, tx Transaction) {
	payload := GobEncode(Tx{nodeAddress, tx.Serialize()})
	request := append(CmdToBytes(TX), payload...)
	SendData(addr, request)
}
