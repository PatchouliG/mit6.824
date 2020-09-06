package main

import (
	"fmt"
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

type Test int

func test(ouput chan int) {
	time.Sleep(time.Second * 5)
	ouput <- 1
	fmt.Printf("ok")
}

func main() {
	for i := 1; i < 10; i++ {
		go func() {
			fmt.Print(i)
		}()
	}
	time.Sleep(time.Second)
}
