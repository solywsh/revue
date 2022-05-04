package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
)

// RevueConfig 根据命令对机器人的一些配置进行动态配置
type RevueConfig struct {
	gorm.Model
	ID          uint `gorm:"primaryKey"`
	ReplyEnable bool // 开启回复
	MusicEnable bool // 开启音乐推荐
}

// KeywordsReply 关键词触发消息
type KeywordsReply struct {
	gorm.Model
	Mode     uint   // 匹配模式 0:不匹配 1:完全匹配 2:存在匹配(find)
	Keywords string // 关键词
	Msg      string // 回复消息
	Userid   string // 设置人的qq
}

//PathExists 判断一个文件或文件夹是否存在
//输入文件路径，根据返回的bool值来判断文件或文件夹是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

//
//  dbInit
//  @Description:
// 			数据库初始化
//			如果不纯在数据库对应的表则自动迁移
//			如果RevueConfig表不纯在信息则创建
//  @param dbPath
//
func dbInit() {
	// 自动迁移,如果数据库不存在对应表，则创建
	globalDB.AutoMigrate(RevueConfig{}, KeywordsReply{})
	//var rc RevueConfig
	// 不存在则自动创建(理论上配置只有一条记录，所以ID只能为1)
	globalDB.Where(RevueConfig{ID: 1}).Attrs(RevueConfig{ID: 1}).FirstOrCreate(&RevueConfig{})
}

func getAllKeywordsReply(res *[]KeywordsReply) {
	//dbPath := yamlConfig.Database.Path
	globalDB.Find(&res)
}

var globalDB, _ = gorm.Open(sqlite.Open("./data.db"), &gorm.Config{})

//func main() {
//	dbInit()
//
//}
