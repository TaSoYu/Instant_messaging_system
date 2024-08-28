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

// 用户的广播单独提供一个接口
func (user *User) DoMessage(msg string) {

	if msg == "who" {
		user.server.Maplock.Lock()
		for _, cli := range user.server.OnlineMap {
			onlineMsg := "[" + cli.Addr + "] " + cli.Name + " : 现在线....\n"
			user.SendMessage(onlineMsg)
		}
		user.server.Maplock.Unlock()

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
