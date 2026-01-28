package serverFinder

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
)

type Config struct {
	Enable     string
	DataLogDir string
	Host       string
	Port       string
}

// 数据字典
var serverFinderMap *sync.Map = &sync.Map{}

// 设置数据
func Set(k string, v any) {
	serverFinderMap.Store(k, v)
	SaveDataToLog(k)
	// 发送通知
	SendNotifyMessage(k)
}

// 设置项目
func SetItem(mainKey, itemKey string, data any) {
	// 获取主库
	mainDB, ok := Get(mainKey)
	// 主库为空
	if !ok {
		Set(mainKey, map[string]any{itemKey: data})
	} else {
		// 已存在数据
		dataOld, ok := mainDB.(map[string]any)
		if ok {
			dataOld[itemKey] = data
			Set(mainKey, dataOld)
		}
	}
}

// 获取数据
func Get(k string) (any, bool) {
	return serverFinderMap.Load(k)
}

// 获取项目数据
func GetItem(mainKey, itemKey string) (any, bool) {
	mainDB, ok := Get(mainKey)
	// 主库为空
	if !ok {
		return nil, false
	}
	data, ok := mainDB.(map[string]any)
	if !ok {
		return nil, false
	}
	item, ok := data[itemKey]
	return item, ok
}

// 删除数据
func Remove(mainKey string) {
	// 获取主库
	_, ok := Get(mainKey)
	if !ok {
		return
	}
	serverFinderMap.Delete(mainKey)
	os.Remove(filepath.Join(GlobalConfig.DataLogDir, mainKey+".json"))
}

// 删除项目
func RemoveItem(mainKey, itemKey string) {
	// 获取主库
	mainDB, ok := Get(mainKey)
	// 主库为空
	if !ok {
		return
	}
	data, ok := mainDB.(map[string]any)
	if !ok {
		return
	}
	delete(data, itemKey)
	Set(mainKey, data)
}

// 保存数据到数据日志文件
func SaveDataToLog(k string) error {
	mapdata, ok := serverFinderMap.Load(k)
	if !ok {
		return errors.New("ServerFinder Error : data not found")
	}
	str, err := json.Marshal(mapdata)
	if err != nil {
		return errors.New("ServerFinder Error : json marshal error")
	}
	filePath := filepath.Join(GlobalConfig.DataLogDir, k+".json")
	err = os.WriteFile(filePath, str, 0777)
	if err != nil {
		return errors.New("ServerFinder Error : write file error")
	}
	return nil
}
