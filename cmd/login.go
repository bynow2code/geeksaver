package cmd

import (
	"log"

	"github.com/bynow2code/geekbangdocsaver/internal/geekbang/config"
	"github.com/spf13/cobra"
)

// LoginFlag 登陆配置 flag
type LoginFlag struct {
	GCID  string // 用户id
	Gcess string // 用户令牌
}

var loginFlg LoginFlag

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "设置极客时间登陆态",
	Long:  `设置极客时间登陆态后，才可以使用此工具`,
	Run: func(cmd *cobra.Command, args []string) {
		// 设置配置项
		config.SetGCID(loginFlg.GCID)
		config.SetGCESS(loginFlg.Gcess)
		newConfig := config.GetConfig()
		// 持久化到配置文件
		err := config.WriteConfig(newConfig)
		if err != nil {
			log.Fatalln(err)
		}
	},
}

func init() {
	loginCmd.Flags().StringVar(&loginFlg.GCID, "gcid", "", "用户id（必填）")
	loginCmd.Flags().StringVar(&loginFlg.Gcess, "gcess", "", "用户令牌（必填）")
	err := loginCmd.MarkFlagRequired("gcid")
	if err != nil {
		log.Fatalln(err)
	}
	loginCmd.MarkFlagsRequiredTogether("gcid", "gcess")
	rootCmd.AddCommand(loginCmd)
}
