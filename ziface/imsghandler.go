package ziface

type IMsgHandler interface {
	DoMsgHandler(req IRequest)
	AddRouter(msgId uint32, router IRouter)
}
