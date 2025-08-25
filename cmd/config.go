package cmd

import (
	"log"
	"os"

	"github.com/bynow2code/geeksaver/internal/geekbang/config"
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
		data = append(data, []string{"user", "gcid", cfg.User.GCID, "极客时间登陆用户id"})
		data = append(data, []string{"user", "gcess", cfg.User.GCESS, "极客时间登陆token"})
		data = append(data, []string{"md", "savepath", cfg.Md.SavePath, "markdown 保存路径"})

		// 添加 table 头部
		table := tablewriter.NewWriter(os.Stdout)
		table.Header([]string{"Group", "Key", "Value", "Remark"})

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
