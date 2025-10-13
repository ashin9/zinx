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

func (this *PingRouter) PreHandler(req ziface.IRequest) {
	fmt.Println("call router prehandler...")
	_, err := req.GetConnection().GetTCPConnection().Write([]byte("before ping...\n"))
	if err != nil {
		fmt.Println("call prehandler err")
	}
}

func (this *PingRouter) Handler(req ziface.IRequest) {
	fmt.Println("call router handler...")
	_, err := req.GetConnection().GetTCPConnection().Write([]byte("ping...\n"))
	if err != nil {
		fmt.Println("call handler err")
	}
}

func (this *PingRouter) PostHandler(req ziface.IRequest) {
	fmt.Println("call router posthandler...")
	_, err := req.GetConnection().GetTCPConnection().Write([]byte("post ping...\n"))
	if err != nil {
		fmt.Println("call posthandler err")
	}
}

func main() {
	s := znet.NewServer("[zinx v0.3]")
	// 添加路由方法
	s.AddRouter(&PingRouter{})
	s.Server()
}
