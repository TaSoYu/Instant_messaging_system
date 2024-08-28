### 1.启动服务
#### 服务器端：
linux: `go build -o server main.go server.go user.go` 

windows: `go build main.go server.go user.go`

#### 客户端:
linux:  `nc 127.0.0.1 8888`

windows: `telnet 127.0.0.1 8888` 

### 2.功能
#### 2.1 用户上线提醒
#### 2.2 用户消息广播(直接输入信息即可)
#### 2.3 在线用户查询(输入who)
#### 2.4 修改用户名(输入rename|新用户名)
#### 2.5 超时强制踢下线功能
#### 2.6 私聊功能(输入 to|目标用户名|消息内容)