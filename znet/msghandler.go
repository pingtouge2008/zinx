package znet

import (
	"fmt"

	"github.com/pingtouge2008/zinx/ziface"
)

type MsgHandler struct {
	Apis map[uint32]ziface.IRouter
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		Apis: make(map[uint32]ziface.IRouter),
	}
}

func (mh *MsgHandler) DoMsgHandler(req ziface.IRequest) {
	handler, exist := mh.Apis[req.GetMsgId()]
	if !exist {
		fmt.Println("Handle for msgId = ", req.GetMsgId(), "not found")
	}
	handler.PreHandle(req)
	handler.Handle(req)
	handler.PostHandle(req)
}

func (mh *MsgHandler) AddRouter(msgId uint32, router ziface.IRouter) {
	if _, exist := mh.Apis[msgId]; exist {
		fmt.Println("handler for msgId = ", msgId, "already exists!")
		return
	}
	mh.Apis[msgId] = router
	fmt.Println("add handler for msgId = ", msgId)
}
