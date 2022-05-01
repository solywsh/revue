package main

import (
	"encoding/json"
	"github.com/ghodss/yaml"
	"github.com/thedevsaddam/gojsonq"
)

type config struct {
	urlHeader string
}

type yamlDecoder struct {
}

var yamlConfig config
var yamlPath = "./config.yaml"

func readConfig() {
	yamlInfo := gojsonq.New(gojsonq.SetDecoder(&yamlDecoder{})).File(yamlPath)
	yamlConfig.urlHeader = yamlInfo.Reset().Find("host.urlHeader").(string)
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
