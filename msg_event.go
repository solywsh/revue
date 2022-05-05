package main

import (
	"github.com/thedevsaddam/gojsonq"
	"strconv"
	"strings"
)

//
//  getGroupId
//  @Description: 从消息id中得到qq群号,用来判断是否在监听群号列表中
//  @param formInfo
//  @return string
//  @return bool
//
func getGroupId(formInfo cqPostFrom) (string, bool) {
	msg, err := getMsg(strconv.Itoa(formInfo.MessageId))
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
//  @param formInfo
//  @param groupId
//
func groupEvent(fromInfo cqPostFrom, groupId string) {
	// demo
	if fromInfo.Message == "叫两声" {
		sendGroupMsg(groupId, "汪汪", "false")
	}
	//fmt.Println(fromInfo.Message)
	if strings.HasPrefix(fromInfo.Message, "查找音乐") {
		if res, musicId := music163(strings.TrimPrefix(fromInfo.Message, "查找音乐")); res {
			sendGroupMsg(groupId, cqCodeMusic("163", musicId), "false")
		} else {
			sendGroupMsg(groupId, "没有找到", "false")
		}
	}
}

//
//  judgmentAdminUser
//  @Description: 判断是否为adminUser并且是否为adminUser命令
//  @param formInfo
//  @return bool
//
func judgmentAdminUser(formInfo cqPostFrom) bool {
	for _, s := range yamlConfig.AdminUser {
		// 判断是否为管理员以及是否为管理员命令头开头
		if strconv.Itoa(formInfo.UserId) == s && formInfo.Message[:1] == yamlConfig.AdminUOH {
			return true
		}
	}
	return false
}

//
//  adminEvent
//  @Description: adminUser 事件
//  @param formInfo
//
func adminEvent(formInfo cqPostFrom) {
	// 群消息
	if formInfo.MessageType == "group" {
		groupId, _ := getGroupId(formInfo) // 得到消息的qq群号
		// demo 复读消息
		sendGroupMsg(groupId, formInfo.Message, "false")
	}
	// 私聊消息
	if formInfo.MessageType == "private" {
		// demo 复读消息
		sendMsg(formInfo.MessageType, strconv.Itoa(formInfo.UserId), "", formInfo.Message, "false")
	}
}

//
//  msgEvent
//  @Description: 对message事件进行相应
//  @param formInfo
//
func msgEvent(formInfo cqPostFrom) {
	// 判断是否为adminUser且为命令
	if judgmentAdminUser(formInfo) {
		adminEvent(formInfo) //执行对应admin命令事件
		return               // 如果是执行之后直接返回，不再继续响应
	}
	// 群消息进行响应
	if formInfo.MessageType == "group" {
		groupId, ok := getGroupId(formInfo)
		if ok {
			// 如果是监听qq群列表的才做出相应
			groupEvent(formInfo, groupId)
			// 发送菜单
			if formInfo.Message == "/help" {
				sendMenu(formInfo, 2)
			}
		}
		// 对不是监听qq群列表的消息做出相应
		// do event
	}
	// 对私聊进行响应
	if formInfo.MessageType == "private" {
		// 发送菜单
		if formInfo.Message == "/help" {
			sendMenu(formInfo, 1)
		}
		// /getToken 创建对应token
		if formInfo.Message == "/getToken" {
			msgAddApiToken(formInfo) // 添加对应apiToken
			return
		}
		// /resetToken 重置对应token
		if formInfo.Message == "/resetToken" {
			msgResetApiToken(formInfo)
			return
		}
		// /deleteToken 删除对应token
		if formInfo.Message == "/deleteToken" {
			msgDeleteApiToken(formInfo)
			return
		}
	}
}

//
//  sendMenu
//  @Description: 发送命令菜单
//  @param formInfo
//  @param model 模式,1的时候为私聊,2的时候为群聊
//
func sendMenu(formInfo cqPostFrom, model int) {
	s := "revue提供以下命令:\n"
	if model == 1 {
		s += "revueApi 相关(私聊执行命令):\n"
		s += "/getToken 获取token\n"
		s += "/resetToken 重置token\n"
		s += "/deleteToken 删除token"
		sendMsg("", strconv.Itoa(formInfo.UserId), "", s, "false")
	} else {
		s += "该环境下暂时为空呢"
		gId, _ := getGroupId(formInfo)
		sendMsg("", "", gId, s, "false")
	}
}

//
//  msgAddApiToken
//  @Description: 添加token
//  @param formInfo
//
func msgAddApiToken(formInfo cqPostFrom) {
	if res, token := insertRevueApiToken(strconv.Itoa(formInfo.UserId), 4); res {
		sendMsg("", strconv.Itoa(formInfo.UserId), "", "创建成功,你的token是:"+token+"\n注意,该token只能给自己发送消息", "false")
	} else {
		sendMsg("", strconv.Itoa(formInfo.UserId), "", "创建失败,你已经创建过了,token是:"+token, "false")
	}
}

//
//  msgResetApiToken
//  @Description: 重置token
//  @param formInfo
//
func msgResetApiToken(formInfo cqPostFrom) {
	if res, token := resetRevueApiToken(strconv.Itoa(formInfo.UserId)); res {
		sendMsg("", strconv.Itoa(formInfo.UserId), "", "重置成功,你的token是:"+token+"\n注意,该token只能给自己发送消息", "false")
	} else {
		sendMsg("", strconv.Itoa(formInfo.UserId), "", "重置失败,请先创建token", "false")
	}
}

func msgDeleteApiToken(formInfo cqPostFrom) {
	if res, token := deleteRevueApiToken(strconv.Itoa(formInfo.UserId)); res {
		sendMsg("", strconv.Itoa(formInfo.UserId), "", token+"删除成功", "false")
	} else {
		sendMsg("", strconv.Itoa(formInfo.UserId), "", "删除失败,可能数据库没有对应的信息", "false")
	}
}
