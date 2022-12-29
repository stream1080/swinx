package znet

import (
	"net"

	"github.com/stream1080/zinx/ziface"
)

type Connect struct {
	Conn      *net.TCPConn      // 当前连接的 TCP 套接字
	ConnId    uint32            // 当前连接的Id
	isClosed  bool              // 当前连接的状态
	handleApi ziface.HandleFunc // 当前连接绑定的业务 api
	ExitChan  chan bool         // 告知当前连接退出的 chan
}

// 初始化连接
func NewConnect(conn *net.TCPConn, connId uint32, callBackApi ziface.HandleFunc) *Connect {
	return &Connect{
		Conn:      conn,
		ConnId:    connId,
		isClosed:  false,
		handleApi: callBackApi,
		ExitChan:  make(chan bool, 1),
	}
}
