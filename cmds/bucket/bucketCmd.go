package bucket

import (
    "github.com/DGHeroin/wgets3/cmds/bucket/bucket_create"
    "github.com/DGHeroin/wgets3/cmds/bucket/bucket_list"
    "github.com/DGHeroin/wgets3/cmds/bucket/bucket_remove"
    "github.com/spf13/cobra"
)

var (
    Cmd = &cobra.Command{
        Use:   "bucket",
        Short: "桶管理命令",
    }
)

var (
    configNode string
)

func init() {
    Cmd.AddCommand(bucket_list.BucketListCmd, bucket_create.BucketCreateCmd, bucket_remove.BucketRemoveCmd)
    Cmd.PersistentFlags().StringVar(&configNode, "n", "s3", "配置节点名称")
}
