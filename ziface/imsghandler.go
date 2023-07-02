package ziface

type IMsgHandler interface {
	DoMsgHandler(req IRequest)
	AddRouter(msgId uint32, router IRouter)
	StartWorkerPool()
	// 将消息交给TaskQueue, 由Worker处理
	SendMsgToTaskQueue(request IRequest)
}
