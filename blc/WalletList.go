package blc

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

const walletsFile = "wallets_%s.dat"

type WalletList struct {
	Wallets map[string]*Wallet
}

func NewWalletList(nodeId string) (*WalletList, error) {
	walletsFile := fmt.Sprintf(walletsFile, nodeId)
	if _, err := os.Stat(walletsFile); os.IsNotExist(err) {
		ws := &WalletList{}
		ws.Wallets = make(map[string]*Wallet)
		return ws, err
	}

	fileContent, err := ioutil.ReadFile(walletsFile)
	if err != nil {
		log.Panic(err)
	}

	var ws WalletList
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&ws)
	if err != nil {
		log.Panic(err)
	}
	return &ws, nil
}

func (ws *WalletList) CreateNewWallet(nodeId string) {
	w := NewWallet()
	fmt.Printf("address:%s\n", w.GetAddress())
	ws.Wallets[string(w.GetAddress())] = w
	// 保存钱包
	ws.SaveWallets(nodeId)
}

func (ws *WalletList) SaveWallets(nodeId string) {
	var content bytes.Buffer

	gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(ws)
	if err != nil {
		log.Panic(err)
	}
	walletsFile := fmt.Sprintf(walletsFile, nodeId)
	err = ioutil.WriteFile(walletsFile, content.Bytes(), 0644)
	if err != nil {
		log.Panic(err)
	}
}

func (ws *WalletList) LoadFromFile(nodeId string) error {
	walletsFile := fmt.Sprintf(walletsFile, nodeId)
	if _, err := os.Stat(walletsFile); os.IsNotExist(err) {
		return err
	}
	fileContent, err := ioutil.ReadFile(walletsFile)
	if err != nil {
		log.Panic(err)
	}
	var wallets WalletList
	gob.Register(elliptic.P256())
	decorder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decorder.Decode(&wallets)
	if err != nil {
		log.Panic(err)
	}
	ws.Wallets = wallets.Wallets
	return nil
}
