package cmd

import (
	"log"
	"os"

	"github.com/bynow2code/geekbangdocsaver/internal/geekbang/config"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "配置信息",
	Long:  `配置信息`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.GetConfig()

		// 组装 table 数据
		var data [][]string
		data = append(data, []string{"user", "gcid", cfg.User.GCID})
		data = append(data, []string{"", "gcess", cfg.User.GCESS})
		data = append(data, []string{"md", "savepath", cfg.Md.SavePath})

		// 添加 table 头部
		table := tablewriter.NewWriter(os.Stdout)
		table.Header([]string{"Group", "Key", "Value"})

		// 添加 table 数据
		err := table.Bulk(data)
		if err != nil {
			log.Fatalln(err)
		}

		// 渲染
		if err = table.Render(); err != nil {
			log.Fatalln(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
