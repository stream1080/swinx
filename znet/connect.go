package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/stream1080/swinx/conf"
	"github.com/stream1080/swinx/face"
	"github.com/stream1080/swinx/pack"
)

type Connect struct {
	TcpServer    face.IServer           // 当前连接隶属的 server
	Conn         *net.TCPConn           // 当前连接的 TCP 套接字
	ConnId       uint32                 // 当前连接的Id
	isClosed     bool                   // 当前连接的状态
	ExitChan     chan bool              // 告知当前连接退出的 chan
	msgChan      chan []byte            // 无缓冲管道，用于读、写两个 goroutine 之间的消息通信
	msgBuffChan  chan []byte            // 有缓冲管道，用于读、写两个 goroutine 之间的消息通信
	MsgHandle    face.IMsgHandle        // 当前连接处理的方法 handle
	propertyMap  map[string]interface{} // 连接属性
	propertyLock sync.RWMutex           // 保护连接属性修改的锁
}

// 初始化连接
func NewConnect(server face.IServer, conn *net.TCPConn, connId uint32, msgHandler face.IMsgHandle) *Connect {
	c := &Connect{
		TcpServer:   server,
		Conn:        conn,
		ConnId:      connId,
		isClosed:    false,
		MsgHandle:   msgHandler,
		ExitChan:    make(chan bool, 1),
		msgChan:     make(chan []byte),
		msgBuffChan: make(chan []byte, conf.ServerConfig.MaxMsgChanLen),
		propertyMap: make(map[string]interface{}),
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
		dp := pack.NewDataPack()

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
	defer c.Stop()

	for {
		select {
		case data := <-c.msgChan:
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("[Connect Writer exit] Send Data error:, ", err)
				return
			}
		case data, ok := <-c.msgBuffChan:
			if ok {
				// 有数据要写给客户端
				if _, err := c.Conn.Write(data); err != nil {
					fmt.Println("[Connect Writer exit] Send Buff Data error:, ", err)
					return
				}
			} else {
				fmt.Println("msgBuffChan is Closed")
				break
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

	// 创建连接后，执行钩子方法
	c.TcpServer.CallOnConnStart(c)
}

// 关闭连接
func (c *Connect) Stop() {
	fmt.Println("Connect Stop()... ConnId: ", c.ConnId)

	// 判断是否关闭
	if c.isClosed {
		return
	}

	c.isClosed = true

	// 关闭连接前，执行钩子方法
	c.TcpServer.CallOnConnStop(c)

	// 关闭连接
	c.Conn.Close()

	// 删除连接管理器的连接
	c.TcpServer.GetConnMgr().Remove(c.GetConnId())

	// 告知 writer 关闭
	c.ExitChan <- true

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
	dp := pack.NewDataPack()
	msg, err := dp.Pack(pack.NewMessage(msgId, data))
	if err != nil {
		fmt.Printf("pack error:%s msgId: %d\n", err, msgId)
		return errors.New("pack error")
	}

	// 回写客户端，发送给 Channel 供 Writer 读取
	c.msgChan <- msg

	return nil
}

// 发送消息, 有缓冲
func (c *Connect) SendBuffMsg(msgId uint32, data []byte) error {
	if c.isClosed {
		return errors.New("connect is close when send buff msg")
	}

	// 封包
	dp := pack.NewDataPack()
	msg, err := dp.Pack(pack.NewMessage(msgId, data))
	if err != nil {
		fmt.Printf("pack error:%s msgId: %d\n", err, msgId)
		return errors.New("pack error")
	}

	// 回写客户端，发送给 Channel 供 Writer 读取
	c.msgBuffChan <- msg

	return nil
}

// 设置链接属性
func (c *Connect) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.propertyMap[key] = value
}

// 获取链接属性
func (c *Connect) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	// 判断是否存在
	if value, ok := c.propertyMap[key]; ok {
		return value, nil
	}

	return nil, errors.New("property not found")
}

// 移除链接属性
func (c *Connect) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.propertyMap, key)
}
