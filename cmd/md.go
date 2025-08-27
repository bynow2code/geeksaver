package cmd

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/bynow2code/geeksaver/internal/geekbang"
	"github.com/bynow2code/geeksaver/internal/geekbang/config"
	"github.com/olekukonko/tablewriter"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

var mdCmd = &cobra.Command{
	Use:   "md",
	Short: "极客时间课程转 Markdown 并本地保存",
	Long:  fmt.Sprintf("极客时间课程转 Markdown 并本地保存，默认保存路径：%s", config.DefaultMdSavePath),
	Run: func(cmd *cobra.Command, args []string) {
		// 获取参数
		cid, err := cmd.Flags().GetString("cid")
		if err != nil {
			log.Fatalln(err)
		}

		processor := &mdProcessor{cid: cid}

		// 获取课程信息
		err = processor.getCourse()
		if err != nil {
			log.Fatalln(err)
		}

		// 获取课表信息
		err = processor.getOutline()
		if err != nil {
			log.Fatalln(err)
		}

		// 保存文章内容
		err = processor.getArticle()
		if err != nil {
			log.Fatalln(err)
		}

		// 生成 SUMMARY 文件
		err = processor.saveSummary()
		if err != nil {
			log.Fatalln(err)
		}
	},
}

func init() {
	mdCmd.Flags().String("cid", "", "课程id（必填）")

	err := mdCmd.MarkFlagRequired("cid")
	if err != nil {
		log.Fatalln(err)
	}

	rootCmd.AddCommand(mdCmd)
}

type mdProcessor struct {
	cid         string                // 课程id
	courseResp  *geekbang.CourseResp  // 课程响应体
	outlineResp *geekbang.OutlineResp // 课表响应体
}

// 获取课程信息
func (p *mdProcessor) getCourse() error {
	// 进度条
	progressBar := progressbar.Default(1, "获取课程信息")

	productId, err := strconv.Atoi(p.cid)
	if err != nil {
		return err
	}

	courseResp, err := geekbang.GetCourse(geekbang.CourseReq{
		ProductId:            int64(productId),
		WithRecommendArticle: false,
	})
	if err != nil {
		return err
	}
	p.courseResp = courseResp

	// 进度条满
	err = progressBar.Finish()
	if err != nil {
		return err
	}

	return nil
}

// 获取课表
func (p *mdProcessor) getOutline() error {
	// 进度条
	progressBar := progressbar.Default(1, "获取课表信息")

	// 请求 api
	outlineResp, err := geekbang.GetOutline(geekbang.OutlineReq{
		Cid:    p.cid,
		Size:   500,
		Prev:   0,
		Order:  "earliest",
		Sample: false,
	})
	if err != nil {
		return err
	}
	p.outlineResp = outlineResp

	// 进度条满
	err = progressBar.Finish()
	if err != nil {
		return err
	}

	// 组装 table 数据
	var data [][]string
	for _, outline := range outlineResp.Data.List {
		data = append(data, []string{strconv.Itoa(outline.Id), outline.ArticleTitle})
	}

	// 添加 table 头部
	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"id", "课表名称"})

	// 添加 table 数据
	err = table.Bulk(data)
	if err != nil {
		return err
	}

	// 添加 table 脚部并渲染
	table.Footer([]string{"total", strconv.Itoa(outlineResp.Data.Page.Count)})
	if err = table.Render(); err != nil {
		return err
	}

	return nil
}

// 保存文章
func (p *mdProcessor) getArticle() error {
	for _, outline := range p.outlineResp.Data.List {
		// 进度条
		progressBar := progressbar.Default(3, "获取文章内容")

		// 开始获取文章内容
		articleResp, err := geekbang.GetArticle(geekbang.ArticleReq{
			Id:               strconv.Itoa(outline.Id),
			IncludeNeighbors: true,
			IsFreelyRead:     true,
		})
		if err != nil {
			return err
		}

		// 获取文章内容完成，进度->1
		err = progressBar.Add(1)
		if err != nil {
			return err
		}

		// 开始转换为 md
		progressBar.Describe("转换为 markdown 格式")
		htmlInput := fmt.Sprintf("<h1>%s</h1>", articleResp.Data.ArticleTitle)
		htmlInput += articleResp.Data.ArticleContent
		mdString, err := htmltomarkdown.ConvertString(htmlInput)
		if err != nil {
			return err
		}

		// 转换为 md 完成，进度->2
		err = progressBar.Add(1)
		if err != nil {
			return err
		}

		// 开始保存 md
		progressBar.Describe("保存 markdown 文件")
		err = p.saveArticle(articleResp, mdString)
		if err != nil {
			return err
		}

		//保存 md 完成，进度->3
		progressBar.Describe(fmt.Sprintf("%d|%s", articleResp.Data.Id, articleResp.Data.ArticleTitle))
		err = progressBar.Add(1)
		if err != nil {
			return err
		}

		// 延迟 1-1.5 秒（随机波动）
		delay := time.Second + time.Duration(rand.Intn(500))*time.Millisecond
		time.Sleep(delay)
	}

	return nil
}

// 保存文章 Markdown
func (p *mdProcessor) saveArticle(article *geekbang.ArticleResp, mdString string) error {
	// 创建文件夹
	cfg := config.GetConfig()
	articleDir := path.Join(cfg.Md.SavePath, p.courseResp.Data.Title)
	mdDir := path.Join(articleDir, "docs")
	err := os.MkdirAll(mdDir, 0755)
	if err != nil {
		return err
	}

	// 创建 md 文件
	mdName := fmt.Sprintf("%d.md", article.Data.Id)
	mdFile := path.Join(mdDir, mdName)
	file, err := os.OpenFile(mdFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
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

// 创建 summary.md
func (p *mdProcessor) saveSummary() error {
	// 进度条
	progressBar := progressbar.Default(1, "创建课程 summary.md")

	cfg := config.GetConfig()
	summaryFile := path.Join(cfg.Md.SavePath, p.courseResp.Data.Title, "summary.md")
	file, err := os.OpenFile(summaryFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer file.Close()

	// 在内存中拼接
	var content strings.Builder
	// todo 使用 content.Grow 预分配内存，估算内容长度进行优化
	for _, outline := range p.outlineResp.Data.List {
		line := fmt.Sprintf("* [%s](./docs/%d.md)\n", outline.ArticleTitle, outline.Id)
		content.WriteString(line)
	}

	// 一次性写入
	_, err = file.WriteString(content.String())
	if err != nil {
		return err
	}

	err = progressBar.Finish()
	if err != nil {
		return err
	}

	return nil
}
