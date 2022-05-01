package main

import (
	"encoding/json"
	"github.com/ghodss/yaml"
	"github.com/thedevsaddam/gojsonq"
)

type config struct {
	urlHeader   string
	adminUser   []string
	listenGroup []string
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
