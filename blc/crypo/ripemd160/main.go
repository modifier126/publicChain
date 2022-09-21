package main

import (
	"fmt"
	"log"

	"golang.org/x/crypto/ripemd160"
)

func main() {
	hasher := ripemd160.New()

	_, err := hasher.Write([]byte("modifier126"))

	if err != nil {
		log.Panic(err)
	}

	bytes := hasher.Sum(nil)
	fmt.Printf("%x\n", bytes)
}
