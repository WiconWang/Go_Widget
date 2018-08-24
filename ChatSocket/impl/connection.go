package impl

import (
	"github.com/gorilla/websocket"
	"time"
	"math/rand"
	"strconv"
	"fmt"
)


//流程
// 客户 -->  WS连接 --> reqChannel
// 客户 <-- WS连接 <-- respChannel


type Connnection struct{
	// 声明ws
	wsConn *websocket.Conn
	// 入的消息
	reqChannel chan []byte
	//要发的消息
	respChannel chan []byte


	//要发的消息
	messageChannel chan []byte
}

// 封装 连接
func InitConnection(wsConn *websocket.Conn)(conn *Connnection, err error) {
	conn = &Connnection{
		wsConn:wsConn, // 连接
		reqChannel:make(chan []byte, 1000), //入消息 最多1000
		respChannel:make(chan []byte,1000),
		messageChannel:make(chan []byte,1000),

	}
	// 启动 接收协程
	go conn.reqLoop()

	//启动 发送协程
	go conn.respLoop()

	//go conn.heartLoop()

	return
}

//// 储存 取消息
//func (conn *Connnection) ReadMessage() (data []byte,err error)  {
//	data = <- conn.reqChannel
//	return
//}
//
//// 储存  写消息
//func (conn *Connnection)  respMessage(data []byte)(err error)  {
//	conn.respChannel <- data
//	return
//}

//关闭线程
func (conn *Connnection) Close()  {
	conn.wsConn.Close()
}


//启动死循环

// 死循环监听，如果ws有内容，则 进入 reqChannel 队列
func (conn *Connnection) reqLoop()  {

	var (
		data  []byte
		err error
		MsgType int
	)

	for  {
		if MsgType,data,err = conn.wsConn.ReadMessage(); err != nil {
			fmt.Println("读取消息失败了")
			goto ERR
		}
		fmt.Println(MsgType)
		fmt.Println(data)
		fmt.Println("GOT")
		if conn.messageChannel <- data; err != nil {
			fmt.Println("失败了")
			goto ERR
		}
		fmt.Println("GOT2")

	}
	ERR:
		conn.Close()
}


// 死循环监听，如果 respChannel 输出队列中有内容，则直接丢给ws
func   (conn *Connnection) respLoop()  {

	var (
		data  []byte
		err error
	)

	for  {
		//输出队列有内容时
		data = <- conn.messageChannel
		fmt.Println("---")
		fmt.Println(data)
		//输出给WS
		if err = conn.wsConn.WriteMessage(websocket.TextMessage,data); err != nil {
			goto ERR
		}

	}
	ERR:
		conn.Close()

}

func (conn *Connnection) heartLoop()  {
	var (
		err error
	)
	for {
		// 每隔一秒发送一次心跳
		str := strconv.Itoa(rand.Int())
		if err = conn.wsConn.WriteMessage(websocket.TextMessage, []byte(str)); err != nil {
			return
		}
		time.Sleep(2 * time.Second)
	}

}