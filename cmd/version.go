package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version = "v0.0.1-alpha"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "当前版本",
	Long:  `当前版本`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("当前版本：%s\n", Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
