package main

import "fmt"

func create() (fs [2]func()) {

	for i := 0; i < 2; i++ {
		fs[i] = func() {
			fmt.Printf("i=%d\n", i)
		}
	}
	return
}

func main() {

	f1 := create()

	for i := 0; i < len(f1); i++ {
		f1[i]()
	}

	fmt.Println("main...")
}