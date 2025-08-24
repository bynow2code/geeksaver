package geekbang

import (
	"net/http"
	"sync"
	"time"
)

// http客户端单例
type httpClient struct {
	once   sync.Once
	Client *http.Client
}

var defaultClient = &httpClient{}

func GetHttpClient() *http.Client {
	defaultClient.once.Do(func() {
		transport := &http.Transport{
			MaxIdleConns:          3,                // 空闲时保留连接数
			IdleConnTimeout:       30 * time.Second, // 空闲连接超时时间
			MaxConnsPerHost:       6,                // 对单个主机最多x个并发连接
			ResponseHeaderTimeout: 5 * time.Second,  // 响应头超时
		}

		defaultClient.Client = &http.Client{
			Timeout:   10 * time.Second,
			Transport: transport,
		}
	})
	return defaultClient.Client
}
