package main

import (
	"flag"
	"fmt"
	"net"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
}

func NewClient(serverIp string, serverPort int) *Client {
	//创建客户端对象
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
	}

	//连接server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial err:", err)
		return nil
	}

	client.conn = conn
	//返回对象
	return client
}

var serverIp string
var serverPort int

//   ./client -ip 127.0.0.1 -porn 8888   作为命令行输入

func init() {
	// 绑定变量   输入变量名   默认值   解释说明
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "server ip address")
	flag.IntVar(&serverPort, "port", 8888, "server port")
}

func main() {

	flag.Parse()

	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println("-------连接服务器失败")
		return
	}
	fmt.Println("-------连接服务器成功")

	for {

	}
}
