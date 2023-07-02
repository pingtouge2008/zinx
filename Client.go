package main

import (
	"fmt"
	"io"
	"math/rand"
	"net"
	"time"

	"github.com/pingtouge2008/zinx/znet"
)

func main() {
	// time.Sleep(2 * time.Second)

	conn, _ := net.Dial("tcp4", "127.0.0.1:6868")
	for {

		dp := znet.NewDataPack()
		msgId := uint32(rand.Intn(6)) % 2
		bytesToBeSent, _ := dp.Pack(znet.NewMsgPacket(msgId, []byte(fmt.Sprintf("这条是msgId=%d的消息", msgId))))
		_, err := conn.Write(bytesToBeSent)
		if err != nil {
			panic(err)
		}

		headData := make([]byte, dp.GetHeadLen())

		_, err = io.ReadFull(conn, headData)
		if err != nil {
			panic(err)
		}
		msgHead, _ := dp.Unpack(headData)
		if msgHead.GetDataLen() > 0 {
			msg := msgHead.(*znet.Message)

			msg.Data = make([]byte, msg.GetDataLen())
			io.ReadFull(conn, msg.Data)
			fmt.Println("recv from server:", string(msg.Data))
		}
		time.Sleep(8 * time.Second)
	}
}
