package main

import "fmt"

type P struct {
	X, Y, Z int
	Name    string
}

type Q struct {
	X, Y *int32
	Name string
}

type A interface {
}

type Test int

func main() {
	var a Test
	a = 3
	fmt.Println(int(a) == 4)

}
