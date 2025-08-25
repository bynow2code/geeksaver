package geekbang

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/bynow2code/geeksaver/internal/geekbang/config"
)

// OutlineReq 课表请求体
type OutlineReq struct {
	Cid    string `json:"cid"`
	Size   int    `json:"size"`
	Prev   int    `json:"prev"`
	Order  string `json:"order"`
	Sample bool   `json:"sample"`
}

// OutlineResp 课表响应体
type OutlineResp struct {
	Error OutlineRespError `json:"error"`
	Data  struct {
		List []Outline `json:"list"`
		Page struct {
			Count int  `json:"count"`
			More  bool `json:"more"`
		} `json:"page"`
	} `json:"data"`
	Code int `json:"code"`
}

// Outline 课表数据
type Outline struct {
	Id           int    `json:"id"`
	ArticleTitle string `json:"article_title"`
}

// OutlineRespError 错误响应体
type OutlineRespError struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
}

// UnmarshalJSON 兼容极客不规范的 json 返回格式
func (l *OutlineRespError) UnmarshalJSON(data []byte) error {
	// 处理空数组情况
	if string(data) == "[]" {
		*l = OutlineRespError{}
		return nil
	}

	// 处理结构体情况
	var tmp OutlineRespError
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	*l = tmp
	return nil
}

// GetOutline 获取课表
func GetOutline(outlineReq OutlineReq) (*OutlineResp, error) {
	reqJson, err := json.Marshal(outlineReq)
	if err != nil {
		return nil, err
	}

	url := "https://time.geekbang.org/serv/v1/column/articles"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqJson))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
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

	var outlineResp OutlineResp
	if err = json.Unmarshal(body, &outlineResp); err != nil {
		return nil, err
	}

	// 检查业务错误
	if outlineResp.Code != 0 {
		return nil, fmt.Errorf("%s", outlineResp.Error.Msg)
	}

	return &outlineResp, nil
}
