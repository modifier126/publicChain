package main

import (
	"fmt"
	"strconv"
)

func main() {

	ch := make(chan bool, 1)

	var m = make(map[string]string)

	go func() {
		for i := 0; i < 5; i++ {
			k := fmt.Sprintf("%s%d", "k", i)
			m[k] = strconv.Itoa(i)
		}
		ch <- true
	}()

	fmt.Println("arrived here")

	flag := <-ch

	for {
		if flag {
			break
		}
	}

	for k, v := range m {
		fmt.Printf("k=%s v=%s\n", k, v)
	}

}
