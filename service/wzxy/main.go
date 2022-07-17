package wzxy

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/go-resty/resty/v2"
	"github.com/thedevsaddam/gojsonq"
	"log"
	"strconv"
	"time"
)

func getDate() string {
	return time.Now().Format("20060102")
}

func getSha256(src string) string {
	sha256Bytes := sha256.Sum256([]byte(src))
	sha256String := hex.EncodeToString(sha256Bytes[:])
	return sha256String
}

func (u UserWzxy) CheckOperate(seq int) int {
	now := time.Now()
	signTime := strconv.FormatInt(now.UnixNano()/1e6, 10) //时间戳精确到毫秒
	content := u.Province + "_" + signTime + "_" + u.City
	signatureHeader := getSha256(content)
	client := resty.New()
	post, err := client.R().SetHeaders(map[string]string{
		"JWSESSION":  u.Jwsession,
		"User-Agent": u.UserAgent,
	}).SetQueryParams(map[string]string{
		"answers":         "[\"0\"]",
		"seq":             strconv.Itoa(seq),
		"temperature":     "36.5",
		"userId":          "",
		"latitude":        "34.108216",
		"longitude":       "108.605084",
		"country":         "中国",
		"city":            "西安市",
		"district":        "鄠邑区",
		"province":        "陕西省",
		"township":        "甘亭街道",
		"street":          "东街",
		"myArea":          "610118",
		"timestampHeader": signTime,
		"signatureHeader": signatureHeader,
	}).Post("https://student.wozaixiaoyuan.com/heat/save.json")
	if err != nil {
		log.Println(u.Name, "打卡失败，网络错误", "seq=", seq)
		return -1
	}
	postJson := gojsonq.New().JSONString(string(post.Body()))
	if int(postJson.Reset().Find("code").(float64)) == 0 {
		log.Println(u.Name, "打卡成功", "seq=", seq)
		// 正常
		return 0
	} else {
		log.Println(u.Name, "打卡失败,jwsession可能失效", "seq=", seq)
		return -2
	}
}

func (u UserWzxy) getSignMessage() (res int, signId, logId string) {
	client := resty.New()
	post, err := client.R().SetHeaders(map[string]string{
		"jwsession": u.Jwsession,
	}).SetQueryParams(map[string]string{
		"page": "1",
		"size": "5",
	}).Post("https://student.wozaixiaoyuan.com/sign/getSignMessage.json")
	if err != nil {
		return -1, "", ""
	}
	//fmt.Println(string(post.Body()))
	pJson := gojsonq.New().JSONString(string(post.Body()))
	if int(pJson.Reset().Find("code").(float64)) == 0 {
		signTimeStart := pJson.Reset().Find("data.[0].start")
		signTimeEnd := pJson.Reset().Find("data.[0].end")
		timeNow := time.Now().Format("2006-01-02 15:04:05")
		if timeNow > signTimeStart.(string) && timeNow < signTimeEnd.(string) {
			// 在签到区间
			signId = pJson.Reset().Find("data.[0].id").(string)
			logId = pJson.Reset().Find("data.[0].logId").(string)
			return 0, signId, logId
		} else {
			// 在签到区间
			signId = pJson.Reset().Find("data.[0].id").(string)
			logId = pJson.Reset().Find("data.[0].logId").(string)
			//fmt.Println(signId, logId)
			// 不在签到区间
			return -2, "", ""
		}
	}
	return -1, "", ""
}

func (u UserWzxy) doEveningCheck(signId, logId string) int {
	url := "https://student.wozaixiaoyuan.com/sign/doSign.json"
	client := resty.New()
	post, err := client.R().SetHeaders(map[string]string{
		"JWSESSION":  u.Jwsession,
		"User-Agent": u.UserAgent,
	}).SetBody(map[string]string{
		"signId":    signId,
		"city":      "西安市",
		"id":        logId,
		"latitude":  "34.10154079861111",
		"longitude": "108.65831163194444",
		"country":   "中国",
		"district":  "鄠邑区",
		"township":  "五竹街道",
		"province":  "陕西省",
	}).Post(url)
	if err != nil {
		log.Println(u.Name, "晚打卡失败，网络错误", err.Error())
		return -1
	}
	pJson := gojsonq.New().JSONString(string(post.Body()))
	if int(pJson.Reset().Find("code").(float64)) == 0 {
		log.Println(u.Name, "晚检签到成功")
		return 0
	} else {
		log.Println(u.Name, "晚检签到失败,返回信息为:", string(post.Body()))
		return -2
	}
}

func (u UserWzxy) EveningCheckOperate() int {
	res, signId, logId := u.getSignMessage()
	switch res {
	case 0:
		// 正常执行签到
		// -2 晚打卡失败, -1 网络错误, 0 正常
		return u.doEveningCheck(signId, logId)
	case -1:
		log.Println(u.Name, "获取晚检信息失败,网络错误")
		return -1
	case -2:
		log.Println(u.Name, "晚检签到失败,不在签到时间范围内")
		return -3 // 不在签到时间范围内
	default:
		log.Println(u.Name, "获取晚检信息失败,未知错误")
		return -4 // 未知错误
	}
}

//
//func operation() {
//	dataNow := getDate()
//	dateTmp := ""
//	var err error
//	eventMap := make(map[string]map[string]int) // 记录今日任务执行flag
//
//	for {
//		dataNow = getDate()
//		// 第二天执行刷新
//		if dataNow != dateTmp {
//			dateTmp = getDate()
//			yamlConfig, err = NewConf("./config.yaml")
//			if err != nil {
//				return // 读取错误退出
//			}
//			// 刷新flag,0为今日未执行
//			for _, user := range yamlConfig.User {
//				eventMap[user.Name] = map[string]int{"morning": 0}
//				eventMap[user.Name] = map[string]int{"afternoon": 0}
//				eventMap[user.Name] = map[string]int{"evening": 0}
//			}
//		}
//		timeNow := time.Now().Format("15:04:05")
//		for _, user := range yamlConfig.User {
//			if user.MorningCheck.Enable &&
//				timeNow < user.MorningCheck.EndTime &&
//				timeNow > user.MorningCheck.CheckTime &&
//				eventMap[user.Name]["morning"] != 1 {
//				eventMap[user.Name]["morning"] = 1 // flag 置为1
//				// 晨检
//				go user.CheckOperate(1)
//			}
//
//			if user.AfternoonCheck.Enable &&
//				timeNow < user.AfternoonCheck.EndTime &&
//				timeNow > user.AfternoonCheck.CheckTime &&
//				eventMap[user.Name]["afternoon"] != 1 {
//				eventMap[user.Name]["afternoon"] = 1 // flag 置为1
//				// 午检
//				go user.CheckOperate(2)
//			}
//
//			if user.EveningCheck.Enable &&
//				timeNow < user.EveningCheck.EndTime &&
//				timeNow > user.EveningCheck.CheckTime &&
//				eventMap[user.Name]["evening"] != 1 {
//				eventMap[user.Name]["evening"] = 1
//				// 晚检
//				go user.EveningCheckOperate()
//			}
//		}
//		time.Sleep(10 * time.Second)
//	}
//}

//func main() {
//	operation()
//}
