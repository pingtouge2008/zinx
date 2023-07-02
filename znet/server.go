package znet

import (
	"fmt"
	"net"

	"github.com/pingtouge2008/zinx/utils"
	"github.com/pingtouge2008/zinx/ziface"
)

type Server struct {
	Name        string
	IPVersion   string
	IP          string
	Port        int
	msgHandler  ziface.IMsgHandler
	ConnMgr     ziface.IConnectionManager
	OnConnStart func(ziface.IConnection)
	OnConnStop  func(ziface.IConnection)
}

func NewServer() ziface.IServer {
	utils.GlobalObject.Reload()
	s := &Server{
		Name:       utils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		msgHandler: NewMsgHandler(),
		ConnMgr:    NewConnectionManager(),
	}

	return s
}

func (s *Server) Start() {

	fmt.Printf("[Zinx] Version: %s, MaxConn: %d, MaxPacketSize: %d\n",
		utils.GlobalObject.Version, utils.GlobalObject.MaxConn, utils.GlobalObject.MaxPacketSize)
	go func() {
		s.msgHandler.StartWorkerPool()
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp add err: ", err)
			return
		}

		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen tcp err: ", err)
			return
		}
		fmt.Printf("[START] server is listening at %s:%d\n", s.IP, s.Port)

		// TODO
		var cid uint32
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("accept err: ", err)
				continue
			}
			if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
				fmt.Println("Too many connections")
				conn.Close()
				continue
			}
			dealConn := NewConnection(s, conn, cid, s.msgHandler)
			cid++
			go dealConn.Start()
		}
	}()
}

func (s *Server) Stop() {
	fmt.Print("[STOP] server is stopping ...")
	s.ConnMgr.ClearConn()
}

func (s *Server) Serve() {
	s.Start()

	select {}
}

func (s *Server) AddRouter(msgId uint32, router ziface.IRouter) {
	s.msgHandler.AddRouter(msgId, router)
}

func (s *Server) GetConnManager() ziface.IConnectionManager {
	return s.ConnMgr
}

func (s *Server) SetOnConnStart(onConnStart func(ziface.IConnection)) {
	s.OnConnStart = onConnStart
}

func (s *Server) SetOnConnStop(onConnStop func(ziface.IConnection)) {
	s.OnConnStop = onConnStop
}

func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("---> CallOnConnStart")
		s.OnConnStart(conn)
	}
}

func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("---> CallOnConnStop")
		s.OnConnStop(conn)
	}
}
