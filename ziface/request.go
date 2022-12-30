package ziface

// 把客户端请求的链接信息 和 请求的数据 包装到了 Request 里
type Request interface {
	GetConn() Connect // 获取当前连接
	GetData() []byte  // 获取请求的数据
}
