package ziface

// 链接管理模块抽象
type IConnManager interface {
	// 添加
	Add(conn IConnection)
	// 删除
	Remove(conn IConnection)
	// 查询链接, 根据 connID
	Get(connID uint32) (IConnection, error)
	// 查询链接总数
	Len() int
	// 清理并终止所有链接
	ClearConn()
}
