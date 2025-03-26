package serverFinder

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有来源的连接，生产环境中应进行更严格的检查
	},
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	connIP := r.URL.Query().Get("ip")
	mainKey := r.URL.Query().Get("mainKey")
	if mainKey == "" || connIP == "" {
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	// 记录监听连接
	conns, ok := ListenClients[mainKey]
	if !ok {
		ListenClients[mainKey] = map[string]*websocket.Conn{connIP: conn}
	} else {
		conns[connIP] = conn
		ListenClients[mainKey] = conns
	}
	// 连接被关闭
	defer func() {
		conn.Close()
		conns, ok := ListenClients[mainKey]
		if ok {
			delete(conns, connIP)
			ListenClients[mainKey] = conns
		}
	}()
	for {
		// 读取消息
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			break
		}
		// 回复消息
		err = conn.WriteMessage(messageType, message)
		if err != nil {
			break
		}
	}
}
func StartListenInServer() {
	go func() {
		http.HandleFunc("/", websocketHandler)
		log.Println("✔ ServerFinder : 监听服务启动，端口:" + GlobalConfig.ListenPort)
		err := http.ListenAndServe(":"+GlobalConfig.ListenPort, nil)
		if err != nil {
			log.Fatal("✘ ServerFinder : 监听服务启动失败，", err.Error())
		}
	}()
}
