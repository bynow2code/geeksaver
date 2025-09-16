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
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		doMdPreRunCheck()
	},
	Run: func(cmd *cobra.Command, args []string) {
		doMd(cmd)
	},
}

func init() {
	mdCmd.Flags().String("cid", "", "课程id（必填）")

	// 必填参数校验
	if err := mdCmd.MarkFlagRequired("cid"); err != nil {
		log.Fatalln(err)
	}

	rootCmd.AddCommand(mdCmd)
}

type mdProcessor struct {
	cid         string                // 课程id
	courseResp  *geekbang.CourseResp  // 课程响应体
	outlineResp *geekbang.OutlineResp // 课表响应体
}

func doMdPreRunCheck() {
	if version == devVersion {
		green := "\033[32m"
		reset := "\033[0m"
		// 输出绿色提示信息，结束后重置颜色
		fmt.Println(green + "当前使用的是 go install 安装，可前往 https://github.com/bynow2code/geeksaver/releases 查看新发布" + reset)
		return
	}

	processor := &upgradeProcessor{}
	if err := processor.getReleaseLatest(); err == nil {
		need, _ := processor.needUpgrade()
		if need {
			green := "\033[32m"
			reset := "\033[0m"
			// 输出绿色升级提示信息，结束后重置颜色
			fmt.Println(green + "有新版本，请执行 geeksaver upgrade 升级到最新版" + reset)
		}
	}
}

// 创建 markdown
func doMd(cmd *cobra.Command) {
	// 获取参数
	cid, err := cmd.Flags().GetString("cid")
	if err != nil {
		log.Fatalln(err)
	}

	processor := &mdProcessor{cid: cid}

	// 获取课程信息
	if err := processor.getCourse(); err != nil {
		log.Fatalln(err)
	}

	// 获取课表信息
	if err := processor.getOutline(); err != nil {
		log.Fatalln(err)
	}

	// 保存文章md
	if err := processor.articleHandle(); err != nil {
		log.Fatalln(err)
	}

	// 生成 SUMMARY 文件
	if err := processor.saveSummary(); err != nil {
		log.Fatalln(err)
	}
}

// 获取课程信息
func (p *mdProcessor) getCourse() error {
	bar := progressbar.NewOptions(1,
		progressbar.OptionSetDescription("获取课程信息"),
		progressbar.OptionShowBytes(false),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetTheme(progressbar.Theme{}),
		progressbar.OptionOnCompletion(func() {
			fmt.Fprint(os.Stderr, "\n")
		}),
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionSetElapsedTime(true),
		progressbar.OptionSetPredictTime(true),
		progressbar.OptionShowElapsedTimeOnFinish(),
	)

	productId, err := strconv.Atoi(p.cid)
	if err != nil {
		return err
	}

	p.courseResp, err = geekbang.GetCourse(geekbang.CourseReq{
		ProductId:            int64(productId),
		WithRecommendArticle: false,
	})
	if err != nil {
		return err
	}

	if err := bar.Finish(); err != nil {
		return err
	}

	return nil
}

// 获取课表
func (p *mdProcessor) getOutline() error {
	bar := progressbar.NewOptions(1,
		progressbar.OptionSetDescription("获取课表信息"),
		progressbar.OptionShowBytes(false),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetTheme(progressbar.Theme{}),
		progressbar.OptionOnCompletion(func() {
			fmt.Fprint(os.Stderr, "\n")
		}),
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionSetElapsedTime(true),
		progressbar.OptionSetPredictTime(true),
		progressbar.OptionShowElapsedTimeOnFinish(),
	)

	var err error
	p.outlineResp, err = geekbang.GetOutline(geekbang.OutlineReq{
		Cid:    p.cid,
		Size:   500,
		Prev:   0,
		Order:  "earliest",
		Sample: false,
	})
	if err != nil {
		return err
	}

	if err := bar.Finish(); err != nil {
		return err
	}

	// 组装 table 数据
	var data [][]string
	for _, outline := range p.outlineResp.Data.List {
		data = append(data, []string{strconv.Itoa(outline.Id), outline.ArticleTitle})
	}

	table := tablewriter.NewWriter(os.Stdout)
	// 添加 table 头部
	table.Header([]string{"id", "课表名称"})
	// 添加 table 数据
	if err := table.Bulk(data); err != nil {
		return err
	}
	// 添加 table 脚部
	table.Footer([]string{"total", strconv.Itoa(p.outlineResp.Data.Page.Count)})
	// 渲染
	if err := table.Render(); err != nil {
		return err
	}

	return nil
}

// 保存文章
func (p *mdProcessor) articleHandle() error {
	for _, outline := range p.outlineResp.Data.List {
		bar := progressbar.NewOptions(3,
			progressbar.OptionSetDescription("获取文章内容"),
			progressbar.OptionShowBytes(false),
			progressbar.OptionFullWidth(),
			progressbar.OptionSetTheme(progressbar.Theme{}),
			progressbar.OptionOnCompletion(func() {
				fmt.Fprint(os.Stderr, "\n")
			}),
			progressbar.OptionSetRenderBlankState(true),
			progressbar.OptionShowCount(),
			progressbar.OptionSetElapsedTime(true),
			progressbar.OptionSetPredictTime(true),
			progressbar.OptionShowElapsedTimeOnFinish(),
		)

		// 获取文章内容
		articleResp, err := geekbang.GetArticle(geekbang.ArticleReq{
			Id:               strconv.Itoa(outline.Id),
			IncludeNeighbors: true,
			IsFreelyRead:     true,
		})
		if err != nil {
			return err
		}
		// 进度+1
		if err := bar.Add(1); err != nil {
			return err
		}

		// 开始转换为 md
		bar.Describe("转换为 markdown 格式")
		htmlInput := fmt.Sprintf("<h1>%s</h1>", articleResp.Data.ArticleTitle)
		htmlInput += articleResp.Data.ArticleContent
		mdString, err := htmltomarkdown.ConvertString(htmlInput)
		if err != nil {
			return err
		}
		// 进度+1
		if err := bar.Add(1); err != nil {
			return err
		}

		// 保存 md
		bar.Describe("保存 markdown 文件")
		if err := p.saveArticle(articleResp, mdString); err != nil {
			return err
		}
		// 进度+1
		bar.Describe(fmt.Sprintf("%d|%s", articleResp.Data.Id, articleResp.Data.ArticleTitle))
		if err := bar.Add(1); err != nil {
			return err
		}

		// 延迟 1-1.5 秒（随机波动）
		delay := time.Second + time.Duration(rand.Intn(500))*time.Millisecond
		time.Sleep(delay)
	}

	return nil
}

// 保存文章
func (p *mdProcessor) saveArticle(article *geekbang.ArticleResp, mdString string) error {
	// 创建文章文件夹
	cfg := config.GetConfig()
	dir := path.Join(cfg.Md.SavePath, p.courseResp.Data.Title)
	saveDir := path.Join(dir, "docs")
	if err := os.MkdirAll(saveDir, 0755); err != nil {
		return err
	}

	// 创建 md 文件
	filename := path.Join(saveDir, fmt.Sprintf("%d.md", article.Data.Id))
	newFile, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer newFile.Close()

	// 写入 md 内容
	if _, err := newFile.WriteString(mdString); err != nil {
		return err
	}

	return nil
}

// 创建 summary.md
func (p *mdProcessor) saveSummary() error {
	bar := progressbar.NewOptions(1,
		progressbar.OptionSetDescription("创建课程 summary.md"),
		progressbar.OptionShowBytes(false),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetTheme(progressbar.Theme{}),
		progressbar.OptionOnCompletion(func() {
			fmt.Fprint(os.Stderr, "\n")
		}),
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionSetElapsedTime(true),
		progressbar.OptionSetPredictTime(true),
		progressbar.OptionShowElapsedTimeOnFinish(),
	)

	// 创建 summary.md 文件
	cfg := config.GetConfig()
	filename := path.Join(cfg.Md.SavePath, p.courseResp.Data.Title, "summary.md")
	newFile, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer newFile.Close()

	// 在内存中拼接
	// 可以使用 content.Grow 估算内容长度，通过预分配内存减少扩容次数
	var content strings.Builder
	for _, outline := range p.outlineResp.Data.List {
		line := fmt.Sprintf("* [%s](./docs/%d.md)\n", outline.ArticleTitle, outline.Id)
		content.WriteString(line)
	}

	// 一次性写入
	if _, err := newFile.WriteString(content.String()); err != nil {
		return err
	}

	if err := bar.Finish(); err != nil {
		return err
	}

	return nil
}
