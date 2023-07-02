package utils

import (
	"encoding/json"
	"io/ioutil"

	"github.com/pingtouge2008/zinx/ziface"
)

type GlobalObj struct {
	TcpServer        ziface.IServer
	Host             string
	TcpPort          int
	Name             string
	Version          string
	MaxPacketSize    uint32
	MaxConn          int
	WorkerPoolSize   uint32 // 业务线程池worker个数
	MaxWorkerTaskLen uint32 // 每个worker对应的消息队列的长度
	ConfFilePath     string
	MaxMsgChanLen    uint32
}

var GlobalObject *GlobalObj

func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile(g.ConfFilePath)
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(data, &GlobalObject); err != nil {
		panic(err)
	}
}

func init() {
	GlobalObject = &GlobalObj{
		Name:             "ZinxServerApp",
		Version:          "v0.4",
		Host:             "0.0.0.0",
		TcpPort:          6868,
		MaxConn:          12000,
		MaxPacketSize:    4096,
		ConfFilePath:     "conf/zinx.json",
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 1024,
		MaxMsgChanLen:    1024,
	}

	GlobalObject.Reload()
}
