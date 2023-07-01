package ziface

// IRequest 把客户端请求的连接和数据包装起来
type IRequest interface {
	GetConnection() IConnection
	GetData() []byte
}
