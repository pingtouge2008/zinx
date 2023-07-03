package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/pingtouge2008/zinx/utils"
	"github.com/pingtouge2008/zinx/ziface"
)

type Connection struct {
	TcpServer    *Server
	Conn         *net.TCPConn
	ConnID       uint32
	isClosed     bool
	ExitBuffChan chan bool
	MsgHandler   ziface.IMsgHandler
	msgChan      chan []byte // 用于读写两个Goroutine的消息通道(无缓冲)
	msgBufChan   chan []byte // 用于读写两个Goroutine的消息通道(有缓冲)
	properties   map[string]interface{}
	propertyLock sync.RWMutex
}

func NewConnection(server *Server, conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandler) ziface.IConnection {

	c := &Connection{
		TcpServer:    server,
		Conn:         conn,
		ConnID:       connID,
		isClosed:     false,
		ExitBuffChan: make(chan bool, 1),
		MsgHandler:   msgHandler,
		msgChan:      make(chan []byte),
		msgBufChan:   make(chan []byte, utils.GlobalObject.MaxMsgChanLen),
		properties:   make(map[string]interface{}),
	}
	c.TcpServer.GetConnManager().Add(c)

	return c

}

func (c *Connection) Start() {
	go c.StartReader()
	go c.StartWriter()
	c.TcpServer.CallOnConnStart(c)
}

func (c *Connection) Stop() {
	if c.isClosed {
		return
	}
	c.isClosed = true

	c.TcpServer.CallOnConnStop(c)

	c.Conn.Close()

	c.ExitBuffChan <- true

	c.TcpServer.GetConnManager().Remove(c)

	close(c.ExitBuffChan)
	close(c.msgChan)
	close(c.msgBufChan)
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running")
	defer fmt.Println(c.RemoteAddr(), "[conn Reader exit]")
	defer c.Stop()

	for {

		dp := NewDataPack()
		headData := make([]byte, dp.GetHeadLen())

		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head err ", err)
			c.ExitBuffChan <- true
			return
		}

		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("Unpack err ", err)
			c.ExitBuffChan <- true
			return
		}

		var data []byte

		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data err ", err)
				c.ExitBuffChan <- true
				return
			}
		}
		msg.SetData(data)
		req := Request{
			conn: c,
			msg:  msg,
		}
		if utils.GlobalObject.WorkerPoolSize > 0 {
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			go c.MsgHandler.DoMsgHandler(&req)
		}

	}
}

func (c *Connection) StartWriter() {
	fmt.Println("Writer Goroutine is running")

	defer fmt.Println(c.RemoteAddr().String(), "[conn Writer exit!]")

	for {
		select {
		case data := <-c.msgChan:
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("send data err", err, "Conn Writer exit")
				return
			}
		case data, ok := <-c.msgBufChan:
			if ok {
				if _, err := c.Conn.Write(data); err != nil {
					fmt.Println("send buf data err", err, "Conn Writer exit")
					return
				}
			} else {
				return
			}
		case <-c.ExitBuffChan:
			return
		}
	}

}

func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed {
		fmt.Println("Connection already closed when sending msg")
		return errors.New("Connection already closed")
	}

	dp := NewDataPack()
	msg, err := dp.Pack(NewMsgPacket(msgId, data))
	if err != nil {
		fmt.Println("Pack data err", err, "msgId", msgId)
		return errors.New("Pack data err")
	}
	c.msgChan <- msg
	return nil
}

func (c *Connection) SendBufMsg(msgId uint32, data []byte) error {
	if c.isClosed {
		fmt.Println("Connection already closed when sending msg")
		return errors.New("Connection already closed")
	}

	dp := NewDataPack()

	msg, err := dp.Pack(NewMsgPacket(msgId, data))
	if err != nil {
		fmt.Println("Pack data err", err, "msgId", msgId)
		return errors.New("Pack data err")
	}
	c.msgBufChan <- msg
	return nil
}

func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	c.properties[key] = value
}

func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()
	if value, ok := c.properties[key]; ok {
		return value, nil
	}
	return nil, errors.New("no property found")
}

func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	delete(c.properties, key)
}
