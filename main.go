package main

import "github.com/bynow2code/geekbangdocsaver/cmd"

func main() {
	cmd.Execute()
}

//	article, err := geekbang.GetArticle(geekbang.ArticleReq{
//		Id:               "600122",
//		IncludeNeighbors: true,
//		IsFreelyRead:     true,
//	})
//	if err != nil {
//		panic(err)
//	}
//	fmt.Println(article)
//}
