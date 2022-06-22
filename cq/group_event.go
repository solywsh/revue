package cq

import (
	"github.com/solywsh/qqBot-revue/db"
	"strconv"
	"strings"
)

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

// KeywordsReplyDeleteEvent 关键词删除事件处理
func (cpf PostForm) KeywordsReplyDeleteEvent() {
	keywords := strings.TrimPrefix(cpf.Message, "删除自动回复:")
	if res, kr := gdb.SearchKeywordsReply(keywords); res {
		gdb.DeleteKeywordsReply(kr.ID)
		_, _ = cpf.SendGroupMsg("删除\"" + keywords + "\"成功")
	} else {
		_, _ = cpf.SendGroupMsg("没有找到对应的关键词")
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

//	AutoGroupMsg 根据群消息自动回复
func (cpf *PostForm) AutoGroupMsg() {
	if res, kr := gdb.SearchKeywordsReply(cpf.Message); res {
		_, _ = cpf.SendGroupMsg(kr.Msg)
	}
}
