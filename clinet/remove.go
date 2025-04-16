package client

import (
	"github.com/cnlesscode/gotool"
	"github.com/cnlesscode/gotool/request"
)

// 删除数据 [ 不需要保持连接 ]
func Remove(addr, mainKey, sonKey string) {
	// 初始化连接地址
	url := "http://" + addr + APIBaseURL + "remove&mainKey=" + mainKey + "&sonKey=" + sonKey + "&addr=" + gotool.GetLocalIP()
	request.GET(url, nil, nil)
}
