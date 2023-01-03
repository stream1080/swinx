package face

// 连接管理
type IConnManager interface {
	Add(conn IConnect)                   // 添加连接
	Remove(connId uint32)                // 删除连接
	Get(connId uint32) (IConnect, error) // 获取连接
	Count() int                          // 获取连接数
	Clear()                              // 清除所有连接
}
