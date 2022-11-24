package store

import (
    "context"
    "github.com/minio/minio-go/v7"
    "github.com/minio/minio-go/v7/pkg/credentials"
    "github.com/spf13/viper"
    "log"
)

func GetConfig(prefixs ...string) (*Store, error) {
    s := &Store{}
    prefix := "s3"
    if len(prefixs) == 1 {
        prefix = prefixs[0]
    }
    s.defaultBucket = viper.GetString(prefix + ".bucket")
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
        return nil, err
    }
    buckets, err := minioClient.ListBuckets(context.Background())
    log.Println("===>", err)
    log.Println("===>", buckets)
    s.client = minioClient
    return s, nil
}
