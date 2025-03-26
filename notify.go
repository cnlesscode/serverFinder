package serverFinder

import (
	"fmt"
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
	fmt.Printf("data: %v\n", data)
}
