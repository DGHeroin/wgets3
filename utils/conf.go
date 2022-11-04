package utils

import (
    "fmt"
    "strconv"
    "strings"
)

type Conf map[string]interface{}

func (c Conf) String(k string) string {
    v, _ := c[k]
    switch val := v.(type) {
    case string:
        return val
    case int8, int16, int32, int64, uint8, uint16, uint32, uint64, int, uint, float32, float64:
        return fmt.Sprint(c)
    default:
        return ""
    }
}

func (c Conf) Int(k string) int64 {
    v, _ := c[k]
    if !isNumber(v) {
        return 0
    }
    val, _ := toInt64(v)
    return val
}
func (c Conf) Float(k string) float64 {
    v, _ := c[k]
    if !isNumber(v) {
        return 0
    }
    val, _ := toFloat(v)
    return val
}
func (c Conf) Bool(k string) bool {
    v, _ := c[k]
    val, _ := toBool(v)
    return val
}
func (c Conf) Object(k string) Conf {
    v, _ := c[k]
    switch val := v.(type) {
    case Conf:
        return val
    case map[string]interface{}:
        return val
    }
    return nil
}

func isNumber(v interface{}) bool {
    switch val := v.(type) {
    case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr:
        return true
    case float32, float64:
        return true
    case string:
        if _, err := strconv.ParseFloat(val, 64); err == nil {
            return true
        }
        return false
    default:
        return false
    }
}
func toInt64(v interface{}) (int64, bool) {
    switch u := v.(type) {
    case int:
        return int64(u), true
    case int8:
        return int64(u), true
    case int16:
        return int64(u), true
    case int32:
        return int64(u), true
    case int64:
        return u, true
    case uint:
        return int64(u), true
    case uint8:
        return int64(u), true
    case uint16:
        return int64(u), true
    case uint32:
        return int64(u), true
    case uint64:
        return int64(u), true
    case uintptr:
        return int64(u), true
    case float32:
        return int64(u), true
    case float64:
        return int64(u), true
    case string:
        if i, err := strconv.ParseInt(u, 10, 64); err == nil {
            return i, true
        }
        return 0, false
    default:
        return 0, false
    }
}
func toFloat(v interface{}) (float64, bool) {
    switch u := v.(type) {
    case int:
        return float64(u), true
    case int8:
        return float64(u), true
    case int16:
        return float64(u), true
    case int32:
        return float64(u), true
    case int64:
        return float64(u), true
    case uint:
        return float64(u), true
    case uint8:
        return float64(u), true
    case uint16:
        return float64(u), true
    case uint32:
        return float64(u), true
    case uint64:
        return float64(u), true
    case uintptr:
        return float64(u), true
    case float32:
        return float64(u), true
    case float64:
        return u, true
    case string:
        if f, err := strconv.ParseFloat(u, 64); err == nil {
            return f, true
        }
        return 0, false
    default:
        return 0, false
    }
}
func toBool(v interface{}) (bool, bool) {
    if isNumber(v) {
        val, _ := toInt64(v)
        return val == 1, true
    }
    switch val := v.(type) {
    case bool:
        return val, true
    case string:
        val = strings.ToLower(val)
        if val == "true" {
            return true, true
        } else if val == "false" {
            return false, true
        } else {
            return false, false
        }
    default:
        return false, false
    }
}
