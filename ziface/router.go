package ziface

// 路由接口，用户自定义的 处理业务 方法
type Router interface {
	PreHandle(request Request)  // 处理业务前的钩子方法 Hook
	Handle(request Request)     // 处理业务的主方法 Hook
	PostHandle(request Request) // 处理业务后的钩子方法 Hook
}
