package main

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/thedevsaddam/gojsonq"
	"strconv"
)

// 发送私聊消息
func sendPrivateMsg(userId, groupId, message, autoEscape string) (string, error) {
	client := resty.New()
	post, err := client.R().SetQueryParams(map[string]string{
		"user_id":     userId,
		"group_id":    groupId,
		"message":     message,
		"auto_escape": autoEscape,
	}).Post(yamlConfig.urlHeader + "/send_private_msg")
	if err != nil {
		return "", err
	}
	rJson := gojsonq.New().JSONString(string(post.Body()))
	if rJson.Reset().Find("retcode") != nil && rJson.Reset().Find("retcode").(float64) != 0.0 {
		return "", fmt.Errorf(string(post.Body()))
	}
	messageId := strconv.Itoa(int(rJson.Reset().Find("data.message_id").(float64)))
	return messageId, err
}

// 发送群消息
func sendGroupMsg(groupId, message, autoEscape string) (string, error) {
	client := resty.New()
	post, err := client.R().SetQueryParams(map[string]string{
		"group_id":    groupId,
		"message":     message,
		"auto_escape": autoEscape,
	}).Post(yamlConfig.urlHeader + "/send_group_msg")
	if err != nil {
		return "", err
	}
	rJson := gojsonq.New().JSONString(string(post.Body()))
	if rJson.Reset().Find("retcode") != nil && rJson.Reset().Find("retcode").(float64) != 0.0 {
		return "", fmt.Errorf(string(post.Body()))
	}
	messageId := strconv.Itoa(int(rJson.Reset().Find("data.message_id").(float64)))
	return messageId, err
}

// 发送消息
func sendMsg(messageType, userId, groupId, message, autoEscape string) error {
	client := resty.New()
	post, err := client.R().SetQueryParams(map[string]string{
		"message_type": messageType,
		"user_id":      userId,
		"group_id":     groupId,
		"message":      message,
		"auto_escape":  autoEscape,
	}).Post(yamlConfig.urlHeader + "/send_msg")
	if err != nil {
		return err
	}
	rJson := gojsonq.New().JSONString(string(post.Body()))
	if rJson.Reset().Find("retcode") != nil && rJson.Reset().Find("retcode").(float64) != 0.0 {
		return fmt.Errorf(string(post.Body()))
	}
	return nil
}

func main() {
	readConfig()
	msgId, _ := sendPrivateMsg("1228014966", "58796599", "这个山谷", "false")
	fmt.Println(msgId)
	//err := send_group_msg("58796599", "msg test", "false")
	//if err != nil {
	//	return
	//}
}
