package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/stream1080/swinx/face"
	"github.com/stream1080/swinx/znet"
)

func main() {

	endpoint := flag.String("t", "", "input endpoint, server、client1 or client2")
	flag.Parse()

	switch *endpoint {
	case "server":
		server()
	case "client1":
		client(1)
	case "client2":
		client(2)
	default:
		flag.Usage()
	}
}

type PingRouter struct {
	znet.BaseRouter
}

// 处理业务的主方法 Hook
func (p *PingRouter) Handle(request face.IRequest) {
	fmt.Println("Call Router Handle...")

	//先读取客户端的数据，再回写客户端
	fmt.Println("[server]==> receive from client : msgId=", request.GetMsgId(), ", data=", string(request.GetData()))

	err := request.GetConn().SendMsg(1, []byte("[ping] receive, complete! \n"))
	if err != nil {
		fmt.Println("Handle call back error: ", err)
	}
}

type PongRouter struct {
	znet.BaseRouter
}

// 处理业务的主方法 Hook
func (p *PongRouter) Handle(request face.IRequest) {
	fmt.Println("Call Router Handle...")

	//先读取客户端的数据，再回写客户端
	fmt.Println("[server]==> receive from client : msgId=", request.GetMsgId(), ", data=", string(request.GetData()))

	err := request.GetConn().SendMsg(2, []byte("[pong] receive, complete! \n"))
	if err != nil {
		fmt.Println("Handle call back error: ", err)
	}
}

// 创建连接的时候执行
func DoConnectBegin(conn face.IConnect) {
	fmt.Println("DoConnecionBegin is Called ... ")

	// 设置连接属性
	fmt.Println("start set connect property...")
	conn.SetProperty("client-name", "Tcp-client")
	conn.SetProperty("conn-id", conn.GetConnId())

	err := conn.SendMsg(2, []byte("DoConnect BEGIN..."))
	if err != nil {
		fmt.Println(err)
	}
}

// 连接断开的时候执行
func DoConnectLost(conn face.IConnect) {
	fmt.Println("DoConneciotnLost is Called ... ")

	if value, err := conn.GetProperty("client-name"); err != nil {
		fmt.Println("property[client-name]: ", value)
	}

	if value, err := conn.GetProperty("conn-id"); err != nil {
		fmt.Println("property[conn-id]: ", value)
	}
}

func server() {
	// 新建一个服务示例
	s := znet.NewServer()

	// 注册 hook 函数
	s.SetOnConnStart(DoConnectBegin)
	s.SetOnConnStop(DoConnectLost)

	// 注册路由
	s.AddRouter(1, &PingRouter{})
	s.AddRouter(2, &PongRouter{})

	// 运行服务
	s.Serve()
}

func client(msgId uint32) {

	conn, err := net.Dial("tcp", "127.0.0.1:8889")
	if err != nil {
		fmt.Printf("client start error: %s \n", err)
		return
	}

	for {

		dp := znet.NewDataPack()

		// 封包
		msg, err := dp.Pack(znet.NewMessage(msgId, []byte("hello tcp server")))
		if err != nil {
			fmt.Println("pack error: ", err)
		}

		// 发送数据
		_, err = conn.Write(msg)
		if err != nil {
			fmt.Println("write conn error:", err)
			return
		}

		// 读取 head
		head := make([]byte, dp.GetHeadLen())
		_, err = io.ReadFull(conn, head)
		if err != nil {
			fmt.Println("read head error:", err)
			return
		}

		// 拆包
		msgHead, err := dp.UnPack(head)
		if err != nil {
			fmt.Println("unpack error:", err)
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
