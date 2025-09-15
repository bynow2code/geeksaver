package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	goversion "github.com/hashicorp/go-version"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "下载安装最新版",
	Long:  `下载安装最新版`,
	Run: func(cmd *cobra.Command, args []string) {
		err := doUpgrade()
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
	},
}

func init() {
	rootCmd.AddCommand(upgradeCmd)
}

type upgradeProcessor struct {
	releaseLatest ReleaseLatest
}

type ReleaseLatest struct {
	TagName     string    `json:"tag_name"`
	PublishedAt time.Time `json:"published_at"`
	Assets      []struct {
		Name               string `json:"name"`
		BrowserDownloadUrl string `json:"browser_download_url"`
	} `json:"assets"`
}

// 程序升级入口
func doUpgrade() error {
	fmt.Printf("当前版本: %s\n", version)

	progressBar := progressbar.NewOptions(1,
		progressbar.OptionSetDescription("检查更新中..."),
		progressbar.OptionShowBytes(false),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetTheme(progressbar.Theme{}),
		progressbar.OptionOnCompletion(func() {
			fmt.Fprint(os.Stderr, "\n")
		}),
	)

	processor := &upgradeProcessor{}
	err := processor.getReleaseLatest()
	if err != nil {
		return fmt.Errorf("获取新版本信息错误：%s", err)
	}
	err = progressBar.Finish()
	if err != nil {
		return err
	}

	fmt.Printf("最新版本: %s\n", processor.releaseLatest.TagName)

	upgrade, err := processor.needUpgrade()
	if err != nil {
		return err
	}
	if !upgrade {
		fmt.Println("当前版本已经是最新版本")
		return nil
	}

	if version == devVersion {
		return errors.New("请使用 go install github.com/bynow2code/geeksaver@latest 的形式安装最新版")
	}

	return nil
}

// 获取最新版本
func (u *upgradeProcessor) getReleaseLatest() error {
	url := "https://api.github.com/repos/bynow2code/geeksaver/releases/latest"
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("检查新版本错误: %s", resp.Status))
	}

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(result, &u.releaseLatest)
	if err != nil {
		return err
	}

	return nil
}

// 是否需要升级
func (u *upgradeProcessor) needUpgrade() (bool, error) {
	newTagName := strings.TrimLeft(u.releaseLatest.TagName, "v")
	newVersion, err := goversion.NewVersion(newTagName)
	if err != nil {
		return false, err
	}

	oldTagName := strings.TrimLeft(version, "v")
	oldVersion, err := goversion.NewVersion(oldTagName)
	if err != nil {
		return false, err
	}

	if oldVersion.LessThan(newVersion) {
		return true, nil
	}

	return false, nil
}

// 获取平台对应的二进制文件名
func getPlatformBinaryName(version string) (string, error) {
	switch runtime.GOOS {
	case "darwin":
		return fmt.Sprintf("geeksaver-%s-macos-%s", version, runtime.GOARCH), nil
	case "linux":
		return fmt.Sprintf("geeksaver-%s-linux-%s", version, runtime.GOARCH), nil
	case "windows":
		return fmt.Sprintf("geeksaver-%s-windows-%s.exe", version, runtime.GOARCH), nil
	default:
		return "", errors.New("不支持的操作系统")
	}
}
