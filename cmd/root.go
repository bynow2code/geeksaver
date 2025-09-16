package cmd

import (
	"fmt"
	"log"

	"github.com/bynow2code/geeksaver/internal/geekbang/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Short: "极客时间课程本地保存工具",
	Long:  "极客时间课程本地保存工具，该程序仅供学习使用",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		initConfig()
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

// 初始化配置文件
func initConfig() {
	vp := config.GetViper()

	// 读取程序配置文件
	if err := vp.ReadInConfig(); err != nil {
		log.Fatalln(err)
	}

	// 绑定配置文件到结构体
	if err := vp.Unmarshal(config.GetConfig()); err != nil {
		log.Fatalln(err)
	}

}
