package geekbang

import (
	"net/http"
	"time"
)

var httpClient *http.Client

func init() {
	transport := &http.Transport{
		MaxConnsPerHost:       4,               // 对单个主机最多x个并发连接
		ResponseHeaderTimeout: 5 * time.Second, // 响应头超时
	}

	httpClient = &http.Client{
		Timeout:   10 * time.Second,
		Transport: transport,
	}
}

// GetClient 返回客户端实例
func GetClient() *http.Client {
	return httpClient
}
