package face

// 消息管理
type MsgHandle interface {
	DoMsgHandler(request Request)          // 以非阻塞方式处理消息
	AddRouter(msgId uint32, router Router) // 为消息添加具体处理逻辑
}
