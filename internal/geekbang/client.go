package geekbang

import (
	"net/http"
	"sync"
	"time"
)

var once sync.Once
var client *Client

type Client struct {
	httpClient *http.Client
}

func GetClient() *Client {
	once.Do(func() {
		client = &Client{httpClient: &http.Client{
			Timeout: 10 * time.Second,
		}}
	})
	return client
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	return c.httpClient.Do(req)
}
