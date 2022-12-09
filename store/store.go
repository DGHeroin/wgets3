package store

import (
    "context"
    "github.com/minio/minio-go/v7"
    "io"
    "strings"
    "time"
)

type (
    Store struct {
        defaultBucket string
        client        *minio.Client
        region        string
    }
)
type (
    ObjectInfo struct {
        Key            string    `json:"key,omitempty"`
        LastModified   time.Time `json:"lastModified,omitempty"`
        Size           int64     `json:"size,omitempty"`
        ContentType    string    `json:"contentType,omitempty"`
        Expires        time.Time `json:"expires,omitempty"`
        ChecksumCRC32  string    `json:"checksumCRC32,omitempty"`
        ChecksumCRC32C string    `json:"checksumCRC32C,omitempty"`
        ChecksumSHA1   string    `json:"checksumSHA1,omitempty"`
        ChecksumSHA256 string    `json:"checksumSHA256,omitempty"`
    }
)

func or(a, b string) string {
    if a != "" {
        return a
    }
    return b
}
func (s *Store) Exist(bucket, key string) (*ObjectInfo, error) {
    bucket = or(bucket, s.defaultBucket)
    st, err := s.client.StatObject(context.Background(), bucket, key, minio.StatObjectOptions{})
    if err != nil {
        if minio.ToErrorResponse(err).StatusCode == 404 {
            return &ObjectInfo{}, nil
        }
        return &ObjectInfo{}, err
    }
    result := &ObjectInfo{
        Key:            st.Key,
        LastModified:   st.LastModified,
        Size:           st.Size,
        ContentType:    st.ContentType,
        Expires:        st.Expires,
        ChecksumCRC32:  st.ChecksumCRC32,
        ChecksumCRC32C: st.ChecksumCRC32C,
        ChecksumSHA1:   st.ChecksumSHA1,
        ChecksumSHA256: st.ChecksumSHA256,
    }

    return result, nil
}

func (s *Store) Upload(bucket, key string, r io.Reader, objectSize int64) error {
    bucket = or(bucket, s.defaultBucket)
    _, err := s.client.PutObject(context.Background(), bucket, key, r, objectSize, minio.PutObjectOptions{
        SendContentMd5: true,
    })
    return err
}

func (s *Store) UploadFile(bucket, key string, filename string) error {
    bucket = or(bucket, s.defaultBucket)
    _, err := s.client.FPutObject(context.Background(), bucket, key, filename, minio.PutObjectOptions{
        SendContentMd5: true,
    })
    return err
}

func (s *Store) List(bucket, prefix string, fn func(ObjectInfo) bool) {
    bucket = or(bucket, s.defaultBucket)
    ch := s.client.ListObjects(context.Background(), bucket, minio.ListObjectsOptions{
        Prefix:    prefix,
        Recursive: true,
        // MaxKeys:   10000,
    })

    for st := range ch {
        if st.Err != nil {
            return
        }
        if strings.HasSuffix(st.Key, "/") {
            continue
        }
        info := ObjectInfo{
            Key:            st.Key,
            LastModified:   st.LastModified,
            Size:           st.Size,
            ContentType:    st.ContentType,
            Expires:        st.Expires,
            ChecksumCRC32:  st.ChecksumCRC32,
            ChecksumCRC32C: st.ChecksumCRC32C,
            ChecksumSHA1:   st.ChecksumSHA1,
            ChecksumSHA256: st.ChecksumSHA256,
        }
        if !fn(info) {
            break
        }
    }

}
func (s *Store) RemoveBucket(bucket string) error {
    bucket = or(bucket, s.defaultBucket)
    return s.client.RemoveBucket(context.Background(), bucket)
}
func (s *Store) Remove(bucket, key string) error {
    bucket = or(bucket, s.defaultBucket)
    return s.client.RemoveObject(context.Background(), bucket, key, minio.RemoveObjectOptions{})
}
func (s *Store) Download(bucket, key string, w io.Writer) error {
    bucket = or(bucket, s.defaultBucket)
    obj, err := s.client.GetObject(context.Background(), bucket, key, minio.GetObjectOptions{})
    if err != nil {
        return err
    }
    _, err = io.Copy(w, obj)
    return err
}
func (s *Store) BucketList() ([]string, error) {
    info, err := s.client.ListBuckets(context.Background())
    if err != nil {
        return nil, err
    }
    var result []string
    for _, bucketInfo := range info {
        result = append(result, bucketInfo.Name)
    }
    return result, nil
}
func (s *Store) BucketCreate(name string) error {
    return s.client.MakeBucket(context.Background(), name, minio.MakeBucketOptions{
        Region: s.region,
    })
}
func (s *Store) BucketRemove(name string) error {
    return s.client.RemoveBucket(context.Background(), name)
}
