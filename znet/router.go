package znet

import "github.com/ashin9/zinx/ziface"

// 实现 Router 时, 嵌入 BaseRouter 作为基类, 然后根据需要对基类方法重写
type BaseRouter struct{}

// 方法先都不实现, 因为有的 Router 可能不需要某些 Router, 直接继承, 只实现需要的
func (br *BaseRouter) PreHandler(request ziface.IRequest) {}

func (br *BaseRouter) Handler(request ziface.IRequest) {}

func (br *BaseRouter) PostHandler(request ziface.IRequest) {}
