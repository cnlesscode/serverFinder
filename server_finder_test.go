package serverFinder

import (
	"fmt"
	"testing"
	"time"

	"github.com/cnlesscode/gotool"
	"github.com/cnlesscode/serverFinder/client"
)

// 启动服务
// go test -v -run=TestStartServer
func TestStartServer(t *testing.T) {
	gotool.SetLogLevel(3)
	config := Config{
		Host: "192.168.0.185",
		Port: "9901",
	}
	Start(config)
	for {
		time.Sleep(time.Second)
	}
}

// 服务注册
// go test -v -run=TestRegist
func TestRegist(t *testing.T) {
	client.Regist(
		"192.168.0.185:9901",
		"test",
		"192.168.0.185:81",
		nil)
	for {
		time.Sleep(time.Second)
	}
}

// 服务注册
// go test -v -run=TestRegistAndListen
func TestRegistAndListen(t *testing.T) {
	client.Regist(
		"192.168.0.185:9901",
		"test",
		"192.168.0.185:80",
		func(msg map[string]int) {
			fmt.Println(msg)
		})
	for {
		time.Sleep(time.Second)
	}
}

// go test -v -run=TestGET
func TestGET(t *testing.T) {
	res, err := client.Get(
		"192.168.0.185:9901", "test")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res)
	}
}

// go test -v -run=TestListen
func TestListen(t *testing.T) {
	client.Listen(
		"192.168.0.185:9901",
		"test",
		func(msg map[string]int) {
			fmt.Println(msg)
		})
	for {
		time.Sleep(time.Second)
	}
}
