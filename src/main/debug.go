package main

import (
	"fmt"
	"time"
)

func main() {
	c1 := make(chan int64)
	select {
	case res := <-c1:
		fmt.Println(res)
	case <-time.After(1 * time.Second):
		fmt.Println("timeout 1")
	}

}
