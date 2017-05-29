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

package module

import (
	"encoding/json"
	"errors"

	"github.com/Huawei/containerops/pilotage/models"

	log "github.com/Sirupsen/logrus"
)

// UserSetting is
type UserSetting struct {
	*models.UserSetting
}

// GetUserSetting is
func GetUserSetting(namespace, repository string) (*UserSetting, error) {
	result := new(UserSetting)
	setting := new(models.UserSetting)

	err := setting.GetUserSetting().Where("namespace = ?", namespace).Where("repository = ?", repository).First(setting).Error
	if err != nil && err.Error() != "record not found" {
		log.Error("[System's GetUserSetting]:errro when get user setting from db:", err.Error())
		return nil, errors.New("error when get user setting from db")
	}

	result.UserSetting = setting
	return result, nil
}

// SetUserSetting is
func (setting *UserSetting) SetUserSetting(namespace, repository string, setMap map[string]interface{}) error {
	userSetMap := make(map[string]interface{})
	for key, value := range setMap {
		if key == "KUBE_APISERVER_IP" || key == "KUBE_NODE_IP" {
			userSetMap[key] = value
		}
	}

	settingBytes, _ := json.Marshal(userSetMap)

	setting.Namespace = namespace
	setting.Repository = repository
	setting.Setting = string(settingBytes)

	err := setting.GetUserSetting().Save(setting).Error
	if err != nil {
		log.Error("[System's SetUserSetting]:error when save user setting:", err.Error())
		return errors.New("error when save user setting")
	}

	return nil
}
