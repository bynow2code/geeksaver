package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var mdCmd = &cobra.Command{
	Use:   "md",
	Short: "下载课程为 markdown 格式",
	Long:  `下载课程为 markdown 格式`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	mdCmd.Flags().StringVar(&loginCfg.GCID, "cid", "", "课程id（必填）")
	err := mdCmd.MarkFlagRequired("cid")
	if err != nil {
		log.Fatalln(err)
	}
	rootCmd.AddCommand(mdCmd)
}
