package serverFinder

import (
	"encoding/json"

	"github.com/gorilla/websocket"
)

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
