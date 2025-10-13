package ziface

import "net"

// 链接模块抽象
type IConnection interface {
	// 启动链接
	Start()
	// 停止链接
	Stop()
	// 获取当前链接的绑定 socket conn
	GetTCPConnection() *net.TCPConn
	// 获取当前链接的链接 id
	GetConnID() uint32
	// 获取远程客户端的 TCP状态 IP Port
	RemoteAddr() net.Addr
	// 发送数据, 将数据发送给远程客户端
	Send(data []byte) error
}

// 处理链接业务的方法
type HandleFunc func(*net.TCPConn, []byte, int) error
