package main

import "fmt"

type A struct {
	Name string
}

func (a A) GetName() string {
	return a.Name
}

func (a *A) SetName() {
	a.Name = "lisi"
}

func main() {
	a := A{Name: "zhangsan"}

	pa := &a

	res1 := pa.GetName() // pa.GetName() 等价 (*pa).GetName()
	res2 := (*pa).GetName()

	fmt.Println("res1=", res1)
	fmt.Println("res2=", res2)

	a.SetName()

	(&A{Name: "wangwu"}).SetName()

	fmt.Println("res2=", a.GetName())
}
