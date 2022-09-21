package blc

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net"
)

func HandleConnection(conn net.Conn, bc *BlockChain) {
	req, err := ioutil.ReadAll(conn)
	defer conn.Close()
	if err != nil {
		log.Panic(err)
	}

	command := BytesToCmd(req[:COMMAND_LENGTH])
	fmt.Printf("Receive %s command\n", command)

	switch command {
	case ADDR:
		HandleAddr(req)
	case BLOCK:
		HandleBlock(req, bc)
	case INV:
		HandleInv(req, bc)
	case GETBLOCKS:
		HandleGetBlocks(req, bc)
	case GETDATA:
		HandleGetData(req, bc)
	case TX:
		HandleTx(req, bc)
	case VERSION:
		HandleVersion(req, bc)
	default:
		fmt.Println("Unknown command")
	}

}

func HandleAddr(req []byte) {
	var buff bytes.Buffer
	var payload Addr

	buff.Write(req[COMMAND_LENGTH:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)

	}

	KnownNodes = append(KnownNodes, payload.AddrList...)
	fmt.Printf("there are %d known nodes\n", len(KnownNodes))

	RequestBlocks()
}

func HandleBlock(req []byte, bc *BlockChain) {
	var buff bytes.Buffer
	var payload BlockData

	buff.Write(req[COMMAND_LENGTH:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	blockData := payload.Block
	block := DeSerializeBlock(blockData)

	fmt.Println("Recevied a new block!")
	bc.SynBlock(block)

	if len(blocksInTransit) > 0 {
		blockHash := blocksInTransit[0]
		SendGetData(payload.AddrFrom, BLOCK, blockHash)
		blocksInTransit = blocksInTransit[1:]
	} else {
		utxoSet := &UTXOSet{bc}
		utxoSet.ResetUTXOSet()
	}
}

func HandleInv(req []byte, bc *BlockChain) {
	var buf bytes.Buffer
	var payload Inv
	buf.Write(req[COMMAND_LENGTH:])
	dec := gob.NewDecoder(&buf)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("Recevied inventory with %d %s\n", len(payload.Items), payload.Type)

	if payload.Type == BLOCK {
		blocksInTransit = payload.Items
		blockHash := payload.Items[0]
		SendGetData(payload.AddrFrom, BLOCK, blockHash)
		newInTransit := [][]byte{}
		for _, b := range blocksInTransit {
			if !bytes.Equal(b, blockHash) {
				newInTransit = append(newInTransit, b)
			}
		}
		blocksInTransit = newInTransit

	} else if payload.Type == TX {
		txID := payload.Items[0]
		if memoryPool[hex.EncodeToString(txID)].TxHash == nil {
			SendGetData(payload.AddrFrom, TX, txID)
		}
	}

}

func HandleGetBlocks(req []byte, bc *BlockChain) {
	var buf bytes.Buffer
	var payload GetBlocks
	_, err := buf.Write(req[COMMAND_LENGTH:])
	if err != nil {
		log.Panic(err)
	}
	dec := gob.NewDecoder(&buf)
	err = dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}
	hashes := bc.GetBlockHashes()
	SendInv(payload.AddrFrom, BLOCK, hashes)
}

func HandleGetData(req []byte, bc *BlockChain) {
	var buf bytes.Buffer
	var payload GetData
	buf.Write(req[COMMAND_LENGTH:])
	dec := gob.NewDecoder(&buf)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	if payload.Type == BLOCK {
		block, err := bc.GetBlock(payload.ID)
		if err != nil {
			log.Panic(err)
		}
		SendBlockData(payload.AddrFrom, block)

	} else if payload.Type == TX {
		txId := hex.EncodeToString(payload.ID)
		tx := memoryPool[txId]
		SendTx(payload.AddrFrom, tx)
	}

}

func HandleTx(req []byte, bc *BlockChain) {
	var buf bytes.Buffer
	var payload Tx

	buf.Write(req[COMMAND_LENGTH:])
	dec := gob.NewDecoder(&buf)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}
	txData := payload.Transaction

	tx := *DeSerializeTransaction(txData)
	memoryPool[hex.EncodeToString(tx.TxHash)] = tx

	if nodeAddress == KnownNodes[0] {
		for _, node := range KnownNodes {
			if node != nodeAddress && node != payload.AddrFrom {
				SendInv(node, TX, [][]byte{tx.TxHash})
			}
		}
	} else {
		if len(memoryPool) >= 2 && len(mineAddress) > 0 {
			MineTx(bc)
		}
	}
}

func HandleVersion(req []byte, bc *BlockChain) {
	var buf bytes.Buffer
	var payload Version
	buf.Write(req[COMMAND_LENGTH:])
	dec := gob.NewDecoder(&buf)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	// 接受方区块高度
	bestHeight := bc.GetBestHeight()
	otherHeight := payload.BestHeight

	if bestHeight < otherHeight {
		SendGetBlocks(payload.AddrFrom)
	} else if bestHeight > otherHeight {
		SendVersion(payload.AddrFrom, bc)
	}

	if !NodeIsKnown(payload.AddrFrom) {
		KnownNodes = append(KnownNodes, payload.AddrFrom)
	}
}

func MineTx(bc *BlockChain) {

	var txs []*Transaction

	for txId := range memoryPool {
		fmt.Printf("tx: %s\n", memoryPool[txId].TxHash)
		tx := memoryPool[txId]

		if bc.VerifyTransaction(&tx, txs) {
			txs = append(txs, &tx)
		}
	}

	if len(txs) == 0 {
		fmt.Println("All Transactions are invalid")
		return
	}

	// 奖励
	tx := NewCoinBaseTrans(mineAddress)
	txs = append(txs, tx)

	newBlock := bc.MineNewBlock(txs)

	utxoSet := UTXOSet{bc}
	utxoSet.Update()
	fmt.Println("New Block mined")

	for _, tx := range txs {
		txId := hex.EncodeToString(tx.TxHash)
		delete(memoryPool, txId)
	}

	for _, node := range KnownNodes {
		if node != nodeAddress {
			SendInv(node, BLOCK, [][]byte{newBlock.Hash})
		}
	}

	if len(memoryPool) > 0 {
		MineTx(bc)
	}
}
