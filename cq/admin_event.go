package cq

import (
	uuid "github.com/satori/go.uuid"
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
	// 创建我在校园token 格式为$wzxyct:alive_days:user:status:times:organization
	case strings.HasPrefix(cpf.Message, yamlConf.AdminUserOrderHeader+"wzxyct"):
		cpf.CreateWzxyToken()
	// 删除我在校园token 格式为$wzxydt:token:user:organization
	case strings.HasPrefix(cpf.Message, yamlConf.AdminUserOrderHeader+"wzxydt"):
		cpf.DeleteWzxyToken()
	// 查找token信息 格式为$wzxyft:token:user:organization
	case strings.HasPrefix(cpf.Message, yamlConf.AdminUserOrderHeader+"wzxyft"):
		cpf.FindWzxyToken()
	}
}

func (cpf *PostForm) AdminHelp() {
	msg := "欢迎使用admin命令\n"
	msg += "[" + yamlConf.AdminUserOrderHeader + "help]显示菜单\n"
	msg += "[" + yamlConf.AdminUserOrderHeader + "bash {command}]执行bash命令\n"
	msg += "[" + yamlConf.AdminUserOrderHeader + "wzxyct]创建我在校园token,输入" + yamlConf.AdminUserOrderHeader + "wzxyct -h显示更多信息\n"
	msg += "[" + yamlConf.AdminUserOrderHeader + "wzxydt]删除我在校园token,输入" + yamlConf.AdminUserOrderHeader + "wzxydt -h显示更多信息\n"
	msg += "[" + yamlConf.AdminUserOrderHeader + "wzxyft]查找我在校园token,输入" + yamlConf.AdminUserOrderHeader + "wzxyft -h显示更多信息\n"
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

// 创建我在校园token 格式为$wzxyct:alive_days:user:status:times:organization

func (cpf *PostForm) CreateWzxyToken() {
	if cpf.Message == yamlConf.AdminUserOrderHeader+"wzxyct -h" {
		msg := "创建我在校园token\n"
		msg += "格式:\n" +
			"\t完整命令:" + yamlConf.AdminUserOrderHeader + "wzxyct:alive_days:user:status:times:organization\n" +
			"\t快速创建:" + yamlConf.AdminUserOrderHeader + "wzxyct:alive_days\n"
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
	cmd := strings.Split(cpf.Message, ":")
	if len(cmd) != 6 && len(cmd) != 2 {
		flag = false
	} else if len(cmd) == 6 {
		if cmd[1] != "" {
			if d, err := strconv.Atoi(cmd[1]); err != nil {
				flag = false
			} else {
				wt.Deadline = time.Now().AddDate(0, 0, d)
			}
		} else {
			wt.Deadline = time.Now().AddDate(0, 0, 1)
		}
		if cmd[2] != "" {
			wt.CreateUser = cmd[2]
		} else {
			wt.CreateUser = strconv.Itoa(cpf.UserId)
		}
		if cmd[3] != "" {
			if s, err := strconv.Atoi(cmd[3]); err != nil {
				flag = false
			} else if s == 0 || s == 2 {
				wt.Status = s
			} else {
				flag = false
			}
		} else {
			wt.Status = 0
		}
		if cmd[4] != "" {
			if t, err := strconv.Atoi(cmd[4]); err != nil {
				flag = false
			} else {
				wt.Times = t
			}
		} else if wt.Status == 2 {
			wt.Times = 30
		} else {
			wt.Times = 1
		}
		if cmd[4] != "" {
			if t, err := strconv.Atoi(cmd[4]); err != nil {
				flag = false
			} else if t >= 0 {
				wt.Times = t
			} else {
				flag = false
			}
		}
		if cmd[5] != "" {
			wt.Organization = cmd[5]
		} else {
			wt.Organization = "default"
		}
	} else if len(cmd) == 2 {
		if cmd[1] != "" {
			if d, err := strconv.Atoi(cmd[1]); err != nil {
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
		cpf.SendMsg("格式错误,格式为$wzxyct:alive_days:user:status:times:organization")
		return
	}
	wt.Token = uuid.NewV4().String()
	count, err := gdb.InsertWzxyTokenOne(wt)
	if err != nil || count <= 0 {
		cpf.SendMsg("创建失败")
	} else {
		cpf.SendMsg("创建成功,信息为:\n" + wt.String())
	}
}

// 删除我在校园token 格式为$wzxydt:token:user:organization

func (cpf *PostForm) DeleteWzxyToken() {
	if cpf.Message == yamlConf.AdminUserOrderHeader+"wzxydt -h" {
		msg := "删除我在校园token\n"
		msg += "格式:\n" +
			"\t" + yamlConf.AdminUserOrderHeader + "wzxydt:token:user:organization\n"
		msg += "参数:\n" +
			"\ttoken:我在校园token\n" +
			"\tuser:用户名\n" +
			"\torganization:组织\n"
		msg += "注意:\n" +
			"\t该命令会对所有符合条件的token进行删除,但确保至少有一个信息有效\n"
		cpf.SendMsg(msg)
		return
	}
	var wt wzxy.TokenWzxy
	flag := true
	tag := 0
	cmd := strings.Split(cpf.Message, ":")
	if len(cmd) != 4 {
		flag = false
	} else {
		if cmd[1] != "" {
			wt.Token = cmd[1]
			tag++
		}
		if cmd[2] != "" {
			wt.CreateUser = cmd[2]
			tag++
		}
		if cmd[3] != "" {
			wt.Organization = cmd[3]
			tag++
		}
	}
	if !flag {
		cpf.SendMsg("格式错误,格式为:$wzxydt:token:user:organization")
		return
	}
	if tag > 0 {
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
	} else {
		cpf.SendMsg("格式错误,请确保至少有一个参数带有值")
		return
	}
}

func (cpf PostForm) FindWzxyToken() {
	if cpf.Message == yamlConf.AdminUserOrderHeader+"wzxyft -h" {
		msg := "查找我在校园token信息\n"
		msg += "格式:\n" +
			"\t完整命令:" + yamlConf.AdminUserOrderHeader + "wzxyft:token:user:organization\n" +
			"\t快速查找:" + yamlConf.AdminUserOrderHeader + "wzxyft:token\n"
		msg += "参数:\n" +
			"\ttoken:token\n" +
			"\tuser:用户名\n" +
			"\torganization:组织\n"
		msg += "注意:\n" +
			"\t确保至少有一个信息有效\n"
		cpf.SendMsg(msg)
		return
	}
	flag := true
	tag := 0
	var wt wzxy.TokenWzxy
	cmd := strings.Split(cpf.Message, ":")
	if len(cmd) != 4 && len(cmd) != 2 {
		flag = false
	} else if len(cmd) == 4 {
		if cmd[1] != "" {
			wt.Token = cmd[1]
			tag++
		}
		if cmd[2] != "" {
			wt.CreateUser = cmd[2]
			tag++
		}
		if cmd[3] != "" {
			wt.Organization = cmd[3]
			tag++
		}
	} else if len(cmd) == 2 {
		wt.Token = cmd[1]
		tag++
	} else {
		flag = false
	}
	if !flag {
		cpf.SendMsg("格式错误,格式为:$wzxyft:token:user:organization")
		return
	}
	if tag > 0 {
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
	} else {
		cpf.SendMsg("格式错误,请确保至少有一个参数带有值")
		return
	}
}
