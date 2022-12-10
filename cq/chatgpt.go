package cq

import (
	"errors"
	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/solywsh/chatgpt"
	"github.com/solywsh/qqBot-revue/db"
	"strconv"
	"strings"
	"time"
)

const ChatGPTName = "chat_gpt"

var (
	chatMap = cmap.New[*Chat]()
)

type Chat struct {
	ChatGPT     *chatgpt.ChatGPT
	UserSession *db.UserSession
}

// ChatGPTEvent status = 1 初始化，status = 2 判断/响应 status = 3 结束
func (cpf *PostForm) ChatGPTEvent(status int) {
	switch status {
	case 1:
		userSession, err := gdb.FindOrCreateUserSession(strconv.Itoa(cpf.UserId), ChatGPTName)
		if err != nil {
			cpf.SendMsg("发生错误" + err.Error())
			return
		}
		userSession.Status = 1 // 激活状态
		userSession.UpdateTime = time.Now()
		ok, err := gdb.UpdateUserSessionMany(userSession, true)
		if err != nil || !ok {
			cpf.SendMsg("发生错误" + err.Error())
			return
		}
		chatMap.Set(strconv.Itoa(cpf.UserId), &Chat{
			ChatGPT:     chatgpt.New(yamlConf.ChatGPT.ApiKey, strconv.Itoa(cpf.UserId), time.Minute*10),
			UserSession: &userSession,
		})
		go func() {
			if v, ok := chatMap.Get(strconv.Itoa(cpf.UserId)); ok {
				select {
				case <-v.ChatGPT.GetDoneChan():
					cpf.SendMsg(GetCqCodeAt(strconv.Itoa(cpf.UserId), "") + " 结束与你的对话")
					chatMap.Remove(strconv.Itoa(cpf.UserId))
				}
			}
		}()
		cpf.SendMsg("请说")
	case 2:
		chat, ok := chatMap.Get(strconv.Itoa(cpf.UserId))
		if !ok || chat.UserSession.AppName != ChatGPTName || chat.UserSession.Status != 1 {
			return
		}
		ans, err := chat.ChatGPT.ChatWithContext(cpf.Message)
		if err != nil {
			switch {
			case strings.Contains(err.Error(), "Timeout"):
				cpf.SendMsg(GetCqCodeAt(strconv.Itoa(cpf.UserId), "") + " 回答超时")
			case strings.Contains(err.Error(), "context canceled"):
				cpf.SendMsg(GetCqCodeAt(strconv.Itoa(cpf.UserId), "") + " 对话已关闭")
			case errors.As(err, &chatgpt.OverMaxSequenceTimes):
				cpf.SendMsg(GetCqCodeAt(strconv.Itoa(cpf.UserId), "") + " 超过最大对话次数，如需继续请重新开始")
			case errors.As(err, &chatgpt.OverMaxTextLength):
				cpf.SendMsg(GetCqCodeAt(strconv.Itoa(cpf.UserId), "") + " 总文本超过最大文本长度，请缩短文本或重新开始")
			default:
				cpf.SendMsg("未知错误" + err.Error())
			}
			return
		}
		chat.UserSession.UpdateTime = time.Now()
		ans = formatAnswer(ans)
		if ans == "" {
			cpf.SendMsg("AI返回了空结果")
		} else {
			cpf.SendMsg(ans)
		}
	case 3:
		chat, ok := chatMap.Get(strconv.Itoa(cpf.UserId))
		if !ok || chat.UserSession.AppName != ChatGPTName || chat.UserSession.Status != 1 {
			return
		}
		chat.UserSession.Status = 2 // 关闭状态
		chat.UserSession.UpdateTime = time.Now()
		ok, err := gdb.UpdateUserSessionMany(*chat.UserSession, true)
		if err != nil || !ok {
			cpf.SendMsg("发生错误" + err.Error())
			return
		}
		if v, ok := chatMap.Get(strconv.Itoa(cpf.UserId)); ok {
			v.ChatGPT.Close()
		}
	}
}

func formatAnswer(answer string) string {
	for len(answer) > 0 {
		switch {
		case answer[0] == '\n':
			answer = answer[1:]
		case answer[0] == ' ':
			answer = answer[1:]
		case answer[0] == '\t':
			answer = answer[1:]
		case answer[0] == '\r':
			answer = answer[1:]
		case answer[0] == '?':
			answer = answer[1:]
		case answer[0] == '.':
			answer = answer[1:]
		case answer[0] == '!':
			answer = answer[1:]
		case strings.HasPrefix(answer, "？"):
			answer = strings.Replace(answer, "？", "", 1)
		case strings.HasPrefix(answer, "。"):
			answer = strings.Replace(answer, "。", "", 1)
		default:
			return answer
		}
	}
	return answer
}
