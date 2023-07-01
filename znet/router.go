package znet

import "github.com/pingtouge2008/zinx/ziface"

type BaseRouter struct{}

func (b *BaseRouter) PreHandle(req ziface.IRequest)  {}
func (b *BaseRouter) Handle(req ziface.IRequest)     {}
func (b *BaseRouter) PostHandle(req ziface.IRequest) {}
