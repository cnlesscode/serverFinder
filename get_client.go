package serverFinder

import (
	"encoding/json"
	"time"

	"github.com/cnlesscode/gotool"
	"github.com/gorilla/websocket"
)

func GetDate(addr, mainKey string, callBack func(msg map[string]any)) {
	go func() {
		// 初始化连接地址
		url := "ws://" + addr + "?action=getData&mainKey=" + mainKey + "&addr=" + gotool.GetLocalIP()
	ListenLoop:
		// 建立连接
		conn, _, err := websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			// 失败重连
			time.Sleep(time.Second)
			goto ListenLoop
		}
		defer conn.Close()
		// 监听消息
		for {
			_, messageByte, err := conn.ReadMessage()
			if err != nil {
				break
			}
			message := make(map[string]any)
			if err := json.Unmarshal(messageByte, &message); err != nil {
				break
			}
			callBack(message)
		}
	}()
}
