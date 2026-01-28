package serverFinder

import (
	"encoding/json"
	"net/http"
	"os"
	"sync"

	"github.com/cnlesscode/gotool"
	"github.com/cnlesscode/gotool/gfs"
	"github.com/cnlesscode/serverFinder/client"
	"github.com/gorilla/websocket"
)

var GlobalConfig = Config{}
var ConnsMu sync.RWMutex

// 监听客户端连接池
var ListenClients map[string]map[*websocket.Conn]int = map[string]map[*websocket.Conn]int{}

// 启动服务
func Start(config Config) {

	// 获取本机IP
	localIP := gotool.GetLocalIP()

	GlobalConfig = config
	if GlobalConfig.Host != localIP {
		return
	}

	if GlobalConfig.DataLogDir == "" {
		GlobalConfig.DataLogDir = "./sf_data_log"
	}

	// 初始化数据目录
	if !gfs.DirExists(GlobalConfig.DataLogDir) {
		err := os.Mkdir(GlobalConfig.DataLogDir, 0777)
		if err != nil {
			gotool.LogFatal(
				"ServerFinder Startup failed. Error : ",
				err.Error(), ".")
		}
	}

	// 加载数据到 syncMap
	res := gfs.ScanDirStruct{
		Path: GlobalConfig.DataLogDir,
	}

	err := gfs.ScanDir(false, &res)
	if err != nil {
		gotool.LogFatal(
			"ServerFinder Startup failed. Error : ",
			err.Error(), ".")
	}

	for _, v := range res.Sons {
		if v.IsDir {
			continue
		}
		// 读取文件内容
		content, err := os.ReadFile(v.Path)
		if err != nil {
			continue
		}
		// 解析数据
		mapData := make(map[string]any, 0)
		err = json.Unmarshal(content, &mapData)
		if err != nil {
			continue
		}
		serverFinderMap.Store(v.Name[0:len(v.Name)-5], mapData)
	}

	// 开启 websocket 监听服务
	http.HandleFunc(client.APIRouteURL, Handler)
	gotool.LogOk(
		"ServerFinder is running on port ",
		GlobalConfig.Port, ".")
	go func() {
		err = http.ListenAndServe(":"+GlobalConfig.Port, nil)
		if err != nil {
			gotool.LogFatal(
				"ServerFinder Startup failed. Error : ",
				err.Error(), ".")
		}
	}()

}
