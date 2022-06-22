package cq

import (
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

// RepeatOperation 对adminUser进行复读操作,用于防止风控(在处于listen状态下的群使用)
func (cpf *PostForm) RepeatOperation() {
	for _, s := range yamlConf.AdminUser {
		if strconv.Itoa(cpf.UserId) == s {
			_, _ = cpf.SendGroupMsg("[复读机](在防止风控中)" + cpf.Message)
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
		// pass
	}
	//  通用
	switch {
	// case "$help" 显示菜单
	case strings.HasPrefix(cpf.Message, yamlConf.AdminUserOrderHeader+"help"):
		//pass
	// case "$bash" 对系统执行命令
	case strings.HasPrefix(cpf.Message, yamlConf.AdminUserOrderHeader+"bash"):
		cpf.BashCommand()
	}

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
	_, _ = cpf.SendMsg(cpf.MessageType, res)
}
