package db

import (
	uuid "github.com/satori/go.uuid"
	"github.com/solywsh/qqBot-revue/conf"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

// RevueConfig 根据命令对机器人的一些配置进行动态配置
type RevueConfig struct {
	ID          uint `gorm:"primaryKey;autoIncrement"`
	ReplyEnable bool // 开启回复
	MusicEnable bool // 开启音乐推荐
}

// KeywordsReply 关键词触发消息
type KeywordsReply struct {
	ID   uint `gorm:"primaryKey;autoIncrement"`
	Mode uint // 匹配模式 0:不匹配 1:完全匹配 2:存在匹配(find)
	Flag uint // 设置标记,
	// 触发“开始添加”时,插入userid记录,其他为空,flag=1
	// 该userid回复的第二条msg设置为keywords,flag=2
	// 该userid回复的第三条msg设置为msg,flag=0,添加完成
	Keywords string // 关键词
	Msg      string // 回复消息
	Userid   string // 设置人的qq
}

// RevueApiToken 对每个qq生成对应token
type RevueApiToken struct {
	ID         uint   `gorm:"primaryKey;autoIncrement"`
	UserId     string //`gorm:"primaryKey"` //qq号
	Token      string //生成的token,这里采用uuid
	Permission uint   //权限
	//  Permission
	//  @Description: 权限
	//	私聊:群聊:通过群私聊(不是好友的情况下) 对应 4:2:1
	//	比如:
	//	Permission==4,则只能私聊
	//	Permission==6,则可以让机器人发送群消息,也可以私聊 //注:目前只有实现私聊接口的打算
	//	Permission==0,则没有权限
}

// GormDb 这里对GormDb重新封装了一下
type GormDb struct {
	DB *gorm.DB
}

// NewDB 重新封装
func NewDB() (gb *GormDb) {
	// 定义gorm日志
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // 慢 SQL 阈值
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  false,         // 禁用彩色打印
		},
	)
	yamlConfig, _ := conf.NewConf("./config.yaml")
	gb = new(GormDb) // 由于定义的是地址,在使用前需要先分配内存
	gb.DB, _ = gorm.Open(
		sqlite.Open(yamlConfig.Database.Path),
		&gorm.Config{Logger: newLogger})
	// 自动迁移,如果数据库不存在对应表，则创建
	err := gb.DB.AutoMigrate(RevueConfig{}, KeywordsReply{}, RevueApiToken{})
	if err != nil {
		return nil
	}
	// 不存在则自动创建(理论上配置只有一条记录，所以ID只能为1)
	gb.DB.Where(RevueConfig{ID: 1}).Attrs(RevueConfig{ID: 1, ReplyEnable: true, MusicEnable: true}).FirstOrCreate(&RevueConfig{})
	return gb
}

// 获取所有的关键词回复
func (gb *GormDb) GetAllKeywordsReply(res *[]KeywordsReply) {
	gb.DB.Find(&res)
}

// 查询并创建,查询到对应的信息,则不创建返回false,没有则创建,返回true
func (gb *GormDb) InsertRevueApiToken(userId string, permission uint) (bool, string) {
	var rt RevueApiToken
	res := gb.DB.Where(RevueApiToken{UserId: userId}).Attrs(RevueApiToken{
		UserId:     userId,
		Token:      uuid.NewV4().String(),
		Permission: permission,
	}).FirstOrCreate(&rt)
	if res.RowsAffected == 1 {
		return true, rt.Token
	} else {
		return false, rt.Token
	}
}

// 根据qq号得到token,注意此函数必须在qq号存在的业务逻辑下使用
func (gb *GormDb) GetRevueApiToken(userId string) string {
	var rt RevueApiToken
	gb.DB.Where(RevueApiToken{UserId: userId}).First(&rt)
	return rt.Token
}

// 根据传入的qq号和token判断是否具有权限,并返回权限
func (gb *GormDb) SearchRevueApiToken(userId, token string) (bool, uint) {
	var rt RevueApiToken
	res := gb.DB.Where(RevueApiToken{UserId: userId, Token: token}).First(&rt)
	if res.RowsAffected >= 1 {
		return true, rt.Permission
	} else {
		return false, rt.Permission
	}
}

// 重新设置token
func (gb *GormDb) ResetRevueApiToken(userId string) (bool, string) {
	var rt RevueApiToken
	res := gb.DB.Where(RevueApiToken{UserId: userId}).First(&rt)
	if res.RowsAffected >= 1 {
		gb.DB.Model(&rt).Update("token", uuid.NewV4().String())
		return true, rt.Token
	} else {
		// 不存在
		return false, rt.Token
	}
}

// 删除token
func (gb *GormDb) DeleteRevueApiToken(userId string) (bool, string) {
	var rt RevueApiToken
	res := gb.DB.Where(RevueApiToken{UserId: userId}).First(&rt)
	if res.RowsAffected >= 1 {
		gb.DB.Unscoped().Delete(&rt)
		return true, rt.Token
	} else {
		// 不存在
		return false, rt.Token
	}
}

// 获取RevueConfig配置
func (gb *GormDb) GetRevueConfig(config *RevueConfig) bool {
	res := gb.DB.Where(RevueConfig{ID: 1}).First(&config)
	if res.RowsAffected >= 1 {
		return true
	}
	return false
}

// 设置music开启和关闭
func (gb *GormDb) SetRevueConfigMusic(enable bool) {
	gb.DB.Model(&RevueConfig{}).Where(RevueConfig{ID: 1}).Update("music_enable", enable)
}

// 设置reply开启和关闭
func (gb *GormDb) SetRevueConfigReply(db *gorm.DB, enable bool) {
	gb.DB.Model(&RevueConfig{}).Where(RevueConfig{ID: 1}).Update("reply_enable", enable)
}

// 插入/更新关键词回复,如果存在对应关键词,则更新,不存在则insert
func (gb *GormDb) UpdateKeywordsReply(info KeywordsReply) {
	var kr KeywordsReply
	if info.Flag == 1 {
		// 第一次创建
		gb.DB.Create(&info)
	} else if info.Flag == 2 {
		// 第二次添加关键词
		gb.DB.Where(KeywordsReply{ID: info.ID}).First(&kr)
		// 如果第二次为更新之前查询,需要更新设置人的qq
		gb.DB.Model(&kr).Updates(KeywordsReply{Keywords: info.Keywords, Userid: info.Userid, Flag: 2})
	} else if info.Flag == 3 {
		// 第三次添加回复和匹配模式
		gb.DB.Where(KeywordsReply{ID: info.ID}).First(&kr)
		gb.DB.Model(&kr).Updates(KeywordsReply{Msg: info.Msg, Mode: info.Mode, Flag: 3})
	}
	//fmt.Println(KR.ID, KR.Keywords, KR.Msg)
}

// 根据userid查找该userid是否存在正在记录中的自动回复
func (gb *GormDb) GetKeywordsReplyFlag(userId string) (bool, KeywordsReply) {
	var kr KeywordsReply
	if res := gb.DB.Where("userid = ? AND flag <> ?", userId, 3).First(&kr); res.RowsAffected >= 1 {
		return true, kr
	} else {
		return false, kr
	}
}

// 根据关键词搜索是否存在记录,如果存在记录则返回对应回复
func (gb *GormDb) SearchKeywordsReply(keywords string) (bool, KeywordsReply) {
	var kr KeywordsReply
	// 存在关键词并且状态flag为0
	if res := gb.DB.Where(KeywordsReply{Keywords: keywords, Flag: 0}).First(&kr); res.RowsAffected >= 1 {
		return true, kr
	}
	return false, kr
}

// 删除对应的关键词回复
func (gb *GormDb) DeleteKeywordsReply(id uint) {
	gb.DB.Delete(&KeywordsReply{}, id)
}
