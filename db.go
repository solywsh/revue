package main

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
)

var globalDB, _ = gorm.Open(sqlite.Open("./data.db"), &gorm.Config{})

// RevueConfig 根据命令对机器人的一些配置进行动态配置
type RevueConfig struct {
	ID          uint `gorm:"primaryKey;autoIncrement"`
	ReplyEnable bool // 开启回复
	MusicEnable bool // 开启音乐推荐
}

// KeywordsReply 关键词触发消息
type KeywordsReply struct {
	ID       uint   `gorm:"primaryKey;autoIncrement"`
	Mode     uint   // 匹配模式 0:不匹配 1:完全匹配 2:存在匹配(find)
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
	globalDB.AutoMigrate(RevueConfig{}, KeywordsReply{}, RevueApiToken{})
	//var rc RevueConfig
	// 不存在则自动创建(理论上配置只有一条记录，所以ID只能为1)
	globalDB.Where(RevueConfig{ID: 1}).Attrs(RevueConfig{ID: 1, ReplyEnable: true, MusicEnable: true}).FirstOrCreate(&RevueConfig{})
}

func getAllKeywordsReply(res *[]KeywordsReply) {
	//dbPath := yamlConfig.Database.Path
	globalDB.Find(&res)
}

//
//  insertRevueApiToken
//  @Description: 采用查询并创建,查询到对应的信息,则不创建返回false,没有则创建,返回true
//  @param userId 生成的qq号
//  @param permission 权限
//  @return bool 是否创建成功
//	@return string 对应的token
//
func insertRevueApiToken(userId string, permission uint) (bool, string) {
	var rt RevueApiToken
	res := globalDB.Where(RevueApiToken{UserId: userId}).Attrs(RevueApiToken{
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

//
//  getRevueApiToken
//  @Description: 根据qq号得到token,注意此函数必须在qq号存在的业务逻辑下使用
//  @param userId qq号
//  @return string token
//
func getRevueApiToken(userId string) string {
	var rt RevueApiToken
	globalDB.Where(RevueApiToken{UserId: userId}).First(&rt)
	return rt.Token
}

//
//  searchRevueApiToken
//  @Description: 根据传入的qq号和token判断是否具有权限,并返回权限
//  @param userId qq号
//  @param token qq号对应的token
//  @return bool 结果是否为真
//  @return uint 对应的权限
//
func searchRevueApiToken(userId, token string) (bool, uint) {
	var rt RevueApiToken
	res := globalDB.Where(RevueApiToken{UserId: userId, Token: token}).First(&rt)
	if res.RowsAffected >= 1 {
		return true, rt.Permission
	} else {
		return false, rt.Permission
	}
}

//
//  resetRevueApiToken
//  @Description: 重新设置token
//  @param userId qq号
//  @return bool
//  @return string
//
func resetRevueApiToken(userId string) (bool, string) {
	var rt RevueApiToken
	res := globalDB.Where(RevueApiToken{UserId: userId}).First(&rt)
	if res.RowsAffected >= 1 {
		globalDB.Model(&rt).Update("token", uuid.NewV4().String())
		return true, rt.Token
	} else {
		// 不存在
		return false, rt.Token
	}
}

//
//  deleteRevueApiToken
//  @Description: 删除token
//  @param userId
//  @return bool
//  @return string
//
func deleteRevueApiToken(userId string) (bool, string) {
	var rt RevueApiToken
	res := globalDB.Where(RevueApiToken{UserId: userId}).First(&rt)
	if res.RowsAffected >= 1 {
		globalDB.Unscoped().Delete(&rt)
		return true, rt.Token
	} else {
		// 不存在
		return false, rt.Token
	}
}

//
//  getRevueConfig
//  @Description: 获取RevueConfig配置
//  @param config
//  @return bool
//
func getRevueConfig(config *RevueConfig) bool {
	res := globalDB.Where(RevueConfig{ID: 1}).First(&config)
	if res.RowsAffected >= 1 {
		return true
	}
	return false
}

//
//  setRevueConfigMusic
//  @Description: 设置music开启和关闭
//  @param enable
//
func setRevueConfigMusic(enable bool) {
	globalDB.Model(&RevueConfig{}).Where(RevueConfig{ID: 1}).Update("music_enable", enable)
}

//
//  setRevueConfigReply
//  @Description: 设置reply开启和关闭
//  @param enable
//
func setRevueConfigReply(enable bool) {
	globalDB.Model(&RevueConfig{}).Where(RevueConfig{ID: 1}).Update("reply_enable", enable)
}

//func main() {
//	dbInit()
//	setRevueConfigMusic(true)
//	setRevueConfigReply(false)
//	var rc RevueConfig
//	if getRevueConfig(&rc) {
//		fmt.Println(rc.MusicEnable, rc.ReplyEnable)
//	}
//}
