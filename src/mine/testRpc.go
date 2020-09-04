package main

import (
	"bytes"
	"fmt"
)
import "../labgob"

type test struct {
	A int
}

func main() {
	w := new(bytes.Buffer)
	t := test{4}
	s := test{2}
	m := test{}
	e := labgob.NewEncoder(w)
	d := labgob.NewDecoder(w)
	e.Encode(&t)
	e.Encode(&s)
	data := w.Bytes()
	println(data)
	d.Decode(&m)
	fmt.Println(m)
	d.Decode(&m)
	fmt.Println(m)
}
