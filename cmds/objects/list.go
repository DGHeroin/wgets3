package objects

import (
    "github.com/DGHeroin/wgets3/store"
    "github.com/DGHeroin/wgets3/utils"
    "github.com/dustin/go-humanize"
    "github.com/spf13/cobra"
)

var listObject = &cobra.Command{
    Use:   "list",
    Short: "列出对象",
    RunE: func(cmd *cobra.Command, args []string) error {
        var key string
        if len(args) >= 1 {
            key = args[0]
        }
        return _listObject(key)
    },
}

func _listObject(prefix string) error {
    s, err := store.GetConfig(configNode)
    if err != nil {
        return err
    }
    i := 0
    s.List(bucket, prefix, func(info store.ObjectInfo) bool {
        i++
        utils.LogI("%d. %v\t\t%s %s\n", i, humanize.Bytes(uint64(info.Size)), info.Key, info.ContentType)
        if maxFiles > 0 && i >= maxFiles {
            return false
        }
        return true
    })
    return nil
}
