package db

import (
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

// ProgrammerAlmanac 程序员黄历
type ProgrammerAlmanac struct {
	ID      uint   `gorm:"primaryKey;autoIncrement"`
	Date    string //时间
	Almanac string //黄历内容
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
	// 读取配置
	yamlConfig, _ := conf.NewConf("./config.yaml")
	gb = new(GormDb) // 由于定义的是地址,在使用前需要先分配内存
	gb.DB, _ = gorm.Open(
		sqlite.Open(yamlConfig.Database.Path),
		&gorm.Config{Logger: newLogger})
	// 自动迁移,如果数据库不存在对应表，则创建
	err := gb.DB.AutoMigrate(RevueConfig{}, KeywordsReply{}, RevueApiToken{}, ProgrammerAlmanac{})
	if err != nil {
		return nil
	}
	// 不存在则自动创建(理论上配置只有一条记录，所以ID只能为1)
	gb.DB.Where(RevueConfig{ID: 1}).Attrs(RevueConfig{ID: 1, ReplyEnable: true, MusicEnable: true}).FirstOrCreate(&RevueConfig{})
	return gb
}
