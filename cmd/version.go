package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const devVersion = "0.0.0-dev"

// 默认值为开发版（go install 直接安装时的状态）
var version = devVersion

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "当前版本",
	Long:  `当前版本`,
	Run: func(cmd *cobra.Command, args []string) {
		doVersion()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

// 获取版本信息
func doVersion() {
	fmt.Printf("当前版本：%s\n", version)
}
