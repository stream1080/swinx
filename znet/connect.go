package znet

import (
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/stream1080/swinx/face"
)

type Connect struct {
	Conn      *net.TCPConn   // 当前连接的 TCP 套接字
	ConnId    uint32         // 当前连接的Id
	isClosed  bool           // 当前连接的状态
	ExitChan  chan bool      // 告知当前连接退出的 chan
	MsgHandle face.MsgHandle // 当前连接处理的方法 handle
}

// 初始化连接
func NewConnect(conn *net.TCPConn, connId uint32, msgHandler face.MsgHandle) *Connect {
	return &Connect{
		Conn:      conn,
		ConnId:    connId,
		isClosed:  false,
		MsgHandle: msgHandler,
		ExitChan:  make(chan bool, 1),
	}
}

func (c *Connect) StartReader() {
	fmt.Println("Reader Goroutine is running... ")
	defer fmt.Println("ConnId: ", c.ConnId, " Reader is exit, remote addr is ", c.RemoteAddr().String())
	defer c.Stop()

	for {

		// 创建数据包对象
		dp := NewDataPack()

		// 获取客户端 msg head
		head := make([]byte, dp.GetHeadLen())
		_, err := io.ReadFull(c.GetTCPConnect(), head)
		if err != nil {
			fmt.Println("read msg head error: ", err)
			c.ExitChan <- true
			continue
		}

		// 拆包，获取 msgId 和 dataLen
		msg, err := dp.UnPack(head)
		if err != nil {
			fmt.Println("unpack error: ", err)
			c.ExitChan <- true
			continue
		}

		// 根据 dataLen 读取 data, 放在 msg.Data 中
		var data []byte
		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())
			_, err = io.ReadFull(c.GetTCPConnect(), data)
			if err != nil {
				fmt.Println("read msg data error: ", err)
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

		// 从绑定好的消息和对应的处理方法中执行对应的Handle方法
		go c.MsgHandle.DoMsgHandler(req)
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

	// 回写客户端
	_, err = c.Conn.Write(msg)
	if err != nil {
		fmt.Printf("write msg error: %s, msgId: %d\n", err, msgId)
		c.ExitChan <- true
		return errors.New("conn write error")
	}

	return nil
}
