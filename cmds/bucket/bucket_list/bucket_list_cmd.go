package bucket_list

import (
    "github.com/DGHeroin/wgets3/store"
    "github.com/DGHeroin/wgets3/utils"
    "github.com/spf13/cobra"
)

var (
    BucketListCmd = &cobra.Command{
        Use:   "list",
        Short: "列出桶",
        RunE: func(cmd *cobra.Command, args []string) error {
            s, err := store.GetConfig(configNode)
            if err != nil {
                return err
            }
            buckets, err := s.BucketList()
            if err != nil {
                return err
            }
            for _, bucket := range buckets {
                utils.LogI("node:%s bucket:%s\n", configNode, bucket)
            }
            return nil
        },
    }
)

var (
    configNode string
)

func init() {
    BucketListCmd.PersistentFlags().StringVar(&configNode, "n", "s3", "配置节点名称")
}
