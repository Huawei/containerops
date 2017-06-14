// Package setting
package setting

import (
	"encoding/json"

	"github.com/spf13/viper"
)

func SetConfig(configPath string) error {
	viper.SetConfigType("toml")
	viper.SetConfigFile(configPath)
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	if err := setDatabaseConfig(viper.GetStringMap("database")); err != nil {
		return err
	}

	if err := setWebConfig(viper.GetStringMap("web")); err != nil {
		return err
	}

	return nil
}

type DatabaseConfig struct {
	Driver   string `json:"driver"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"db"`
}

type WebConfig struct {
	Mode    string `json:"mode" description:"Listen mode, 'https' or 'unix'"`
	Address string `json:"address" description:"The host address when mode is 'https', or socket file path when mode is 'unix'"`
	Port    int    `json:"port"`
	Key     string `json:"key"`
	Cert    string `json:"cert"`
}

var Database DatabaseConfig
var Web WebConfig

func setDatabaseConfig(config map[string]interface{}) error {
	bs, err := json.Marshal(&config)
	if err != nil {
		return err
	}

	return json.Unmarshal(bs, &Database)
}

func setWebConfig(config map[string]interface{}) error {
	bs, err := json.Marshal(&config)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bs, &Web)
	if err != nil {
		return err
	}

	return nil
}
