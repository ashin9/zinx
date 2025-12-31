package ziface

// 服务器接口
type IServer interface {
	// 启动
	Start()
	// 停止
	Stop()
	// 运行
	Server()
	// 给当前服务注册一个路由方法, 供客户端的链接使用
	AddRouter(msgID uint32, router IRouter)
	// 获取当前 server 的链接管理器
	GetConnMgr() IConnManager
	// 注册 OnConnStart 钩子函数方法
	SetOnConnStart(func(connection IConnection))
	// 注册 OnConnStop 钩子函数方法
	SetOnConnStop(func(connection IConnection))
	// 调用 OnConnStart 钩子函数方法
	CallOnConnStart(connection IConnection)
	// 调用 OnConnStop 钩子函数方法
	CallOnConnStop(connection IConnection)
}
