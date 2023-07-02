package znet

import (
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/pingtouge2008/zinx/ziface"
)

type Connection struct {
	Conn         *net.TCPConn
	ConnID       uint32
	isClosed     bool
	ExitBuffChan chan bool
	MsgHandler   ziface.IMsgHandler
	msgChan      chan []byte // 用于读写两个Goroutine的消息通道
}

func NewConnection(conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandler) ziface.IConnection {

	c := &Connection{
		Conn:         conn,
		ConnID:       connID,
		isClosed:     false,
		ExitBuffChan: make(chan bool, 1),
		MsgHandler:   msgHandler,
		msgChan:      make(chan []byte),
	}

	return c

}

func (c *Connection) Start() {

	go c.StartReader()
	go c.StartWriter()
	for {
		select {
		case <-c.ExitBuffChan:
			return
		}
	}
}

func (c *Connection) Stop() {
	if c.isClosed {
		return
	}
	c.isClosed = true

	c.Conn.Close()

	c.ExitBuffChan <- true

	close(c.ExitBuffChan)
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
	defer fmt.Printf("%s conn reader exit", c.RemoteAddr())
	defer c.Stop()

	for {

		dp := NewDataPack()
		headData := make([]byte, dp.GetHeadLen())

		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head err ", err)
			c.ExitBuffChan <- true
			continue
		}

		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("Unpack err ", err)
			c.ExitBuffChan <- true
			continue
		}

		var data []byte

		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data err ", err)
				c.ExitBuffChan <- true
				continue
			}
		}
		msg.SetData(data)
		req := Request{
			conn: c,
			msg:  msg,
		}
		go c.MsgHandler.DoMsgHandler(&req)

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
