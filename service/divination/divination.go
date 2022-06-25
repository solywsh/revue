// Package divination 求签
package divination

import (
	"math/rand"
	"time"
)

func NewDivination() string {
	tag := []string{
		"大吉", "吉", "小吉", "平", "小凶", "大凶",
		"吉", "小吉", "小凶", "平", "小吉", "小凶", "平", "凶",
		"小吉", "吉",
	}
	rand.Seed(time.Now().Unix())
	return tag[rand.Intn(len(tag))]
}
