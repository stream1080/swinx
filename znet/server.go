package znet

import (
	"fmt"
	"net"

	"github.com/stream1080/swinx/conf"
	"github.com/stream1080/swinx/face"
	"github.com/stream1080/swinx/pack"
)

// Server 的服务接口实现
type Server struct {
	Name        string                   // 名称
	IpVersion   string                   // ip版本
	IP          string                   // ip地址
	Port        int                      // 端口
	msgHandle   face.IMsgHandle          // 消息管理器
	ConnMgr     face.IConnManager        // 连接管理器
	OnConnStart func(conn face.IConnect) // 创建连接后自动调用的 hook 函数
	OnConnStop  func(conn face.IConnect) // 销毁连接后自动调用的 hook 函数
	packet      face.IDataPack           // 封包拆包类实例
}

// 初始化 Server 方法
func NewServer() face.IServer {
	return &Server{
		Name:      conf.ServerConfig.Name,
		IpVersion: "tcp4",
		IP:        conf.ServerConfig.Host,
		Port:      conf.ServerConfig.Port,
		msgHandle: NewMsgHandle(),
		ConnMgr:   NewConnManager(),
		packet:    pack.NewDataPack(),
	}
}

// 启动服务器
func (s *Server) Start() {
	fmt.Printf("[Start] Server Listenner at IP: %s, Port: %d, is starting \n", s.IP, s.Port)

	go func() {

		// 0. 启动 worker 工作池机制
		s.msgHandle.StartWorkerPool()

		// 1. 获取一个 TCP 的 Addr
		addr, err := net.ResolveTCPAddr(s.IpVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error: ", err)
			return
		}

		// 2. 监听服务器的地址
		listener, err := net.ListenTCP(s.IpVersion, addr)
		if err != nil {
			fmt.Println("listen: ", s.IpVersion, "error: ", err)
			return
		}

		fmt.Printf("start [%s] server successful, Listenning...\n", s.Name)

		// 连接 id
		connId := uint32(0)

		// 3. 阻塞等待客户端连接，处理客户端连接
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("listener.AcceptTCP() error: ", err)
				continue
			}

			// 是否大于最大连接数
			if s.ConnMgr.Count() >= conf.ServerConfig.MaxConn {
				fmt.Println("too many conn, maxConn:", conf.ServerConfig.MaxConn)
				conn.Close()
				continue
			}

			// 将处理连接到业务方法与 conn 进行绑定
			dealConn := NewConnect(s, conn, connId, s.msgHandle)
			connId++

			// 启动 conn 的业务处理
			go dealConn.Start()
		}
	}()
}

// 停止服务器
func (s *Server) Stop() {
	fmt.Println("[stop server]", s.Name)
	// 回收资源
	s.GetConnMgr().Clear()
}

// 运行服务器
func (s *Server) Serve() {
	// 启动服务
	s.Start()

	// 阻塞
	select {}
}

// 注册路由
func (s *Server) AddRouter(msgId uint32, router face.IRouter) {
	s.msgHandle.AddRouter(msgId, router)
	fmt.Println("Add Router Success! ")
}

// 获取连接管理器
func (s *Server) GetConnMgr() face.IConnManager {
	return s.ConnMgr
}

// 注册创建连接的 hook 函数
func (s *Server) SetOnConnStart(hookFunc func(conn face.IConnect)) {
	s.OnConnStart = hookFunc
}

// 注册销毁连接的 hook 函数
func (s *Server) SetOnConnStop(hookFunc func(conn face.IConnect)) {
	s.OnConnStop = hookFunc
}

// 调用创建连接的 hook 函数
func (s *Server) CallOnConnStart(conn face.IConnect) {
	if s.OnConnStart != nil {
		fmt.Print("[CallOnConnStart] ====> ")
		s.OnConnStart(conn)
	}
}

func (s *Server) Packet() face.IDataPack {
	return s.packet
}

// 调用销毁连接的 hook 函数
func (s *Server) CallOnConnStop(conn face.IConnect) {
	if s.OnConnStop != nil {
		fmt.Print("[CallOnConnStop] ====> ")
		s.OnConnStop(conn)
	}
}
