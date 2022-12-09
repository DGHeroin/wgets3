package cmds

import (
    "github.com/DGHeroin/wgets3/cmds/bucket"
    "github.com/DGHeroin/wgets3/cmds/download"
    "github.com/DGHeroin/wgets3/cmds/objects"
    "github.com/DGHeroin/wgets3/cmds/s2ssync"
    "github.com/DGHeroin/wgets3/cmds/upload"
    "github.com/DGHeroin/wgets3/utils"
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
    "os"
)

var (
    rootCmd = &cobra.Command{
        Use: "wgets3",
        Long: `wget + s3
`,
    }
)

var (
    confFile string
    // bucketName string
)

func init() {
    rootCmd.PersistentFlags().StringVar(&confFile, "c", "", "配置文件")
    // rootCmd.PersistentFlags().StringVar(&bucketName, "bucket", "", "配置文件")

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

func Run() {
    rootCmd.AddCommand(download.Cmd, upload.Cmd, s2ssync.Cmd)
    rootCmd.AddCommand(bucket.Cmd)
    rootCmd.AddCommand(objects.Cmd)
    if err := rootCmd.Execute(); err != nil {
        utils.LogE("%v", err)
        os.Exit(-1)
    }
}
