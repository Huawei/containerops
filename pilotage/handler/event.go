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

	"github.com/Huawei/containerops/pilotage/models"

	"gopkg.in/macaron.v1"
)

//PostEventV1Handler is
func PostEventV1Handler(ctx *macaron.Context) (int, []byte) {
	result, _ := json.Marshal(map[string]string{"message": ""})
	return http.StatusOK, result
}

//GetEventV1Handler is
func GetEventV1Handler(ctx *macaron.Context) (int, []byte) {
	result, _ := json.Marshal(map[string]string{"message": ""})
	return http.StatusOK, result
}

//PutEventV1Handler is
func PutEventV1Handler(ctx *macaron.Context) (int, []byte) {
	result, _ := json.Marshal(map[string]string{"message": ""})
	return http.StatusOK, result
}

//DeleteEventV1Handler is
func DeleteEventV1Handler(ctx *macaron.Context) (int, []byte) {
	result, _ := json.Marshal(map[string]string{"message": ""})
	return http.StatusOK, result
}

// GetEventDefineJsonV1Handler is
func GetEventDefineJsonV1Handler(ctx *macaron.Context) (int, []byte) {
	result, _ := json.Marshal(map[string]string{"message": ""})

	siteName := ctx.Params(":site")
	if siteName == "" {
		result, _ = json.Marshal(map[string]string{"errMsg": "site name can't be empty!"})
		return http.StatusBadRequest, result
	}

	eventName := ctx.Params(":event")
	if eventName == "" {
		result, _ = json.Marshal(map[string]string{"errMsg": "event name can't be empty!"})
		return http.StatusBadRequest, result
	}

	eventDefine := new(models.EventJson)
	err := eventDefine.GetEventJson().Where("site = ?", siteName).Where("type = ?", eventName).Find(&eventDefine).Error
	if err != nil {
		result, _ = json.Marshal(map[string]string{"errMsg": "error when get eventDefine from db: " + err.Error()})
		return http.StatusBadRequest, result
	}

	if eventDefine.ID == 0 {
		result, _ = json.Marshal(map[string]string{"errMsg": "can't get " + siteName + "'s : " + eventName + " define json"})
		return http.StatusBadRequest, result
	}

	outputMap := make(map[string]interface{}, 0)
	json.Unmarshal([]byte(eventDefine.Output), &outputMap)

	result, _ = json.Marshal(map[string]interface{}{"output": outputMap})
	return http.StatusOK, result
}
