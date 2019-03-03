package sockects

import (
	"github.com/gorilla/websocket"
	"net/http"
	"time"
	"../common"
)

//全局静态连接
var conns [] *websocket.Conn
var wsupgrader = websocket.Upgrader{
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	HandshakeTimeout: 5 * time.Second,
	// 取消ws跨域校验
	//CheckOrigin: func(r *http.Request) bool {
	//	return true
	//},
}

// 处理ws请求
func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	var conn *websocket.Conn
	var err error
	//Log("进入web socket"+time.Now().Format("2006-01-02 15:04:05"))
	conn, err = wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	conns = append(conns, conn)

	//gin通过协程调用该handler函数，一旦退出函数，ws会被主动销毁
	for {
		mt, message, err := conn.ReadMessage()
		common.Log("接收消息："+string(message))
		if err != nil {
			break
		}
		for _, item := range conns {
			common.Log("发送消息："+string(message))
			err = item.WriteMessage(mt, message)
			if err != nil {
				continue
			}
		}
	}
}

func Send(msg string) {
	for _, item := range conns {
		err := item.WriteMessage(1, []byte (msg))
		if err != nil {
			continue
		}
	}
}
