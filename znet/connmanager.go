package znet

import (
	"errors"
	"fmt"
	"sync"
	"zinx/ziface"
)

// 链接管理模块
// 为什么不用全局锁?
// 全局锁: 程序只有一份链接管理; 结构体锁: 链接管理是可复用, 可实例化, 线程安全的组件
// 理论可以, 但实践不够工程, 因为看到结构体时不知道他是不是线程安全的, 应该能谁拥有数据, 谁有应该拥有保护它的锁
// 同时若多个实例公用一把锁, 会过度同步, 结构体锁可以更细粒度的并发控制, 符合高性能网络服务器设计思路
// 全局锁适用的场景: 全局唯一的资源如日志文件, 进程级配置, 全局 ID 生成器
type ConnManager struct {
	connections map[uint32]ziface.IConnection // 管理的链接集合
	connLock    sync.RWMutex                  // 保护链接集合的读写锁, 每个实例独有一把锁
}

// 添加链接
func (cm *ConnManager) Add(conn ziface.IConnection) {
	// map 不是线程安全的
	// 保护共享资源, 加写锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	// 将 conn 加入 connManager 中
	cm.connections[conn.GetConnID()] = conn
	fmt.Println("connID = ", conn.GetConnID(), " add to ConnManager success: conn num = ", cm.Len())
}

func (cm *ConnManager) Remove(conn ziface.IConnection) {
	// 保护共享资源, 加写锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	// 将 conn 从 connManager 中删除
	delete(cm.connections, conn.GetConnID())
	fmt.Println("connID = ", conn.GetConnID(), " delete to ConnManager success: conn num = ", cm.Len())
}

func (cm *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	// 加读锁
	cm.connLock.RLock()
	defer cm.connLock.RUnlock()

	if conn, ok := cm.connections[connID]; ok {
		// 找到 conn
		return conn, nil
	} else {
		return nil, errors.New("connection Not Found!")
	}
}

func (cm *ConnManager) Len() int {
	return len(cm.connections)
}

func (cm *ConnManager) ClearConn() {
	// 保护共享资源, 加写锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	for connID, conn := range cm.connections {
		// 停止
		conn.Stop()
		// 删除
		delete(cm.connections, connID)
	}

	fmt.Println("clear all connections sucess! conn num = ", cm.Len())
}

// 创建当前链接的方法
func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
		connLock:    sync.RWMutex{},
	}
}
