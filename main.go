package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"io/ioutil"
	"net/http"
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

// GinAuthentication
// @Description: gin中间件,如果开启反向鉴权(reverseAuthentication)时,对数据进行验证
// @return gin.HandlerFunc
//
func GinAuthentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		if yamlConfig.RAuth.Enable {
			body, _ := ioutil.ReadAll(c.Request.Body)
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body)) // 重设body
			headerXSignature := c.Request.Header.Get("X-Signature")  // sha1签名
			//headerXSelfId := c.Request.Header.Get("X-Self-Id")      // 发送消息的qq
			//fmt.Println(headerXSignature[len("sha1="):], headerXSelfId)
			if headerXSignature[len("sha1="):] != hmacSHA1Encrypt(yamlConfig.RAuth.Secret, body) {
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
//  getSHA256
//  @Description: 得到SHA256之后的密钥
//  @param str
//  @return string
//
func getSHA256(str string) string {
	sha256Bytes := sha256.Sum256([]byte(str))
	return hex.EncodeToString(sha256Bytes[:])
}

func listenFromSendPrivateMsg(c *gin.Context) {
	// 如果revue没开启直接结束
	if !yamlConfig.Revue.Enable {
		c.Abort()
	}
	var form revueSendMsgApi
	// token正确时
	if c.ShouldBind(&form) == nil &&
		getSHA256(yamlConfig.Revue.Secret) == form.Token {
		fmt.Printf("%#v\n", form)
	}
}

func main() {
	yamlConfig.getConf() // 读取配置
	router := gin.Default()

	router.POST("/send_private_msg", listenFromSendPrivateMsg)
	//router.Use(GinAuthentication()) // 下面的路由都使用这个中间件
	// 监听动作并做出反应
	router.POST("/", listenFromCqhttp)
	router.Run(":5000")
}
