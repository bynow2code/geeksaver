package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"time"

	"github.com/spf13/cobra"
)

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "下载安装最新版",
	Long:  `下载安装最新版`,
	Run: func(cmd *cobra.Command, args []string) {
		processor := &upgradeProcessor{}
		processor.getReleaseLatest()
		fmt.Printf("最新版本: %s\n", processor.releaseLatest.TagName)
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

func getPlatform() string {
	return fmt.Sprintf("%s_%s", runtime.GOOS, runtime.GOARCH)
}

func doUpdate() {

}
