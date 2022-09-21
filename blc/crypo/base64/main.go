package main

import (
	"encoding/base64"
	"fmt"
	"log"
)

func main() {
	str := base64.StdEncoding.EncodeToString([]byte("Hello,世界"))
	fmt.Printf("%s\n", str)

	bytes, err := base64.StdEncoding.DecodeString(str)

	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("%s\n", bytes)
}
