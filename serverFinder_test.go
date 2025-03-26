package serverFinder

import (
	"fmt"
	"testing"
	"time"
)

var config Config = Config{
	Enable:     "on",
	DataLogDir: "D:\\githubApps\\serverFinder\\dataLogs",
	Port:       "8001",
	ListenPort: "8002",
}

// 测试命令 :
// go test -v -run=TestRun
func TestRun(t *testing.T) {
	go func() {
		Start(config)
	}()
	SetItem("key1", "skey1", "value1")
	SetItem("key1", "skey2", "value2")
	SetItem("key2", "skey1", "value3")
	res, ok := GetItem("key1", "skey1")
	if ok {
		fmt.Printf("res: %v\n", res)
	}
	serverFinderMap.Range(func(key, value any) bool {
		fmt.Printf("key: %v, value: %v\n", key, value)
		return true
	})
	for {
		time.Sleep(1 * time.Second)
	}
}
