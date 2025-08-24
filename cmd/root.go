package cmd

import (
	"fmt"
	"log"

	"github.com/bynow2code/geekbangdocsaver/internal/geekbang/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "geek",
	Short: "极客时间课程保存工具",
	Long:  "极客时间课程保存工具，该脚本仅供学习使用",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// 读取程序配置文件
		err := config.ReadInConfig()
		if err != nil {
			log.Fatalln(err)
		}

		//绑定配置文件到结构体
		err = config.Unmarshal()
		if err != nil {
			log.Fatalln(err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(cmd.Usage())
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
