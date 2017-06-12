// Package setting
package setting

import (
	"encoding/json"
	"fmt"

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

	if err := setListenMode(viper.GetStringMap("listenmode")); err != nil {
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

type ListenModeConfig struct {
	Mode    string `json:"mode"`
	Address string `json:"address"`
	Port    int    `json:"port"`
	CertKey string `json:"key"`
	Cert    string `json:"cert"`
}

var Database DatabaseConfig
var ListenMode ListenModeConfig

func setDatabaseConfig(config map[string]interface{}) error {
	bs, err := json.Marshal(&config)
	if err != nil {
		return err
	}

	return json.Unmarshal(bs, &Database)
}

func setListenMode(config map[string]interface{}) error {
	bs, err := json.Marshal(&config)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bs, &ListenMode)
	if err != nil {
		return err
	}

	if ListenMode.Mode != "http" && ListenMode.Mode != "https" && ListenMode.Mode != "unix" {
		return fmt.Errorf("Invalid listen mode: %s", ListenMode.Mode)
	}

	return nil
}
