package main
// 用户和服务器交互websockets 简单例子
//用户发起 ，并发送数据到服务器
//服务器接收数据并处理
//然后发送回用户


import (
	"fmt"
	"net/http"
	"github.com/gorilla/websocket"
	"socket/SingleSocket/impl"
)

var wsUpgrade = websocket.Upgrader{
	// 读取存储空间大小
	ReadBufferSize:1024,
	// 写入存储空间大小
	WriteBufferSize:1024,
	// 允许所有CORS跨域请求
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main()  {
	fmt.Println("Start")
	//声明路由
	http.HandleFunc("/chat",wsHandler)
	//启用服务器监听
	http.ListenAndServe("0.0.0.0:20181",nil)
}

//访问ws时的方法
func wsHandler(resp http.ResponseWriter, req *http.Request) {
	var(
		wsConn *websocket.Conn
		conn *impl.Connnection
		err error
	)
	if wsConn,err = wsUpgrade.Upgrade(resp, req, nil); err != nil {
		return
	}
	if conn,err = impl.InitConnection(wsConn); err != nil {
		goto ERR
	}

	fmt.Println(conn)

ERR:
	conn.Close()
		//关闭

}


