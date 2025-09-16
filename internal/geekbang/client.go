package geekbang

import (
	"net/http"
	"time"

	"github.com/bynow2code/geeksaver/internal/geekbang/config"
)

type client struct {
	*http.Client
}

// 发起 http 请求
func (c *client) do(req *http.Request) (*http.Response, error) {
	setCommonRequestParams(req)
	return c.Do(req)
}

// 设置通用 req 参数
func setCommonRequestParams(req *http.Request) {
	req.Header.Set("Origin", "https://time.geekbang.org")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36 Edg/139.0.0.0")

	cfg := config.GetConfig()
	req.AddCookie(&http.Cookie{
		Name:  "GCID",
		Value: cfg.User.GCID,
	})
	req.AddCookie(&http.Cookie{
		Name:  "GCESS",
		Value: cfg.User.GCESS,
	})
}

// defaultClient
var defaultClient *client

func init() {
	defaultClient = &client{}
	transport := &http.Transport{
		MaxConnsPerHost:       4,                // 对单个主机最多x个并发连接
		ResponseHeaderTimeout: 10 * time.Second, // 响应头超时
	}

	defaultClient.Client = &http.Client{
		Timeout:   30 * time.Second,
		Transport: transport,
	}
}

// 返回客户端实例
func getClient() *client {
	return defaultClient
}
