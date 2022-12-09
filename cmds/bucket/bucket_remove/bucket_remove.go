package bucket_remove

import (
    "github.com/DGHeroin/wgets3/store"
    "github.com/spf13/cobra"
)

var (
    BucketRemoveCmd = &cobra.Command{
        Use:   "remove",
        Short: "删除桶",
        RunE: func(cmd *cobra.Command, args []string) error {
            s, err := store.GetConfig(configNode)
            if err != nil {
                return err
            }
            return s.BucketRemove(name)
        },
    }
)

var (
    configNode string
    name       string
)

func init() {
    BucketRemoveCmd.PersistentFlags().StringVar(&configNode, "n", "s3", "配置节点名称")
    BucketRemoveCmd.PersistentFlags().StringVar(&name, "b", "bucket_name", "桶名称")
}
