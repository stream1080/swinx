package znet

import (
	"fmt"
	"net"

	"github.com/stream1080/zinx/ziface"
)

// Service 的服务接口实现
type Service struct {
	Name      string        // 名称
	IpVersion string        // ip版本
	IP        string        // ip地址
	Port      int           // 端口
	Router    ziface.Router // 路由
}

// 初始化 Service 方法
func NewService(name string) ziface.Service {
	return &Service{
		Name:      name,
		IpVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8888,
		Router:    nil,
	}
}

// 启动服务器
func (s *Service) Start() {
	fmt.Printf("[Start] Server Listenner at IP: %s, Port: %d, is starting \n", s.IP, s.Port)

	go func() {
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

		fmt.Println("start Zinx server successful ", s.Name, " Listenning...")

		// 连接 id
		connId := uint32(0)

		// 3. 阻塞等待客户端连接，处理客户端连接
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("listener.AcceptTCP() error: ", err)
				continue
			}

			// 将处理连接到业务方法与 conn 进行绑定
			dealConn := NewConnect(conn, connId, s.Router)
			connId++

			// 启动 conn 的业务处理
			go dealConn.Start()
		}
	}()
}

// 停止服务器
func (s *Service) Stop() {
	// TODO 回收资源
}

// 运行服务器
func (s *Service) Serve() {
	// 启动服务
	s.Start()

	// 阻塞
	select {}
}

// 注册路由
func (s *Service) AddRouter(router ziface.Router) {
	s.Router = router
	fmt.Println("Add Router Success! ")
}
