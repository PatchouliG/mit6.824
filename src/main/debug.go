package main

import (
	"fmt"
)

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

type Test struct {
	i int
	m int
}

func testtest(a interface{}) {
	s, ok := a.(Q)
	if ok {
		fmt.Println(s.Name)
	}
	m, ok := a.(Test)
	if ok {
		fmt.Println(m.i)
	}
}

func (t *Test) test() {
	t.m += 3
}
func (t *Test) test1() int {
	return t.i
}

func main() {
	a := 2
	for {
		switch a {
		case 1:
			fmt.Print(1)
			return
		case 2:
		case 3:
			fmt.Print(3)
			return
		}
		fmt.Println(2)
	}
}
