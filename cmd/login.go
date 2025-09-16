package cmd

import (
	"log"

	"github.com/bynow2code/geeksaver/internal/geekbang/config"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "极客时间登录态设置",
	Long:  `极客时间登录态设置，设置好后工具才可以正常使用`,
	Run: func(cmd *cobra.Command, args []string) {
		doLogin(cmd)
	},
}

func init() {
	loginCmd.Flags().String("gcid", "", "用户id（必填）")
	loginCmd.Flags().String("gcess", "", "用户令牌（必填）")

	// 添加必填参数
	if err := loginCmd.MarkFlagRequired("gcid"); err != nil {
		log.Fatalln(err)
	}
	// 添加必填参数组合
	loginCmd.MarkFlagsRequiredTogether("gcid", "gcess")

	rootCmd.AddCommand(loginCmd)
}

// 登录
func doLogin(cmd *cobra.Command) {
	gcid, err := cmd.Flags().GetString("gcid")
	if err != nil {
		log.Fatalln(err)
	}
	gcess, err := cmd.Flags().GetString("gcess")
	if err != nil {
		log.Fatalln(err)
	}

	// 设置配置项
	config.SetGCID(gcid)
	config.SetGCESS(gcess)

	// 持久化到配置文件
	if err := config.WriteConfig(config.GetConfig()); err != nil {
		log.Fatalln(err)
	}
}
