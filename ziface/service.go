package ziface

// 服务器接口
type Service interface {
	// 启动服务
	Start()
	// 停止服务
	Stop()
	// 运行服务
	Serve()
}
