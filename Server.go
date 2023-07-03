package main

import (
	"fmt"

	"github.com/pingtouge2008/zinx/ziface"
	"github.com/pingtouge2008/zinx/znet"
)

func main() {
	s := znet.NewServer()
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloZinxRouter{})
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)
	s.Serve()
}

type PingRouter struct {
	znet.BaseRouter
}

func (p *PingRouter) Handle(req ziface.IRequest) {
	fmt.Println("PingRouter Handle")
	fmt.Println("recv from client, msgId:", req.GetMsgId(), "data: ", string(req.GetData()))
	err := req.GetConnection().SendBufMsg(0, []byte("ping...\n"))
	if err != nil {
		fmt.Println("PingRouter Handle err", err)
	}
}

type HelloZinxRouter struct {
	znet.BaseRouter
}

func (p *HelloZinxRouter) Handle(req ziface.IRequest) {
	fmt.Println("HelloZinxRouter Handle")
	fmt.Println("recv from client, msgId:", req.GetMsgId(), "data: ", string(req.GetData()))
	err := req.GetConnection().SendBufMsg(1, []byte("hello zinx...\n"))
	if err != nil {
		fmt.Println("HelloZinxRouter Handle err", err)
	}
}

func DoConnectionBegin(conn ziface.IConnection) {
	fmt.Println("DoConnectionBegin", conn)
	conn.SetProperty("name", "connection-ptg")
	conn.SendMsg(2, []byte("DoConnectionBegin"))
}

func DoConnectionLost(conn ziface.IConnection) {
	fmt.Println("DoConnectionLost is called")
	p, _ := conn.GetProperty("name")
	fmt.Println(p)
}
