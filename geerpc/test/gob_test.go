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
	err := eg.Encode(s)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(buf.String())

	dg := gob.NewDecoder(&buf)
	p := S{}
	err = dg.Decode(&p)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(p)
}
