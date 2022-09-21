package blc

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"net"
)

var (
	KnownNodes      = []string{"localhost:3000"} // 默认主节点
	nodeAddress     string                       // 开启服务的节点地址
	mineAddress     string
	blocksInTransit = [][]byte{}
	memoryPool      = make(map[string]Transaction)
)

func StartServer(nodeId, minerAdd string) {
	nodeAddress = fmt.Sprintf("localhost:%s", nodeId)
	mineAddress = minerAdd
	fmt.Printf("开启节点服务地址:%s\n", nodeAddress)
	ln, err := net.Listen(PROTOCOL, nodeAddress)
	if err != nil {
		log.Panic(err)
	}
	defer ln.Close()

	bc := BlockChainObj(nodeId)
	defer bc.DB.Close()

	if nodeAddress != KnownNodes[0] {
		SendVersion(KnownNodes[0], bc)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Panic(err)
		}
		go HandleConnection(conn, bc)
	}

}

func CmdToBytes(cmd string) []byte {
	var bytes [COMMAND_LENGTH]byte
	for i, c := range cmd {
		bytes[i] = byte(c)
	}
	return bytes[:]
}

func BytesToCmd(bytes []byte) string {
	var cmd []byte

	for _, b := range bytes {
		if b != 0x0 {
			cmd = append(cmd, b)
		}
	}
	return string(cmd)
}

func GobEncode(data interface{}) []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(data)
	if err != nil {
		log.Panic(err)
	}
	return buf.Bytes()
}

func NodeIsKnown(addr string) bool {
	for _, node := range KnownNodes {
		if addr == node {
			return true
		}
	}
	return false
}

func RequestBlocks() {
	for _, node := range KnownNodes {
		SendGetBlocks(node)
	}
}
