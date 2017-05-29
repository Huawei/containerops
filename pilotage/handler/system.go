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

package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Huawei/containerops/pilotage/module"

	"gopkg.in/macaron.v1"
)

//GetSettingV1Handler is get user's system default setting.
func GetSettingV1Handler(ctx *macaron.Context) (int, []byte) {
	result, _ := json.Marshal(map[string]string{"message": ""})

	namespace := ctx.Params(":namespace")
	if namespace == "" {
		result, _ = json.Marshal(map[string]string{"errMsg": "namespace can't be empty"})
		return http.StatusBadRequest, result
	}

	repository := ctx.Params(":repository")
	if repository == "" {
		result, _ = json.Marshal(map[string]string{"errMsg": "repository can't be empty"})
		return http.StatusBadRequest, result
	}

	setting, err := module.GetUserSetting(namespace, repository)
	if err != nil {
		result, _ = json.Marshal(map[string]string{"errMsg": "error when get user setting"})
		return http.StatusBadRequest, result
	}

	settingMap := make(map[string]interface{})
	settingMap["KUBE_APISERVER_IP"] = ""
	settingMap["KUBE_NODE_IP"] = ""

	if setting.Setting != "" {
		err = json.Unmarshal([]byte(setting.Setting), &settingMap)
		if err != nil {
			result, _ = json.Marshal(map[string]string{"errMsg": "error when get setting info"})
			return http.StatusBadRequest, result
		}
	}

	result, _ = json.Marshal(map[string]interface{}{"setting": settingMap})
	return http.StatusOK, result
}

//PutSettingV1Handler is set user's system default setting.
func PutSettingV1Handler(ctx *macaron.Context) (int, []byte) {
	result, _ := json.Marshal(map[string]string{"message": ""})

	body := new(struct {
		Setting map[string]interface{} `json:"setting"`
	})

	namespace := ctx.Params(":namespace")
	if namespace == "" {
		result, _ = json.Marshal(map[string]string{"errMsg": "namespace can't be empty"})
		return http.StatusBadRequest, result
	}

	repository := ctx.Params(":repository")
	if repository == "" {
		result, _ = json.Marshal(map[string]string{"errMsg": "repository can't be empty"})
		return http.StatusBadRequest, result
	}

	reqBody, err := ctx.Req.Body().Bytes()
	if err != nil {
		result, _ = json.Marshal(map[string]string{"errMsg": "error when get request body:" + err.Error()})
		return http.StatusBadRequest, result
	}

	err = json.Unmarshal(reqBody, &body)
	if err != nil {
		result, _ = json.Marshal(map[string]string{"errMsg": "error when unmarshal request body:" + err.Error()})
		return http.StatusBadRequest, result
	}

	setting, err := module.GetUserSetting(namespace, repository)
	if err != nil {
		result, _ = json.Marshal(map[string]string{"errMsg": "error when get user setting"})
		return http.StatusBadRequest, result
	}

	err = setting.SetUserSetting(namespace, repository, body.Setting)
	if err != nil {
		result, _ = json.Marshal(map[string]string{"errMsg": "error when save user setting"})
		return http.StatusBadRequest, result
	}

	result, _ = json.Marshal(map[string]string{"message": "success"})
	return http.StatusOK, result
}
