package znet

import (
	"fmt"
	"strconv"
	"zinx/utils"
	"zinx/ziface"
)

// 消息处理模块

type MsgHandle struct {
	// 存放每个 msgID 对应的 Router 方法
	Apis map[uint32]ziface.IRouter

	// worker 取任务的消息队列
	TaskQueue []chan ziface.IRequest
	// worker 池数量
	WorkerPollSize uint32
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]ziface.IRouter),
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObj.WorkerPoolSize),
		WorkerPollSize: utils.GlobalObj.WorkerPoolSize,
	}
}

// 执行对应的 Router 消息处理方法
func (mh *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	// 1 从 request 中找到 msgID
	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgID = ", request.GetMsgID(), "is not found! need register")
	}
	// 2 根据 msgID 调用对应的 router 业务
	handler.PreHandler(request)
	handler.Handler(request)
	handler.PostHandler(request)
}

// 为消息添加具体的处理逻辑
func (mh *MsgHandle) AddRouter(msgID uint32, router ziface.IRouter) {
	// 1 判断当前 msg 绑定的 api 处理方法是否已经存在
	if _, ok := mh.Apis[msgID]; ok {
		// id 已经注册
		panic("repeat api, msgID = " + strconv.Itoa(int(msgID)))
	}
	// 2 添加 msg 与 api 的绑定关系
	mh.Apis[msgID] = router
	fmt.Println("add api msgID = ", msgID, " succ!")
}

// 开辟 Worker 工作池, 只启动一次
func (mh *MsgHandle) StartWorkerPool() {
	// 根据 size 启动 worker goroutine
	for i := 0; i < int(mh.WorkerPollSize); i++ {
		// go 每个 worker 对应的 channel 消息队列
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObj.MaxWorkerTaskLen)
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}

// 启动一个 Worker
func (mh *MsgHandle) StartOneWorker(workerID int, taskQueue chan ziface.IRequest) {
	fmt.Println("Worker ID = ", workerID, "is started...")
	// 不断阻塞等待消息队列的消息
	for {
		select {
		// 每个消息是个客户端的 request, 执行其绑定的业务
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

// 分发请求到不同的消息队列中, 轮询均等分配, 分布式场景可以根据 ip 等分配
func (mh *MsgHandle) SendMsgToTaskQueue(request ziface.IRequest) {
	// 1 将消息平均分配给不同 worker
	// 根据 client 建立的 ConnID 进行分配
	workerID := request.GetConnection().GetConnID() % mh.WorkerPollSize
	fmt.Println("Add ConnID = ", request.GetConnection().GetConnID(),
		"request MsgID = ", request.GetMsgID(), "to WorkerID = ", workerID)
	// 2 将消息发送给对应 worker 的 TaskQueue
	mh.TaskQueue[workerID] <- request
}
