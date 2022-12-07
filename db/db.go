package db

import (
	"fmt"
	"github.com/solywsh/qqBot-revue/conf"
	"github.com/solywsh/qqBot-revue/service/wzxy"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"sync"
	"time"
)

var (
	dbOnce sync.Once
	gb     *GormDb
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

// Divination 求签
type Divination struct {
	ID     uint   `gorm:"primaryKey;autoIncrement"`
	Date   string //日期
	Tag    string //求签结果
	UserId string //用户id,QQ号
}

// ListenGroup 监听的群
type ListenGroup struct {
	ID     uint   `gorm:"primaryKey;autoIncrement"`
	Date   string //日期
	UserId string //设置用户id,QQ号
	Group  string //群号
}

// UserSession 用户会话控制
type UserSession struct {
	ID         uint      `gorm:"primaryKey;autoIncrement"`
	UserId     string    //用户id,QQ号
	AppName    string    //应用名称
	Status     uint      //状态,各个服务区分,0为初始
	UpdateTime time.Time //更新时间
}

// GormDb 这里对GormDb重新封装了一下
type GormDb struct {
	DB *gorm.DB
}

// NewDB 重新封装
func NewDB() *GormDb {
	dbOnce.Do(func() {
		// 读取配置
		yamlConfig := conf.NewConf()
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
		gb = new(GormDb) // 由于定义的是地址,在使用前需要先分配内存
		if yamlConfig.Database.Sqlite.Enable {
			log.Println("检测到使用sqlite数据库")
			// 使用sqlite数据库
			gb.DB, _ = gorm.Open(
				sqlite.Open(yamlConfig.Database.Sqlite.Path),
				&gorm.Config{Logger: newLogger})
		} else if yamlConfig.Database.Mysql.Enable {
			// 使用mysql数据库
			mysqlConf := yamlConfig.Database.Mysql
			dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=True&loc=Local",
				mysqlConf.Username, mysqlConf.Password, mysqlConf.Address, mysqlConf.Dbname, mysqlConf.Charset)
			log.Println("检测到使用了mysql数据库,链接dsn为:", dsn)
			gb.DB, _ = gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: newLogger})
		} else {
			// 如果没有启用数据库直接退出程序
			log.Printf("请在config.yaml启用一个数据库!")
			os.Exit(0)
		}
		// 自动迁移,如果数据库不存在对应表，则创建
		err := gb.DB.AutoMigrate(
			RevueConfig{}, KeywordsReply{}, RevueApiToken{},
			ProgrammerAlmanac{}, Divination{}, wzxy.UserWzxy{},
			wzxy.TokenWzxy{}, wzxy.MonitorWzxy{},
			wzxy.ClassStudentWzxy{}, PostForm{}, ListenGroup{}, UserSession{})
		if err != nil {
			log.Printf("数据库迁移失败:%s", err)
			return
		}
		// 不存在则自动创建(理论上配置只有一条记录，所以ID只能为1)
		gb.DB.Where(RevueConfig{ID: 1}).Attrs(RevueConfig{ID: 1, ReplyEnable: true, MusicEnable: true}).FirstOrCreate(&RevueConfig{})
		log.Println("数据库连接成功")
	})
	return gb
}

func (lg *ListenGroup) String() string {
	msg := "监听群号:" + lg.Group + "\n"
	msg += "设置时间:" + lg.Date + "\n"
	msg += "设置用户:" + lg.UserId + "\n"
	return msg
}
