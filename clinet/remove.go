package client

import (
	"time"

	"github.com/cnlesscode/gotool"
	"github.com/gorilla/websocket"
)

// 删除数据 [ 不需要保持连接 ]
func Remove(addr, mainKey, sonKey string) {
	tryCount := 0
	// 初始化连接地址
	url := "ws://" + addr + APIBaseURL + "remove&mainKey=" + mainKey + "&sonKey=" + sonKey + "&addr=" + gotool.GetLocalIP()
GetDataLoop:
	// 建立连接
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		// 失败重连
		tryCount++
		if tryCount > 3 {
			return
		}
		time.Sleep(time.Second)
		goto GetDataLoop
	}
	defer conn.Close()
	// 监听消息
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}
