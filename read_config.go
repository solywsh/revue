package main

import (
	"encoding/json"
	"github.com/ghodss/yaml"
	"github.com/thedevsaddam/gojsonq"
)

type access struct {
	enable        bool
	tokenOrSecret string
}

type config struct {
	urlHeader   string
	adminUser   []string
	listenGroup []string
	fAuth       access // 正向鉴权 forward authentication
	rAuth       access // 反向鉴权 reverse authentication
}

type yamlDecoder struct {
}

var yamlConfig config
var yamlPath = "./config.yaml"

func readConfig() {
	yamlInfo := gojsonq.New(gojsonq.SetDecoder(&yamlDecoder{})).File(yamlPath)
	yamlConfig.urlHeader = yamlInfo.Reset().Find("apiHost").(string)
	for _, v := range yamlInfo.Reset().Find("adminUser").([]interface{}) {
		yamlConfig.adminUser = append(yamlConfig.adminUser, v.(string))
	}
	for _, v := range yamlInfo.Reset().Find("listenGroup").([]interface{}) {
		yamlConfig.listenGroup = append(yamlConfig.listenGroup, v.(string))
	}
	// 正向鉴权
	yamlConfig.fAuth.enable = yamlInfo.Reset().Find("forwardAuthentication.enable").(bool)
	yamlConfig.fAuth.tokenOrSecret = yamlInfo.Reset().Find("forwardAuthentication.token").(string)
	// 反向鉴权
	yamlConfig.rAuth.enable = yamlInfo.Reset().Find("reverseAuthentication.enable").(bool)
	yamlConfig.rAuth.tokenOrSecret = yamlInfo.Reset().Find("reverseAuthentication.secret").(string)
}

// Decode 实现gojsonq.Decoder
func (i *yamlDecoder) Decode(data []byte, v interface{}) error {
	bb, err := yaml.YAMLToJSON(data)
	if err != nil {
		return err
	}
	return json.Unmarshal(bb, &v)
}

//func main() {
//	readConfig()
//}
