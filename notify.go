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
//

var ListenClients map[string]map[string]*websocket.Conn = make(map[string]map[string]*websocket.Conn)

func SendNotifyMessage(mainKey string) {
	// 获取对应 mainKey 的数据
	mainDB, ok := Get(mainKey)
	if !ok {
		return
	}

	data, ok := mainDB.(map[string]any)
	if !ok {
		return
	}
	message, err := json.Marshal(data)
	if err != nil {
		return
	}
	// 获取监听客户端连接，发送通知
	if conns, ok := ListenClients[mainKey]; ok {
		for _, conn := range conns {
			conn.WriteMessage(websocket.TextMessage, message)
		}
	}
}

func Listen(addr, mainKey string, onChange func(msg map[string]any)) {
	go func() {
		// 初始化 websocket 客户端连接
		addr = "ws://" + addr + "?mainKey=" + mainKey + "&ip=" + gotool.GetLocalIP()
	ListenLoop:
		conn, _, err := websocket.DefaultDialer.Dial(addr, nil)
		if err != nil {
			// 失败重连
			time.Sleep(time.Second)
			goto ListenLoop
		}
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
			message := make(map[string]any)
			if err := json.Unmarshal(messageByte, &message); err != nil {
				return
			}
			onChange(message)
		}
	}()
}
