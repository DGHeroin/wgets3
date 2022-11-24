package list

import (
    "github.com/DGHeroin/wgets3/store"
    "github.com/DGHeroin/wgets3/utils"
    "github.com/spf13/cobra"
)

var (
    Cmd = &cobra.Command{
        Use:   "list command <args>",
        Short: "从s3列出文件",
        RunE: func(cmd *cobra.Command, args []string) error {
            var key string
            if len(args) >= 1 {
                key = args[0]
            }
            return listS3(key)
        },
    }
)

var (
    configNode string
    maxFiles   int
    bucket     string
)

func init() {
    Cmd.PersistentFlags().StringVar(&bucket, "b", "", "桶")
    Cmd.PersistentFlags().StringVar(&configNode, "n", "s3", "配置节点名称")
    Cmd.PersistentFlags().IntVar(&maxFiles, "max", 10, "列出最大文件数量")
}

func listS3(prefix string) error {
    s, err := store.GetConfig(configNode)
    if err != nil {
        return err
    }
    i := 0
    s.List(bucket, prefix, func(info store.ObjectInfo) bool {
        i++
        utils.LogI("%d. %v\t%s %s\n", i, utils.HumanSize(float64(info.Size)), info.Key, info.ContentType)
        if i >= maxFiles {
            return false
        }
        return true
    })
    return nil
}
