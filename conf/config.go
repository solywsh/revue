package conf

import (
	"fmt"
	"github.com/ghodss/yaml"
	"io/ioutil"
)

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

func NewConf(yamlPath string) (conf *Config, err error) {
	yamlFile, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		fmt.Println("文件打开错误,请传入正确的文件路径!", err)
		return conf, err
	}
	fmt.Println(string(yamlFile))
	err = yaml.Unmarshal(yamlFile, &conf)
	//err = yaml.UnmarshalStrict(yamlFile, kafkaCluster)
	if err != nil {
		fmt.Println("文件解析错误,请配置正确的yaml格式!", err)
		return conf, err
	}
	return conf, nil
}
