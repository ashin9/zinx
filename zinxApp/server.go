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
	err := req.GetConnection().SendMsg(200, []byte("ping..ping..ping"))
	if err != nil {
		fmt.Println(err)
		return
	}
}

type HelloZinxRouter struct {
	znet.BaseRouter
}

func (this *HelloZinxRouter) Handler(req ziface.IRequest) {
	fmt.Println("call hello zinx router handler...")
	// 先读取 Client 数据, 再 ping
	fmt.Println("recv from client msgID = ", req.GetMsgID(), ", data = ", string(req.GetData()))
	err := req.GetConnection().SendMsg(201, []byte("hello welcome to zinx"))
	if err != nil {
		fmt.Println(err)
		return
	}
}

// 创建链接后执行的 hook 函数
func DoConnBegin(conn ziface.IConnection) {
	fmt.Println("===> DoConnBegin is called...")
	if err := conn.SendMsg(202, []byte("DoConnection Begin")); err != nil {
		fmt.Println(err)
		return
	}

	// 给当前的链接设置一些属性
	fmt.Println("Set conn Name, Hoe ...")
	conn.SetProperty("Name", "Ashing")
	conn.SetProperty("Github", "https://github.com/ashin9")
	conn.SetProperty("Blog", "https://ashin9.github.io/")
}

// 销毁链接前执行的 hook 函数
func DoConnLost(conn ziface.IConnection) {
	fmt.Println("===> DoConnLost is called...")
	fmt.Println("conn ID", conn.GetConnID(), "is Lost...")

	// 获取链接属性
	//name, err := conn.GetProperty("Name")
	//if err != nil {
	//	return
	//}
	//fmt.Println(name)
	if name, err := conn.GetProperty("Name"); err == nil {
		fmt.Println("Name = ", name)
	}

	if github, err := conn.GetProperty("Github"); err == nil {
		fmt.Println("Github = ", github)
	}

	if blog, err := conn.GetProperty("Blog"); err == nil {
		fmt.Println("Blog = ", blog)
	}

}

func main() {
	s := znet.NewServer("[zinx v0.9]")
	// 注册链接的 hook 方法
	s.SetOnConnStart(DoConnBegin)
	s.SetOnConnStop(DoConnLost)
	// 添加路由方法
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloZinxRouter{})
	s.Server()
}
