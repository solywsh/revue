package wzxy

import (
	"github.com/go-resty/resty/v2"
	"github.com/thedevsaddam/gojsonq"
	"log"
	"strconv"
	"strings"
	"time"
)

type UserWzxy struct {
	ID              uint   `gorm:"primaryKey;autoIncrement"`
	Jwsession       string // JWSESSION
	JwsessionStatus bool   // JWSESSION是否有效
	Token           string // token,有效应验证
	UserId          string // 用户ID/QQ
	Name            string // 用户名

	MorningCheckEnable   bool   // 晨检打卡是否开启
	MorningCheckTime     string // 晨检打卡时间
	MorningLastCheckDate string // 晨检打卡日期

	AfternoonCheckEnable   bool   // 午检打卡是否开启
	AfternoonCheckTime     string // 午检打卡时间
	AfternoonLastCheckDate string // 午检打卡日期

	EveningCheckEnable   bool   // 晚检打卡是否开启
	EveningCheckTime     string // 晚检打卡时间
	EveningLastCheckDate string // 晚检打卡日期

	Province  string // 省份
	City      string // 城市
	UserAgent string // UserAgent

	Deadline time.Time // 过期时间
}

type TokenWzxy struct {
	ID           uint      `gorm:"primaryKey;autoIncrement"`
	Token        string    // token,用于注册
	Deadline     time.Time // 过期时间
	CreateUser   string    // 创建人,默认只能管理员
	Status       int       // 状态,0未使用,1已经使用,2可多次使用(0和1针对单次使用)
	Times        int       // 可使用次数	默认1次
	Organization string    // 组织机构,默认为private,多次使用时需要修改
}

func (u UserWzxy) CheckOperate(seq int) (res int, message string) {
	client := resty.New()
	postStr := ""
	payload := strings.NewReader(`{"location":"中国/陕西省/西安市/鄠邑区/五竹街道/长庆石化路/156/610118/156610100/610118003","t1":"是","t2":"绿色","t3":"是","type":0,"locationType":0}`)
	post, err := client.R().SetHeaders(map[string]string{
		"JWSESSION":  u.Jwsession,
		"User-Agent": u.UserAgent,
		"Cookie":     "JWSESSION=" + u.Jwsession,
		"Referer":    "https://gw.wozaixiaoyuan.com/h5/mobile/health/index/health/detail?id=170000" + strconv.Itoa(seq),
	}).SetBody(payload).Post("https://gw.wozaixiaoyuan.com/health/mobile/health/save?batch=170000" + strconv.Itoa(seq))
	if err != nil {
		log.Println(u.Name, "打卡失败，网络错误", "seq=", seq, err.Error())
		return -1, "网络错误"
	}
	postStr = string(post.Body())
	postJson := gojsonq.New().JSONString(postStr)
	if int(postJson.Reset().Find("code").(float64)) == 0 {
		log.Println(u.Name, "打卡成功", "seq=", seq)
		// 正常
		return 0, ""
	} else {
		log.Println(u.Name, "打卡失败", "seq=", seq, postStr)
		res = int(postJson.Reset().Find("code").(float64))
		message = postJson.Reset().Find("message").(string)
		return res, message
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
		"towncode":  "",
		"citycode":  "",
		"areacode":  "",
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

func (u UserWzxy) String() string {
	msg := "我在校园用户信息\n"
	msg += "打卡任务:\n"
	msg += "id：" + u.UserId + "\n"
	msg += "name：" + u.Name + "\n"
	msg += "token：" + u.Token + "\n"
	msg += "jwsession：" + u.Jwsession + "\n"

	msg += "jwsession有效性："
	if u.JwsessionStatus {
		msg += "有效\n"
	} else {
		msg += "无效\n"
	}
	msg += "晨检打卡状态："
	if u.MorningCheckEnable {
		msg += "开启\n"
	} else {
		msg += "关闭\n"
	}
	msg += "晨检打卡时间：" + u.MorningCheckTime + "\n"
	msg += "午检打卡状态："
	if u.AfternoonCheckEnable {
		msg += "开启\n"
	} else {
		msg += "关闭\n"
	}
	msg += "午检打卡时间：" + u.AfternoonCheckTime + "\n"
	msg += "签到(晚检)打卡状态："
	if u.EveningCheckEnable {
		msg += "开启\n"
	} else {
		msg += "关闭\n"
	}
	msg += "签到(晚检)打卡时间：" + u.EveningCheckTime + "\n"
	msg += "有效期：" + u.Deadline.Format("2006-01-02 15:04:05") + "\n"
	return msg
}

func (wt TokenWzxy) String() string {
	msg := "token：" + wt.Token + "\n"
	msg += "有效期至：" + wt.Deadline.Format("2006-01-02 15:04") + "\n"
	msg += "用户：" + wt.CreateUser + "\n"
	if wt.Status == 0 {
		msg += "token状态：未使用\n"
	} else if wt.Status == 1 {
		msg += "token状态：已使用\n"
	} else if wt.Status == 2 {
		msg += "token状态：多次使用\n"
	} else {
		msg += "token状态：未知\n"
	}
	msg += "剩余次数：" + strconv.Itoa(wt.Times) + "\n"
	msg += "组织：" + wt.Organization + "\n"
	return msg
}
