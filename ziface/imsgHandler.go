package ziface

// 消息管理抽象
type IMsgHandle interface {
	// 调度 Router 的方法
	DoMsgHandler(request IRequest)
	// 添加 Router 的方法
	AddRouter(msgID uint32, router IRouter)
	StartWorkerPool()
	SendMsgToTaskQueue(request IRequest)
}
