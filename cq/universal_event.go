package cq

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/solywsh/qqBot-revue/service/wzxy"
	"github.com/thedevsaddam/gojsonq"
	"strconv"
	"strings"
	"time"
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
	s += "\t[wzxy -h] 我在校园打卡相关"
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
		cpf.FindMusicEvent() // 查找音乐
		return true
	case cpf.Message == "程序员黄历":
		cpf.GetProgramAlmanac() // 发送程序员黄历
		return true
	case cpf.Message == "求签":
		cpf.GetDivination() // 求签
		return true
	case cpf.Message == "无内鬼来点涩图":
		cpf.HImgEvent(0, "") // 涩图事件,非r18
		return true
	case cpf.Message == "无内鬼来点色图":
		cpf.HImgEvent(1, "") // 色图事件,r18
		return true
	case strings.HasPrefix(cpf.Message, "无内鬼来点"):
		cpf.HImgEvent(2, strings.TrimPrefix(cpf.Message, "无内鬼来点")) // 涩图事件,搜索标签tag
		return true
	case strings.HasPrefix(cpf.Message, "搜索答案"):
		cpf.GetAnswer() // 搜索答案
		return true
	case strings.HasPrefix(cpf.Message, "wzxy"):
		cpf.HandleUserWzxy()
	}
	return false
}

func (cpf PostForm) HandleUserWzxy() {
	if cpf.Message == "wzxy -h" {
		msg := "我在校园打卡\n"
		msg += "使用方法:\n"
		msg += "\twzxy -a [token] [jwsession]\t注册我在校园打卡任务\n"
		msg += "\twzxy -d\t删除我在校园打卡任务\n"
		msg += "\twzxy -m morning/afternoon/check 15:04\t修改[晨检/午检/签到]打卡时间,时间格式为15:04\n"
		msg += "\twzxy -r [jwsession]\t修改我在校园jwsession\n"
		msg += "\twzxy -l\t查看我在校园打卡任务\n"
		msg += "\twzxy -alive 2006-01-02 15:04\t修改打卡任务结束时间,时间格式为2006-01-02 15:04,注意不能超过token的时间权限\n"
		msg += "\twzxy -on/off morning/afternoon/check/all\t开启/关闭打卡任务\n"
		msg += "\twzxy -do morning/afternoon/check\t手动打卡\n"
		msg += "\twzxy -h\t查看帮助\n"
		cpf.SendMsg(msg)
		return
	}
	cmd := strings.Split(cpf.Message, " ")
	if len(cmd) < 2 {
		cpf.SendMsg("参数错误")
		return
	}
	// 提前查找用户任务
	userWzxy := wzxy.UserWzxy{}
	flag := false
	if strings.HasPrefix(cpf.Message, "wzxy -") {
		manyUser, i, err := gdb.FindWzxyUserMany(wzxy.UserWzxy{UserId: strconv.Itoa(cpf.UserId)})
		if err != nil || i != 1 {
			flag = false
		} else if len(manyUser) == 1 {
			userWzxy = manyUser[0]
			flag = true
		}
	}

	// 注册任务
	if strings.HasPrefix(cpf.Message, "wzxy -a") && len(cmd) > 3 {
		if flag {
			cpf.SendMsg("您已经注册过任务了")
			return
		}
		token := cmd[2]
		wzxyTokens, i, err := gdb.FindWzxyTokenMany(wzxy.TokenWzxy{Token: token})
		wt := wzxyTokens[0]
		if err != nil || i != 1 || wt.Times <= 0 || wt.Status == 1 {
			cpf.SendMsg("注册失败,请输入一个有效的token")
			return
		}
		if len(cmd) == 4 {
			_, err := gdb.InsertWzxyUserOne(wzxy.UserWzxy{
				Jwsession:              cmd[3],
				JwsessionStatus:        true,
				Token:                  cmd[2],
				UserId:                 strconv.Itoa(cpf.UserId),
				Name:                   cpf.Sender.Nickname,
				MorningCheckEnable:     true,
				MorningCheckTime:       "08:00",
				MorningLastCheckDate:   "2006-01-02",
				AfternoonCheckEnable:   true,
				AfternoonCheckTime:     "13:00",
				AfternoonLastCheckDate: "2006-01-02",
				EveningCheckEnable:     true,
				EveningCheckTime:       "21:30",
				EveningLastCheckDate:   "2006-01-02",
				Province:               "陕西省",
				City:                   "西安市",
				UserAgent:              "Mozilla/5.0 (iPhone; CPU iPhone OS 14_2_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 MicroMessenger/8.0.18(0x18001236) NetType/WIFI Language/zh_CN",
				Deadline:               wt.Deadline,
			})
			if err != nil {
				cpf.SendMsg("注册失败")
				return
			}
			cpf.SendMsg("注册成功,输入wzxy -l查看具体打卡任务信息")
			wt.Times--
			if wt.Status == 0 {
				wt.Status = 1
			}
			_, err = gdb.UpdateWzxyTokenOne(wt, true)
		} else {
			cpf.SendMsg("注册失败,请按照wzxy -a [token] [jwsession]格式重新输入")
		}
		return
	}

	if !flag {
		cpf.SendMsg("您还没有注册任务")
		return
	}

	// 打印任务
	if cpf.Message == "wzxy -l" {
		cpf.SendMsg(userWzxy.String())
		return
	}

	// 删除任务
	if cpf.Message == "wzxy -d" {
		i, err := gdb.DeleteWzxyUserOne(userWzxy)
		if err != nil || i != 1 {
			cpf.SendMsg("删除失败")
			return
		}
		cpf.SendMsg("删除成功")
		return
	}

	// 修改jwsession
	if strings.HasPrefix(cpf.Message, "wzxy -r") {
		if len(cmd) == 3 {
			userWzxy.Jwsession = cmd[2]
			userWzxy.JwsessionStatus = true
			one, err := gdb.UpdateWzxyUserOne(userWzxy, flag)
			if err != nil || one != 1 {
				cpf.SendMsg("修改失败")
				return
			}
			cpf.SendMsg("修改成功")
		} else {
			cpf.SendMsg("修改失败,请按照wzxy -r [jwsession]格式重新输入")
		}
		return
	}
	// 修改打卡时间
	if strings.HasPrefix(cpf.Message, "wzxy -m") {
		if len(cmd) == 4 && (cmd[2] == "morning" || cmd[2] == "afternoon" || cmd[2] == "check") {
			// todo check cmd[3] is time
			switch cmd[2] {
			case "morning":
				userWzxy.MorningCheckTime = cmd[3]
			case "afternoon":
				userWzxy.AfternoonCheckTime = cmd[3]
			case "check":
				userWzxy.EveningCheckTime = cmd[3]
			default:
				cpf.SendMsg("修改失败,请按照wzxy -m morning/afternoon/check 15:04格式输入")
				return
			}
			one, err := gdb.UpdateWzxyUserOne(userWzxy, false)
			if err != nil || one != 1 {
				cpf.SendMsg("修改失败")
				return
			}
			cpf.SendMsg("修改成功")
		} else {
			cpf.SendMsg("修改失败,请按照wzxy -m morning/afternoon/check 15:04格式输入")
		}
		return
	}
	// 修改打卡任务存活时间
	if strings.HasPrefix(cpf.Message, "wzxy -alive") {
		if len(cmd) == 4 {
			t, err := time.Parse("2006-01-02 15:04:05", cmd[2]+" "+cmd[3])
			if err != nil {
				cpf.SendMsg("修改失败,请按照wzxy -alive 2006-01-02 15:04格式输入")
				return
			}
			tokens, i, err := gdb.FindWzxyTokenMany(wzxy.TokenWzxy{Token: userWzxy.Token})
			if err != nil || i != 1 {
				cpf.SendMsg("修改失败,没有找到相关token信息")
				return
			}
			if tokens[0].Deadline.Before(t) {
				cpf.SendMsg("修改失败,超过了token的时间权限")
				return
			} else {
				userWzxy.Deadline = t
				one, err := gdb.UpdateWzxyUserOne(userWzxy, false)
				if err != nil || one != 1 {
					cpf.SendMsg("修改失败")
					return
				}
				cpf.SendMsg("修改成功")
			}
		} else {
			cpf.SendMsg("修改失败,请按照wzxy -alive 2006-01-02 15:04格式输入")
		}
		return
	}
	// 开启/关闭打卡任务
	if strings.HasPrefix(cpf.Message, "wzxy -on") ||
		strings.HasPrefix(cpf.Message, "wzxy -off ") {
		if len(cmd) == 3 {
			var taskStatus bool
			if cmd[1] == "-on" {
				taskStatus = true
			} else {
				taskStatus = false
			}
			switch cmd[2] {
			case "morning":
				userWzxy.MorningCheckEnable = taskStatus
			case "afternoon":
				userWzxy.AfternoonCheckEnable = taskStatus
			case "check":
				userWzxy.EveningCheckEnable = taskStatus
			case "all":
				userWzxy.MorningCheckEnable = taskStatus
				userWzxy.AfternoonCheckEnable = taskStatus
				userWzxy.EveningCheckEnable = taskStatus
			default:
				cpf.SendMsg("修改失败,请按照wzxy -on morning/afternoon/check格式输入")
				return
			}
			one, err := gdb.UpdateWzxyUserOne(userWzxy, true)
			if err != nil || one != 1 {
				cpf.SendMsg("修改失败")
				return
			}
			cpf.SendMsg("修改成功")
		} else {
			cpf.SendMsg("修改失败,请按照wzxy -on morning/afternoon/check格式输入")
		}
		return
	}
	// 手动执行打卡
	if strings.HasPrefix(cpf.Message, "wzxy -do") {

		if len(cmd) == 3 {
			var status int
			var msg string
			switch cmd[2] {
			case "morning":
				status, msg = userWzxy.CheckOperate(1)
			case "afternoon":
				status, msg = userWzxy.CheckOperate(2)
			case "check":
				status = userWzxy.EveningCheckOperate()
			default:
				cpf.SendMsg("执行失败,请按照wzxy -on morning/afternoon/check格式输入")
				return
			}
			if status == 0 {
				cpf.SendMsg("执行成功")
			} else if msg != "" {
				cpf.SendMsg("执行失败," + msg)
			} else {
				cpf.SendMsg("执行失败")
			}
		} else {
			cpf.SendMsg("执行失败,请按照wzxy -do morning/afternoon/check格式输入")
		}
		return
	}
}
