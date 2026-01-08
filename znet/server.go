package znet

import (
	"fmt"
	"github.com/ashin9/zinx/utils"
	"github.com/ashin9/zinx/ziface"
	"net"
)

// iServer 接口的实现
type Server struct {
	Name      string
	IPVersion string
	IP        string
	Port      int
	// 当前 server 的消息管理模块, 绑定 msgID 和对应处理业务 API 关系
	MsgHandler ziface.IMsgHandle
	// 该 server 的链接管理器
	ConnMgr ziface.IConnManager
	// 链接创建后和销毁前 Hook
	OnConnStart func(conn ziface.IConnection)
	OnConnStop  func(conn ziface.IConnection)
}

// TCP 服务器最核心的三步: 解析 addr, 创建 listen, accept
func (s *Server) Start() {
	fmt.Printf("[*] %s Server Listener at IP: %s, Port %d, is starting\n", s.Name, s.IP, s.Port)
	fmt.Printf("[+] [Zinx] Version: %s, MaxConn: %d, MaxPktSize: %d\n", s.IPVersion, utils.GlobalObj.MaxConn, utils.GlobalObj.MaxPackageSize)

	go func() {
		// 开启 worker goroutine 池和对应的 taskqueue
		s.MsgHandler.StartWorkerPool()

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

			// 判断是否达到最大链接数量
			if s.ConnMgr.Len() >= utils.GlobalObj.MaxConn {
				// todo: 给 client 响应一个超出最大链接的错误包
				fmt.Println("Too Many Connections, MaxConn = ", utils.GlobalObj.MaxConn)
				err := conn.Close()
				if err != nil {
					return
				}
				continue
			}

			// 将处理链接的业务方法与 conn 绑定, 得到自定义的链接模块
			dealConn := NewConnection(s, conn, cid, s.MsgHandler)
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
	fmt.Println("[STOP] Zinx server name", s.Name)
	s.ConnMgr.ClearConn()
}

func (s *Server) Server() {
	// 启动服务
	s.Start()
	// todo : 做一些启动服务后的额外业务
	// 阻塞, 否则 start() 是异步的, 执行完就释放了
	select {}
}

func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgID, router)
	fmt.Println("Add Router Success!")
}

func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}

func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:       utils.GlobalObj.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObj.Host,
		Port:       utils.GlobalObj.Port,
		MsgHandler: NewMsgHandle(),
		ConnMgr:    NewConnManager(),
	}
	return s
}

// 注册 OnConnStart 钩子函数方法
func (s *Server) SetOnConnStart(hookFunc func(connection ziface.IConnection)) {
	s.OnConnStart = hookFunc
}

// 注册 OnConnStop 钩子函数方法
func (s *Server) SetOnConnStop(hookFunc func(connection ziface.IConnection)) {
	s.OnConnStop = hookFunc
}

// 调用 OnConnStart 钩子函数方法
func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("-> Call OnConnStart()...")
		s.OnConnStart(conn)
	}
}

// 调用 OnConnStop 钩子函数方法
func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("-> Call OnConnStop()...")
		s.OnConnStop(conn)
	}
}
