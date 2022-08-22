package cq

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/solywsh/qqBot-revue/conf"
	"github.com/solywsh/qqBot-revue/db"
	"github.com/thedevsaddam/gojsonq"
	"strconv"
)

// PostForm 用于接收go-cqhttp消息
type PostForm struct {
	Anonymous     interface{} `json:"anonymous"`
	GroupId       int         `json:"group_id"`
	Interval      int         `json:"interval"`
	MetaEventType string      `json:"meta_event_type"`
	Font          int         `json:"font"`
	Message       string      `json:"message"`
	MessageId     int         `json:"message_id"`
	MessageSeq    int         `json:"message_seq"`
	MessageType   string      `json:"message_type"`
	PostType      string      `json:"post_type"`
	RawMessage    string      `json:"raw_message"`
	SelfId        int64       `json:"self_id"`
	Sender        struct {
		Age      int    `json:"age"`
		Area     string `json:"area"`
		Card     string `json:"card"`
		Level    string `json:"level"`
		Nickname string `json:"nickname"`
		Role     string `json:"role"`
		Sex      string `json:"sex"`
		Title    string `json:"title"`
		UserId   int    `json:"user_id"`
	} `json:"sender"`
	Status struct {
		AppEnabled     bool        `json:"app_enabled"`
		AppGood        bool        `json:"app_good"`
		AppInitialized bool        `json:"app_initialized"`
		Good           bool        `json:"good"`
		Online         bool        `json:"online"`
		PluginsGood    interface{} `json:"plugins_good"`
		Stat           struct {
			PacketReceived  int `json:"PacketReceived"`
			PacketSent      int `json:"PacketSent"`
			PacketLost      int `json:"PacketLost"`
			MessageReceived int `json:"MessageReceived"`
			MessageSent     int `json:"MessageSent"`
			LastMessageTime int `json:"LastMessageTime"`
			DisconnectTimes int `json:"DisconnectTimes"`
			LostTimes       int `json:"LostTimes"`
		} `json:"stat"`
	} `json:"status"`
	SubType  string `json:"sub_type"`
	TargetId int64  `json:"target_id"`
	Time     int    `json:"time"`
	UserId   int    `json:"user_id"`
}

// 定义全局变量
var (
	yamlConf *conf.Config
	gdb      *db.GormDb
)

// 初始化配置
func init() {
	yamlConf = conf.NewConf() // 得到配置文件
	//gdb = new(db.GormDb)
	gdb = db.NewDB() // 初始化操作数据库
}

//  格式化url，根据配置文件是否开启鉴权格式化
func formatAccessUrl(str string) string {
	if yamlConf.ForwardAuthentication.Enable {
		return yamlConf.UrlHeader + str + "?access_token=" + yamlConf.ForwardAuthentication.Token
	} else {
		return yamlConf.UrlHeader + str
	}
}

//
//// PublicSendPrivateMsg 对外的发送私聊消息接口，给其他包调用
//func PublicSendPrivateMsg(userId, msg string) error {
//	client := resty.New()
//	post, err := client.R().SetQueryParams(map[string]string{
//		"user_id":     userId,
//		"group_id":    "",
//		"message":     msg,
//		"auto_escape": "false",
//	}).Post(formatAccessUrl("/send_private_msg"))
//	if err != nil {
//		return err
//	}
//	rJson := gojsonq.New().JSONString(string(post.Body()))
//	if rJson.Reset().Find("retcode") != nil && rJson.Reset().Find("retcode").(float64) != 0.0 {
//		return fmt.Errorf(string(post.Body()))
//	}
//	return nil
//}

// SendPrivateMsg 发送私聊消息
func (cpf *PostForm) SendPrivateMsg(msg string) (string, error) {
	client := resty.New()
	post, err := client.R().SetQueryParams(map[string]string{
		"user_id":     strconv.Itoa(cpf.UserId),
		"group_id":    strconv.Itoa(cpf.GroupId),
		"message":     msg,
		"auto_escape": "false",
	}).Post(formatAccessUrl("/send_private_msg"))
	if err != nil {
		return "", err
	}
	rJson := gojsonq.New().JSONString(string(post.Body()))
	if rJson.Reset().Find("retcode") != nil && rJson.Reset().Find("retcode").(float64) != 0.0 {
		return "", fmt.Errorf(string(post.Body()))
	}
	messageId := strconv.Itoa(int(rJson.Reset().Find("data.message_id").(float64)))
	return messageId, err
}

// SendGroupMsg 发送群消息
func (cpf *PostForm) SendGroupMsg(msg string) (string, error) {
	client := resty.New()
	post, err := client.R().SetQueryParams(map[string]string{
		"group_id":    strconv.Itoa(cpf.GroupId),
		"message":     msg,
		"auto_escape": "false",
	}).Post(formatAccessUrl("/send_group_msg"))
	if err != nil {
		return "", err
	}
	rJson := gojsonq.New().JSONString(string(post.Body()))
	if rJson.Reset().Find("retcode") != nil && rJson.Reset().Find("retcode").(float64) != 0.0 {
		return "", fmt.Errorf(string(post.Body()))
	}
	messageId := strconv.Itoa(int(rJson.Reset().Find("data.message_id").(float64)))
	return messageId, err
}

// SendMsg
//  @Description: 发送消息
//  @param userId private时对方的qq
//  @param groupId group时群号
//  @param message 消息
//  @param autoEscape 是否解析CQ码
//  @return string 返回message_id
//
func (cpf *PostForm) SendMsg(message string) string {
	postData := map[string]string{
		"message_type": cpf.MessageType,
		"user_id":      strconv.Itoa(cpf.UserId),
		"group_id":     strconv.Itoa(cpf.GroupId),
		"message":      message,
		"auto_escape":  "false",
	}
	client := resty.New()
	post, err := client.R().SetQueryParams(postData).Post(formatAccessUrl("/send_msg"))
	if err != nil {
		return "POST ERROR"
	}
	rJson := gojsonq.New().JSONString(string(post.Body()))
	if rJson.Reset().Find("retcode") != nil && rJson.Reset().Find("retcode").(float64) != 0.0 {
		return rJson.Reset().Find("msg").(string)
	}
	messageId := strconv.Itoa(int(rJson.Reset().Find("data.message_id").(float64)))
	return messageId
}

// DeleteMsg
//  @Description: 撤回消息
//  @param messageId 需要撤回的消息Id
//
func (cpf *PostForm) DeleteMsg() {
	client := resty.New()
	_, _ = client.R().SetQueryParams(map[string]string{
		"message_id": strconv.Itoa(cpf.MessageId),
	}).Post(formatAccessUrl("/delete_msg"))
}

// DeleteFriend
//  @Description: 删除好友
//  @param friendId 好友qq号
//
func (cpf *PostForm) DeleteFriend() {
	client := resty.New()
	_, _ = client.R().SetQueryParams(map[string]string{
		"friend_id": strconv.Itoa(cpf.UserId),
	}).Post(formatAccessUrl("/delete_friend"))
}
