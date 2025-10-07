package main

import (
	"fmt"
	"net"
	"time"
)

// 模拟客户端
func main() {
	fmt.Printf("client start...")
	time.Sleep(1 * time.Second)

	//	1, 链接远程服务, 得到 conn
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("connect err:", err, ", exit")
		return
	}
	//	2, 调用 write 写数据
	for {
		if _, err := conn.Write([]byte("Hello Zinx")); err != nil {
			fmt.Println("write conn err", err)
			return
		}

		// 接受回显
		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read buf err:", err)
			return
		}
		fmt.Printf("server call back: %s, cnt = %d \n", buf, cnt)

		// 阻塞 1s
		time.Sleep(1 * time.Second)
	}
}
