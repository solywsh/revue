package cq

import (
	"github.com/solywsh/qqBot-revue/db"
	"strconv"
	"strings"
)

// 对adminUser进行复读操作,用于防止风控(在处于listen状态下的群使用)
func (cpf *PostForm) RepeatOperation() {
	for _, s := range yamlConf.AdminUser {
		if strconv.Itoa(cpf.UserId) == s {
			_, _ = cpf.SendGroupMsg("[复读机](在防止风控中)" + cpf.Message)
		}
	}
}

// JudgeListenGroup 判断该群消息是否在监听群号列表中
func (cpf *PostForm) JudgeListenGroup() bool {
	groupId := strconv.Itoa(cpf.GroupId)
	for _, s := range yamlConf.ListenGroup {
		if s == groupId {
			return true
		}
	}
	return false
}

// GroupEvent 群消息事件
func (cpf *PostForm) GroupEvent() {
	//cpf.RepeatOperation() // 对adminUSer复读防止风控
	//fmt.Println("收到群消息:", cpf.Message, cpf.UserId)
	switch {
	// demo
	case cpf.Message == "叫两声":
		_, _ = cpf.SendGroupMsg("汪汪")
	case strings.HasPrefix(cpf.Message, "查找音乐"):
		// 查找音乐
		cpf.FindMusicEvent()
	case cpf.Message == "开始添加":
		// 触发添加自动回复
		cpf.KeywordsReplyAddEvent(1, 0)
	case strings.HasPrefix(cpf.Message, "删除自动回复:"):
		// 删除自动回复
		cpf.KeywordsReplyDeleteEvent(strings.TrimPrefix(cpf.Message, "删除自动回复:"))
	default:
		// 添加自动回复(关键词/回复内容)
		if res, kr := gdb.GetKeywordsReplyFlag(strconv.Itoa(cpf.UserId)); res {
			cpf.KeywordsReplyAddEvent(kr.Flag+1, kr.ID)
		} else {
			// 自动回复
			cpf.AutoGroupMsg()
		}
	}

}

// FindMusicEvent 查找音乐事件处理
func (cpf PostForm) FindMusicEvent() {
	if res, musicId := Music163(strings.TrimPrefix(cpf.Message, "查找音乐")); res {
		_, _ = cpf.SendGroupMsg(GetCqCodeMusic("163", musicId))
	} else {
		_, _ = cpf.SendGroupMsg("没有找到")
	}
}

//  关键词删除事件处理
func (cpf PostForm) KeywordsReplyDeleteEvent(msg string) {
	if res, kr := gdb.SearchKeywordsReply(msg); res {
		gdb.DeleteKeywordsReply(kr.ID)
		_, _ = cpf.SendGroupMsg("删除\"" + msg + "\"成功")
	} else {
		_, _ = cpf.SendGroupMsg("没有找到对应的关键词")
	}
}

// KeywordsReplyAddEvent 关键词添加事件处理
func (cpf *PostForm) KeywordsReplyAddEvent(rate uint, krId uint) {
	if rate == 1 {
		if res, kr := gdb.GetKeywordsReplyFlag(strconv.Itoa(cpf.UserId)); res {
			if kr.Flag == 1 {
				_, _ = cpf.SendGroupMsg("你有一个正在添加的任务,不能重复发送\"开始添加\",请设置触发关键词")
			} else if kr.Flag == 2 {
				_, _ = cpf.SendGroupMsg("你有一个正在添加的任务,不能重复发送\"开始添加\",请为:\"" + kr.Keywords + "\"设置回复")
			}
		} else {
			gdb.UpdateKeywordsReply(db.KeywordsReply{Flag: 1, Userid: strconv.Itoa(cpf.UserId)})
			_, _ = cpf.SendGroupMsg("开始添加,请设置关键词")
		}
	} else if rate == 2 {
		if res, kr := gdb.SearchKeywordsReply(cpf.Message); res {
			// 这个关键词已经存在了,覆盖
			gdb.DeleteKeywordsReply(krId) // 删除暂存的
			kr.Userid = strconv.Itoa(cpf.UserId)
			kr.Flag = 2
			gdb.UpdateKeywordsReply(kr)
		} else {
			gdb.UpdateKeywordsReply(db.KeywordsReply{ID: krId, Flag: 2, Userid: strconv.Itoa(cpf.UserId), Keywords: cpf.Message})
		}
		_, _ = cpf.SendGroupMsg("请为\"" + cpf.Message + "\"设置回复")
	} else if rate == 3 {
		gdb.UpdateKeywordsReply(db.KeywordsReply{ID: krId, Flag: 3, Msg: cpf.Message, Mode: 1})
		_, _ = cpf.SendGroupMsg("添加完成")
	}
}

// JudgmentAdminUser 判断是否为adminUser并且是否为adminUser命令
func (cpf *PostForm) JudgmentAdminUser() bool {
	for _, s := range yamlConf.AdminUser {
		// 判断是否为管理员以及是否为管理员命令头开头
		if strconv.Itoa(cpf.UserId) == s && cpf.Message[:1] == yamlConf.AdminUserOrderHeader {
			return true
		}
	}
	return false
}

// AdminEvent AdminUser事件
func (cpf *PostForm) AdminEvent() {
	// 群消息
	if cpf.MessageType == "group" {
		// demo 复读消息
		//_, _ = cpf.SendGroupMsg(cpf.Message)
	}
	// 私聊消息
	if cpf.MessageType == "private" {
		// demo 复读消息
		//_, _ = cpf.SendMsg(cpf.MessageType, cpf.Message)
	}
}

// MsgEvent 对message事件进行相应
func (cpf *PostForm) MsgEvent() {
	// 判断是否为adminUser且为命令
	if cpf.JudgmentAdminUser() {
		cpf.AdminEvent() //执行对应admin命令事件
		return           // 如果是执行之后直接返回，不再继续响应
	}
	// 群消息进行响应
	if cpf.MessageType == "group" {
		if ok := cpf.JudgeListenGroup(); ok {
			// 如果是监听qq群列表的才做出相应
			cpf.GroupEvent()
			// 发送菜单
			if cpf.Message == "/help" {
				cpf.SendMenu()
				return
			}
		}
		// 对不是监听qq群列表的消息做出相应

		// do event
	}
	// 对私聊进行响应
	if cpf.MessageType == "private" {
		// 发送菜单
		if cpf.Message == "/help" {
			cpf.SendMenu()
		}
		// /getToken 创建对应token
		if cpf.Message == "/getToken" {
			cpf.MsgAddApiToken() // 添加对应apiToken
			return
		}
		// /resetToken 重置对应token
		if cpf.Message == "/resetToken" {
			cpf.MsgResetApiToken()
			return
		}
		// /deleteToken 删除对应token
		if cpf.Message == "/deleteToken" {
			cpf.MsgDeleteApiToken()
			return
		}
	}
}

// SendMenu 发送命令菜单
func (cpf *PostForm) SendMenu() {
	s := "revue提供以下命令:\n"
	if cpf.MessageType == "private" {
		s += "revueApi 相关(私聊执行命令):\n"
		s += "[/getToken] 获取token\n"
		s += "[/resetToken] 重置token\n"
		s += "[/deleteToken] 删除token"
		_, _ = cpf.SendMsg(cpf.MessageType, s)
	} else if cpf.MessageType == "group" {
		s += "群聊菜单:\n"
		s += "[开始添加] 添加自动回复\n"
		s += "[删除自动回复:{关键词}] 删除自动回复\n"
		s += "[查找音乐{关键词}] 查找音乐(暂时只支持163)\n"
		_, _ = cpf.SendMsg(cpf.MessageType, s)
	}
}

// MsgAddApiToken 添加token
func (cpf *PostForm) MsgAddApiToken() {
	//gdb := db.NewDB()
	if res, token := gdb.InsertRevueApiToken(strconv.Itoa(cpf.UserId), 4); res {
		_, _ = cpf.SendMsg(cpf.MessageType, "创建成功,你的token是:"+token+"\n注意,该token只能给自己发送消息")
	} else {
		_, _ = cpf.SendMsg(cpf.MessageType, "创建失败,你已经创建过了,token是:"+token)
	}
}

// MsgResetApiToken 重置token
func (cpf *PostForm) MsgResetApiToken() {
	if res, token := gdb.ResetRevueApiToken(strconv.Itoa(cpf.UserId)); res {
		_, _ = cpf.SendMsg(cpf.MessageType, "重置成功,你的token是:"+token+"\n注意,该token只能给自己发送消息")
	} else {
		_, _ = cpf.SendMsg(cpf.MessageType, "重置失败,请先创建token")
	}
}

// MsgDeleteApiToken 删除token
func (cpf *PostForm) MsgDeleteApiToken() {
	if res, token := gdb.DeleteRevueApiToken(strconv.Itoa(cpf.UserId)); res {
		_, _ = cpf.SendMsg(cpf.MessageType, token+"删除成功")
	} else {
		_, _ = cpf.SendMsg(cpf.MessageType, "删除失败,可能数据库没有对应的信息")
	}
}

//	AutoGroupMsg 根据群消息自动回复
func (cpf *PostForm) AutoGroupMsg() {
	if res, kr := gdb.SearchKeywordsReply(cpf.Message); res {
		_, _ = cpf.SendGroupMsg(kr.Msg)
	}
}
