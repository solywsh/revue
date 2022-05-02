package main

import (
	"fmt"
	"github.com/ghodss/yaml"
	"io/ioutil"
)

var yamlConfig Config
var yamlPath = "./config.yaml"

// Config 配置相关
type Config struct {
	AdminUOH    string                `yaml:"adminUserOrderHeader"`  // 管理员命令头 adminUserOrderHeader
	ListenGroup []string              `yaml:"listenGroup"`           // 监听群列表
	FAuth       ForwardAuthentication `yaml:"forwardAuthentication"` // 正向鉴权 forward authentication
	RAuth       ReverseAuthentication `yaml:"reverseAuthentication"` // 反向鉴权 reverse authentication
	Revue       Revue                 `yaml:"revue"`                 // revue相关
	UrlHeader   string                `yaml:"urlHeader"`             // url
	SelfId      string                `yaml:"selfId"`                // 机器人的qq
	AdminUser   []string              `yaml:"adminUser"`             // 管理员列表
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
	Enable bool   `yaml:"enable"`
	Secret string `yaml:"secret"`
}

//  getConf
//  @Description: 读取Yaml配置文件,并转换成conf对象
//  @receiver conf
//  @return *Config
//
func (conf *Config) getConf() *Config {
	//应该是 绝对地址
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
//}
