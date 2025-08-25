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
	Long:  fmt.Sprintf("将极客时间的课程用 markdown 的形式保存到本地，默认保存路径：%s", config.DefaultMdSavePath),
	Run: func(cmd *cobra.Command, args []string) {
		processor := &MdProcessor{cid: cid}
		// 获取课程信息
		err := processor.getCourse()
		if err != nil {
			log.Fatalln(err)
		}

		// 获取课表信息
		err = processor.getOutline()
		if err != nil {
			log.Fatalln(err)
		}

		// 保存文章内容
		err = processor.saveArticle()
		if err != nil {
			log.Fatalln(err)
		}

		// 生成 SUMMARY 文件
		err = processor.createMdSummary()
		if err != nil {
			log.Fatalln(err)
		}
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

type MdProcessor struct {
	cid         string                // 课程id
	courseResp  *geekbang.CourseResp  // 课程响应体
	outlineResp *geekbang.OutlineResp // 课表响应体
}

// 获取课程信息
func (p *MdProcessor) getCourse() error {
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
func (p *MdProcessor) getOutline() error {
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
func (p *MdProcessor) saveArticle() error {
	for _, outline := range p.outlineResp.Data.List {
		// 进度条
		progressBar := progressbar.NewOptions(3)

		// 获取文章内容
		progressBar.Describe("获取文章内容")
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

		// 转换为 md
		progressBar.Describe("转换到 markdown")
		htmlInput := articleResp.Data.ArticleContent
		mdString, err := htmltomarkdown.ConvertString(htmlInput)
		if err != nil {
			return err
		}

		// 转换为 md 完成，进度->2
		err = progressBar.Add(1)
		if err != nil {
			return err
		}

		// 保存 md
		progressBar.Describe("保存 markdown 文件")
		err = p.saveMdArticle(articleResp, mdString)
		if err != nil {
			return err
		}

		//保存 md 完成，进度->3
		progressBar.Describe(fmt.Sprintf("%d|%s", articleResp.Data.Id, articleResp.Data.ArticleTitle))
		err = progressBar.Add(1)
		if err != nil {
			return err
		}

		// 换行
		fmt.Println()
	}

	return nil
}

// 保存文章 md
func (p *MdProcessor) saveMdArticle(article *geekbang.ArticleResp, mdString string) error {
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
func (p *MdProcessor) createMdSummary() error {
	// 进度条
	progressBar := progressbar.Default(1, "创建课程 summary.md")

	cfg := config.GetConfig()
	summaryFile := path.Join(cfg.Md.SavePath, p.courseResp.Data.Title, "summary.md")
	file, err := os.OpenFile(summaryFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, outline := range p.outlineResp.Data.List {
		_, err = file.WriteString(fmt.Sprintf("* [%s](./docs/%d.md)\n", outline.ArticleTitle, outline.Id))
		if err != nil {
			return err
		}
	}

	err = progressBar.Finish()
	if err != nil {
		return err
	}

	return nil
}
