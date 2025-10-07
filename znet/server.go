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
		// 3 阻塞等待客户端链接, 处理客户端业务
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}

			// 已经建立连接, 处理业务
			go func() {
				for {
					buf := make([]byte, 512)
					// 读取客户端数据 c -> s, 放到 buf 中
					cnt, err := conn.Read(buf)
					if err != nil {
						fmt.Println("recv buf err", err)
						continue
					}

					// 回显给客户端数据 s -> c, 取 buf 前 cnt 字节
					if _, err := conn.Write(buf[:cnt]); err != nil {
						fmt.Println("write back buf err", err)
						continue
					}
				}
			}()
		}
	}()
}

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

func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8999,
	}
	return s
}
