package main

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/stream1080/swinx/face"
	"github.com/stream1080/swinx/znet"
)

func main() {

	// 启动服务端
	go server()
	time.Sleep(2 * time.Second)

	for {
		// 启动客户端
		go client()
		time.Sleep(10 * time.Millisecond)
	}
}

type PingRouter struct {
	znet.BaseRouter
}

// 处理业务的主方法 Hook
func (p *PingRouter) Handle(request face.Request) {
	fmt.Println("Call Router Handle...")

	//先读取客户端的数据，再回写客户端
	fmt.Println("[server]==> receive from client : msgId=", request.GetMsgId(), ", data=", string(request.GetData()))

	err := request.GetConn().SendMsg(200, []byte("receive, complete! \n"))
	if err != nil {
		fmt.Println("Handle call back error: ", err)
	}
}

func server() {
	s := znet.NewServer()
	s.AddRouter(&PingRouter{})
	s.Serve()
}

func client() {

	conn, err := net.Dial("tcp", "127.0.0.1:8889")
	if err != nil {
		fmt.Printf("client start error: %s \n", err)
		return
	}

	for {

		dp := znet.NewDataPack()

		// 封包
		msg, err := dp.Pack(znet.NewMessage(1, []byte("hello server v0.5")))
		if err != nil {
			fmt.Println("pack error: ", err)
		}

		// 发送数据
		_, err = conn.Write(msg)
		if err != nil {
			fmt.Println("write conn error: ", err)
			return
		}

		// 读取 head
		head := make([]byte, dp.GetHeadLen())
		_, err = io.ReadFull(conn, head)
		if err != nil {
			fmt.Println("read head error: ", err)
			return
		}

		// 拆包
		msgHead, err := dp.UnPack(head)
		if err != nil {
			fmt.Println("unpack error: ", err)
			return
		}

		if msgHead.GetDataLen() > 0 {
			msg := msgHead.(*znet.Message)
			msg.Data = make([]byte, msg.GetDataLen())

			_, err = io.ReadFull(conn, msg.Data)
			if err != nil {
				fmt.Println("read data err:", err)
				return
			}
			fmt.Println("[client]==> recvive msgId:", msg.Id, ", len:", msg.DataLen, ", data:", string(msg.Data))
		}

		time.Sleep(1 * time.Second)
	}
}
