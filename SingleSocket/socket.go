package main
// 用户和服务器交互websockets 简单例子
//用户发起 ，并发送数据到服务器
//服务器接收数据并处理
//然后发送回用户

//流程
// 客户 -->  WS连接 --> messageChannel   服务器可以对此队列做处理
// 客户 <-- WS连接 <-- messageChannel


import (
	"fmt"
	"net/http"
	"github.com/gorilla/websocket"
	"./impl"
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
		fmt.Println("升级为WebSockets失败")
		return
	}
	if conn,err = impl.InitConnection(wsConn); err != nil {
		fmt.Println("初始化连接出错")
		conn.Close()
		return
	}
	fmt.Println("--结束-")


}




