package objects

import (
    "github.com/DGHeroin/wgets3/store"
    "github.com/DGHeroin/wgets3/utils"
    "github.com/dustin/go-humanize"
    "github.com/spf13/cobra"
    "time"
)

var countObject = &cobra.Command{
    Use:   "count",
    Short: "列出对象",
    RunE: func(cmd *cobra.Command, args []string) error {
        var key string
        if len(args) >= 1 {
            key = args[0]
        }
        return _countObject(key)
    },
}

func _countObject(prefix string) error {
    s, err := store.GetConfig(configNode)
    if err != nil {
        return err
    }
    i := 0
    sz := uint64(0)
    go func() {
        for {
            time.Sleep(time.Second)
            utils.LogI("counting object number:%v total size:%v\n", i, humanize.Bytes(sz))
        }
    }()

    s.List(bucket, prefix, func(info store.ObjectInfo) bool {
        i++
        sz += uint64(info.Size)
        if maxFiles > 0 && i >= maxFiles {
            return false
        }
        return true
    })
    utils.LogI("finish count object number:%v total size:%v\n", i, humanize.Bytes(sz))
    return nil
}
