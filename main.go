package main

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/gojsonq"
	"io"
	"strconv"
)

type cqPostFrom struct {
	Time        int    `json:"time"`
	SelfId      int    `json:"self_id"`
	PostType    string `json:"post_type"`
	MessageType string `json:"message_type"`
	SubType     string `json:"sub_type"`
	MessageId   int    `json:"message_id"`
	UserId      int    `json:"user_id"`
	Message     string `json:"message"`
	RawMessage  string `json:"raw_message"`
	Font        int    `json:"font"`
	Sender      struct {
		Nickname string `json:"nickname"`
		Sex      string `json:"sex"`
		Age      int    `json:"age"`
	} `json:"sender"`
}

func getGroupId(fromInfo cqPostFrom) (string, bool) {
	msg, err := getMsg(strconv.Itoa(fromInfo.MessageId))
	if err != nil {
		return "", false
	}
	fJson := gojsonq.New().JSONString(msg)
	groupId := strconv.Itoa(int(fJson.Reset().Find("data.group_id").(float64)))
	//fmt.Println("---------------------------------------")
	//fmt.Println(groupId)
	//fmt.Println("---------------------------------------")
	for _, s := range yamlConfig.listenGroup {
		if s == groupId {
			return groupId, true
		}
	}
	// 返回id,但是不在监听群列表中
	return groupId, false
}

func groupEvent(fromInfo cqPostFrom, groupId string) {
	// 只是demo
	if fromInfo.Message == "叫两声" {
		sendGroupMsg(groupId, "汪汪", "false")
	}
}

func event(fromInfo cqPostFrom) {
	if fromInfo.MessageType == "group" {
		groupId, ok := getGroupId(fromInfo)
		if ok {
			groupEvent(fromInfo, groupId)
		}
	}
	//fmt.Println("---------------------------------------")
	//fmt.Println(getMsg(strconv.Itoa(fromInfo.MessageId)))
	//fmt.Println("---------------------------------------")
	//if fromInfo.Message == "[CQ:at,qq="+strconv.Itoa(fromInfo.SelfId)+"] 叫两声" {
	//
	//}
}

func getSHA1(data string) string {
	t := sha1.New()
	_, err := io.WriteString(t, data)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%x", t.Sum(nil))
}

func listenFromCqhttp(c *gin.Context) {
	headerXSignature := c.Request.Header.Get("X-Signature") // sha1
	headerXSelfId := c.Request.Header.Get("X-Self-Id")      // 发送qq

	fmt.Println(headerXSignature[len("sha1="):], headerXSelfId)

	var form cqPostFrom
	if c.ShouldBind(&form) == nil {
		//fmt.Printf("%#v", form)
		info, _ := json.Marshal(form)
		fmt.Println(string(info))
		pJson := gojsonq.New().JSONString(string(info))
		if pJson.Reset().Find("post_type").(string) != "meta_event" {
			// do event
			event(form)
		}
	}
}

func main() {
	readConfig() // 读取配置
	//msg, err := sendMsg("private", "1228014966", "", "test", "false")
	//if err != nil {
	//	return
	//}
	router := gin.Default()
	router.POST("/", listenFromCqhttp)
	router.Run(":5000")
}
