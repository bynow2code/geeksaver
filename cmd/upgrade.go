package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	goversion "github.com/hashicorp/go-version"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "升级到最新版",
	Long:  `升级到最新版`,
	Run: func(cmd *cobra.Command, args []string) {
		doUpgrade()
	},
}

func init() {
	rootCmd.AddCommand(upgradeCmd)
}

type upgradeProcessor struct {
	releaseLatest *ReleaseLatest
	assert        *Asset
}

type ReleaseLatest struct {
	TagName     string    `json:"tag_name"`
	PublishedAt time.Time `json:"published_at"`
	Assets      []Asset   `json:"assets"`
}

type Asset struct {
	Name               string `json:"name"`
	BrowserDownloadUrl string `json:"browser_download_url"`
}

// 程序升级入口
func doUpgrade() {
	fmt.Printf("当前版本: %s\n", version)

	processor := &upgradeProcessor{}

	// 检查更新
	if err := processor.checkForUpdate(); err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("最新版本: %s\n", processor.releaseLatest.TagName)

	// 检查是否需要升级
	if err := processor.checkNeedUpgrade(); err != nil {
		log.Fatalln(err)
	}

	// 升级
	if err := processor.upgrade(); err != nil {
		log.Fatalln(err)
	}

	fmt.Println("升级完成")
}

// 检查更新并显示进度条
func (u *upgradeProcessor) checkForUpdate() error {
	bar := progressbar.NewOptions(1,
		progressbar.OptionSetDescription("检查更新中..."),
		progressbar.OptionShowBytes(false),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetTheme(progressbar.Theme{}),
		progressbar.OptionOnCompletion(func() {
			fmt.Fprint(os.Stderr, "\n")
		}),
		progressbar.OptionSetRenderBlankState(true),
	)

	// 获取最新版本信息
	if err := u.getReleaseLatest(); err != nil {
		return errors.New(fmt.Sprintf("获取新版本信息错误：%s", err))
	}

	if err := bar.Finish(); err != nil {
		return err
	}

	return nil
}

// 获取最新版本
func (u *upgradeProcessor) getReleaseLatest() error {
	url := "https://api.github.com/repos/bynow2code/geeksaver/releases/latest"
	client := http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("检查新版本错误: %s", resp.Status))
	}

	newBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(newBytes, &u.releaseLatest); err != nil {
		return err
	}

	return nil
}

// 检查是否需要升级
func (u *upgradeProcessor) checkNeedUpgrade() error {
	need, err := u.needUpgrade()
	if err != nil {
		return err
	}
	if !need {
		return errors.New("当前版本已经是最新版本")
	}

	// go install 安装或 dev 不允许命令行更新
	if version == devVersion {
		return errors.New("请使用 go install github.com/bynow2code/geeksaver@latest 的形式安装最新版")
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

// 匹配升级文件下载地址
func (u *upgradeProcessor) matchUpgradeUrl() error {
	upgradeFile, err := u.getPlatformBinaryName()
	if err != nil {
		return err
	}

	for i, asset := range u.releaseLatest.Assets {
		if upgradeFile == asset.Name {
			u.assert = &u.releaseLatest.Assets[i]
			return nil
		}
	}

	return errors.New("没有匹配的升级文件")
}

// 获取平台对应的二进制文件名
func (u *upgradeProcessor) getPlatformBinaryName() (string, error) {
	switch runtime.GOOS {
	case "darwin":
		return fmt.Sprintf("geeksaver-%s-macos-%s", u.releaseLatest.TagName, runtime.GOARCH), nil
	case "linux":
		return fmt.Sprintf("geeksaver-%s-linux-%s", u.releaseLatest.TagName, runtime.GOARCH), nil
	case "windows":
		return fmt.Sprintf("geeksaver-%s-windows-%s.exe", u.releaseLatest.TagName, runtime.GOARCH), nil
	default:
		return "", errors.New("不支持的操作系统")
	}
}

// 升级
func (u *upgradeProcessor) upgrade() error {
	bar := progressbar.NewOptions(1,
		progressbar.OptionSetDescription("升级中..."),
		progressbar.OptionShowBytes(false),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetTheme(progressbar.Theme{}),
		progressbar.OptionOnCompletion(func() {
			fmt.Fprint(os.Stderr, "\n")
		}),
		progressbar.OptionSetRenderBlankState(true),
	)

	// 匹配升级文件下载地址
	if err := u.matchUpgradeUrl(); err != nil {
		return err
	}

	// 升级文件
	if err := u.upgradeFile(); err != nil {
		return err
	}

	if err := bar.Finish(); err != nil {
		return err
	}

	return nil
}

// 升级文件
func (u *upgradeProcessor) upgradeFile() error {
	// 下载升级文件
	client := http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(u.assert.BrowserDownloadUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	newBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	exePath, err := os.Executable()
	if err != nil {
		return err
	}
	targetPath, err := filepath.EvalSymlinks(exePath)
	if err != nil {
		return err
	}

	updateDir := filepath.Dir(targetPath)
	filename := filepath.Base(targetPath)

	// 创建新升级文件
	newPath := filepath.Join(updateDir, fmt.Sprintf(".%s.new", filename))
	newFile, err := os.OpenFile(newPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer newFile.Close()

	if _, err := io.Copy(newFile, bytes.NewReader(newBytes)); err != nil {
		return err
	}

	// windows 打开文件后，需要先关闭，然后才能才能移动这个文件
	newFile.Close()

	// 删除旧文件
	// Windows 如果重命名操作的目标文件已存在，重命名操作将会失败
	oldPath := filepath.Join(updateDir, fmt.Sprintf(".%s.old", filename))
	_ = os.Remove(oldPath)
	// 备份旧文件
	if err := os.Rename(targetPath, oldPath); err != nil {
		return err
	}

	// 替换新文件
	if err := os.Rename(newPath, targetPath); err != nil {
		// 升级没有成功，恢复旧文件
		_ = os.Rename(oldPath, targetPath)
		return fmt.Errorf("升级失败（已回滚）：%w", err)
	}

	// 删除旧文件
	// Windows 更新后进程仍在运行导致无法删除旧文件
	_ = os.Remove(oldPath)

	return nil
}
