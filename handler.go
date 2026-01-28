package serverFinder

import (
	"encoding/json"
	"net/http"
	"time"

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
	action := r.URL.Query().Get("action")
	if mainKey == "" || action == "" || addr == "" {
		return
	}

	// 服务注册
	switch action {

	case "register":
		// 升级协议
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		SetItem(mainKey, addr, time.Now().Unix())
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

	case "get":
		data, ok := Get(mainKey)
		if ok {
			messageByte, _ := json.Marshal(data)
			w.Write(messageByte)
		}

	case "listen":
		// 升级协议
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		// 记录监听连接
		ConnsMu.Lock()
		if _, ok := ListenClients[mainKey]; ok {
			ListenClients[mainKey][conn] = 1
		} else {
			ListenClients[mainKey] = map[*websocket.Conn]int{}
			ListenClients[mainKey][conn] = 1
		}
		ConnsMu.Unlock()
		// 连接被关闭
		defer func() {
			// 删除连接
			ConnsMu.Lock()
			if _, ok := ListenClients[mainKey]; ok {
				delete(ListenClients[mainKey], conn)
			}
			ConnsMu.Unlock()
		}()
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				break
			}
		}

	}

}
