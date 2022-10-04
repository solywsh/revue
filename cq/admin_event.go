package cq

import (
	uuid "github.com/satori/go.uuid"
	"github.com/solywsh/qqBot-revue/db"
	"github.com/solywsh/qqBot-revue/service/wzxy"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// RepeatOperation 对adminUser进行复读操作,用于防止风控(在处于listen状态下的群使用)
func (cpf *PostForm) RepeatOperation() {
	for _, s := range yamlConf.AdminUser {
		if strconv.Itoa(cpf.UserId) == s {
			cpf.SendGroupMsg("[复读机](在防止风控中)" + cpf.Message)
		}
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
	// 只对群消息响应
	if cpf.MessageType == "group" {
		//// demo 复读消息
		//if cpf.Message == "$any"{
		//	_, _ = cpf.SendGroupMsg(cpf.Message)
		//	return
		//}
	}
	// 只对私聊消息响应
	if cpf.MessageType == "private" {
		//switch {
		//
		//}
	}
	//  通用
	switch {
	// case "$help" 显示菜单
	case strings.HasPrefix(cpf.Message, yamlConf.AdminUserOrderHeader+"help"):
		cpf.AdminHelp()
	// case "$bash" 对系统执行命令
	case strings.HasPrefix(cpf.Message, yamlConf.AdminUserOrderHeader+"bash"):
		cpf.BashCommand()
	// case "$wzxy" 我在校园token相关
	case strings.HasPrefix(cpf.Message, yamlConf.AdminUserOrderHeader+"wzxy"):
		cpf.HandleAdminWzxy()
	// case "$listen" 监听群消息
	case strings.HasPrefix(cpf.Message, yamlConf.AdminUserOrderHeader+"lg"):
		cpf.HandleAdminListenGroup()
	}
}

func (cpf *PostForm) AdminHelp() {
	msg := "欢迎使用admin命令\n"
	msg += "[" + yamlConf.AdminUserOrderHeader + "help]显示菜单\n"
	msg += "[" + yamlConf.AdminUserOrderHeader + "bash {command}]执行Linux bash命令\n"
	msg += "[" + yamlConf.AdminUserOrderHeader + "wzxy]我在校园token相关,输入" + yamlConf.AdminUserOrderHeader + "wzxy -h显示更多信息\n"
	msg += "[" + yamlConf.AdminUserOrderHeader + "lg]监听群消息相关,输入" + yamlConf.AdminUserOrderHeader + "lg -h显示更多信息\n"
	cpf.SendMsg(msg)
}

// BashCommand 对admin执行bash命令
func (cpf *PostForm) BashCommand() {
	cmd := strings.TrimPrefix(cpf.Message, yamlConf.AdminUserOrderHeader+"bash ")
	res := func(cmd string) string {
		if runtime.GOOS == "windows" {
			return "不能在Windows下执行"
		}
		out, err := exec.Command("bash", "-c", cmd).Output()
		if err != nil {
			return "执行错误"
		}
		return string(out[:])
	}(cmd)
	cpf.SendMsg(res)
}

func (cpf *PostForm) HandleAdminWzxy() {
	if cpf.Message == yamlConf.AdminUserOrderHeader+"wzxy -h" {
		msg := "我在校园Token管理\n"
		msg += "使用方法:\n"
		msg += "\t" + yamlConf.AdminUserOrderHeader + "wzxy -c\t注册我在校园Token,输入" + yamlConf.AdminUserOrderHeader + "wzxy -c -h显示更多信息\n"
		msg += "\t" + yamlConf.AdminUserOrderHeader + "wzxy -d\t删除我在校园Token,输入" + yamlConf.AdminUserOrderHeader + "wzxy -d -h显示更多信息\n"
		msg += "\t" + yamlConf.AdminUserOrderHeader + "wzxy -f\t查找我在校园Token,输入" + yamlConf.AdminUserOrderHeader + "wzxy -f -h显示更多信息\n"
		msg += "\t" + yamlConf.AdminUserOrderHeader + "wzxy -h\t查看帮助\n"
		cpf.SendMsg(msg)
		return
	}
	switch {
	case strings.HasPrefix(cpf.Message, yamlConf.AdminUserOrderHeader+"wzxy -c"):
		cpf.createWzxyToken()
	case strings.HasPrefix(cpf.Message, yamlConf.AdminUserOrderHeader+"wzxy -d"):
		cpf.deleteWzxyToken()
	case strings.HasPrefix(cpf.Message, yamlConf.AdminUserOrderHeader+"wzxy -f"):
		cpf.findWzxyToken()
	}

}

// createWzxyToken 创建我在校园token 格式为$wzxy -c <alive_days> <user> <status> <times> <organization>
func (cpf *PostForm) createWzxyToken() {
	cmd := strings.Split(cpf.Message, " ")
	if cpf.Message == yamlConf.AdminUserOrderHeader+"wzxy -c -h" {
		msg := "创建我在校园token\n"
		msg += "格式:\n" +
			"\t完整命令:" + yamlConf.AdminUserOrderHeader + "wzxy -c <alive_days> <user> <status> <times> <organization>\n" +
			"\t快速创建:" + yamlConf.AdminUserOrderHeader + "wzxy -c <alive_days>\n"
		msg += "参数:\n" +
			"\talive_days:有效天数,默认1天\n" +
			"\tuser:用户名,默认管理员qq号\n" +
			"\tstatus:状态,0,1区分单次使用状态,2为多次使用,默认0(单词未使用)\n" +
			"\ttimes:次数,默认1\n" +
			"\torganization:组织,默认default\n"
		cpf.SendMsg(msg)
		return
	}
	flag := true
	var wt wzxy.TokenWzxy
	if len(cmd) != 7 && len(cmd) != 3 {
		flag = false
	} else if len(cmd) == 7 {
		if cmd[2] != "" {
			if d, err := strconv.Atoi(cmd[2]); err != nil {
				flag = false
			} else {
				wt.Deadline = time.Now().AddDate(0, 0, d)
			}
		} else {
			wt.Deadline = time.Now().AddDate(0, 0, 1)
		}
		if cmd[3] != "" {
			wt.CreateUser = cmd[3]
		} else {
			wt.CreateUser = strconv.Itoa(cpf.UserId)
		}
		if cmd[4] != "" {
			if s, err := strconv.Atoi(cmd[4]); err != nil {
				flag = false
			} else if s == 0 || s == 2 {
				wt.Status = s
			} else {
				flag = false
			}
		} else {
			wt.Status = 0
		}
		if cmd[5] != "" {
			if t, err := strconv.Atoi(cmd[5]); err != nil {
				flag = false
			} else {
				wt.Times = t
			}
		} else if wt.Status == 2 {
			wt.Times = 30
		} else {
			wt.Times = 1
		}
		if cmd[5] != "" {
			if t, err := strconv.Atoi(cmd[5]); err != nil {
				flag = false
			} else if t >= 0 {
				wt.Times = t
			} else {
				flag = false
			}
		}
		if cmd[6] != "" {
			wt.Organization = cmd[6]
		} else {
			wt.Organization = "default"
		}
	} else if len(cmd) == 3 {
		if cmd[2] != "" {
			if d, err := strconv.Atoi(cmd[2]); err != nil {
				flag = false
			} else {
				wt.Deadline = time.Now().AddDate(0, 0, d)
			}
		} else {
			wt.Deadline = time.Now().AddDate(0, 0, 1)
		}
		wt.CreateUser = strconv.Itoa(cpf.UserId)
		wt.Status = 0
		wt.Times = 1
		wt.Organization = "default"
	} else {
		flag = false
	}
	if !flag {
		cpf.SendMsg("格式错误,输入" + yamlConf.AdminUserOrderHeader + "wzxy -c -h查看帮助")
		return
	}
	wt.Token = uuid.NewV4().String()
	count, err := gdb.InsertWzxyTokenOne(wt)
	if err != nil || count <= 0 {
		cpf.SendMsg("创建失败")
	} else {
		cpf.SendMsg("创建成功,信息为:\n" + wt.String())
	}
	return
}

// deleteWzxyToken 删除我在校园token
func (cpf *PostForm) deleteWzxyToken() {
	if cpf.Message == yamlConf.AdminUserOrderHeader+"wzxy -d -h" {
		msg := "删除我在校园token\n"
		msg += "格式:\n" +
			"\t" + yamlConf.AdminUserOrderHeader + "wzxy -d -t <token>\t删除对应token\n" +
			"\t" + yamlConf.AdminUserOrderHeader + "wzxy -d -u <user>\t删除对应用户的所有token\n" +
			"\t" + yamlConf.AdminUserOrderHeader + "wzxy -d -u <organization>\t删除对应组织的所有token\n"
		cpf.SendMsg(msg)
		return
	}
	var wt wzxy.TokenWzxy
	flag := true
	cmd := strings.Split(cpf.Message, " ")
	if len(cmd) != 4 {
		flag = false
	} else {
		switch cmd[2] {
		case "-t":
			wt.Token = cmd[3]
		case "-u":
			wt.CreateUser = cmd[3]
		case "-o":
			wt.Organization = cmd[3]
		default:
			flag = false
		}
	}
	if !flag {
		cpf.SendMsg("格式错误,输入" + yamlConf.AdminUserOrderHeader + "wzxy -d -h查看帮助")
		return
	}
	many, i, err := gdb.FindWzxyTokenMany(wt)
	if err != nil {
		cpf.SendMsg("删除失败")
		return
	}
	for _, tokenWzxy := range many {
		one, err := gdb.DeleteWzxyTokenOne(tokenWzxy)
		if err != nil || one <= 0 {
			cpf.SendMsg(tokenWzxy.Token + "删除失败")
		}
	}
	cpf.SendMsg("影响了" + strconv.Itoa(int(i)) + "条")
	return
}

// findWzxyToken 查询我在校园token
func (cpf PostForm) findWzxyToken() {
	if cpf.Message == yamlConf.AdminUserOrderHeader+"wzxy -f -h" {
		msg := "查找我在校园token信息\n"
		msg += "格式:\n" +
			"\t" + yamlConf.AdminUserOrderHeader + "wzxy -f -t <token>\t查找token对应信息\n" +
			"\t" + yamlConf.AdminUserOrderHeader + "wzxy -f -u <user>\t查找user对应token信息\n" +
			"\t" + yamlConf.AdminUserOrderHeader + "wzxy -f -o <organization>\t查找organization对应token信息\n"
		cpf.SendMsg(msg)
		return
	}
	flag := true
	var wt wzxy.TokenWzxy
	cmd := strings.Split(cpf.Message, " ")
	if len(cmd) != 4 {
		flag = false
	} else {
		switch cmd[2] {
		case "-t":
			wt.Token = cmd[3]
		case "-u":
			wt.CreateUser = cmd[3]
		case "-o":
			wt.Organization = cmd[3]
		default:
			flag = false
		}
	}
	if !flag {
		cpf.SendMsg("格式错误,输入" + yamlConf.AdminUserOrderHeader + "wzxy -f -h查看帮助")
		return
	}
	many, i, err := gdb.FindWzxyTokenMany(wt)
	if i == 0 || err != nil {
		cpf.SendMsg("没有找到相关结果")
		return
	}
	msg := "查找成功,找到了" + strconv.Itoa(int(i)) + "条\n"
	for _, tokenWzxy := range many {
		msg += "=========================\n"
		msg += tokenWzxy.String()
		msg += "=========================\n"
	}
	cpf.SendMsg(msg)
}

func (cpf *PostForm) HandleAdminListenGroup() {
	if cpf.Message == yamlConf.AdminUserOrderHeader+"lg -h" {
		msg := "监听群消息\n"
		msg += "使用方法:\n"
		msg += "\t" + yamlConf.AdminUserOrderHeader + "lg -c\t注册监听群消息,输入" + yamlConf.AdminUserOrderHeader + "lg -c -h显示更多信息\n"
		msg += "\t" + yamlConf.AdminUserOrderHeader + "lg -d\t删除监听群消息,输入" + yamlConf.AdminUserOrderHeader + "lg -d -h显示更多信息\n"
		msg += "\t" + yamlConf.AdminUserOrderHeader + "lg -f\t查找监听群消息,输入" + yamlConf.AdminUserOrderHeader + "lg -f -h显示更多信息\n"
		msg += "\t" + yamlConf.AdminUserOrderHeader + "lg -h\t查看帮助\n"
		cpf.SendMsg(msg)
		return
	}
	switch {
	case strings.HasPrefix(cpf.Message, yamlConf.AdminUserOrderHeader+"lg -c"):
		cpf.createListenGroup()
	case strings.HasPrefix(cpf.Message, yamlConf.AdminUserOrderHeader+"lg -d"):
		cpf.deleteListenGroup()
	case strings.HasPrefix(cpf.Message, yamlConf.AdminUserOrderHeader+"lg -f"):
		cpf.findListenGroup()
	}
}

func (cpf *PostForm) createListenGroup() {
	if cpf.Message == yamlConf.AdminUserOrderHeader+"lg -c -h" {
		msg := "添加监听群组\n"
		msg += "格式:\n" +
			"\t" + yamlConf.AdminUserOrderHeader + "lg -c <qq群号>\n"
		cpf.SendMsg(msg)
		return
	}
	cmd := strings.Split(cpf.Message, " ")
	flag := true
	var lg db.ListenGroup
	if len(cmd) != 3 {
		flag = false
	} else {
		lg.Group = cmd[2]
		lg.Date = time.Now().Format("2006-01-02 15:04:05")
		lg.UserId = strconv.Itoa(cpf.UserId)
	}
	if !flag {
		cpf.SendMsg("格式错误,输入" + yamlConf.AdminUserOrderHeader + "lg -c -h查看帮助")
		return
	}

	count, err := gdb.InsertListenGroupOne(lg)
	if err != nil || count <= 0 {
		cpf.SendMsg("创建失败")
	} else {
		cpf.SendMsg("创建成功,信息为:\n" + lg.String())
	}
	return
}

func (cpf *PostForm) deleteListenGroup() {
	if cpf.Message == yamlConf.AdminUserOrderHeader+"lg -c -h" {
		msg := "删除监听群组\n"
		msg += "格式:\n" +
			"\t" + yamlConf.AdminUserOrderHeader + "lg -d <qq群号>\n"
		cpf.SendMsg(msg)
		return
	}
	cmd := strings.Split(cpf.Message, " ")
	flag := true
	var lg db.ListenGroup
	if len(cmd) != 3 {
		flag = false
	} else {
		lg.Group = cmd[2]
	}
	if !flag {
		cpf.SendMsg("格式错误,输入" + yamlConf.AdminUserOrderHeader + "lg -d -h查看帮助")
		return
	}

	many, i, err := gdb.FindListenGroupMany(lg)
	if err != nil {
		cpf.SendMsg("删除失败")
		return
	}
	miss := 0
	for _, listenGroup := range many {
		one, err := gdb.DeleteListenGroupOne(listenGroup)
		if err != nil || one <= 0 {
			cpf.SendMsg(listenGroup.Group + "删除失败")
			miss++
		}
	}
	cpf.SendMsg("影响了" + strconv.Itoa(int(i)-miss) + "条")
	return
}

func (cpf *PostForm) findListenGroup() {
	if cpf.Message == yamlConf.AdminUserOrderHeader+"lg -f -h" {
		msg := "查找监听群组\n"
		msg += "格式:\n" +
			"\t" + yamlConf.AdminUserOrderHeader + "lg -f -g <qq群号>\t查找监听群组对应信息\n" +
			"\t" + yamlConf.AdminUserOrderHeader + "lg -f -u <用户>\t查找该用户创建的监听群组对应信息\n"
		cpf.SendMsg(msg)
		return
	}
	cmd := strings.Split(cpf.Message, " ")
	flag := true
	var lg db.ListenGroup
	if len(cmd) != 4 {
		flag = false
	} else if cmd[2] == "-g" {
		lg.Group = cmd[3]
	} else if cmd[2] == "-u" {
		lg.UserId = cmd[3]
	} else {
		flag = false
	}
	if !flag {
		cpf.SendMsg("格式错误,输入" + yamlConf.AdminUserOrderHeader + "lg -f -h查看帮助")
		return
	}

	many, i, err := gdb.FindListenGroupMany(lg)
	if err != nil {
		cpf.SendMsg("查找失败")
		return
	}
	msg := "查找成功,找到了" + strconv.Itoa(int(i)) + "条\n"
	for _, listenGroup := range many {
		msg += "=========================\n"
		msg += listenGroup.String()
		msg += "=========================\n"
	}
	cpf.SendMsg(msg)
	return
}
