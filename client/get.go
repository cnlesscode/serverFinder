package client

import (
	"github.com/cnlesscode/gotool"
	"github.com/cnlesscode/gotool/request"
)

// 获取数据
func Get(addr, mainKey string) (string, error) {
	url := "http://" + addr + APIBaseURL + "get&mainKey=" + mainKey + "&addr=" + gotool.GetLocalIP()
	return request.GET(url, nil, nil)
}
