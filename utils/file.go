package utils

import (
    "bufio"
    "io"
    "os"
)

func FileReadLine(path string, fn func(line int, content string) bool) error {
    f, err := os.Open(path)
    if err != nil {
        return err
    }
    defer func() {
        _ = f.Close()
    }()
    return ReadLine(f, fn)
}
func ReadLine(r io.Reader, fn func(line int, content string) bool) error {
    buf := bufio.NewReader(r)
    i := 0
    for {
        line, _, err := buf.ReadLine()
        if err != nil {
            if err == io.EOF {
                break
            }
            return err
        }
        i++
        if !fn(i, string(line)) {
            break
        }
    }
    return nil
}
