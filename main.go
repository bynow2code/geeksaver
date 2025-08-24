package main

import (
	"fmt"

	"github.com/bynow2code/geekbangdocsaver/internal/geekbang"
)

func main() {
	article, err := geekbang.GetArticle(geekbang.ArticleReq{
		Id:               "600122",
		IncludeNeighbors: true,
		IsFreelyRead:     true,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(article)
}
