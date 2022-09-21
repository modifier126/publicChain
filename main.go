package main

import (
	"blockDemo/blc"
)

func main() {
	cli := blc.CLI{}
	cli.Run()

	a := make(map[string]string)

	a["a"] = "1"
}
