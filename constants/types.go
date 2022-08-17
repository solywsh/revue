package constants

// 我在校园相关

type UserWzxy struct {
	Jwsession       string // JWSESSION
	JwsessionStatus bool   // JWSESSION是否有效
	Status          int    // 状态码 0 默认定时执行 1 手动执行
	Token           string // token,有效应验证
	UserId          string // 用户ID/QQ
	Name            string // 用户名

	MorningCheckEnable   bool   // 晨检打卡是否开启
	MorningCheckTime     string // 晨检打卡时间
	MorningCheckLastDate string // 晨检打卡最后打卡时间

	AfternoonCheckEnable   bool   // 午检打卡是否开启
	AfternoonCheckTime     string // 午检打卡时间
	AfternoonCheckLastDate string // 午检打卡最后打卡时间

	EveningCheckEnable   bool   // 晚检打卡是否开启
	EveningCheckTime     string // 晚检打卡时间
	EveningCheckLastTime string // 晚检打卡最后打卡时间

	Province  string // 省份
	City      string // 城市
	UserAgent string // UserAgent

	Deadline string // 过期时间
}

type TokenWzxy struct {
	Token        string // token,用于注册
	Deadline     string // 过期时间
	CreateUser   string // 创建人,默认只能管理员
	Status       int    // 状态,0未使用,1使用,2可多次使用
	Times        int    // 可使用次数
	Organization string // 组织机构,默认为private,多次使用时需要修改
}
