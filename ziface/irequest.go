package ziface

// 包装链接和数据到 Request
type IRequest interface {
	// 返回当前链接
	GetConnection() IConnection
	// 返回请求的数据
	GetData() []byte
	// 返回请求消息的 ID
	GetMsgID() uint32
}
