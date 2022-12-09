package bucket_create

import (
    "github.com/DGHeroin/wgets3/store"
    "github.com/spf13/cobra"
)

var (
    BucketCreateCmd = &cobra.Command{
        Use:   "create <args>",
        Short: "创建桶",
        RunE: func(cmd *cobra.Command, args []string) error {
            s, err := store.GetConfig(configNode)
            if err != nil {
                return err
            }
            return s.BucketCreate(name)
        },
    }
)

var (
    configNode string
    name       string
)

func init() {
    BucketCreateCmd.PersistentFlags().StringVar(&configNode, "n", "s3", "配置节点名称")
    BucketCreateCmd.PersistentFlags().StringVar(&name, "b", "bucket_name", "桶名称")
}
