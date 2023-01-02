package face

// 服务器接口
type Server interface {
	Start()                                // 启动服务
	Stop()                                 // 停止服务
	Serve()                                // 运行服务
	AddRouter(msgId uint32, router Router) // 注册路由
	GetConnMgr() ConnManager               // 获取链接管理器
	SetOnConnStart(func(conn Connect))     // 注册创建连接的 hook 函数
	SetOnConnStop(func(conn Connect))      // 注册销毁连接的 hook 函数
	CallOnConnStart(conn Connect)          // 调用创建连接的 hook 函数
	CallOnConnStop(conn Connect)           // 调用销毁连接的 hook 函数
}
