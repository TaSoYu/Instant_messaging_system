package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn

	flag int
}

func NewClient(serverIp string, serverPort int) *Client {
	//创建客户端对象
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       1,
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

	go client.DealResponse()

	client.run()
}

func (client *Client) UpdateName() bool {
	fmt.Println("\n请您的新用户名：")
	fmt.Scanln(&client.Name)

	sendMsg := "rename|" + client.Name

	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("client.conn.Write err:", err)
		return false
	}
	return true
}

// 处理server的回应消息
func (client *Client) DealResponse() {
	//一但client有数据，就拷贝到标准输出上，且永久阻塞 运行
	io.Copy(os.Stdout, client.conn)

	/*
		等效于
		for{
			buf := make()
			client.conn.Read(buf)
			fmt.Printf(buf)
		}
	*/
}

func (client *Client) Menu() bool {
	var flag int

	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更新当前用户名")
	fmt.Println("0.退出")

	fmt.Scanln(&flag)

	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println("请输入正确范围内的数字")
		return false
	}
}

func (client *Client) run() {
	for client.flag != 0 {
		for client.Menu() != true {
			//不返回 true就一直读
		}
		//根据flag 处理不同的业务
		switch client.flag {
		case 1:
			//公聊
			fmt.Println("请输入你要公聊的信息")

		case 2:
			//私聊
			fmt.Println("请按要求输入信息进行私聊   to|张三|消息内容")
		case 3:
			//更新用户名
			client.UpdateName()

		}
	}
}
