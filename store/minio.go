package store

import (
    "github.com/minio/minio-go/v7"
    "github.com/minio/minio-go/v7/pkg/credentials"
    "github.com/spf13/viper"
    "log"
)

func GetConfig(prefixs ...string) *Store {
    s := &Store{}
    prefix := "s3"
    if len(prefixs) == 1 {
        prefix = prefixs[0]
    }
    s.bucket = viper.GetString(prefix + ".bucket")
    var endpoint = viper.GetString(prefix + ".endpoint")
    var accessKeyId = viper.GetString(prefix + ".keyId")
    var accessKeySecret = viper.GetString(prefix + ".keySecret")
    var region = viper.GetString(prefix + ".region")
    if region == "" {
        region = "auto"
    }

    minioClient, err := minio.New(endpoint, &minio.Options{
        Region: region,
        Creds:  credentials.NewStaticV4(accessKeyId, accessKeySecret, ""),
        Secure: true,
    })
    if err != nil {
        log.Fatal(err)
    }
    s.client = minioClient
    return s
}
