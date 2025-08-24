package geekbang

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// 文章详情请求参数
type ArticleReq struct {
	Id               string `json:"id"`
	IncludeNeighbors bool   `json:"include_neighbors"`
	IsFreelyRead     bool   `json:"is_freelyread"`
}

// 文章详情返回参数
type ArticleResp struct {
	Error ArticleRespError `json:"error"`
	Data  struct {
		Id             int    `json:"id"`
		ArticleTitle   string `json:"article_title"`
		ArticleContent string `json:"article_content"`
		Neighbors      struct {
			Left struct {
				ArticleTitle string `json:"article_title"`
				Id           int    `json:"id"`
			} `json:"left"`
			Right struct {
				ArticleTitle string `json:"article_title"`
				Id           int    `json:"id"`
			} `json:"right"`
		} `json:"neighbors"`
	} `json:"data"`
	Code int `json:"code"`
}

// 文章详情错误返回
type ArticleRespError struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
}

// 兼容极客不规范的 json 返回格式
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

// 获取文章详情数据
func GetArticle(articleReq ArticleReq) (*ArticleResp, error) {
	reqJson, err := json.Marshal(articleReq)
	if err != nil {
		return nil, err
	}

	url := "https://time.geekbang.org/serv/v1/article"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqJson))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", "https://time.geekbang.org")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36 Edg/139.0.0.0")
	req.AddCookie(&http.Cookie{
		Name:  "GCID",
		Value: "e185b7b-cbe0e85-281047b-6227ff1",
	})
	req.AddCookie(&http.Cookie{
		Name:  "GCESS",
		Value: "Bg0BAQcEIVN2hwoEAAAAAAIE8bqqaAUEAAAAAAME8bqqaAkBAQgBAwQEAI0nAAEI4LspAAAAAAALAgYADAEBBgR66IGN",
	})
	resp, err := GetClient().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var articleResp *ArticleResp
	if err := json.Unmarshal(body, &articleResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w, response body: %s", err, string(body))
	}
	return articleResp, nil
}
