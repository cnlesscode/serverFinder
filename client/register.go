package client

import (
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
)

func Register(addr, mainKey, registerAddr string, callback func(msg map[string]any)) {
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
			_, messageByte, err := conn.ReadMessage()
			// 服务端断开连接
			if err != nil {
				conn.Close()
				break
			}
			// 数据变化消息
			message := make(map[string]any)
			if err := json.Unmarshal(messageByte, &message); err != nil {
				continue
			}
			callback(message)
		}

		// 断线重连
		time.Sleep(time.Second)
		goto RegisterLoop

	}()
}
