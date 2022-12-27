package test

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

func doBadSomething(ch chan bool) {
	//time.Sleep(time.Second)
	select {
	case ch <- true:
	default:
		return
	}
	time.Sleep(time.Second)
	fmt.Println("work 2")

}

func doBad() {
	ch := make(chan bool)
	go doBadSomething(ch)
	select {
	case <-time.After(1):
		fmt.Println("timeout")
	case <-ch:
		fmt.Println("done")
	}
	fmt.Println("do bad end")
}

func TestBad(t *testing.T) {
	for i := 0; i < 10; i++ {
		doBad()
	}
	time.Sleep(time.Second * 2)
	fmt.Println(runtime.NumGoroutine())
}
