package serverFinder

import (
	"time"

	"github.com/gorilla/websocket"
)

func Register(action, addr, mainKey, registerAddr string) {
	go func() {
		// 初始化连接地址
		url := "ws://" + addr + "?action=" + action + "&mainKey=" + mainKey + "&addr=" + registerAddr

	RegisterLoop:
		// 创建连接
		conn, _, err := websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			// 失败重连
			time.Sleep(time.Millisecond * 500)
			goto RegisterLoop
		}

		// 间隔 10 秒，向服务端发送心跳信号
		go func(connIn *websocket.Conn) {
			for {
				err = conn.WriteMessage(websocket.TextMessage, []byte("ping"))
				if err != nil {
					break
				}
				time.Sleep(time.Second * 10)
			}
		}(conn)

		// 读取消息保持连接
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				// 断开连接
				conn.Close()
				// 断线重连
				time.Sleep(time.Millisecond * 500)
				goto RegisterLoop
			}
			// fmt.Printf("string(message): %v\n", string(message))
		}
	}()
}
