package main

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/pingtouge2008/zinx/znet"
)

func main() {
	time.Sleep(2 * time.Second)

	conn, _ := net.Dial("tcp4", "127.0.0.1:6868")
	for {

		dp := znet.NewDataPack()
		bytesToBeSent, _ := dp.Pack(znet.NewMsgPacket(0, []byte("zinx v0.5 message sent by client")))

		conn.Write(bytesToBeSent)
		headData := make([]byte, dp.GetHeadLen())

		io.ReadFull(conn, headData)
		msgHead, _ := dp.Unpack(headData)
		if msgHead.GetDataLen() > 0 {
			msg := msgHead.(*znet.Message)

			msg.Data = make([]byte, msg.GetDataLen())
			io.ReadFull(conn, msg.Data)
			fmt.Println("recv from server:", string(msg.Data))
		}
		time.Sleep(2 * time.Second)
	}
}
