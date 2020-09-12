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

type Test struct {
	i int
	m int
}

func (t *Test) test() {
	t.m += 3
}
func (t *Test) test1() int {
	return t.i
}

func main() {
	t := Test{1, 2}
	for i := 1; i < 1000; i++ {
		go t.test()
		x := t.test1()
		fmt.Println(x)
	}

}
