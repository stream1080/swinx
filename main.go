package main

import (
	"fmt"
	"net"
	"time"

	"github.com/stream1080/zinx/znet"
)

func main() {

	// 启动服务端
	go server()
	time.Sleep(2 * time.Second)

	// 启动客户端
	client()
}

func server() {
	s := znet.NewService("[server v0.2]")
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
