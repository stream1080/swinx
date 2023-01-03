package face

// 把客户端请求的链接信息 和 请求的数据 包装到了 IRequest 里
type IRequest interface {
	GetConn() IConnect // 获取当前连接
	GetData() []byte   // 获取请求的数据
	GetMsgId() uint32  // 获取消息Id
}
