package store

import (
    "context"
    "github.com/minio/minio-go/v7"
    "io"
    "time"
)

type (
    Store struct {
        bucket string
        client *minio.Client
    }
)
type (
    ExistResult struct {
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

func (s *Store) Exist(key string) (*ExistResult, error) {
    st, err := s.client.StatObject(context.Background(), s.bucket, key, minio.StatObjectOptions{})
    if err != nil {
        if minio.ToErrorResponse(err).StatusCode == 404 {
            return &ExistResult{}, nil
        }
        return &ExistResult{}, err
    }
    result := &ExistResult{
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

func (s *Store) Upload(key string, r io.Reader, objectSize int64) error {
    _, err := s.client.PutObject(context.Background(), s.bucket, key, r, objectSize, minio.PutObjectOptions{
        SendContentMd5: true,
    })
    return err
}

func (s *Store) UploadFile(key string, filename string) error {
    _, err := s.client.FPutObject(context.Background(), s.bucket, key, filename, minio.PutObjectOptions{
        SendContentMd5: true,
    })
    return err
}
