package main

import (
	"fmt"
	"zinx/ziface"
	"zinx/znet"
)

// App 自定义路由
type PingRouter struct {
	znet.BaseRouter
}

func (this *PingRouter) Handler(req ziface.IRequest) {
	fmt.Println("call router handler...")
	// 先读取 Client 数据, 再 ping
	fmt.Println("recv from client msgID = ", req.GetMsgID(), ", data = ", string(req.GetData()))
	err := req.GetConnection().SendMsg(1, []byte("ping..ping..ping"))
	if err != nil {
		fmt.Println(err)
		return
	}
}

func main() {
	s := znet.NewServer("[zinx v0.5]")
	// 添加路由方法
	s.AddRouter(&PingRouter{})
	s.Server()
}
