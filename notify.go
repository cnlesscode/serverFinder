package serverFinder

import (
	"encoding/json"
	"net"
	"time"

	"github.com/cnlesscode/gotool"
)

type NotifyMessage struct {
	Action int
	Data   []byte
}

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
	dataByte, err := json.Marshal(data)
	if err != nil {
		return
	}
	for _, v := range data {
		addr := v.(string)
		if addr == "none" {
			continue
		}
		conn, err := net.DialTimeout("tcp", addr, 2*time.Second)
		if err != nil {
			continue
		}
		message := NotifyMessage{
			Action: 100,
			Data:   dataByte,
		}
		messageByte, _ := json.Marshal(message)
		gotool.WriteTCPResponse(conn, messageByte)
		conn.Close()
	}
}
