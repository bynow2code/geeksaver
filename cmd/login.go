package cmd

import (
	"log"

	"github.com/bynow2code/geekbangdocsaver/internal/geekbang"
	"github.com/spf13/cobra"
)

type LoginConfig struct {
	// 极客时间登陆后用户id
	GCID string
	// 极客时间登陆后用户令牌
	Gcess string
}

var loginCfg LoginConfig

func init() {
	loginCmd.Flags().StringVar(&loginCfg.GCID, "gcid", "", "极客时间登陆后用户id（必填）")
	loginCmd.Flags().StringVar(&loginCfg.Gcess, "gcess", "", "极客时间登陆后用户令牌（必填）")
	err := loginCmd.MarkFlagRequired("gcid")
	if err != nil {
		log.Fatalln(err)
	}
	loginCmd.MarkFlagsRequiredTogether("gcid", "gcess")
	rootCmd.AddCommand(loginCmd)
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "设置极客时间登陆态",
	Long:  `设置极客时间登陆态`,
	Run: func(cmd *cobra.Command, args []string) {
		config := &geekbang.Config{User: geekbang.UserConfig{
			GCID:  loginCfg.GCID,
			GCESS: loginCfg.Gcess,
		}}
		geekbang.GetViper().WriteConfig(config)
	},
}
