package main

import "fmt"

// func A1(a *int) {
// 	fmt.Printf("a=%d\n", *a)
// }

// func B1() {
// 	a, b := 1, 2
// 	defer A1(&a)
// 	a = a + b
// 	fmt.Printf("a=%d b=%d\n", a, b)
// }

// func A() {
// 	a, b := 1, 2
// 	defer func(b int) {
// 		a = a + b
// 		fmt.Printf("a=%d b=%d\n", a, b)
// 	}(b)

// 	a = a + b
// 	fmt.Printf("a=%d b=%d\n", a, b)
// }

// func A() {
// 	defer A1()
// 	defer A2()
// }

// func A1() {
// 	fmt.Println("A1")
// }

// func A2() {
// 	defer B1()
// 	defer B2()
// 	fmt.Println("A2")
// }

// func B1() {
// 	fmt.Println("B1")
// }

// func B2() {
// 	fmt.Println("B2")
// }

func A() {
	defer A1()
	defer A2()
	panic("panic A")
}

func A1() {
	fmt.Println("A1")
}

func A2() {

	r := recover()

	fmt.Printf("r=%s\n", r)

	panic("panic A2")
}

func main() {
	//B1()
	//A()
	A()
}
