package ziface

// 服务器接口
type Service interface {
	Start() // 启动服务
	Stop()  // 停止服务
	Serve() // 运行服务
}
