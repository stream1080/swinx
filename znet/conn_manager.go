package znet

import (
	"errors"
	"fmt"
	"sync"

	"github.com/stream1080/swinx/face"
)

// 连接管理
type ConnManager struct {
	connMap map[uint32]face.Connect // 管理的连接信息
	lock    sync.RWMutex            // 读写锁
}

// 创建连接管理器
func NewConnManager() *ConnManager {
	return &ConnManager{
		connMap: make(map[uint32]face.Connect),
	}
}

// 添加连接
func (cm *ConnManager) Add(conn face.Connect) {

	// 加写锁, 保护共享资源 Map
	cm.lock.Lock()
	defer cm.lock.Unlock()

	// 将 conn 连接添加到 ConnMananger 中
	cm.connMap[conn.GetConnId()] = conn

	fmt.Println("connect:", conn.GetConnId(), "add to ConnManager successfully: connMap count:", cm.Count())
}

// 删除连接
func (cm *ConnManager) Remove(connId uint32) {

	// 加写锁, 保护共享资源 Map
	cm.lock.Lock()
	defer cm.lock.Unlock()

	// 删除连接信息
	delete(cm.connMap, connId)

	fmt.Println("connect:", connId, "remove successfully: connMap count:", cm.Count())
}

// 获取连接
func (cm *ConnManager) Get(connId uint32) (face.Connect, error) {

	// 加锁, 保护共享资源 Map
	cm.lock.RLock()
	defer cm.lock.RUnlock()

	// 删除连接信息
	if conn, ok := cm.connMap[connId]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connect not found")
	}
}

// 获取连接数
func (cm *ConnManager) Count() int {
	return len(cm.connMap)
}

// 清除所有连接
func (cm *ConnManager) Clear() {

	// 加写锁, 保护共享资源 Map
	cm.lock.Lock()
	defer cm.lock.Unlock()

	// 停止并删除全部的连接信息
	for connId, conn := range cm.connMap {
		conn.Stop()
		delete(cm.connMap, connId)
	}

	fmt.Println("clear all connect successfully!")
}
