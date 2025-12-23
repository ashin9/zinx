package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

// 测试 datapack 封包拆包的单元测试
func TestDataPack_Pack(t *testing.T) {
	// 模拟服务

	// 1 创建 socketTCP
	listener, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("server listen err: ", err)
		return
	}

	// goroutine 负责从客户端处理业务
	go func() {
		// 2 从客户端读取数据, 拆包处理
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("server accept error: ", err)
			}

			go func(conn net.Conn) {
				// 处理客户端请求
				// 拆包过程
				// 定义拆包对象
				dp := NewDataPack()
				for {
					// 1 先把包的 head 读出来
					headData := make([]byte, dp.GetHeadLen())
					_, err := io.ReadFull(conn, headData)
					if err != nil {
						fmt.Println("read head err: ", err)
						return
					}

					msgHead, err := dp.UnPack(headData)
					if err != nil {
						fmt.Println("server unpack err: ", err)
						return
					}
					if msgHead.GetMsgLen() > 0 {
						// msg 是有数据的, 需要进行二次读取
						// 2 第二次读取, 根据 head 中的 datalen 读取 data 内容
						msg := msgHead.(*Message)
						msg.Data = make([]byte, msg.GetMsgLen())
						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							fmt.Println("server unpack data err: ", err)
							return
						}

						// 完整的一个消息已经读取完毕
						fmt.Println("-> Reve MsgID: ", msg.MsgId, ", datalen = ", msg.MsgLen, "data = ", string(msg.Data))
					}
				}

			}(conn)
		}
	}()

	// 模拟 Client
	dialer, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client dial err: ", err)
		return
	}

	dp := NewDataPack()
	// 模拟粘包过程, 封装两个包同时发送
	msg1 := &Message{
		MsgId:  1,
		MsgLen: 4,
		Data:   []byte{'z', 'i', 'n', 'x'},
	}
	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("pack msg1 err: ", err)
		return
	}

	msg2 := &Message{
		MsgId:  2,
		MsgLen: 7,
		Data:   []byte{'h', 'l', 'l', 'o', 'z', 'i', 'n', 'x'},
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("pack msg2 err: ", err)
		return
	}

	// 两个包合在一起一起发送
	sendData := append(sendData1, sendData2...) // 将切片的所有元素逐个展开添加

	dialer.Write(sendData)

	// 客户端阻塞
	select {}
}
