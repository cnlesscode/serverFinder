package client

import (
	"time"

	"github.com/cnlesscode/gotool"
	"github.com/cnlesscode/gotool/request"
)

// 获取数据
func Get(addr, mainKey string) (string, error) {
	tryNumber := 0
	url := "http://" + addr + APIBaseURL + "get&mainKey=" + mainKey + "&addr=" + gotool.GetLocalIP()
retry:
	res, err := request.GET(url, nil, nil)
	if err != nil || res == "" {
		if tryNumber < 5 {
			tryNumber++
			time.Sleep(time.Second)
			goto retry
		}
		return "", err
	}
	return res, nil
}
