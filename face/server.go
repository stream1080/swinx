package face

// 服务器接口
type IServer interface {
	Start()                                 // 启动服务
	Stop()                                  // 停止服务
	Serve()                                 // 运行服务
	AddRouter(msgId uint32, router IRouter) // 注册路由
	GetConnMgr() IConnManager               // 获取链接管理器
	SetOnConnStart(func(conn IConnect))     // 注册创建连接的 hook 函数
	SetOnConnStop(func(conn IConnect))      // 注册销毁连接的 hook 函数
	CallOnConnStart(conn IConnect)          // 调用创建连接的 hook 函数
	CallOnConnStop(conn IConnect)           // 调用销毁连接的 hook 函数
}
