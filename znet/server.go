package znet

import (
	"errors"
	"fmt"
	"net"

	"github.com/pingtouge2008/zinx/utils"
	"github.com/pingtouge2008/zinx/ziface"
)

type Server struct {
	Name      string
	IPVersion string
	IP        string
	Port      int
	Router    ziface.IRouter
}

func CallBackToClient(conn *net.TCPConn, data []byte, cnt int) error {

	fmt.Println("[Conn Handle] CallBackToClient ...")

	if _, err := conn.Write(data[:cnt]); err != nil {
		fmt.Println("write back buf err", err)
		return errors.New("CallBackToClient error")
	}

	return nil
}

func NewServer() ziface.IServer {
	utils.GlobalObject.Reload()
	s := &Server{
		Name:      utils.GlobalObject.Name,
		IPVersion: "tcp4",
		IP:        utils.GlobalObject.Host,
		Port:      utils.GlobalObject.TcpPort,
		Router:    nil,
	}

	return s
}

func (s *Server) Start() {

	fmt.Printf("[Zinx] Version: %s, MaxConn: %d, MaxPacketSize: %d\n",
		utils.GlobalObject.Version, utils.GlobalObject.MaxConn, utils.GlobalObject.MaxPacketSize)
	go func() {
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
			dealConn := NewConnection(conn, cid, s.Router)
			cid++
			go dealConn.Start()
		}
	}()
}

func (s *Server) Stop() {

}

func (s *Server) Serve() {
	s.Start()

	select {}
}

func (s *Server) AddRouter(router ziface.IRouter) {
	s.Router = router

	fmt.Println("Add Router succ")
}
