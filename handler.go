package serverFinder

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/cnlesscode/serverFinder/client"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有来源的连接
	},
}

// 监听客户端连接池
var ConnsMu sync.RWMutex
var ListenClients = make(map[string]map[string]*websocket.Conn)

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
		connUUID := uuid.New().String()
		// 升级协议
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		// 保存注册节点数据
		SetItem(mainKey, addr, time.Now().Unix())
		// 注册连接是否同时用于监听
		if listen == "true" {
			AddListener(mainKey, connUUID, conn)
		}
		// 连接被关闭
		defer func() {
			RemoveItem(mainKey, addr)
			if listen == "true" {
				RemoveListener(mainKey, connUUID)
			}
		}()
		webSocketReadLoopHandle(conn)
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
		connUUID := uuid.New().String()
		// 升级协议
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		// 记录监听连接
		AddListener(mainKey, connUUID, conn)
		// 连接被关闭
		defer func() {
			RemoveListener(mainKey, connUUID)
		}()
		webSocketReadLoopHandle(conn)

	}
}

// 保存连接到监听连接池
func AddListener(mainKey, id string, conn *websocket.Conn) {
	// 记录监听连接
	ConnsMu.Lock()
	defer ConnsMu.Unlock()
	if _, ok := ListenClients[mainKey]; ok {
		ListenClients[mainKey][id] = conn
	} else {
		ListenClients[mainKey] = map[string]*websocket.Conn{}
		ListenClients[mainKey][id] = conn
	}
}

// 删除监听连接
func RemoveListener(mainKey, id string) {
	ConnsMu.Lock()
	defer ConnsMu.Unlock()

	if conns, exists := ListenClients[mainKey]; exists {
		delete(conns, id)
		// 清理空的 mainKey 条目，防止内存泄漏
		if len(conns) == 0 {
			delete(ListenClients, mainKey)
		}
	}
}

func webSocketReadLoopHandle(conn *websocket.Conn) {
	defer conn.Close()

	conn.SetReadDeadline(time.Now().Add(client.ReadDeadlineTimer * time.Second))
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
