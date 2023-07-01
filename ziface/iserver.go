package ziface

type IServer interface {
	Start()
	Stop()
	Serve()
	AddRouter(route IRouter)
}
