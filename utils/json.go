package utils

import "encoding/json"

func JSONEncode(v interface{}) string {
    data, _ := json.Marshal(v)
    return string(data)
}

func JSONEncodePretty(v interface{}) string {
    data, _ := json.MarshalIndent(v, "", "  ")
    return string(data)
}
