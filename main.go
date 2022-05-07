package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/thedevsaddam/gojsonq"
)

// 用于接收go-cqhttp消息
type cqPostFrom struct {
	Font        int    `json:"font"`
	Message     string `json:"message"`
	MessageId   int    `json:"message_id"`
	MessageType string `json:"message_type"`
	PostType    string `json:"post_type"`
	RawMessage  string `json:"raw_message"`
	SelfId      int64  `json:"self_id"`
	Sender      struct {
		Age      int    `json:"age"`
		Nickname string `json:"nickname"`
		Sex      string `json:"sex"`
		UserId   int    `json:"user_id"`
	} `json:"sender"`
	SubType  string `json:"sub_type"`
	TargetId int64  `json:"target_id"`
	Time     int    `json:"time"`
	UserId   int    `json:"user_id"`
}

// revue发送消息api,用于接收
type revueSendMsgApi struct {
	Token      string `json:"token"`       // 加密后的密钥
	UserId     string `json:"user_id"`     // qq号
	Message    string `json:"message"`     // 消息内容
	AutoEscape string `json:"auto_escape"` // 消息内容是否作为纯文本发送(即不解析CQ码)
}

//
//  listenFromCqhttp
//  @Description: 监听go-cqhttp动作并以此做出反应
//  @param c
//
func listenFromCqhttp(c *gin.Context) {
	var form cqPostFrom
	if c.ShouldBindBodyWith(&form, binding.JSON) == nil {
		// 对message事件进行响应
		if form.PostType == "message" {
			msgEvent(form)
		}
	}
}

//  hmacSHA1Encrypt
//  @Description: SHA1 加密进行鉴权
//  @param encryptKey 密钥
//  @param encryptText 签名主体
//  @return string
//
func hmacSHA1Encrypt(encryptKey string, encryptText []byte) string {
	key := []byte(encryptKey)
	mac := hmac.New(sha1.New, key)
	mac.Write(encryptText)
	var str = hex.EncodeToString(mac.Sum(nil))
	return str
}

// GinReverseAuthentication
// @Description: gin中间件,如果开启反向鉴权(reverseAuthentication)时,对数据进行验证
// @return gin.HandlerFunc
//
func GinReverseAuthentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		if yamlConfig.ReverseAuthentication.Enable {
			body, _ := ioutil.ReadAll(c.Request.Body)
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body)) // 重设body
			headerXSignature := c.Request.Header.Get("X-Signature")  // sha1签名
			//headerXSelfId := c.Request.Header.Get("X-Self-Id")      // 发送消息的qq
			if headerXSignature[len("sha1="):] != hmacSHA1Encrypt(yamlConfig.ReverseAuthentication.Secret, body) {
				c.JSON(
					http.StatusForbidden,
					gin.H{
						"code": http.StatusForbidden,
						"msg":  "密钥错误",
					},
				)
				c.Abort()
			}
		}
	}
}

//
//  listenFromSendPrivateMsg
//  @Description: 监听revue发送私聊消息的接口
//  @param c
//
func listenFromSendPrivateMsg(c *gin.Context) {
	// 如果revue没开启直接结束
	if !yamlConfig.Revue.Enable {
		c.Abort()
	}
	var form revueSendMsgApi
	if c.ShouldBindBodyWith(&form, binding.JSON) == nil {
		// do event
		autoEscape := ""
		if strings.ToUpper(form.AutoEscape) == "TRUE" {
			autoEscape = "true"
		} else {
			autoEscape = "false"
		}
		msg, err := sendMsg("private", form.UserId, "", form.Message, autoEscape)
		if err != nil {
			c.JSON(
				http.StatusInternalServerError,
				gin.H{
					"code": http.StatusInternalServerError,
					"msg":  msg, // 返回错误信息
				},
			)
		} else {
			c.JSON(
				http.StatusOK,
				gin.H{
					"code": http.StatusOK,
					"msg":  msg, // 正确返回msg id
				},
			)
		}

	} else {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"code": http.StatusBadRequest,
				"msg":  "请求参数不能识别",
			},
		)
	}
}

//
// GinRevueAuthentication
// @Description: revue接口中间件,对发送的token进行验证
// @return gin.HandlerFunc
//
func GinRevueAuthentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		if yamlConfig.ReverseAuthentication.Enable {
			body, _ := ioutil.ReadAll(c.Request.Body)
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body)) // 重设body
			pJson := gojsonq.New().JSONString(string(body))
			pToken := pJson.Reset().Find("token").(string)
			pUserId := pJson.Reset().Find("user_id").(string)
			if res, permission := searchRevueApiToken(pUserId, pToken); !res || permission < 4 {
				c.JSON(
					http.StatusForbidden,
					gin.H{
						"code": http.StatusForbidden,
						"msg":  "密钥错误",
					},
				)
				c.Abort() // 结束会话
			}
		}
	}
}

func main() {
	gin.DisableConsoleColor()
	dbInit()                            // 初始化数据库
	yamlConfig.getConf("./config.yaml") // 读取配置
	router := gin.Default()
	// 监听动作并做出反应
	router.POST("/", GinReverseAuthentication(), listenFromCqhttp)
	// 监听revue提供发送消息的接口
	router.POST("/send_private_msg", GinRevueAuthentication(), listenFromSendPrivateMsg)
	router.Run("0.0.0.0:" + yamlConfig.ListenPort)
}
