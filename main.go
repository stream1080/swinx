package main

import (
	"fmt"
	"net"
	"time"

	"github.com/stream1080/zinx/ziface"
	"github.com/stream1080/zinx/znet"
)

type PingRouter struct {
	znet.BaseRouter
}

// 处理业务前的钩子方法 Hook
func (p *PingRouter) PreHandle(request ziface.Request) {
	fmt.Println("Call Router PreHandle...")
	_, err := request.GetConn().GetTCPConnect().Write([]byte(" before ping ....\n"))
	if err != nil {
		fmt.Println("PreHandle call back error: ", err)
	}
}

// 处理业务的主方法 Hook
func (p *PingRouter) Handle(request ziface.Request) {
	fmt.Println("Call Router Handle...")
	_, err := request.GetConn().GetTCPConnect().Write([]byte(" before ping ....\n"))
	if err != nil {
		fmt.Println("Handle call back error: ", err)
	}
}

// 处理业务后的钩子方法 Hook
func (p *PingRouter) PostHandle(request ziface.Request) {
	fmt.Println("Call Router PostHandle...")
	_, err := request.GetConn().GetTCPConnect().Write([]byte(" after ping ....\n"))
	if err != nil {
		fmt.Println("PostHandle call back error: ", err)
	}
}

func main() {

	// 启动服务端
	go server()
	time.Sleep(2 * time.Second)

	// 启动客户端
	client()
}

func server() {
	s := znet.NewService("[server v0.3]")
	s.AddRouter(&PingRouter{})
	s.Serve()
}

func client() {

	conn, err := net.Dial("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Printf("client satrt error: %s \n", err)
		return
	}

	for {
		_, err := conn.Write([]byte("hello server v0.2"))
		if err != nil {
			fmt.Println("write conn error: ", err)
			return
		}

		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("conn.Read error: ", err)
			return
		}

		fmt.Printf("server call back: %s, cnt: %d \n", buf, cnt)

		time.Sleep(1 * time.Second)
	}
}
