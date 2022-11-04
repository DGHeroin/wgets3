package utils

import "testing"

func TestConf(t *testing.T) {
    c := Conf{"name": "my name",
        "age":    1,
        "height": 1.83,
        "score":  "4.9",
        "money":  "2_000_000",
        "isMe":   "True",
        "address": Conf{
            "country": "country",
            "street":  "street",
        }}
    t.Log("name", c.String("name"))
    t.Log("address:", c.String("address"))
    t.Logf("age:%v", c.Int("age"))
    t.Logf("height:%v", c.Float("height"))
    t.Logf("score:%v", c.Float("score"))
    t.Logf("money:%+v,%v", c.Float("money"), c.Int("money"))
    t.Logf("isMe:%v", c.Bool("isMe"))
}
