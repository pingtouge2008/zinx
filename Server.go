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

func (p *PingRouter) Handle(req ziface.IRequest) {
	fmt.Println("PingRouter Handle")
	fmt.Println("recv from client, msgId:", req.GetMsgId(), "data: ", string(req.GetData()))
	err := req.GetConnection().Send(1, []byte("ping...\n"))
	if err != nil {
		fmt.Println("PingRouter Handle err", err)
	}
}
