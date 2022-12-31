package znet

import "github.com/stream1080/swinx/face"

type Request struct {
	conn face.Connect // 当前连接
	data []byte       // 请求的数据
}

// 获取当前的连接
func (r *Request) GetConn() face.Connect {
	return r.conn
}

// 获取请求数据
func (r *Request) GetData() []byte {
	return r.data
}
