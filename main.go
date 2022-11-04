package main

import (
    "bytes"
    "flag"
    "fmt"
    "github.com/DGHeroin/ActorSystem/actor"
    "github.com/DGHeroin/wgets3/store"
    "github.com/DGHeroin/wgets3/utils"
    "github.com/spf13/viper"
    "io"
    "io/fs"
    "io/ioutil"
    "log"
    "net/http"
    "net/url"
    "path"
    "path/filepath"
    "strings"
    "time"
)

var (
    InputFile  string
    InputDir   string
    OutputPath string
    confFile   string
)

func init() {
    flag.StringVar(&InputFile, "i", "", "download input file")
    flag.StringVar(&InputDir, "dir", "", "download input dir")
    flag.StringVar(&confFile, "c", "", "config file")
    flag.StringVar(&OutputPath, "d", "/", "output to")

    flag.Parse()

    viper.SetConfigName("wgets3")
    viper.SetConfigType("toml")
    viper.AddConfigPath(".wgets3")
    viper.AddConfigPath("$HOME/.wgets3") // call multiple times to add many search paths
    viper.AddConfigPath("/etc/wgets3")
    viper.WatchConfig()
    if confFile != "" {
        viper.SetConfigFile(confFile)
    }
    if err := viper.ReadInConfig(); err != nil {
        switch err.(type) {
        case viper.ConfigFileNotFoundError:
            return
        default:
            panic(err)
        }
    }
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
        log.Println("parse url error:", err)
        return
    }
    filename := path.Base(u.Path)
    key := path.Join(OutputPath, filename)
    // 先检查s3是否存在, 存在就不下载
    st, err := s.Exist(key)
    if st.Size > 0 {
        return // 无需下载
    }
    if err != nil {
        log.Println("check exist error:", err)
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
            if err := s.Upload(path.Join(OutputPath, filename), bytes.NewBuffer(data), size); err != nil {
                fmt.Println(time.Since(startTime), filename, err)
                return
            } else {
                fmt.Println(time.Since(startTime), filename, utils.HumanBytesSize(float64(size)))
                break
            }
        }
    })
}
func (d *SyncDirMsg) Execute(actor.Context) {
    startTime := time.Now()
    key := path.Join(OutputPath, d.filepath)
    // 先检查s3是否存在, 存在就不下载
    st, err := s.Exist(key)
    if st.Size > 0 {
        return // 无需下载
    }
    if err != nil {
        log.Println("check exist error:", err)
        return
    }

    for {
        if err := s.UploadFile(key, d.fullPath); err != nil {
            fmt.Println("上传失败", key, time.Since(startTime), err)
            return
        } else {
            fmt.Println("上传成功", key, time.Since(startTime))
            break
        }
    }
}

var (
    s *store.Store
)

func main() {
    s = store.GetConfig()
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

}
