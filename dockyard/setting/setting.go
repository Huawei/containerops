/*
Copyright 2016 - 2017 Huawei Technologies Co., Ltd. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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

	if err := setStorageConfig(viper.GetStringMap("storage")); err != nil {
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

type StorageConfig struct {
	DockerV2 string `json:"dockerv2" description:"Docker V2 image storage path in the host."`
	BinaryV1 string `json:"binaryv1" description:"Binary V1 file storage path in the host"`
}

var Database DatabaseConfig
var Web WebConfig
var Storage StorageConfig

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

func setStorageConfig(config map[string]interface{}) error {
	bs, err := json.Marshal(&config)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bs, &Storage)
	if err != nil {
		return err
	}

	return nil
}
