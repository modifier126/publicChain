package main

import (
	"fmt"
	"reflect"
)

type A struct {
	name string
}

func (a A) Name() {
	fmt.Printf("Hi %s\n", a.name)
}

func NameOfA(a A) {
	fmt.Printf("Hi %s\n", a.name)
}

func main() {
	a := A{name: "zhangsan"}

	a.Name()
	NameOfA(a)

	t1 := reflect.TypeOf(A.Name)
	t2 := reflect.TypeOf(NameOfA)

	fmt.Println(t1 == t2)

}
