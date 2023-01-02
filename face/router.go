package face

// 路由接口，用户自定义的 处理业务 方法
type IRouter interface {
	PreHandle(request IRequest)  // 处理业务前的钩子方法 Hook
	Handle(request IRequest)     // 处理业务的主方法 Hook
	PostHandle(request IRequest) // 处理业务后的钩子方法 Hook
}
