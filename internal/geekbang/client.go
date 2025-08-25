package geekbang

import (
	"net/http"
	"time"
)

var httpClient *http.Client

func init() {
	transport := &http.Transport{
		MaxIdleConns:          2,                // 空闲时保留连接数
		IdleConnTimeout:       30 * time.Second, // 空闲连接超时时间
		MaxConnsPerHost:       2,                // 对单个主机最多x个并发连接
		ResponseHeaderTimeout: 5 * time.Second,  // 响应头超时
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
