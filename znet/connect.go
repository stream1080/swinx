package znet

import (
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/stream1080/swinx/face"
)

type Connect struct {
	TcpServer face.Server    // 当前连接隶属的 server
	Conn      *net.TCPConn   // 当前连接的 TCP 套接字
	ConnId    uint32         // 当前连接的Id
	isClosed  bool           // 当前连接的状态
	ExitChan  chan bool      // 告知当前连接退出的 chan
	msgChan   chan []byte    // 无缓冲管道，用于读、写两个 goroutine 之间的消息通信
	MsgHandle face.MsgHandle // 当前连接处理的方法 handle
}

// 初始化连接
func NewConnect(server face.Server, conn *net.TCPConn, connId uint32, msgHandler face.MsgHandle) *Connect {
	c := &Connect{
		TcpServer: server,
		Conn:      conn,
		ConnId:    connId,
		isClosed:  false,
		MsgHandle: msgHandler,
		ExitChan:  make(chan bool, 1),
		msgChan:   make(chan []byte),
	}

	// 将新创建的 Conn 添加到链接管理器中
	c.TcpServer.GetConnMgr().Add(c)

	return c
}

func (c *Connect) StartReader() {
	fmt.Println("[Reader Goroutine is running]")
	defer fmt.Println("[Connect Reader exit] , RemoteAddr:", c.RemoteAddr().String(), "connId:", c.ConnId)
	defer c.Stop()

	for {

		// 创建数据包对象
		dp := NewDataPack()

		// 获取客户端 msg head
		head := make([]byte, dp.GetHeadLen())
		_, err := io.ReadFull(c.GetTCPConnect(), head)
		if err != nil {
			fmt.Println("read msg head error:", err)
			c.ExitChan <- true
			continue
		}

		// 拆包，获取 msgId 和 dataLen
		msg, err := dp.UnPack(head)
		if err != nil {
			fmt.Println("unpack error:", err)
			c.ExitChan <- true
			continue
		}

		// 根据 dataLen 读取 data, 放在 msg.Data 中
		var data []byte
		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())
			_, err = io.ReadFull(c.GetTCPConnect(), data)
			if err != nil {
				fmt.Println("read msg data error:", err)
				c.ExitChan <- true
				continue
			}
		}

		msg.SetData(data)

		// 客户端请求的 Request 数据
		req := &Request{
			conn: c,
			msg:  msg,
		}

		// 将消息交给 Worker 处理
		c.MsgHandle.SendMsgTaskQueue(req)
	}
}

// 写消息 Goroutine， 用户将数据发送给客户端
func (c *Connect) StartWriter() {
	fmt.Println("[Writer Goroutine is running]")
	defer fmt.Println("[Connect Writer exit] , RemoteAddr:", c.RemoteAddr().String(), "connId:", c.ConnId)

	for {
		select {
		case data := <-c.msgChan:
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("[Connect Writer exit] Send Data error:, ", err)
				return
			}
		case <-c.ExitChan:
			// conn 关闭
			return
		}
	}
}

// 启动连接
func (c *Connect) Start() {
	fmt.Println("Connect Start()... ConnId: ", c.ConnId)

	// 开启用户从客户端读取数据流程的 Goroutine
	go c.StartReader()
	// 开启用于写回客户端数据流程的 Goroutine
	go c.StartWriter()

	for range c.ExitChan {
		return
	}
}

// 关闭连接
func (c *Connect) Stop() {
	fmt.Println("Connect Stop()... ConnId: ", c.ConnId)

	// 判断是否关闭
	if c.isClosed {
		return
	}

	c.isClosed = true

	// 关闭连接
	c.Conn.Close()

	// 告知 writer 关闭
	c.ExitChan <- true

	// 删除连接管理器的连接
	c.TcpServer.GetConnMgr().Remove(c.GetConnId())

	// 回收资源
	close(c.ExitChan)
	close(c.msgChan)
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
func (c *Connect) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed {
		return errors.New("connect is close when send msg")
	}

	// 封包
	dp := NewDataPack()
	msg, err := dp.Pack(NewMessage(msgId, data))
	if err != nil {
		fmt.Printf("pack error:%s msgId: %d\n", err, msgId)
		return errors.New("pack error")
	}

	// 回写客户端，发送给 Channel 供 Writer 读取
	c.msgChan <- msg

	return nil
}
