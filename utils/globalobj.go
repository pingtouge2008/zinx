package utils

import (
	"encoding/json"
	"io/ioutil"

	"github.com/pingtouge2008/zinx/ziface"
)

type GlobalObj struct {
	TcpServer     ziface.IServer
	Host          string
	TcpPort       int
	Name          string
	Version       string
	MaxPacketSize uint32
	MaxConn       int
}

var GlobalObject *GlobalObj

func (*GlobalObj) Reload() {
	data, err := ioutil.ReadFile("conf/zinx.json")
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(data, &GlobalObject); err != nil {
		panic(err)
	}
}

func init() {
	GlobalObject = &GlobalObj{
		Host:          "0.0.0.0",
		TcpPort:       6868,
		Name:          "ZinxServerApp",
		Version:       "v0.4",
		MaxPacketSize: 4096,
		MaxConn:       12000,
	}

	GlobalObject.Reload()
}