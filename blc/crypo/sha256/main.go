package main

import (
	"crypto/sha256"
	"fmt"
	"log"
)

func main() {

	hasher := sha256.New()

	_, err := hasher.Write([]byte("modifier126"))

	if err != nil {
		log.Panic(err)
	}

	bytes := hasher.Sum(nil)
	fmt.Printf("%x\n", bytes)
}
