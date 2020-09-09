package main

import (
	"fmt"
	"sort"
	"time"
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
}

func test(ouput chan int) {
	time.Sleep(time.Second * 5)
	ouput <- 1
	fmt.Printf("ok")
}

func main() {
	a := []int{1, 2, 1, 2, 3, 5, 3, 4}

	sort.Ints(a)

}
