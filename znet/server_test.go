package znet

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/pingtouge2008/zinx/ziface"
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
	s := NewServer("v0.1")

	go ClientTest()

	s.Serve()
}

// server test v0.3
type PingRouter struct {
	BaseRouter
}

func (p *PingRouter) PreHandle(req ziface.IRequest) {
	// fmt.Println("PingRouter PreHandle")
	_, err := req.GetConnection().GetTCPConnection().Write([]byte("before ping...\n"))
	if err != nil {
		fmt.Println("PingRouter PreHandle err", err)
	}
}

func (p *PingRouter) Handle(req ziface.IRequest) {
	// fmt.Println("PingRouter Handle")
	_, err := req.GetConnection().GetTCPConnection().Write([]byte("ping...\n"))
	if err != nil {
		fmt.Println("PingRouter Handle err", err)
	}
}

func (p *PingRouter) PostHandle(req ziface.IRequest) {
	// fmt.Println("PingRouter PostHandle")
	_, err := req.GetConnection().GetTCPConnection().Write([]byte("post ping...\n"))
	if err != nil {
		fmt.Println("PingRouter PostHandle err", err)
	}
}

func TestServerV0_3(t *testing.T) {
	s := NewServer("v1.3")
	s.AddRouter(&PingRouter{})
	go ClientTest()
	s.Serve()
}
