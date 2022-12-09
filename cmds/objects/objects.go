package objects

import "github.com/spf13/cobra"

var (
    Cmd = &cobra.Command{
        Use:   "objects",
        Short: "对象管理",
    }
)
var (
    configNode string
    maxFiles   int
    bucket     string
)

func init() {
    Cmd.AddCommand(listObject, countObject, deleteObject)
}

func init() {
    Cmd.PersistentFlags().StringVar(&bucket, "b", "", "桶")
    Cmd.PersistentFlags().StringVar(&configNode, "n", "s3", "配置节点名称")
    Cmd.PersistentFlags().IntVar(&maxFiles, "max", 10, "列出最大文件数量")
}
