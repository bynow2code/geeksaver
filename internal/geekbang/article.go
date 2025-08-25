package geekbang

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/bynow2code/geekbangdocsaver/internal/geekbang/config"
)

// ArticleReq 文章详情请求体
type ArticleReq struct {
	Id               string `json:"id"`
	IncludeNeighbors bool   `json:"include_neighbors"`
	IsFreelyRead     bool   `json:"is_freelyread"`
}

// ArticleResp 文章详情响应体
type ArticleResp struct {
	Error ArticleRespError `json:"error"`
	Data  Article          `json:"data"`
	Code  int              `json:"code"`
}

// Article 文章详情
type Article struct {
	Id             int    `json:"id"`
	ArticleTitle   string `json:"article_title"`
	ArticleContent string `json:"article_content"`
}

// ArticleRespError 错误响应体
type ArticleRespError struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
}

// UnmarshalJSON 兼容极客不规范的 json 返回格式
func (a *ArticleRespError) UnmarshalJSON(data []byte) error {
	// 处理空数组情况
	if string(data) == "[]" {
		*a = ArticleRespError{}
		return nil
	}

	// 处理结构体情况
	var tmp ArticleRespError
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	*a = tmp

	return nil
}

// GetArticle 获取文章数据
func GetArticle(articleReq ArticleReq) (*ArticleResp, error) {
	reqJson, err := json.Marshal(articleReq)
	if err != nil {
		return nil, err
	}

	url := "https://time.geekbang.org/serv/v1/article"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqJson))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", "https://time.geekbang.org")
	req.Header.Set("User-Agent", "")

	cfg := config.GetConfig()
	req.AddCookie(&http.Cookie{
		Name:  "GCID",
		Value: cfg.User.GCID,
	})
	req.AddCookie(&http.Cookie{
		Name:  "GCESS",
		Value: cfg.User.GCESS,
	})

	resp, err := GetClient().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var articleResp ArticleResp
	if err = json.Unmarshal(body, &articleResp); err != nil {
		return nil, err
	}
	return &articleResp, nil
}
