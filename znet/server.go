package znet

import (
	"fmt"
	"net"

	"github.com/stream1080/swinx/conf"
	"github.com/stream1080/swinx/face"
)

// Server 的服务接口实现
type Server struct {
	Name      string           // 名称
	IpVersion string           // ip版本
	IP        string           // ip地址
	Port      int              // 端口
	msgHandle face.MsgHandle   // 消息管理器
	ConnMgr   face.ConnManager // 连接管理器
}

// 初始化 Server 方法
func NewServer() face.Server {
	return &Server{
		Name:      conf.ServerConfig.Name,
		IpVersion: "tcp4",
		IP:        conf.ServerConfig.Host,
		Port:      conf.ServerConfig.Port,
		msgHandle: NewMsgHandle(),
		ConnMgr:   NewConnManager(),
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
func (s *Server) AddRouter(msgId uint32, router face.Router) {
	s.msgHandle.AddRouter(msgId, router)
	fmt.Println("Add Router Success! ")
}

// 获取连接管理器
func (s *Server) GetConnMgr() face.ConnManager {
	return s.ConnMgr
}
