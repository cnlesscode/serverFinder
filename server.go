package serverFinder

import (
	"encoding/json"
	"log"
	"net"
	"os"

	"github.com/cnlesscode/gotool"
	"github.com/cnlesscode/gotool/gfs"
)

var GlobalConfig = Config{}

// TCP服务器结构
type TCPServer struct {
	listener net.Listener
}

// 创建TCP服务器
func NewTCPServer(addr string) *TCPServer {
	// 创建 Socket 端口监听
	// listener 是一个用于面向流的网络协议的公用网络监听器接口，
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	log.Println("✔ ServerFinder : 服务启动成功, 端口" + addr)
	// 返回实例
	return &TCPServer{listener: listener}
}

// Accept 等待客户端连接
func (t *TCPServer) Accept() {
	// 关闭接口解除阻塞的 Accept 操作并返回错误
	defer t.listener.Close()
	// 循环等待客户端连接
	for {
		// 等待客户端连接
		conn, err := t.listener.Accept()
		if err == nil {
			// 处理客户端连接
			go t.Handle(conn)
		}
	}
}

// Handle 处理客户端连接
func (t *TCPServer) Handle(conn net.Conn) {
	for {
		// 创建字节切片
		buf, err := gotool.ReadTCPResponse(conn)
		if err != nil {
			// 退出协程
			conn.Close()
			break
		}
		// 处理消息
		Handle(conn, buf)
	}
}

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
	StartListenInServer()
	// 开启 TCP 服务
	tcpServer := NewTCPServer(":" + GlobalConfig.Port)
	tcpServer.Accept()
}
