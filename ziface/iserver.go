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
	AddRouter(router IRouter)
}
