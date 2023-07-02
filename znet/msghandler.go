package znet

import (
	"fmt"

	"github.com/pingtouge2008/zinx/utils"
	"github.com/pingtouge2008/zinx/ziface"
)

type MsgHandler struct {
	Apis           map[uint32]ziface.IRouter
	WorkerPoolSize uint32
	TaskQueue      []chan ziface.IRequest
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		Apis:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
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

func (mh *MsgHandler) StartWorkerPool() {
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}

func (mh *MsgHandler) StartOneWorker(workerId int, taskQueue chan ziface.IRequest) {
	fmt.Println("Worker ID = ", workerId, " is started.")
	for {
		select {
		case req := <-taskQueue:
			mh.DoMsgHandler(req)
		}
	}
}

func (mh *MsgHandler) SendMsgToTaskQueue(request ziface.IRequest) {
	workerId := request.GetConnection().GetConnID() % mh.WorkerPoolSize

	fmt.Println("Add ConnID = ", request.GetConnection().GetConnID(), " request msgID = ",
		request.GetMsgId(), " to workerID = ", workerId)

	mh.TaskQueue[workerId] <- request
}
