package cq

import (
	"github.com/go-resty/resty/v2"
	"github.com/solywsh/qqBot-revue/db"
	"github.com/solywsh/qqBot-revue/mongo_service"
	"github.com/thedevsaddam/gojsonq"
	"strconv"
	"strings"
	"time"
)

// KeywordsReplyAddEvent 关键词添加事件处理
func (cpf *PostForm) KeywordsReplyAddEvent(rate uint, krId uint) {
	if rate == 1 {
		if res, kr := gdb.GetKeywordsReplyFlag(strconv.Itoa(cpf.UserId)); res {
			if kr.Flag == 1 {
				cpf.SendGroupMsg("你有一个正在添加的任务,不能重复发送\"开始添加\",请设置触发关键词")
			} else if kr.Flag == 2 {
				cpf.SendGroupMsg("你有一个正在添加的任务,不能重复发送\"开始添加\",请为:\"" + kr.Keywords + "\"设置回复")
			}
		} else {
			gdb.UpdateKeywordsReply(db.KeywordsReply{Flag: 1, Userid: strconv.Itoa(cpf.UserId)})
			cpf.SendGroupMsg("开始添加,请设置关键词")
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
		cpf.SendGroupMsg("请为\"" + cpf.Message + "\"设置回复")
	} else if rate == 3 {
		gdb.UpdateKeywordsReply(db.KeywordsReply{ID: krId, Flag: 3, Msg: cpf.Message, Mode: 1})
		cpf.SendGroupMsg("添加完成")
	}
}

// KeywordsReplyDeleteEvent 关键词删除事件处理
func (cpf PostForm) KeywordsReplyDeleteEvent() {
	keywords := strings.TrimPrefix(cpf.Message, "删除自动回复:")
	if res, kr := gdb.SearchKeywordsReply(keywords); res {
		gdb.DeleteKeywordsReply(kr.ID)
		cpf.SendGroupMsg("删除\"" + keywords + "\"成功")
	} else {
		cpf.SendGroupMsg("没有找到对应的关键词")
	}
}

// FindMusicEvent 查找音乐事件处理
func (cpf PostForm) FindMusicEvent() {
	if res, musicId := Music163(strings.TrimPrefix(cpf.Message, "查找音乐")); res {
		cpf.SendMsg(GetCqCodeMusic("163", musicId))
	} else {
		cpf.SendMsg("没有找到")
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
		_, err := cpf.SendGroupMsg(kr.Msg)
		if err != nil {
			return
		}
	}
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

// GetProgramAlmanac 得到今天黄历
func (cpf *PostForm) GetProgramAlmanac() {
	cpf.SendMsg(gdb.GetProgrammerAlmanac())
}

func (cpf *PostForm) GetDivination() {
	ok, tag := gdb.GetDivination(strconv.Itoa(cpf.UserId))
	if ok {
		cpf.SendMsg("婚丧嫁娶亲友疾病编程测试\n升职跳槽陨石核弹各类吉凶\n\n请心中默念所求之事,4s后发送结果...")
		time.Sleep(4 * time.Second)
		res := GetCqCodeAt(strconv.Itoa(cpf.UserId), "") + " 所求运势为:" + tag
		cpf.SendMsg(res)
	} else {
		res := GetCqCodeAt(strconv.Itoa(cpf.UserId), "") + " 今日已求过签了,再求就不灵了"
		cpf.SendMsg(res)
	}
}

// randomHImg 获取随机图片 r18 = 0 非r18,r18 = 1 为r18, r18g = 2 混合
func randomHImg(r18 int, tag string) (bool, string) {
	client := resty.New().SetTimeout(time.Second * 5) // 设置超时时间
	post, err := client.R().SetQueryParams(map[string]string{
		"r18": strconv.Itoa(r18),
		"tag": tag,
	}).Post("https://api.lolicon.app/setu/v2")
	if err != nil {
		return false, "请求失败,可能是服务器高峰期(｡ ́︿ ̀｡)"
	}
	postJson := gojsonq.New().JSONString(post.String())
	imgUrl := postJson.Reset().Find("data.[0].urls.original")
	if imgUrl != nil {
		return true, imgUrl.(string)
	}
	return false, "解析url失败,可能是服务器高峰期"
}

func (cpf *PostForm) HImgEvent(r18 int, tag string) {

	// 私有数据库启用时,并且tag为空时
	if yamlConf.Database.Mongo.HImgDB.Enable && tag == "" {
		var flag bool
		if r18 == 0 {
			flag = false
		} else {
			flag = true
		}
		if mongo, ok := mongo_service.NewMongo(); ok {
			res, err := mongo.GetLoLiCon(flag)
			if err != nil {
				cpf.SendMsg(GetCqCodeAt(strconv.Itoa(cpf.UserId), "") + "发生错误:" + err.Error())
			} else {
				cpf.SendMsg(GetCqCodeImg(res.Urls.Original))
			}
		} else {
			cpf.SendMsg("mongo调用失败")
		}
		return
	}

	cpf.SendMsg(GetCqCodeAt(strconv.Itoa(cpf.UserId), "") + " 排队搜索中...")
	if ok, res := randomHImg(r18, tag); ok {
		cpf.SendMsg(GetCqCodeImg(res))
	} else {
		// 发生错误,从其他图床拿一张非涩图
		client := resty.New()
		get, err := client.R().Get("https://api.ixiaowai.cn/mcapi/mcapi.php?return=json")
		if err != nil {
			return
		}
		getJson := gojsonq.New().JSONString(string(get.Body()))
		if getJson.Reset().Find("code").(string) == "200" {
			cpf.SendMsg(GetCqCodeAt(strconv.Itoa(cpf.UserId), "") + " " +
				res + " " + GetCqCodeImg(getJson.Reset().Find("imgurl").(string)))
		} else {
			cpf.SendMsg(GetCqCodeAt(strconv.Itoa(cpf.UserId), "") + " " + res)
		}

	}
}

// GroupEvent 群消息事件
func (cpf *PostForm) GroupEvent() {
	switch {
	// demo
	case cpf.Message == "叫两声":
		cpf.SendGroupMsg("汪汪")
	case cpf.Message == "开始添加":
		// 触发添加自动回复
		cpf.KeywordsReplyAddEvent(1, 0)
	case strings.HasPrefix(cpf.Message, "删除自动回复:"):
		// 删除自动回复
		cpf.KeywordsReplyDeleteEvent()
	case strings.HasPrefix(cpf.Message, "搜索答案"):
		// 搜索答案
		cpf.GetAnswer()
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
