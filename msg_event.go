package main

import (
	"github.com/thedevsaddam/gojsonq"
	"strconv"
)

//
//  getGroupId
//  @Description: 从消息id中得到qq群号,用来判断是否在监听群号列表中
//  @param fromInfo
//  @return string
//  @return bool
//
func getGroupId(fromInfo cqPostFrom) (string, bool) {
	msg, err := getMsg(strconv.Itoa(fromInfo.MessageId))
	if err != nil {
		return "", false
	}
	fJson := gojsonq.New().JSONString(msg)
	groupId := strconv.Itoa(int(fJson.Reset().Find("data.group_id").(float64)))
	for _, s := range yamlConfig.ListenGroup {
		if s == groupId {
			return groupId, true
		}
	}
	// 返回id,但是不在监听群列表中
	return groupId, false
}

//
//  groupEvent
//  @Description: 群消息事件
//  @param fromInfo
//  @param groupId
//
func groupEvent(fromInfo cqPostFrom, groupId string) {
	// demo
	if fromInfo.Message == "叫两声" {
		sendGroupMsg(groupId, "汪汪", "false")
	}
}

//
//  judgmentAdminUser
//  @Description: 判断是否为adminUser并且是否为adminUser命令
//  @param fromInfo
//  @return bool
//
func judgmentAdminUser(fromInfo cqPostFrom) bool {
	for _, s := range yamlConfig.AdminUser {
		// 判断是否为管理员以及是否为管理员命令头开头
		if strconv.Itoa(fromInfo.UserId) == s && fromInfo.Message[:1] == yamlConfig.AdminUOH {
			return true
		}
	}
	return false
}

//
//  adminEvent
//  @Description: adminUser 事件
//  @param fromInfo
//
func adminEvent(fromInfo cqPostFrom) {
	// 群消息
	if fromInfo.MessageType == "group" {
		groupId, _ := getGroupId(fromInfo) // 得到消息的qq群号
		// demo 复读消息
		sendGroupMsg(groupId, fromInfo.Message, "false")
	}
	// 私聊消息
	if fromInfo.MessageType == "private" {
		// demo 复读消息
		sendMsg(fromInfo.MessageType, strconv.Itoa(fromInfo.UserId), "", fromInfo.Message, "false")
	}
}

//
//  msgEvent
//  @Description: 对message事件进行相应
//  @param fromInfo
//
func msgEvent(fromInfo cqPostFrom) {
	// 判断是否为adminUser
	if judgmentAdminUser(fromInfo) {
		adminEvent(fromInfo) //执行对应admin命令事件
		return               // 如果是执行之后直接返回，不再继续响应
	}
	// 群消息进行响应
	if fromInfo.MessageType == "group" {
		groupId, ok := getGroupId(fromInfo)
		if ok {
			// 如果是监听qq群列表的才做出相应
			groupEvent(fromInfo, groupId)
		}
		// 对不是监听qq群列表的消息做出相应
		// do event
	}
	// 对私聊进行响应
	if fromInfo.MessageType == "private" {
		// do event
	}
}
