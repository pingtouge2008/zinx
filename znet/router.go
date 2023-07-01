package znet

import "github.com/pingtouge2008/zinx/ziface"

type BaseRouter struct{}

func (br *BaseRouter) PreHandle(req *ziface.IRequest)  {}
func (br *BaseRouter) Handle(req *ziface.IRequest)     {}
func (br *BaseRouter) PostHandle(req *ziface.IRequest) {}
