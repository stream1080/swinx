package znet

import "github.com/stream1080/zinx/face"

// 实现 router 时，先嵌入这个基类，然后根据需要对这个基类的方法进行重写
type BaseRouter struct{}

// 这里之所以 BaseRouter 的方法都为空，
// 是因为有的 Router 不希望有所有方法都实现
// 所以 Router 全部继承 BaseRouter 的好处是，不需要实现全部方法也可以实例化

// 处理业务前的钩子方法 Hook
func (b *BaseRouter) PreHandle(request face.Request) {}

// 处理业务的主方法 Hook
func (b *BaseRouter) Handle(request face.Request) {}

// 处理业务后的钩子方法 Hook
func (b *BaseRouter) PostHandle(request face.Request) {}
