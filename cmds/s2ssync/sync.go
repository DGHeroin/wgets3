package s2ssync

import (
    "github.com/DGHeroin/wgets3/store"
    "github.com/DGHeroin/wgets3/utils"
    "github.com/dustin/go-humanize"
    "github.com/spf13/cobra"
    "io"
    "os"
    "path"
    "sync"
    "sync/atomic"
)

var (
    Cmd = &cobra.Command{
        Use:   "sync command <args>",
        Short: "s3同步到s3",
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
    bucket     string
    sourceNode string
    targetNode string
    workerNum  int
)

func init() {
    Cmd.PersistentFlags().StringVar(&bucket, "b", "", "桶")
    Cmd.PersistentFlags().StringVar(&sourceNode, "s", "s3", "源节点名称")
    Cmd.PersistentFlags().StringVar(&targetNode, "t", "r2", "目的节点名称")
    Cmd.PersistentFlags().IntVar(&workerNum, "worker", 100, "工作者数量")
}

var (
    count int32
    s     *store.Store
    d     *store.Store
)

func keyHandler(ch chan *store.ObjectInfo, wg *sync.WaitGroup) {
    for v := range ch {
        wg.Add(1)
        go func(info *store.ObjectInfo) {
            key := info.Key
            defer wg.Done()
            f, err := os.CreateTemp(os.TempDir(), "wgets3_sync_")
            if err != nil {
                utils.LogE("创建临时文件错误:%v\n", err)
                return
            }
            defer func() {
                f.Close()
                os.Remove(path.Join(os.TempDir(), f.Name()))
            }()
            err = s.Download(bucket, key, f)
            if err != nil {
                utils.LogE("下载错误:%v\n", key, err)
                return
            }
            _, err = f.Seek(0, io.SeekStart)
            if err != nil {
                utils.LogE("seek start错误:%v\n", err)
                return
            }
            err = d.Upload(bucket, key, f, info.Size)
            if err != nil {
                utils.LogE("上传错误:%v\n", err)
                return
            }
            utils.LogI("上传成功 %d %v\t%s\n", atomic.AddInt32(&count, 1), humanize.Bytes(uint64(info.Size)), key)
        }(v)

    }
}

func listS3(prefix string) error {
    // 启动worker
    var wg sync.WaitGroup
    ch := make(chan *store.ObjectInfo, 500)

    for i := 0; i < workerNum; i++ {
        go keyHandler(ch, &wg)
    }
    var err error
    s, err = store.GetConfig(sourceNode)
    if err != nil {
        return err
    }
    d, err = store.GetConfig(targetNode)
    if err != nil {
        return err
    }
    i := 0
    s.List(bucket, prefix, func(info store.ObjectInfo) bool {
        i++
        utils.LogI("遍历:%v %v\n", i, info.Key)
        ch <- &info
        return true
    })
    utils.LogI("结束遍历,文件数量:%v\n", i)
    close(ch)
    wg.Wait()
    return nil
}
