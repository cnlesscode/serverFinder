package serverFinder

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有来源的连接
	},
}

func handler(w http.ResponseWriter, r *http.Request) {

	// 初始化 url 参数
	addr := r.URL.Query().Get("addr")
	mainKey := r.URL.Query().Get("mainKey")
	action := r.URL.Query().Get("action")
	if mainKey == "" || action == "" || addr == "" {
		return
	}

	// 升级协议
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	// 监听模式
	if action == "listen" {
		// 记录监听连接
		conns, ok := ListenClients[mainKey]
		if !ok {
			ListenClients[mainKey] = map[string]*websocket.Conn{addr: conn}
		} else {
			conns[addr] = conn
			ListenClients[mainKey] = conns
		}
		// 连接被关闭
		defer func() {
			conn.Close()
			conns, ok := ListenClients[mainKey]
			if ok {
				delete(conns, addr)
				ListenClients[mainKey] = conns
			}
		}()
	} else if action == "register" { // 注册模式
		SetItem(mainKey, addr, addr)
		// 连接被关闭
		defer func() {
			conn.Close()
			RemoveItem(mainKey, addr)
		}()
	} else if action == "getData" {
		data, ok := Get(mainKey)
		if ok {
			messageByte, _ := json.Marshal(data)
			conn.WriteMessage(websocket.TextMessage, messageByte)
			conn.Close()
		}
	} else {
		conn.Close()
	}

	// 循环读取消息
	for {
		// 读取消息
		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}
		// 回复消息
		err = conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			break
		}
	}
}
func StartServer() {
	go func() {
		http.HandleFunc("/", handler)
		log.Println("✔ ServerFinder : 监听服务启动，端口:" + GlobalConfig.Port)
		err := http.ListenAndServe(":"+GlobalConfig.Port, nil)
		if err != nil {
			log.Fatal("✘ ServerFinder : 监听服务启动失败，", err.Error())
		}
	}()
}
