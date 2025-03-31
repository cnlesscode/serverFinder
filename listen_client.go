package serverFinder

import (
	"encoding/json"
	"time"

	"github.com/cnlesscode/gotool"
	"github.com/gorilla/websocket"
)

// 监听连接池
// 格式 :
//
//	map[string]map[string]*websocket.Conn = {
//		"fisrtMQfirstMQServers":{"192.168.1.11":conn}
//	}
var ListenClients map[string]map[string]*websocket.Conn = make(map[string]map[string]*websocket.Conn)

func Listen(addr, mainKey string, onChange func(msg map[string]any)) {
	go func() {
		// 初始化连接地址
		url := "ws://" + addr + "/ServerFinder?action=listen&mainKey=" + mainKey + "&addr=" + gotool.GetLocalIP()
	ListenLoop:
		// 建立连接
		conn, _, err := websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			// 失败重连
			time.Sleep(time.Second)
			goto ListenLoop
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

		// 监听消息
		for {
			_, messageByte, err := conn.ReadMessage()
			if err != nil {
				// 断开连接
				conn.Close()
				// 断线重连
				time.Sleep(time.Second)
				goto ListenLoop
			}
			// ping 消息
			messageString := string(messageByte)
			if messageString == "ping" {
				continue
			}
			// 数据变化消息
			message := make(map[string]any)
			if err := json.Unmarshal(messageByte, &message); err != nil {
				return
			}
			onChange(message)
		}
	}()
}
