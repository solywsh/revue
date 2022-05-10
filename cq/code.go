package cq

import (
	"github.com/go-resty/resty/v2"
	"github.com/thedevsaddam/gojsonq"
	"strconv"
)

// GetCqCodeFace qq表情
func GetCqCodeFace(faceId string) string {
	return "[CQ:face,qq=" + faceId + "]"
}

// GetCqCodeAt @某人
func GetCqCodeAt(qq, name string) string {
	// qq 为空或者为all时,@全体成员
	if name == "" {
		// 在本群时
		return "[CQ:at,qq=" + qq + "]"
	} else {
		// 不在本群时
		return "[CQ:at,qq=123,name=" + name + "]"
	}
}

// GetCqCodePoke 戳一戳
func GetCqCodePoke(qq string) string {
	return "[CQ:poke,qq=" + qq + "]"
}

// GetCqCodeMusic 分享音乐(标准)
func GetCqCodeMusic(musicType, musicId string) string {
	// musicType : qq 163 xm
	return "[CQ:music,type=" + musicType + ",id=" + musicId + "]"
}

// Music163 根据关键词获取网易云id
func Music163(keywords string) (bool, string) {
	url := "http://music.cyrilstudio.top/search"
	client := resty.New()
	get, err := client.R().SetQueryParams(map[string]string{
		"keywords": keywords,
		"limit":    "1",
	}).Get(url)
	if err != nil {
		return false, ""
	} else {
		rJson := gojsonq.New().JSONString(string(get.Body()))
		if res := rJson.Reset().Find("result.songs.[0].id"); res != nil {
			return true, strconv.Itoa(int(res.(float64)))
		} else {
			return false, ""
		}
	}
}
