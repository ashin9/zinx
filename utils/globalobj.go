package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"zinx/ziface"
)

type GlobalObject struct {
	// Server
	Host      string
	Port      int
	Name      string
	TcpServer ziface.IServer

	// Zinx
	Version        string
	MaxConn        int
	MaxPackageSize uint32
}

var GlobalObj *GlobalObject

func (g *GlobalObject) ReLoad() {
	data, err := os.ReadFile("conf/zinx.json")
	if err != nil {
		fmt.Println("read config file err:", err)
		return
	}
	if err := json.Unmarshal(data, &GlobalObj); err != nil {
		panic(err)
		return
	}
}

func init() {
	// 默认值
	GlobalObj = &GlobalObject{
		Host:           "0.0.0.0",
		Port:           8999,
		Name:           "ZinxServer",
		TcpServer:      nil,
		Version:        "v0.4",
		MaxConn:        1000,
		MaxPackageSize: 4096,
	}
	GlobalObj.ReLoad()
}
