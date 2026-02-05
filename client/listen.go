package client

import (
	"encoding/json"
	"time"

	"github.com/cnlesscode/gotool"
	"github.com/gorilla/websocket"
)

// 监听服务器数据变化
func Listen(addr, mainKey string, onChange func(msg map[string]int)) {
	// 第一次执行时先查询服务数据并执行监听函数
	data, err := Get(addr, mainKey)
	if err != nil {
		gotool.LogFatal("Failed to listen the serverFinder. Error Code 100001")
		return
	}
	message := make(map[string]int)
	if err := json.Unmarshal([]byte(data), &message); err != nil {
		gotool.LogFatal("Failed to listen the serverFinder. Error Code 100002")
	}
	onChange(message)

	// 开启一个协程进行监听
	go func() {
		// 初始化连接地址
		url := "ws://" + addr + APIBaseURL + "listen&mainKey=" +
			mainKey + "&addr=" + gotool.GetLocalIP()
	ListenLoop:
		// 建立连接
		conn, _, err := websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			// 失败重连
			time.Sleep(time.Second)
			goto ListenLoop
		}

		// 启动心跳协程
		stopHeartbeat := make(chan struct{})
		go func() {
			ticker := time.NewTicker(HeartbeatInterval * time.Second)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
						return
					}
				case <-stopHeartbeat:
					return
				}
			}
		}()

		// 监听 Pong 消息
		conn.SetPongHandler(func(appData string) error {
			return nil
		})

		// 监听消息
		for {
			_, messageByte, err := conn.ReadMessage()
			if err != nil {
				// 断开连接
				close(stopHeartbeat)
				conn.Close()
				break
			}
			// 数据变化消息
			message := make(map[string]int)
			if err := json.Unmarshal(messageByte, &message); err != nil {
				continue
			}
			onChange(message)
		}
		// 断线重连
		time.Sleep(time.Second)
		goto ListenLoop
	}()
}
