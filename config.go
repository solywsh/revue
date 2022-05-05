package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/ghodss/yaml"
	"io/ioutil"
)

var yamlConfig Config // 定义全局变量用于配置

// Config 配置相关
type Config struct {
	ListenPort            string                `yaml:"listenPort"`            // 监听端口
	AdminUOH              string                `yaml:"adminUserOrderHeader"`  // 管理员命令头 adminUserOrderHeader
	ListenGroup           []string              `yaml:"listenGroup"`           // 监听群列表
	ForwardAuthentication ForwardAuthentication `yaml:"forwardAuthentication"` // 正向鉴权 forward authentication
	ReverseAuthentication ReverseAuthentication `yaml:"reverseAuthentication"` // 反向鉴权 reverse authentication
	Revue                 Revue                 `yaml:"revue"`                 // revue相关
	UrlHeader             string                `yaml:"urlHeader"`             // url
	SelfId                string                `yaml:"selfId"`                // 机器人的qq
	AdminUser             []string              `yaml:"adminUser"`             // 管理员列表
	Database              Database              `yaml:"Database"`              // 数据库相关
}

// ForwardAuthentication 正向鉴权相关
type ForwardAuthentication struct {
	Enable bool   `yaml:"enable"`
	Token  string `yaml:"token"`
}

// ReverseAuthentication 反向鉴权相关
type ReverseAuthentication struct {
	Enable bool   `yaml:"enable"`
	Secret string `yaml:"secret"`
}

// Revue 相关
type Revue struct {
	Enable bool `yaml:"enable"`
}

// Database 数据库相关
type Database struct {
	Path string `yaml:"path"`
}

//
//  getSHA256
//  @Description: 得到SHA256之后的密钥
//  @param str
//  @return string
//
func getSHA256(str string) string {
	sha256Bytes := sha256.Sum256([]byte(str))
	return hex.EncodeToString(sha256Bytes[:])
}

//  getConf
//  @Description: 读取Yaml配置文件,并转换成conf对象
//  @receiver conf
//  @return *Config
//
func (conf *Config) getConf(yamlPath string) *Config {
	yamlFile, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		fmt.Println(err.Error())
	}
	err = yaml.Unmarshal(yamlFile, conf)
	//err = yaml.UnmarshalStrict(yamlFile, kafkaCluster)
	if err != nil {
		fmt.Println(err.Error())
	}
	return conf
}

//func main() {
//	yamlConfig.getConf()
//	fmt.Printf("%#v", yamlConfig.UrlHeader)
//}
