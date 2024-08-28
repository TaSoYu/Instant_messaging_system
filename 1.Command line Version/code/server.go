package main

import (
	"fmt"
	"net"
)

type Server struct {
	//先定义 ip和端口号的属性
	Ip   string
	Port int
}

// NewServer 创建一个server对象的接口
func NewServer(ip string, port int) *Server {
	//初始化一个变量传回去
	server := &Server{
		Ip:   ip,
		Port: port,
	}
	return server
}

// Handle 处理连接业务
func (this *Server) Handle(conn net.Conn) {
	fmt.Println("与服务器连接成功")
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

	for {
		//accept
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}

		go this.Handle(conn)
	}
}
