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

type Test struct {
	i int
}

func test(ouput chan int) {
	time.Sleep(time.Second * 5)
	ouput <- 1
	fmt.Printf("ok")
}

func main() {
	fmt.Println(time.Now().Nanosecond())
	fmt.Println(time.Now().Nanosecond())
	fmt.Println(time.Now().Nanosecond())
}
