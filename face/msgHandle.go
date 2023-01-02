package face

// 消息管理
type IMsgHandle interface {
	DoMsgHandler(request IRequest)          // 以非阻塞方式处理消息
	AddRouter(msgId uint32, router IRouter) // 为消息添加具体处理逻辑
	StartWorkerPool()                       // 启动 worker 工作池
	SendMsgTaskQueue(request IRequest)      // 将消息交给 TaskQueue ,由 worker 进行处理
}
