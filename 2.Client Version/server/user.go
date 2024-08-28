package main

import (
	"fmt"
	"net"
	"strings"
)

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn

	//user与server绑定
	server *Server
}

func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		conn: conn,

		server: server,
	}
	//	在创建完user对象后，要让它一直监听msg
	go user.ListenMessage()

	return user
}

func (user *User) ListenMessage() {
	for {
		msg := <-user.C

		_, err := user.conn.Write([]byte(msg + "\n"))
		if err != nil {
			fmt.Println("write error:", err)
		}
	}
}

func (user *User) Online() {

	//用户上线后 需要进行的操作
	//先把他加到map里面
	user.server.Maplock.Lock()
	user.server.OnlineMap[user.Addr] = user
	user.server.Maplock.Unlock()

	//然后进行广播
	user.server.Broadcast(user, "已上线")
}

func (user *User) Offline() {

	//下线就从map中删去
	user.server.Maplock.Lock()
	delete(user.server.OnlineMap, user.Name)
	user.server.Maplock.Unlock()
	//下线广播
	user.server.Broadcast(user, "下线")
}

// 用户的广播 或 私聊操作提供一个接口
func (user *User) DoMessage(msg string) {

	if msg == "who" {
		user.server.Maplock.Lock()
		for _, cli := range user.server.OnlineMap {
			onlineMsg := "[" + cli.Addr + "] " + cli.Name + " : 现在线....\n"
			user.SendMessage(onlineMsg)
		}
		user.server.Maplock.Unlock()

	} else if len(msg) > 4 && msg[:3] == "to|" {
		//输入格式为 to|张三|消息内容

		//1.获取目标的name
		remoteName := strings.Split(msg, "|")[1]
		if remoteName == "" {
			user.SendMessage("输入格式有误，请使用\"to|张三|你好呀\"的格式.")
			return
		}

		//2.根据name获得对方的user对象
		remoteUser, ok := user.server.OnlineMap[remoteName]
		if !ok {
			user.SendMessage("目标用户不存在")
			return
		} else if remoteUser == user {
			user.SendMessage("私聊对象不能为自己")
			return
		}

		//3.发送消息
		content := strings.Split(msg, "|")[2]
		if content == "" {
			user.SendMessage("消息不可为空")
			return
		}
		remoteUser.SendMessage(user.Name + "对您说: " + content)

	} else if len(msg) > 7 && msg[:7] == "rename|" {
		newName := strings.Split(msg, "|")[1]

		//先检测新name是否已经存在
		_, ok := user.server.OnlineMap[newName]
		if ok {
			user.SendMessage("这个用户名已经被使用\n")
		} else {
			user.server.Maplock.Lock()

			delete(user.server.OnlineMap, user.Name)
			user.server.OnlineMap[newName] = user
			user.Name = newName

			user.server.Maplock.Unlock()

			user.SendMessage("您已经成功修改用户名为: " + newName + " \n")
		}
	} else {
		user.server.Broadcast(user, msg)
	}

}

func (user *User) SendMessage(msg string) {
	user.conn.Write([]byte(msg + "\n"))
}
