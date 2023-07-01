package znet

import (
	"fmt"
	"net"

	"github.com/pingtouge2008/zinx/ziface"
)

type Connection struct {
	Conn         *net.TCPConn
	ConnID       uint32
	isClosed     bool
	handleApi    ziface.HandFunc
	ExitBuffChan chan bool
}

func NewConnection(conn *net.TCPConn, connID uint32, handleApi ziface.HandFunc) ziface.IConnection {

	c := &Connection{
		Conn:         conn,
		ConnID:       connID,
		isClosed:     false,
		handleApi:    handleApi,
		ExitBuffChan: make(chan bool, 1),
	}

	return c

}

func (c *Connection) Start() {

	go c.StartReader()
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
		buf := make([]byte, 512)
		cnt, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("recv buf err ", err)
			c.ExitBuffChan <- true
			continue
		}

		if err := c.handleApi(c.Conn, buf, cnt); err != nil {
			fmt.Printf("%d handle err ", c.ConnID)
			c.ExitBuffChan <- true
			return
		}
	}
}
