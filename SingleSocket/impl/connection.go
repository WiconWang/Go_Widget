package impl

import (
	"github.com/gorilla/websocket"
	"fmt"
)



type Connnection struct{
	// 声明ws
	wsConn *websocket.Conn
	//要发的消息
	messageChannel chan []byte
}

// 封装 连接
func InitConnection(wsConn *websocket.Conn)(conn *Connnection, err error) {
	conn = &Connnection{
		wsConn:wsConn, // 连接
		messageChannel:make(chan []byte,1000),

	}
	// 启动 协程
	go conn.Loop()

	return
}


//关闭线程
func (conn *Connnection) Close()  {
	conn.wsConn.Close()
}


// 死循环监听，如果ws有内容，则接收
// 进入服务器处理
// 写回ws
func (conn *Connnection) Loop()  {

	var (
		data  []byte
		err error
		info string
	)

	for  {

		//读取ws
		if _,data,err = conn.wsConn.ReadMessage(); err != nil {
			fmt.Println("读取消息失败了")
			goto ERR
		}

		//对接收数据做处理
		info = "Your Msg : " + string(data)

		//发送回ws
		if err = conn.wsConn.WriteMessage(websocket.TextMessage,[]byte(info)); err != nil {
			goto ERR
		}

	}
	ERR:
		fmt.Println("循环出错")
		conn.Close()
}

