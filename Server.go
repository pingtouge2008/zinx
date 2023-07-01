package main

import (
	"fmt"

	"github.com/pingtouge2008/zinx/ziface"
	"github.com/pingtouge2008/zinx/znet"
)

func main() {
	s := znet.NewServer()
	s.AddRouter(&PingRouter{})
	s.Serve()
}

type PingRouter struct {
	znet.BaseRouter
}

func (p *PingRouter) PreHandle(req ziface.IRequest) {
	fmt.Println("PingRouter PreHandle")
	_, err := req.GetConnection().GetTCPConnection().Write([]byte("before ping...\n"))
	if err != nil {
		fmt.Println("PingRouter PreHandle err", err)
	}
}

func (p *PingRouter) Handle(req ziface.IRequest) {
	fmt.Println("PingRouter Handle")
	_, err := req.GetConnection().GetTCPConnection().Write([]byte("ping...\n"))
	if err != nil {
		fmt.Println("PingRouter Handle err", err)
	}
}

func (p *PingRouter) PostHandle(req ziface.IRequest) {
	fmt.Println("PingRouter PostHandle")
	_, err := req.GetConnection().GetTCPConnection().Write([]byte("post ping...\n"))
	if err != nil {
		fmt.Println("PingRouter PostHandle err", err)
	}
}
