package serverFinder

import (
	"encoding/json"
	"os"

	"github.com/cnlesscode/gotool/gfs"
)

var GlobalConfig = Config{}

// 开启 TCP 服务
func Start(config Config) {
	GlobalConfig = config
	if GlobalConfig.Enable != "on" {
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
	StartServer()
}
