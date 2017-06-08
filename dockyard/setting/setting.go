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
	return nil
}

type DatabaseConfig struct {
	Driver   string `json:"driver"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"db"`
}

var DBConfig DatabaseConfig

func setDatabaseConfig(config map[string]interface{}) error {
	bs, err := json.Marshal(&config)
	if err != nil {
		return err
	}

	return json.Unmarshal(bs, &DBConfig)
}
