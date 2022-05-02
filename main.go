package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/gojsonq"
	"net/http"
	"strconv"
)

//type cqPostFrom struct {
//	Time        int    `json:"time"`
//	SelfId      int    `json:"self_id"`
//	PostType    string `json:"post_type"`
//	MessageType string `json:"message_type"`
//	SubType     string `json:"sub_type"`
//	MessageId   int    `json:"message_id"`
//	UserId      int    `json:"user_id"`
//	Message     string `json:"message"`
//	RawMessage  string `json:"raw_message"`
//	Font        int    `json:"font"`
//	Sender      struct {
//		Nickname string `json:"nickname"`
//		Sex      string `json:"sex"`
//		Age      int    `json:"age"`
//	} `json:"sender"`
//}

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

//
//  getGroupId
//  @Description: 从消息id中得到qq群号,用来判断是否在监听群号列表中
//  @param fromInfo
//  @return string
//  @return bool
//
func getGroupId(fromInfo cqPostFrom) (string, bool) {
	msg, err := getMsg(strconv.Itoa(fromInfo.MessageId))
	if err != nil {
		return "", false
	}
	fJson := gojsonq.New().JSONString(msg)
	groupId := strconv.Itoa(int(fJson.Reset().Find("data.group_id").(float64)))
	for _, s := range yamlConfig.listenGroup {
		if s == groupId {
			return groupId, true
		}
	}
	// 返回id,但是不在监听群列表中
	return groupId, false
}

//
//  groupEvent
//  @Description: 群消息事件
//  @param fromInfo
//  @param groupId
//
func groupEvent(fromInfo cqPostFrom, groupId string) {
	// 只是demo
	if fromInfo.Message == "叫两声" {
		sendGroupMsg(groupId, "汪汪", "false")
	}
}

//
//  msgEvent
//  @Description: 对message事件进行相应
//  @param fromInfo
//
func msgEvent(fromInfo cqPostFrom) {
	// 群消息进行响应
	if fromInfo.MessageType == "group" {
		groupId, ok := getGroupId(fromInfo)
		if ok {
			// 如果是监听qq群列表的才做出相应
			groupEvent(fromInfo, groupId)
		}
	}
	//
	if fromInfo.MessageType == "private" {

	}
}

//
//  listenFromCqhttp
//  @Description: 监听go-cqhttp动作并以此做出反应
//  @param c
//
func listenFromCqhttp(c *gin.Context) {
	var form cqPostFrom
	if c.ShouldBind(&form) == nil {
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
// @Description: gin中间件,如果反向鉴权(reverseAuthentication)时,对数据进行验证
// @return gin.HandlerFunc
//
func GinAuthentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		if yamlConfig.rAuth.enable {
			buf := make([]byte, 1024)
			rBody, _ := c.Request.Body.Read(buf)
			headerXSignature := c.Request.Header.Get("X-Signature") // sha1签名
			//headerXSelfId := c.Request.Header.Get("X-Self-Id")      // 发送消息的qq
			//fmt.Println(headerXSignature[len("sha1="):], headerXSelfId)
			if headerXSignature[len("sha1="):] != hmacSHA1Encrypt(yamlConfig.rAuth.tokenOrSecret, buf[0:rBody]) {
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

func main() {
	readConfig() // 读取配置
	router := gin.Default()
	router.Use(GinAuthentication()) // 下面的路由都使用这个中间件
	// 监听动作并做出反应
	router.POST("/", listenFromCqhttp)
	router.Run(":5000")
}
