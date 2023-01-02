package face

import "net"

// 连接模块的接口
type Connect interface {
	Start()                                      // 启动连接
	Stop()                                       // 关闭连接
	GetTCPConnect() *net.TCPConn                 // 获取当前连接绑定的 socket conn
	GetConnId() uint32                           // 获取当前连接的 Id
	RemoteAddr() net.Addr                        // 获取远程客户端的 TCP 状态 Ip，Port
	SendMsg(msgId uint32, data []byte) error     // 发送数据
	SetProperty(key string, value interface{})   // 设置链接属性
	GetProperty(key string) (interface{}, error) // 获取链接属性
	RemoveProperty(key string)                   // 移除链接属性
}

// 处理连接业务的方法
type HandleFunc func(*net.TCPConn, []byte, int) error
