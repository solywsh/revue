package cq

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/thedevsaddam/gojsonq"
	"strings"
)

// MsgEvent 对message事件进行相应
func (cpf *PostForm) MsgEvent() {

	// 判断是否为adminUser且为命令
	if cpf.JudgmentAdminUser() {
		cpf.AdminEvent() //执行对应admin命令事件
		return           // 如果是执行之后直接返回，不再继续响应
	}
	// 响应通用消息
	if cpf.CommonEvent() {
		return
	}
	// 群消息进行响应
	if cpf.MessageType == "group" {
		if ok := cpf.JudgeListenGroup(); ok {
			// 如果是监听qq群列表的才做出相应
			cpf.GroupEvent()
			return
		}
		// 对不是监听qq群列表的消息做出相应
		// do event
	}
	// 对私聊进行响应
	if cpf.MessageType == "private" {
		cpf.PrivateEvent()
	}
}

// SendMenu 发送命令菜单
func (cpf *PostForm) SendMenu() {
	s := "revue提供以下命令:\n"
	s += "通用菜单:\n"
	s += "\t[搜索答案{关键词}] 搜索答案\n"
	s += "\t[查找音乐{关键词}] 查找音乐(暂时只支持163)\n"
	s += "\t[搜索答案{关键词}] 搜索答案\n"
	s += "\t[程序员黄历] 显示今天黄历\n"
	s += "\t[求签] 今日运势\n"
	s += "\t[无内鬼来点{关键词}] 发送二刺螈图片\n"
	if cpf.MessageType == "private" {
		s += "私聊菜单:\n"
		s += "revueApi 相关(私聊执行命令):\n"
		s += "\t[/getToken] 获取token\n"
		s += "\t[/resetToken] 重置token\n"
		s += "\t[/deleteToken] 删除token\n"
		cpf.SendMsg(s)
	} else if cpf.MessageType == "group" {
		s += "群聊菜单:\n"
		s += "\t[开始添加] 添加自动回复\n"
		s += "\t[删除自动回复:{关键词}] 删除自动回复\n"
		cpf.SendMsg(s)
	}
}

// ProblemRepository 搜索题库
func ProblemRepository(question string) string {
	client := resty.New()
	post, err := client.R().SetQueryParams(map[string]string{
		"q": question,
	}).Post("http://api.902000.xyz:88/wkapi.php")
	if err != nil {
		return "题目请求失败"
	}
	postJSON := gojsonq.New().JSONString(post.String())
	if postJSON.Reset().Find("code").(float64) == float64(1) {
		return fmt.Sprintf("问题:" + postJSON.Reset().Find("tm").(string) + "\n" +
			"答案:" + postJSON.Reset().Find("answer").(string))
	} else {
		return "没有找到相关问题"
	}
}

// GetAnswer 搜索题目答案
func (cpf *PostForm) GetAnswer() {
	question := strings.TrimPrefix(cpf.Message, "搜索答案")
	ans := ProblemRepository(question)
	cpf.SendMsg(ans)
}

// CommonEvent 对通用消息响应
func (cpf *PostForm) CommonEvent() bool {

	// cpf.RepeatOperation() // 对adminUSer复读防止风控
	//fmt.Println("收到群消息:", cpf.Message, cpf.UserId)
	switch {
	case cpf.Message == "/help":
		cpf.SendMenu()
		return true
	case strings.HasPrefix(cpf.Message, "查找音乐"):
		// 查找音乐
		cpf.FindMusicEvent()
		return true
	case cpf.Message == "程序员黄历":
		// 发送程序员黄历
		cpf.GetProgramAlmanac()
		return true
	case cpf.Message == "求签":
		// 求签
		cpf.GetDivination()
		return true
	case cpf.Message == "无内鬼来点涩图":
		// 涩图事件,非r18
		cpf.HImgEvent(0, "")
		return true
	case cpf.Message == "无内鬼来点色图":
		// 色图事件,r18
		cpf.HImgEvent(1, "")
		return true
	case strings.HasPrefix(cpf.Message, "无内鬼来点"):
		// 涩图事件,搜索标签tag
		cpf.HImgEvent(2, strings.TrimPrefix(cpf.Message, "无内鬼来点"))
		return true
	case strings.HasPrefix(cpf.Message, "搜索答案"):
		// 搜索答案
		cpf.GetAnswer()
		return true
	}
	return false
}
