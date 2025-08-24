package cmd

import (
	"log"

	"github.com/bynow2code/geekbangdocsaver/internal/geekbang/config"
	"github.com/spf13/cobra"
)

var (
	gcid  string // 用户id
	gcess string // 用户id
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "设置极客时间登陆态",
	Long:  `设置极客时间登陆态后，才可以使用此工具`,
	Run: func(cmd *cobra.Command, args []string) {
		// 设置配置项
		config.SetGCID(gcid)
		config.SetGCESS(gcess)

		// 持久化到配置文件
		newConfig := config.GetConfig()
		err := config.WriteConfig(newConfig)
		if err != nil {
			log.Fatalln(err)
		}
	},
}

func init() {
	loginCmd.Flags().StringVar(&gcid, "gcid", "", "用户id（必填）")
	loginCmd.Flags().StringVar(&gcess, "gcess", "", "用户令牌（必填）")

	err := loginCmd.MarkFlagRequired("gcid")
	if err != nil {
		log.Fatalln(err)
	}
	loginCmd.MarkFlagsRequiredTogether("gcid", "gcess")

	rootCmd.AddCommand(loginCmd)
}
