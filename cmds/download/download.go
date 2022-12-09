package download

import (
    "github.com/DGHeroin/wgets3/store"
    "github.com/DGHeroin/wgets3/utils"
    "github.com/spf13/cobra"
    "os"
    "path"
    "sync"
    "time"
)

var (
    Cmd = &cobra.Command{
        Use:   "download command <args>",
        Short: "从s3下载",
        RunE: func(cmd *cobra.Command, args []string) error {
            return doDownload()
        },
    }
)
var (
    bucket    string
    saveDir   string
    s3prefix  string
    NodeName  string
    workerNum int
)
var (
    s *store.Store
)

func init() {
    Cmd.PersistentFlags().StringVar(&bucket, "b", "", "桶")
    Cmd.PersistentFlags().StringVar(&saveDir, "dir", "", "保存到本地dir")
    Cmd.PersistentFlags().StringVar(&s3prefix, "prefix", "", "s3前缀")
    Cmd.PersistentFlags().StringVar(&NodeName, "n", "s3", "配置节点名")
    Cmd.PersistentFlags().IntVar(&workerNum, "worker", 100, "配置节点名")
}

func doDownload() error {
    var (
        wg sync.WaitGroup
        ch chan *store.ObjectInfo
    )
    err := os.MkdirAll(saveDir, os.ModePerm)
    if err != nil {
        return err
    }
    s, err = store.GetConfig(NodeName)
    if err != nil {
        return err
    }

    for i := 0; i < workerNum; i++ {
        go downloadFile(ch, &wg)
    }

    s.List(bucket, s3prefix, func(info store.ObjectInfo) bool {
        ch <- &info
        return true
    })
    wg.Wait()
    return nil
}

func downloadFile(ch chan *store.ObjectInfo, wg *sync.WaitGroup) {
    for v := range ch {
        wg.Add(1)
        go func(info *store.ObjectInfo) {
            defer wg.Done()

            startTime := time.Now()
            savePath := path.Join(saveDir, info.Key)
            dir := path.Dir(savePath)
            err := os.MkdirAll(dir, os.ModePerm)
            if err != nil {
                utils.LogE("创建文件夹失败:%v\n", err)
                return
            }
            file, err := os.Create(savePath)
            if err != nil {
                utils.LogE("创建文件失败:%v\n", err)
                return
            }
            defer file.Close()
            err = s.Download(bucket, info.Key, file)
            if err != nil {
                utils.LogE("下载文件失败:%v\n", err)
                return
            }
            utils.LogI("下载成功 耗时%v 文件:%v\n", time.Since(startTime), savePath)
        }(v)
    }
}
