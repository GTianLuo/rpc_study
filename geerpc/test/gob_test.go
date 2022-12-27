package test

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"testing"
)

type S struct {
	A int
	B int
}

func TestGob(t *testing.T) {
	s := S{A: 1, B: 2}
	var buf bytes.Buffer
	eg := gob.NewEncoder(&buf)
	eg.Encode(s)
	fmt.Println(buf.String())

	dg := gob.NewDecoder(&buf)
	//p := S{}
	var p *S
	//dg.Decode()
	err := dg.Decode(&p)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(p)
}

func TestSlice(t *testing.T) {
	s := make([]int, 0, 10)
	fmt.Println("len:", len(s))
	fmt.Println("cap:", cap(s))
}
