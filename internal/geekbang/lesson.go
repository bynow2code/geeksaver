package geekbang

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// LessonReq 课表请求
type LessonReq struct {
	Cid    string `json:"cid"`
	Size   int    `json:"size"`
	Prev   int    `json:"prev"`
	Order  string `json:"order"`
	Sample bool   `json:"sample"`
}

// LessonResp 课表返回
type LessonResp struct {
	Error []interface{} `json:"error"`
	Data  struct {
		List []struct {
			Id           int    `json:"id"`
			ArticleTitle string `json:"article_title"`
		} `json:"list"`
		Page struct {
			Count int  `json:"count"`
			More  bool `json:"more"`
		} `json:"page"`
	} `json:"data"`
	Code int `json:"code"`
}

// LessonRespError 课表错误返回
type LessonRespError struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
}

// UnmarshalJSON 兼容极客不规范的 json 返回格式
func (l *LessonRespError) UnmarshalJSON(data []byte) error {
	// 处理空数组情况
	if string(data) == "[]" {
		*l = LessonRespError{}
		return nil
	}

	// 处理结构体情况
	var tmp LessonRespError
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	*l = tmp
	return nil
}

func GetLessons(lessonReq LessonReq) (*LessonResp, error) {
	reqJson, err := json.Marshal(lessonReq)
	if err != nil {
		return nil, err
	}

	url := "https://time.geekbang.org/serv/v1/column/articles"
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
	resp, err := GetHttpClient().Do(req)
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

	var lessonResp *LessonResp
	if err := json.Unmarshal(body, &lessonResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w, response body: %s", err, string(body))
	}
	return lessonResp, nil
}
