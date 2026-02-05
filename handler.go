package serverFinder

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/cnlesscode/serverFinder/client"
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
	listen := r.URL.Query().Get("listen")
	if mainKey == "" || action == "" || addr == "" {
		return
	}

	switch action {
	// 服务注册
	case "register":
		// 升级协议
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}

		conn.SetReadDeadline(time.Now().Add(client.ReadDeadlineTimer * time.Second))

		// 保存注册节点数据
		SetItem(mainKey, addr, time.Now().Unix())

		// 注册连接是否同时用于监听
		if listen == "true" {
			AddListener(mainKey, conn)
		}

		// 连接被关闭
		defer func() {
			conn.Close()
			RemoveItem(mainKey, addr)
			if listen == "true" {
				RemoveListener(mainKey, conn)
			}
		}()

		// 监听 ping
		conn.SetPingHandler(func(appData string) error {
			conn.SetReadDeadline(time.Now().Add(client.ReadDeadlineTimer * time.Second))
			return conn.WriteMessage(websocket.PongMessage, []byte(appData))
		})

		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				break
			}
		}
	// 获取数据
	case "get":
		data, ok := Get(mainKey)
		if ok {
			messageByte, err := json.Marshal(data)
			if err != nil {
				return
			}
			w.Write(messageByte)
		}
	// 监听数据变化
	case "listen":
		// 升级协议
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		conn.SetReadDeadline(time.Now().Add(client.ReadDeadlineTimer * time.Second))

		// 记录监听连接
		AddListener(mainKey, conn)
		// 连接被关闭
		defer func() {
			// 删除监听连接
			RemoveListener(mainKey, conn)
		}()

		// 监听 ping
		conn.SetPingHandler(func(appData string) error {
			conn.SetReadDeadline(time.Now().Add(client.ReadDeadlineTimer * time.Second))
			return conn.WriteMessage(websocket.PongMessage, []byte(appData))
		})

		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				break
			}
		}

	}

}

// 保存连接到监听连接池
func AddListener(mainKey string, conn *websocket.Conn) {
	// 记录监听连接
	ConnsMu.Lock()
	if _, ok := ListenClients[mainKey]; ok {
		ListenClients[mainKey][conn] = 1
	} else {
		ListenClients[mainKey] = map[*websocket.Conn]int{}
		ListenClients[mainKey][conn] = 1
	}
	ConnsMu.Unlock()
}

// 删除监听连接
func RemoveListener(mainKey string, conn *websocket.Conn) {
	// 删除连接
	ConnsMu.Lock()
	if _, ok := ListenClients[mainKey]; ok {
		delete(ListenClients[mainKey], conn)
	}
	ConnsMu.Unlock()
}
