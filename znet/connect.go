package znet

import (
	"fmt"
	"net"

	"github.com/stream1080/zinx/ziface"
)

type Connect struct {
	Conn     *net.TCPConn  // 当前连接的 TCP 套接字
	ConnId   uint32        // 当前连接的Id
	isClosed bool          // 当前连接的状态
	ExitChan chan bool     // 告知当前连接退出的 chan
	Router   ziface.Router // 当前连接处理的方法 handle
}

// 初始化连接
func NewConnect(conn *net.TCPConn, connId uint32, router ziface.Router) *Connect {
	return &Connect{
		Conn:     conn,
		ConnId:   connId,
		isClosed: false,
		Router:   router,
		ExitChan: make(chan bool, 1),
	}
}

func (c *Connect) StartReader() {
	fmt.Println("Reader Goroutine is running... ")
	defer fmt.Println("ConnId: ", c.ConnId, " Reader is exit, remote addr is ", c.RemoteAddr().String())
	defer c.Stop()

	for {
		// 读取客户端的数据到 buf
		buf := make([]byte, 512)
		_, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("recv buf error ", err)
			continue
		}

		// 客户端请求的 Request 数据
		req := &Request{
			conn: c,
			data: buf,
		}

		// 从路由 Routers 中找到注册绑定 Conn 的对应 Handle
		go func(request ziface.Request) {
			// 执行注册的路由方法
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(req)

	}
}

// 启动连接
func (c *Connect) Start() {
	fmt.Println("Connect Start()... ConnId: ", c.ConnId)

	// 从当前连接读数据
	go c.StartReader()

	// TODO 从当前连接写数据
}

// 关闭连接
func (c *Connect) Stop() {
	fmt.Println("Connect Stop()... ConnId: ", c.ConnId)

	// 判断是否关闭
	if c.isClosed {
		return
	}

	c.isClosed = false

	// 关闭连接
	c.Conn.Close()

	// 回收资源
	close(c.ExitChan)
}

// 获取当前连接绑定的 socket conn
func (c *Connect) GetTCPConnect() *net.TCPConn {
	return c.Conn
}

// 获取当前连接的 Id
func (c *Connect) GetConnId() uint32 {
	return c.ConnId
}

// 获取远程客户端的 TCP 状态 Ip，Port
func (c *Connect) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// 发送数据
func (c *Connect) Send(data []byte) error {
	return nil
}
