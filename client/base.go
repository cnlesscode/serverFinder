package client

import "time"

var APIRouteURL string = "/ServerFinder"
var APIBaseURL string = APIRouteURL + "?action="

// 心跳时间间隔
var HeartbeatInterval time.Duration = 30

// 读超时时间
var ReadDeadlineTimer time.Duration = 50
