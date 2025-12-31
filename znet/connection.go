package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"zinx/utils"
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
	// 告知当前链接已经退出/停止的 channel
	ExitChan chan bool
	// 无缓冲管道, 用于读写 goroutine 之间的通信
	msgChan chan []byte
	// 消息管理, 根据 msgID 选择不同处理 API
	MsgHandler ziface.IMsgHandle
}

// 初始化链接模块的方法
func NewConnection(conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandle) *Connection {
	c := &Connection{
		Conn:       conn,
		ConnID:     connID,
		isClosed:   false,
		ExitChan:   make(chan bool, 1),
		msgChan:    make(chan []byte),
		MsgHandler: msgHandler,
	}
	return c
}

// 链接的读业务方法
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running...")
	defer fmt.Println("[Reader is exit], connID = ", c.ConnID, "remote addr is", c.RemoteAddr().String())
	defer c.Stop()

	for {
		// 解包拆包
		dp := NewDataPack()
		// 读取 Msg Head
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head err: ", err)
			break
		}
		// 拆包, 得到 msgID 和 msgDataLen 放在 msg 消息中
		msg, err := dp.UnPack(headData)
		if err != nil {
			fmt.Println("unpack err: ", err)
			break
		}
		// 根据 dataLen 再次读取 Data 放在 msg.Data 中
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data err: ", err)
				break
			}
		}
		msg.SetData(data)

		req := Request{
			conn: c,
			msg:  msg,
		}

		if utils.GlobalObj.WorkerPoolSize > 0 {
			// 已经开启了 worker goroutine 池, 将消息发送给 worker goroutine 池处理
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			// 从路由中, 找到注册绑定的 Conn 对应的 router 调用
			// 根据绑定好的 msgID, 找到对应的 api 业务处理
			go c.MsgHandler.DoMsgHandler(&req)
		}
	}
}

// 写消息 Goroutine, 专门发送给客户端消息的模块
func (c *Connection) StartWriter() {
	fmt.Println("[Writer Goroutine is running...]")
	defer fmt.Println("[conn Writer exit!]", c.RemoteAddr().String())

	// 阻塞不断等待 channel 的消息, 写给客户端
	for {
		select {
		case data := <-c.msgChan:
			// 有数据发送给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("send data err: ", err)
				return
			}
		case <-c.ExitChan:
			// 表示 Reader 已退出, Writer 也要退出
			return
		}
	}
}

func (c *Connection) Start() {
	fmt.Println("Conn start... ConnID:", c.ConnID)
	// 启动从当前链接的读数据业务
	go c.StartReader()
	// 启动从当前链接写数据的业务
	go c.StartWriter()
}

func (c *Connection) Stop() {
	fmt.Println("Conn stop... ConnID:", c.ConnID)
	if c.isClosed == true {
		return
	}
	c.isClosed = true
	// 关闭 socket 链接
	err := c.Conn.Close()
	if err != nil {
		return
	}
	// 告知 Writer 关闭
	c.ExitChan <- true
	// 回收资源
	close(c.ExitChan)
	close(c.msgChan)
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

// 发送封包数据
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("connection closed when send msg")
	}
	// 将 data 封包
	dp := NewDataPack()
	binMsg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("Pack err. msg id = ", msgId)
		return errors.New("Pack err")
	}
	// 将数据发送给客户端
	c.msgChan <- binMsg
	return nil
}
