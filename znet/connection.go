package znet

import (
	"fmt"
	"net"
	"zinx/ziface"
)

// 链接模块
type Connection struct {
	// 当前链接的 socket TCP 套接字
	Conn *net.TCPConn
	// 当前链接的 ID
	ConnID uint32
	// 当前链接的状态
	isClosed bool
	// 当前链接绑定的处理业务方法 API
	// handlerAPI ziface.HandleFunc
	// 告知当前链接已经退出/停止的 channel
	ExitChan chan bool
	Router   ziface.IRouter
}

// 链接的读业务方法
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running...")
	defer fmt.Println("connID = ", c.ConnID, "Reader is exit, remote addr is", c.RemoteAddr().String())
	defer c.Stop()

	for {
		// 读取客户端的数据到 buf
		buf := make([]byte, 512)
		_, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("recv buf err:", err)
			continue
		}

		// 调用绑定的业务逻辑方法
		//if err := c.handlerAPI(c.Conn, buf, cnt); err != nil {
		//	fmt.Println("ConnID", c.ConnID, "handle is error", err)
		//	break
		//}

		req := Request{
			conn: c,
			data: buf,
		}

		// 执行注册的路由方法
		go func(request ziface.IRequest) {
			c.Router.PreHandler(request)
			c.Router.Handler(request)
			c.Router.PostHandler(request)
		}(&req)
	}
}

func (c *Connection) Start() {
	fmt.Println("Conn start... ConnID:", c.ConnID)
	// 启动从当前链接的读数据业务
	go c.StartReader()
	// todo 启动从当前链接写数据的业务

}

func (c *Connection) Stop() {
	fmt.Println("Conn stop... ConnID:", c.ConnID)
	if c.isClosed == true {
		return
	}
	c.isClosed = true
	// 关闭 socket 链接
	c.Conn.Close()
	// 回收资源
	close(c.ExitChan)
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) Send(data []byte) error {
	//TODO implement me
	panic("implement me")
}

// 初始化链接模块的方法
func NewConnection(conn *net.TCPConn, connID uint32, router ziface.IRouter) *Connection {
	c := &Connection{
		Conn:     conn,
		ConnID:   connID,
		isClosed: false,
		ExitChan: make(chan bool, 1),
		Router:   router,
	}
	return c
}
