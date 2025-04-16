package client

import (
	"time"

	"github.com/gorilla/websocket"
)

func Register(addr, mainKey, registerAddr string) {
	go func() {
		// 初始化连接地址
		url := "ws://" + addr + APIBaseURL + "register&mainKey=" + mainKey + "&addr=" + registerAddr

	RegisterLoop:
		// 创建连接
		conn, _, err := websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			// 失败重连
			time.Sleep(time.Second)
			goto RegisterLoop
		}

		// 读取消息保持连接
		for {
			_, _, err := conn.ReadMessage()
			// 服务端断开连接
			if err != nil {
				conn.Close()
				// 断线重连
				time.Sleep(time.Second)
				goto RegisterLoop
			}
		}
	}()
}
