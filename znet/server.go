package znet

import (
	"fmt"
	"net"
	"zinx/ziface"
)

// iServer 接口的实现
type Server struct {
	Name      string
	IPVersion string
	IP        string
	Port      int
	Router    ziface.IRouter
}

// TCP 服务器最核心的三步: 解析 addr, 创建 listen, accept
func (s *Server) Start() {
	fmt.Printf("[*] Server Listener at IP: %s, Port %d, is starting\n", s.IP, s.Port)

	go func() {
		// 1 获取 TCP 的 addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("reslove tcp addr error : ", err)
			return
		}
		// 2 监听地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen ", s.IPVersion, " err ", err)
		}
		fmt.Println("start Zinx server success.", s.Name, "sucess listenning..")

		var cid uint32
		cid = 0
		// 3 阻塞等待客户端链接, 处理客户端业务
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}

			// 将处理链接的业务方法与 conn 绑定, 得到自定义的链接模块
			dealConn := NewConnection(conn, cid, s.Router)
			cid++

			// 启动当前的链接业务处理
			go dealConn.Start()
		}
	}()
}

// 定义当前客户端链接所绑定的 handle api, 暂时写死, 后续优化应由 app 自定义实现
//func CallBackToClient(conn *net.TCPConn, data []byte, cnt int) error {
//	// 回显的业务
//	fmt.Println("[Conn Handle] CallBackToClient...")
//	if _, err := conn.Write(data[:cnt]); err != nil {
//		fmt.Println("write back buf err", err)
//		return errors.New("CallBackToClient error")
//	}
//	return nil
//}

func (s *Server) Stop() {
	// todo: 释放资源
}

func (s *Server) Server() {
	// 启动服务
	s.Start()
	// todo : 做一些启动服务后的额外业务
	// 阻塞, 否则 start() 是异步的, 执行完就释放了
	select {}
}

func (s *Server) AddRouter(router ziface.IRouter) {
	s.Router = router
	fmt.Println("Add Router Success!")
}

func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8999,
		Router:    nil,
	}
	return s
}
