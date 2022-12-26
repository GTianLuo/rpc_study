package example_1

import (
	"fmt"
	"testing"
)

func TestServer1(t *testing.T) {
	myServer()
}

func TestClient(t *testing.T) {
	myClient()
}

func TestAsyncClient(t *testing.T) {
	asyncClient()
}

func TestChanel(t *testing.T) {
	ch := make(chan int, 1)
	ch <- 1
	close(ch)
	i, ok := <-ch
	fmt.Println(ok, "  ", i)
	i, ok = <-ch
	fmt.Println(ok, "  ", i)
}
