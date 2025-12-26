package znet

import (
	"fmt"
	"strconv"
	"zinx/ziface"
)

// 消息处理模块

type MsgHandle struct {
	// 存放每个 msgID 对应的 Router 方法
	Apis map[uint32]ziface.IRouter
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis: make(map[uint32]ziface.IRouter),
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
