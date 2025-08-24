package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "geek",
	Short: "极客时间课程保存工具，该脚本仅供学习使用",
	Long:  "极客时间课程保存工具，该脚本仅供学习使用",
	Run: func(cmd *cobra.Command, args []string) {
		//lessons, err := geekbang.GetLessons(geekbang.LessonReq{
		//	Cid:    "100093501",
		//	Size:   500,
		//	Prev:   0,
		//	Order:  "earliest",
		//	Sample: false,
		//})
		//if err != nil {
		//	panic(err)
		//}
		//fmt.Println(lessons)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
