package tests

import (
	"fmt"
	"net"
	"testing"
	"time"
)

func TestNet(t *testing.T) {
	l, _ := net.Listen("tcp", ":8888")
	go func() {
		for {
			conn, _ := l.Accept()
			fmt.Println("server", conn.LocalAddr().String(), conn.RemoteAddr().String())
		}
	}()
	time.Sleep(2 * time.Second)
	conn, _ := net.Dial("tcp", ":8888")

	fmt.Println(conn.LocalAddr().String(), conn.RemoteAddr().String())
	time.Sleep(10 * time.Second)
}
