package main

import (
	"github.com/go-resty/resty/v2"
	"github.com/thedevsaddam/gojsonq"
	"strconv"
)

//
//  music163
//  @Description: 根据关键词获取网易云id
//  @param keywords
//  @return bool
//  @return string
//
func music163(keywords string) (bool, string) {
	url := "http://music.cyrilstudio.top/search"
	client := resty.New()
	get, err := client.R().SetQueryParams(map[string]string{
		"keywords": keywords,
		"limit":    "1",
	}).Get(url)
	if err != nil {
		return false, ""
	} else {
		//fmt.Println(string(get.Body()))
		rJson := gojsonq.New().JSONString(string(get.Body()))
		if res := rJson.Reset().Find("result.songs.[0].id"); res != nil {
			return true, strconv.Itoa(int(res.(float64)))
		} else {
			return false, ""
		}
	}
}

//func main() {
//	music163("北宇治高校吹奏楽部")
//}
