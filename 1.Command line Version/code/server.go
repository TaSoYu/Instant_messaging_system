package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	//先定义 ip和端口号的属性
	Ip   string
	Port int
	//	map储存在线用户信息
	OnlineMap map[string]*User
	//加锁保证进程安全
	Maplock sync.RWMutex
	//message管道储存所有需广播的信息
	Messages chan string
}

// NewServer 创建一个server对象的接口
func NewServer(ip string, port int) *Server {
	//初始化一个变量传回去
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Messages:  make(chan string),
	}
	return server
}

// Handle 处理连接业务
func (this *Server) Handle(conn net.Conn) {
	fmt.Println("与服务器连接成功")

	user := NewUser(conn, this)

	//用户的上线服务接口
	user.Online()

	//使用一个channel监听用户是否活跃
	isLive := make(chan bool)

	//接受客户端的信息并进行广播
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				//下线服务接口
				user.Offline()
				return
			}

			if err != nil && err != io.EOF {
				fmt.Println("conn read err:", err)
				return
			}
			//提取用户信息 消除最后的\n
			msg := string(buf[:n-1])
			this.Messages <- msg

			//用户广播信息的接口
			user.DoMessage(msg)

			//只要有消息 就塞true
			isLive <- true
		}
	}()

	//此handle先暂时阻塞
	for {
		select {
		case <-isLive:
			//当前用户活跃 应重置定时器
			//不做如何事情 只是为了激活select

		case <-time.After(time.Second * 30):
			//过了时间就是超时了
			//此时需要将这个user关闭

			user.SendMessage("因长时间未响应，您已被踢出系统")
			//销毁资源
			close(user.C)
			err := conn.Close()
			if err != nil {
				fmt.Println(err)
			}
			return
		}
	}
}

// 给予用户，向其他用户广播的接口
func (this *Server) Broadcast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "] " + user.Name + " : " + msg

	//msg存储到server的message中
	this.Messages <- sendMsg
}

func (this *Server) MessageListen() {
	for {
		msg := <-this.Messages

		//一但从msg读到消息，就向所有用户发送
		this.Maplock.Lock()
		for _, cli := range this.OnlineMap {
			cli.C <- msg
		}
		this.Maplock.Unlock()
	}
}

// Start 启动Server服务的接口
func (this *Server) Start() {

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	//因为调用了panic抛出异常，以此我们需要用recover来捕捉panic
	//这里需要注意一点：父进程的recover只能捕捉父进程的panic
	//如果父进程又开了一个子进程，那么子进程的panic，此recover无法捕捉
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	if err != nil {
		// 使用panic抛出异常
		panic(err)
	}

	//创建成功listener要关闭
	defer listener.Close()

	//需要一个进程实时监听msg管道
	go this.MessageListen()

	for {
		//accept
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}

		go this.Handle(conn)
	}
}
