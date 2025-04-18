package client

import (
	"time"

	"github.com/cnlesscode/gotool"
	"github.com/cnlesscode/gotool/request"
)

// 删除数据 [ 不需要保持连接 ]
func Remove(addr, mainKey, sonKey string) {
	// 初始化连接地址
	url := "http://" + addr + APIBaseURL + "remove&mainKey=" + mainKey + "&sonKey=" + sonKey + "&addr=" + gotool.GetLocalIP()
	tryNumber := 0
retry:
	_, err := request.GET(url, nil, nil)
	if err != nil {
		if tryNumber < 5 {
			tryNumber++
			time.Sleep(time.Second)
			goto retry
		}
		return
	}
}
