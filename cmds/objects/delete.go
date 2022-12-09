package objects

import (
    "fmt"
    "github.com/DGHeroin/wgets3/store"
    "github.com/DGHeroin/wgets3/utils"
    "github.com/dustin/go-humanize"
    "github.com/spf13/cobra"
)

var deleteObject = &cobra.Command{
    Use:   "delete",
    Short: "删除对象",
    RunE: func(cmd *cobra.Command, args []string) error {
        var key string
        if len(args) >= 1 {
            key = args[0]
        }
        return _deleteObject(key)
    },
}

func _deleteObject(prefix string) error {
    s, err := store.GetConfig(configNode)
    if err != nil {
        return err
    }
    i := 0
    sz := uint64(0)
    s.List(bucket, prefix, func(info store.ObjectInfo) bool {
        fmt.Println(info.Key)
        err = s.Remove(bucket, info.Key)
        if err != nil {
            sz += uint64(info.Size)
            utils.LogE("delete [%s] error:%v", info.Key, err)
            return true
        }
        i++
        return true
    })
    utils.LogI("deleted %d total size:%v\n", i, humanize.Bytes(sz))
    return nil
}
