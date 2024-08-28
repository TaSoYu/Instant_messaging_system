package main

import "net"

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

		user.conn.Write([]byte(msg + "\n"))
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
	user.server.Broadcast(user, msg)
}
