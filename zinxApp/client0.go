package main

import (
	"fmt"
	"github.com/ashin9/zinx/znet"
	"io"
	"net"
	"time"
)

// 模拟客户端
func main() {
	fmt.Println("client start...")
	time.Sleep(1 * time.Second)

	//	1, 链接远程服务, 得到 conn
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("connect err:", err, ", exit")
		return
	}
	//	2, 调用 write 写数据
	for {
		// 发送封包的 msg
		dp := znet.NewDataPack()
		binMsg, err := dp.Pack(znet.NewMsgPackage(0, []byte("zinx v0.9 client0 test message")))
		if err != nil {
			fmt.Println("Pack err: ", err)
			return
		}
		if _, err := conn.Write(binMsg); err != nil {
			fmt.Println("write err: ", err)
			return
		}

		// 服务器应该回复一个 msg 数据

		// 1 先读取 tcp 流中的 head 部分得到 id 和 dataLen
		binHead := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn, binHead); err != nil {
			fmt.Println("read head err: ", err)
			break
		}
		// 将 binHead 拆包到 msg 结构体中
		msgHead, err := dp.UnPack(binHead)
		if err != nil {
			fmt.Println("client unpack msgHead err: ", err)
			break
		}
		if msgHead.GetMsgLen() > 0 {
			// 2 根据 dataLen 读取 data
			msg := msgHead.(*znet.Message)
			msg.Data = make([]byte, msg.GetMsgLen())

			if _, err := io.ReadFull(conn, msg.Data); err != nil {
				fmt.Println("read msg data err: ", err)
				return
			}

			fmt.Println("-> Recv Server Msg : ID = ", msg.GetMsgId(), ", len = ", msg.GetMsgLen(), ", data = ", string(msg.GetData()))
		}

		// 阻塞 1s
		time.Sleep(1 * time.Second)
	}
}
