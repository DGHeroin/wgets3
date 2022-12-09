package utils

import (
    "io"
    "net/http"
)

func getClient() *http.Client {
    return http.DefaultClient
}

func HTTPDownload(downloadUrl string, fn func(statusCode int, header http.Header, r io.Reader)) error {
    cli := getClient()
    req, err := http.NewRequest(http.MethodGet, downloadUrl, nil)
    if err != nil {
        return err
    }
    resp, err := cli.Do(req)
    if err != nil {
        return err
    }
    fn(resp.StatusCode, resp.Header, resp.Body)
    return nil
}
