package cmd

import (
	"fmt"
	"log"
	"os"
	"path"
	"strconv"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/bynow2code/geekbangdocsaver/internal/geekbang"
	"github.com/bynow2code/geekbangdocsaver/internal/geekbang/config"
	"github.com/olekukonko/tablewriter"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

var cid string // 课程id

var mdCmd = &cobra.Command{
	Use:   "md",
	Short: "将课程用 markdown 的形式保存到本地",
	Long:  `将极客时间的课程用 markdown 的形式保存到本地`,
	Run: func(cmd *cobra.Command, args []string) {
		runMdTask()
	},
}

func init() {
	mdCmd.Flags().StringVar(&cid, "cid", "", "课程id（必填）")

	err := mdCmd.MarkFlagRequired("cid")
	if err != nil {
		log.Fatalln(err)
	}

	rootCmd.AddCommand(mdCmd)
}

// 保存文章任务人口
func runMdTask() {
	// 获取课表
	lessons, err := getLessons("100093501")
	if err != nil {
		log.Fatalln(err)
	}

	// 获取下载文章
	for _, lesson := range lessons.Data.List {
		fmt.Printf("处理：（%d %s）\n", lesson.Id, lesson.ArticleTitle)

		err = saveArticle(lesson)
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println()

		break
	}
}

// 获取课表
func getLessons(cid string) (*geekbang.LessonResp, error) {
	// 进度条
	lessonBar := progressbar.Default(1, "获取课表")

	// 请求 api
	lessons, err := geekbang.GetLessons(geekbang.LessonReq{
		Cid:    cid,
		Size:   500,
		Prev:   0,
		Order:  "earliest",
		Sample: false,
	})
	if err != nil {
		return nil, err
	}

	// 进度条满
	err = lessonBar.Finish()
	if err != nil {
		return nil, err
	}

	// 组装 table 数据
	var data [][]string
	for _, lesson := range lessons.Data.List {
		data = append(data, []string{strconv.Itoa(lesson.Id), lesson.ArticleTitle})
	}

	// 添加 table 头部
	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"id", "课表名称"})

	// 添加 table 数据
	err = table.Bulk(data)
	if err != nil {
		return nil, err
	}

	// 添加 table 脚部并渲染
	table.Footer([]string{"total", strconv.Itoa(lessons.Data.Page.Count)})
	if err = table.Render(); err != nil {
		return nil, err
	}

	return lessons, nil
}

// 保存文章
func saveArticle(lesson geekbang.Lesson) error {
	// 进度条初始化
	articleBar := progressbar.NewOptions(3,
		progressbar.OptionShowCount(),
	)

	// 获取文章内容
	articleBar.Describe("获取文章内容")
	article, err := geekbang.GetArticle(geekbang.ArticleReq{
		Id:               strconv.Itoa(lesson.Id),
		IncludeNeighbors: true,
		IsFreelyRead:     true,
	})
	if err != nil {
		return err
	}

	// 获取文章内容，进度+1
	err = articleBar.Add(1)
	if err != nil {
		return err
	}

	// 转换为 md
	articleBar.Describe("转换到 markdown")
	htmlInput := article.Data.ArticleContent
	mdString, err := htmltomarkdown.ConvertString(htmlInput)
	if err != nil {
		return err
	}

	// 转换为 md，进度+2
	err = articleBar.Add(1)
	if err != nil {
		return err
	}

	// 保存 md
	articleBar.Describe("保存 markdown 文件")
	err = saveMd(article, mdString)
	if err != nil {
		return err
	}

	//保存 md，进度+3
	err = articleBar.Add(1)
	if err != nil {
		return err
	}

	return nil
}

func saveMd(article *geekbang.ArticleResp, mdString string) error {
	// 创建文件夹
	cfg := config.GetConfig()
	err := os.MkdirAll(cfg.Md.SavePath, 0666)
	if err != nil {
		return err
	}

	// 创建 md 文件
	mdName := fmt.Sprintf("%d.md", article.Data.Id)
	mdFile := path.Join(cfg.Md.SavePath, mdName)
	file, err := os.OpenFile(mdFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	// 写入 md 内容
	_, err = file.WriteString(mdString)
	if err != nil {
		return err
	}

	return nil
}
