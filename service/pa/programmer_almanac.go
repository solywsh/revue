package pa

import (
	"math/rand"
	"time"
)

// 得到黄历
func getOld() string {
	activities := []map[string]string{
		{"name": "写单元测试", "good": "写单元测试将减少出错", "bad": "写单元测试会降低你的开发效率"},
		{"name": "洗澡", "good": "你几天没洗澡了？", "bad": "会把设计方面的灵感洗掉"},
		{"name": "锻炼一下身体", "good": "瘦10斤", "bad": "能量没消耗多少，吃得却更多"},
		{"name": "抽烟", "good": "抽烟有利于提神，增加思维敏捷", "bad": "除非你活够了，死得早点没关系"},
		{"name": "白天上线", "good": "今天白天上线是安全的", "bad": "可能导致灾难性后果"},
		{"name": "重构", "good": "代码质量得到提高", "bad": "你很有可能会陷入泥潭"},
		{"name": "使用%t", "good": "你看起来更有品位", "bad": "别人会觉得你在装逼"},
		{"name": "跳槽", "good": "该放手时就放手", "bad": "鉴于当前的经济形势，你的下一份工作未必比现在强"},
		{"name": "招人", "good": "你遇到千里马的可能性大大增加", "bad": "你只会招到一两个混饭吃的外行"},
		{"name": "面试", "good": "面试官今天心情很好", "bad": "面试官不爽，会拿你出气"},
		{"name": "提交辞职申请", "good": "公司找到了一个比你更能干更便宜的家伙，巴不得你赶快滚蛋", "bad": "鉴于当前的经济形势，你的下一份工作未必比现在强"},
		{"name": "申请加薪", "good": "老板今天心情很好", "bad": "公司正在考虑裁员"},
		{"name": "晚上加班", "good": "晚上是程序员精神最好的时候", "bad": "身心憔悴，早点休息"},
		{"name": "在妹子面前吹牛", "good": "改善你矮穷挫的形象", "bad": "会被识破"},
		{"name": "stack overflow", "good": "避免缓冲区溢出", "bad": "灰飞烟灭"},
		{"name": "浏览非法网站", "good": "重拾对生活的信心", "bad": "你会心神不宁"},
		{"name": "命名变量", "good": "变量名萌萌哒", "bad": "这个变量永远引用不到"},
		{"name": "写超过%l行的方法", "good": "你的代码组织的很好，长一点没关系", "bad": "你的代码将混乱不堪，你自己都看不懂"},
		{"name": "提交代码", "good": "遇到冲突的几率是最低的", "bad": "你遇到的一大堆冲突会让你觉得自己是不是时间穿越了"},
		{"name": "代码复审", "good": "发现重要问题的几率大大增加", "bad": "你什么问题都发现不了，白白浪费时间"},
		{"name": "开会", "good": "写代码之余放松一下打个盹，有益健康", "bad": "你会被扣屎盆子背黑锅"},
		{"name": "打雀魂", "good": "你将有如神助", "bad": "你会被虐的很惨"},
		{"name": "晚上上线", "good": "晚上是程序员精神最好的时候", "bad": "你白天已经筋疲力尽了"},
		{"name": "修复BUG", "good": "你今天对BUG的嗅觉大大提高", "bad": "新产生的BUG将比修复的更多"},
		{"name": "设计评审", "good": "设计评审会议将变成头脑风暴", "bad": "人人筋疲力尽，评审就这么过了"},
		{"name": "需求评审", "good": "这个需求很简单", "bad": "公司需要一个能根据手机外壳变化APP皮肤的功能"},
		{"name": "上微博", "good": "今天发生的事不能错过", "bad": "会被老板看到"},
		{"name": "上AB站", "good": "还需要理由吗？", "bad": "会被老板看到"},
		{"name": "打守望先锋", "good": "你将有如神助", "bad": "你会被虐的很惨"},
		{"name": "在维基萌抽卡", "good": "大概率抽到了自己心仪的卡", "bad": "垃圾卡片满天飞"},
		{"name": "写技术文章", "good": "新的水文即将诞生", "bad": "你的博文会被抄袭"},
	}
	var good, bad string
	// 随机种子
	rand.Seed(time.Now().Unix())
	for i := 0; i < 4; i++ {
		index := rand.Intn(len(activities))
		if i%2 == 0 {
			good += "\t" + activities[index]["name"] + ":" + activities[index]["good"] + "\n"
		} else {
			bad += "\t" + activities[index]["name"] + ":" + activities[index]["bad"] + "\n"
		}
		activities = append(activities[:index], activities[index:]...) // 删除对应的选项
	}
	return "宜:\n" + good + "不宜:\n" + bad
}

// 得到日期
func getDate() string {
	t := time.Now()
	weekMap := []string{"日", "一", "二", "三", "四", "五", "六"}
	return t.Format("今天是2006年01月02号") + " 星期" + weekMap[t.Weekday()]
}

// 得到方位
func getDirections() string {
	directions := []string{"北方", "东北方", "东方", "东南方", "南方", "西南方", "西方", "西北方"}
	// 随机种子
	rand.Seed(time.Now().Unix())
	return directions[rand.Intn(len(directions))]
}

// 得到宜饮
func getDrink() string {
	drinks := []string{
		"水", "茶", "红茶", "绿茶", "咖啡", "奶茶", "可乐",
		"牛奶", "豆奶", "果汁", "果味汽水", "苏打水", "运动饮料",
		"酸奶", "燕京", "崂山", "雪花",
		"大乌苏", "二锅头", "五粮液", "茅台", "剑南春", "青岛",
		"大黑啤", "哈尔滨啤酒", "喜力啤酒", "威士忌🥃",
	}
	return drinks[rand.Intn(len(drinks))]
}

// 得到app
func getApp() string {
	App := []string{
		"上鼠鼠的B站吧,李沐可能更新了", "看看AcFan吧，看下摇曳露营更新了没", "去telegram看看有啥黑料",
		"你有多久没去蓝鸟了", "小黑盒上你愿望单的游戏正在更新", "你需要去学习强国刷积分", "去知乎键政一波",
		"晚餐用美团点个外卖", "是不是该去淘宝剁手一波了", "到点了,请打开网易云", "今天也要keep哟",
	}
	// 随机种子
	rand.Seed(time.Now().Unix())
	return App[rand.Intn(len(App))]
}

// 得到香烟
func getSmoke() string {
	smokes := []string{
		"利群",
		"砖石",
		"白沙",
		"中南海",
		"新石家庄",
		"万宝路",
		"紫云",
		"玉溪",
		"一根华子",
	}
	// 随机种子
	rand.Seed(time.Now().Unix())
	return smokes[rand.Intn(len(smokes))]
}

// NewCalendar 得到日历的集合
func NewCalendar() (res string) {
	// 随机种子
	rand.Seed(time.Now().Unix())
	event := []string{getDrink(), getApp(), getSmoke()}
	index := rand.Intn(len(event))
	res += getDate() + "\n" + getOld()
	res += "\n座位朝向:面向" + getDirections() + "写代码,bug最少\n"
	switch index {
	case 0:
		res += "今日宜饮:" + event[index]
	case 1:
		res += "今日宜刷app:" + event[index]
	case 2:
		res += "今日宜吸烟:" + event[index]
	}
	return res
}
