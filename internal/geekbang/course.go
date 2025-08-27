package geekbang

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// CourseReq 课程信息请求体
type CourseReq struct {
	ProductId            int64 `json:"product_id"`
	WithRecommendArticle bool  `json:"with_recommend_article"`
}

// CourseResp 课程信息响应体
type CourseResp struct {
	Code int `json:"code"`
	Data struct {
		Id      int    `json:"id"`
		Title   string `json:"title"`
		Modules []struct {
			Name    string `json:"name"`
			Title   string `json:"title"`
			Content string `json:"content"`
			Type    string `json:"type"`
			IsTop   bool   `json:"is_top"`
			Audio   struct {
				Title       string `json:"title"`
				Md5         string `json:"md5"`
				Size        int    `json:"size"`
				Dubber      string `json:"dubber"`
				Time        string `json:"time"`
				DownloadUrl string `json:"download_url"`
				Url         string `json:"url"`
			} `json:"audio"`
		} `json:"modules"`
	} `json:"data"`
	Error struct {
		Msg  string `json:"msg"`
		Code int    `json:"code"`
	} `json:"error"`
}

func GetCourse(courseReq CourseReq) (*CourseResp, error) {
	reqJson, err := json.Marshal(courseReq)
	if err != nil {
		return nil, err
	}

	url := "https://time.geekbang.org/serv/v3/column/info"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqJson))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := GetClient().do(req)
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

	var courseResp CourseResp
	if err = json.Unmarshal(body, &courseResp); err != nil {
		return nil, err
	}

	// 检查业务错误
	if courseResp.Code != 0 {
		return nil, fmt.Errorf("%s", courseResp.Error.Msg)
	}

	return &courseResp, nil
}
