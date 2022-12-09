package store

import (
    "github.com/minio/minio-go/v7"
    "github.com/minio/minio-go/v7/pkg/credentials"
    "github.com/spf13/viper"
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
    s.region = region

    minioClient, err := minio.New(endpoint, &minio.Options{
        Region: region,
        Creds:  credentials.NewStaticV4(accessKeyId, accessKeySecret, ""),
        Secure: true,
    })
    if err != nil {
        return nil, err
    }
    s.client = minioClient
    return s, nil
}
