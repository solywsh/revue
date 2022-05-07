package main

import (
	"fmt"
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
func groupEvent(formInfo cqPostFrom, groupId string) {
	switch {
	case formInfo.Message == "叫两声":
		_, _ = sendGroupMsg(groupId, "汪汪", "false")
	case strings.HasPrefix(formInfo.Message, "查找音乐"):
		// 查找音乐
		findMusicEvent(groupId, formInfo)
	case formInfo.Message == "开始添加":
		// 触发添加自动回复
		keywordsReplyAddEvent(groupId, formInfo, 1, 0)
	case strings.HasPrefix(formInfo.Message, "删除自动回复:"):
		// 删除自动回复
		keywordsReplyDeleteEvent(strings.TrimPrefix(formInfo.Message, "删除自动回复:"), groupId)
	default:
		// 添加自动回复(关键词/回复内容)
		if res, kr := getKeywordsReplyFlag(strconv.Itoa(formInfo.UserId)); res {
			keywordsReplyAddEvent(groupId, formInfo, kr.Flag+1, kr.ID)
		} else {
			// 自动回复
			autoGroupMsg(groupId, formInfo)
		}
	}

}

//
//  findMusicEvent
//  @Description: 查找音乐事件处理
//  @param groupId
//  @param formInfo
//
func findMusicEvent(groupId string, formInfo cqPostFrom) {
	if res, musicId := music163(strings.TrimPrefix(formInfo.Message, "查找音乐")); res {
		_, _ = sendGroupMsg(groupId, cqCodeMusic("163", musicId), "false")
	} else {
		_, _ = sendGroupMsg(groupId, "没有找到", "false")
	}
}

//
//  keywordsReplyDeleteEvent
//  @Description: 关键词删除事件处理
//  @param msg
//  @param groupId
//
func keywordsReplyDeleteEvent(msg, groupId string) {
	fmt.Println("删除:", msg)
	if res, kr := searchKeywordsReply(msg); res {
		deleteKeywordsReply(kr.ID)
		_, _ = sendGroupMsg(groupId, "删除\""+msg+"\"成功", "false")
	} else {
		_, _ = sendGroupMsg(groupId, "没有找到对应的关键词", "false")
	}
}

//
//  keywordsReplyAddEvent
//  @Description: 关键词添加事件处理
//  @param groupId
//  @param formInfo
//  @param rate
//  @param krId
//
func keywordsReplyAddEvent(groupId string, formInfo cqPostFrom, rate uint, krId uint) {
	if rate == 1 {
		if res, kr := getKeywordsReplyFlag(strconv.Itoa(formInfo.UserId)); res {
			if kr.Flag == 1 {
				_, _ = sendGroupMsg(groupId, "你有一个正在添加的任务,不能重复发送\"开始添加\",请设置触发关键词", "false")
			} else if kr.Flag == 2 {
				_, _ = sendGroupMsg(groupId, "你有一个正在添加的任务,不能重复发送\"开始添加\",请为:\""+kr.Keywords+"\"设置回复", "false")
			}
		} else {
			updateKeywordsReply(KeywordsReply{Flag: 1, Userid: strconv.Itoa(formInfo.UserId)})
			_, _ = sendGroupMsg(groupId, "开始添加,请设置关键词", "false")
		}
	} else if rate == 2 {
		// 这个关键词已经存在了,覆盖
		if res, kr := searchKeywordsReply(formInfo.Message); res {
			deleteKeywordsReply(krId) // 删除暂存的
			kr.Userid = strconv.Itoa(formInfo.UserId)
			kr.Flag = 2
			updateKeywordsReply(kr)
		} else {
			updateKeywordsReply(KeywordsReply{ID: krId, Flag: 2, Userid: strconv.Itoa(formInfo.UserId), Keywords: formInfo.Message})
		}
		_, _ = sendGroupMsg(groupId, "请为\""+formInfo.Message+"\"设置回复", "false")
	} else if rate == 3 {
		updateKeywordsReply(KeywordsReply{ID: krId, Flag: 3, Msg: formInfo.Message, Mode: 1})
		_, _ = sendGroupMsg(groupId, "添加完成", "false")
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
		_, _ = sendGroupMsg(groupId, formInfo.Message, "false")
	}
	// 私聊消息
	if formInfo.MessageType == "private" {
		// demo 复读消息
		_, _ = sendMsg(formInfo.MessageType, strconv.Itoa(formInfo.UserId), "", formInfo.Message, "false")
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
				return
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
		s += "[/getToken] 获取token\n"
		s += "[/resetToken] 重置token\n"
		s += "[/deleteToken] 删除token"
		_, _ = sendMsg("", strconv.Itoa(formInfo.UserId), "", s, "false")
	} else {
		s += "群聊菜单:\n"
		s += "[开始添加] 添加自动回复\n"
		s += "[删除自动回复:{关键词}] 删除自动回复\n"
		s += "[查找音乐{关键词}] 查找音乐(暂时只支持163)\n"
		gId, _ := getGroupId(formInfo)
		_, _ = sendMsg("", "", gId, s, "false")
	}
}

//
//  msgAddApiToken
//  @Description: 添加token
//  @param formInfo
//
func msgAddApiToken(formInfo cqPostFrom) {
	if res, token := insertRevueApiToken(strconv.Itoa(formInfo.UserId), 4); res {
		_, _ = sendMsg("", strconv.Itoa(formInfo.UserId), "", "创建成功,你的token是:"+token+"\n注意,该token只能给自己发送消息", "false")
	} else {
		_, _ = sendMsg("", strconv.Itoa(formInfo.UserId), "", "创建失败,你已经创建过了,token是:"+token, "false")
	}
}

//
//  msgResetApiToken
//  @Description: 重置token
//  @param formInfo
//
func msgResetApiToken(formInfo cqPostFrom) {
	if res, token := resetRevueApiToken(strconv.Itoa(formInfo.UserId)); res {
		_, _ = sendMsg("", strconv.Itoa(formInfo.UserId), "", "重置成功,你的token是:"+token+"\n注意,该token只能给自己发送消息", "false")
	} else {
		_, _ = sendMsg("", strconv.Itoa(formInfo.UserId), "", "重置失败,请先创建token", "false")
	}
}

func msgDeleteApiToken(formInfo cqPostFrom) {
	if res, token := deleteRevueApiToken(strconv.Itoa(formInfo.UserId)); res {
		_, _ = sendMsg("", strconv.Itoa(formInfo.UserId), "", token+"删除成功", "false")
	} else {
		_, _ = sendMsg("", strconv.Itoa(formInfo.UserId), "", "删除失败,可能数据库没有对应的信息", "false")
	}
}

//
//  autoGroupMsg
//  @Description: 根据群消息自动回复
//  @param groupId qq群号
//  @param formInfo
//
func autoGroupMsg(groupId string, formInfo cqPostFrom) {
	if res, kr := searchKeywordsReply(formInfo.Message); res {
		_, _ = sendGroupMsg(groupId, kr.Msg, "false")
	}
}
