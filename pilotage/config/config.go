package config

import (
	"encoding/json"

	"github.com/spf13/viper"
)

type WebHookConfig struct {
	Host         string `json:"host"`
	Namespace    string `json:"namespace"`
	Repository   string `json:"repository"`
	Binary       string `json:"binary"`
	Tag          string `json:"tag"`
	FlowFilePath string `json:"flowFilePath"`
}

var WebHook WebHookConfig

func InitConfig(cfgFile string) error {
	viper.SetConfigFile(cfgFile)
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	hookMap := viper.GetStringMap("hook")
	bs, err := json.Marshal(hookMap)
	if err != nil {
		return err
	}

	return json.Unmarshal(bs, &WebHook)
}
