package znet

import (
	"fmt"
	"net"
	"testing"
	"time"
)

func ClientTest() {
	time.Sleep(3 * time.Second)

	conn, _ := net.Dial("tcp4", "127.0.0.1:6868")
	for {
		conn.Write([]byte("hello\n"))
		buf := make([]byte, 512)
		cnt, _ := conn.Read(buf)
		fmt.Print(string(buf[:cnt]))
		time.Sleep(3 * time.Second)
	}
}

func TestServerV0_1(t *testing.T) {
	// s := NewServer("v0.1")
	s := NewServer()

	go ClientTest()

	s.Serve()
}
