package upload

import (
    "bytes"
    "github.com/DGHeroin/ActorSystem/actor"
    "github.com/DGHeroin/wgets3/store"
    "github.com/DGHeroin/wgets3/utils"
    "github.com/spf13/cobra"
    "io"
    "io/fs"
    "io/ioutil"
    "net/http"
    "net/url"
    "path"
    "path/filepath"
    "strings"
    "time"
)

var (
    Cmd = &cobra.Command{
        Use:   "upload command <args>",
        Short: "下载并上传到s3",
        RunE: func(cmd *cobra.Command, args []string) error {
            return doUpload()
        },
    }
)

var (
    bucket     string
    InputFile  string
    InputDir   string
    OutputPath string
    NodeName   string
)
var (
    s *store.Store
)

func init() {
    Cmd.PersistentFlags().StringVar(&bucket, "b", "", "桶")
    Cmd.PersistentFlags().StringVar(&InputFile, "file", "", "wget 输入文件")
    Cmd.PersistentFlags().StringVar(&InputDir, "dir", "", "wget 上传dir")
    Cmd.PersistentFlags().StringVar(&OutputPath, "prefix", "", "s3保存目录前缀")
    Cmd.PersistentFlags().StringVar(&NodeName, "node", "s3", "配置节点名")
}

func doUpload() error {
    var err error
    s, err = store.GetConfig(NodeName)
    if err != nil {
        return err
    }
    if InputFile != "" {
        sys := actor.NewSystem("download_sys", &actor.Config{
            MinActor:          3,
            MaxActor:          20,
            ActorQueueSize:    10,
            DispatchQueueSize: 500,
            DispatchBlocking:  true,
        })
        sys.Start()

        _ = utils.FileReadLine(InputFile, func(line int, content string) bool {
            _ = sys.Dispatch(&DownloadMsg{
                id:  line,
                url: content,
            })
            return true
        })
        sys.Stop()
    }
    if InputDir != "" {
        sys := actor.NewSystem("download_sys", &actor.Config{
            MinActor:          3,
            MaxActor:          20,
            ActorQueueSize:    10,
            DispatchQueueSize: 500,
            DispatchBlocking:  true,
        })
        sys.Start()
        _ = filepath.WalkDir(InputDir, func(path string, d fs.DirEntry, err error) error {
            if d.IsDir() {
                return nil
            }
            absPath := strings.ReplaceAll(path, InputDir, "")
            if absPath == "" {
                return nil
            }
            _ = sys.Dispatch(&SyncDirMsg{
                filepath: absPath,
                fullPath: path,
            })
            return nil
        })
        sys.Stop()
    }
    return nil
}

type DownloadMsg struct {
    id  int
    url string
}
type SyncDirMsg struct {
    filepath string
    fullPath string
}

func (d *DownloadMsg) Execute(actor.Context) {
    startTime := time.Now()
    content := d.url
    u, err := url.Parse(content)
    if err != nil {
        utils.LogE("解析错误:%v\n", err)
        return
    }
    filename := path.Base(u.Path)
    key := path.Join(OutputPath, filename)
    // 先检查s3是否存在, 存在就不下载
    st, err := s.Exist(bucket, key)
    if st.Size > 0 {
        return // 无需下载
    }
    if err != nil {
        utils.LogE("检查失败:%v\n", err)
        return
    }

    _ = utils.HTTPDownload(d.url, func(header http.Header, r io.Reader) {
        size := int64(0)
        data, err := ioutil.ReadAll(r)
        if err != nil {
            return
        }
        size = int64(len(data))
        for {
            if err := s.Upload(bucket, path.Join(OutputPath, filename), bytes.NewBuffer(data), size); err != nil {
                utils.LogE("上传失败:%v 错误:%v\n", filename, err)
                return
            } else {
                utils.LogI("上传成功: %v %v 文件:%v \n", time.Since(startTime), utils.HumanBytesSize(float64(size)), filename)
                break
            }
        }
    })
}
func (d *SyncDirMsg) Execute(actor.Context) {
    startTime := time.Now()
    key := path.Join(OutputPath, d.filepath)
    // 先检查s3是否存在, 存在就不下载
    st, err := s.Exist(bucket, key)
    if st.Size > 0 {
        return // 无需下载
    }
    if err != nil {
        utils.LogE("check exist error:%v\n", err)
        return
    }

    for {
        if err := s.UploadFile(bucket, key, d.fullPath); err != nil {
            utils.LogE("上传失败 key:% 耗时:%v 错误:%v\n", key, time.Since(startTime), err)
            return
        } else {
            utils.LogI("上传成功 key:%v 耗时:%v\n", key, time.Since(startTime))
            break
        }
    }
}
