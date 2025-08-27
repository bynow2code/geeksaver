package geekbang

import (
	"net/http"
	"time"

	"github.com/bynow2code/geeksaver/internal/geekbang/config"
)

type Client struct {
	httpClient *http.Client
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	buildCommonReqParams(req)
	return c.httpClient.Do(req)
}

func buildCommonReqParams(req *http.Request) {
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

var client *Client

func init() {
	client = &Client{}
	transport := &http.Transport{
		MaxConnsPerHost:       4,               // 对单个主机最多x个并发连接
		ResponseHeaderTimeout: 5 * time.Second, // 响应头超时
	}

	client.httpClient = &http.Client{
		Timeout:   10 * time.Second,
		Transport: transport,
	}
}

// GetClient 返回客户端实例
func GetClient() *Client {
	return client
}
