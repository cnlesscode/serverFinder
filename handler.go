package serverFinder

import (
	"encoding/json"
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

func Handler(w http.ResponseWriter, r *http.Request) {

	// 初始化 url 参数
	addr := r.URL.Query().Get("addr")
	mainKey := r.URL.Query().Get("mainKey")
	sonKey := r.URL.Query().Get("sonKey")
	action := r.URL.Query().Get("action")
	if mainKey == "" || action == "" || addr == "" {
		return
	}

	// 监听模式
	if action == "listen" {
		// 升级协议
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
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
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				break
			}
		}
	} else if action == "register" { // 注册模式
		// 升级协议
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		SetItem(mainKey, addr, addr)
		// 连接被关闭
		defer func() {
			conn.Close()
			RemoveItem(mainKey, addr)
		}()
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				break
			}
		}
	} else if action == "get" {
		data, ok := Get(mainKey)
		if ok {
			messageByte, _ := json.Marshal(data)
			w.Write(messageByte)
		}
	} else if action == "remove" {
		RemoveItem(mainKey, sonKey)
		w.Write([]byte("ok"))
	}
}
