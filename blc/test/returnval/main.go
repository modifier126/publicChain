package main

import "fmt"

func incr(a int) (b int) {
	defer func() {
		a++
		b++
	}()
	a++
	b = a

	return b
}

func main() {

	var a, b int

	b = incr(a)

	fmt.Printf("b=%d\n", b)

	fmt.Println("main...")
}
