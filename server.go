package serverFinder

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/cnlesscode/gotool"
	"github.com/cnlesscode/gotool/gfs"
	"github.com/cnlesscode/serverFinder/client"
)

var GlobalConfig = Config{}

// 启动服务
func Start(config Config) {

	// 获取本机IP
	localIP := gotool.GetLocalIP()

	GlobalConfig = config
	if GlobalConfig.Host != localIP {
		return
	}

	// 初始化数据目录
	if !gfs.DirExists(GlobalConfig.DataLogDir) {
		err := os.Mkdir(GlobalConfig.DataLogDir, 0777)
		if err != nil {
			panic("ServerFinder Error : 数据目录创建失败: " + err.Error() + "\n")
		}
	}

	// 加载数据到 syncMap
	res := gfs.ScanDirStruct{
		Path: GlobalConfig.DataLogDir,
	}

	err := gfs.ScanDir(false, &res)
	if err != nil {
		panic("ServerFinder Error : 数据目录扫描失败: " + err.Error() + "\n")
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
	log.Println("✔ ServerFinder : 监听服务启动，端口:" + GlobalConfig.Port)
	go func() {
		err = http.ListenAndServe(":"+GlobalConfig.Port, nil)
		if err != nil {
			log.Fatal("✘ ServerFinder : 监听服务启动失败，", err.Error())
		}
	}()
}
