/*
Copyright 2014 Huawei Technologies Co., Ltd. All rights reserved.

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
	"strconv"
	"strings"

	"gopkg.in/macaron.v1"

	"github.com/Huawei/containerops/pilotage/module"
)

//PostStageV1Handler is
func PostStageV1Handler(ctx *macaron.Context) (int, []byte) {
	result, _ := json.Marshal(map[string]string{"message": ""})
	return http.StatusOK, result
}

//GetStageV1Handler is
func GetStageV1Handler(ctx *macaron.Context) (int, []byte) {
	result, _ := json.Marshal(map[string]string{"message": ""})
	return http.StatusOK, result
}

//PutStageV1Handler is
func PutStageV1Handler(ctx *macaron.Context) (int, []byte) {
	result, _ := json.Marshal(map[string]string{"message": ""})
	return http.StatusOK, result
}

//DeleteStageV1Handler is
func DeleteStageV1Handler(ctx *macaron.Context) (int, []byte) {
	result, _ := json.Marshal(map[string]string{"message": ""})
	return http.StatusOK, result
}

func GetStageHistoryInfoV1Handler(ctx *macaron.Context) (int, []byte) {
	result, _ := json.Marshal(map[string]string{"message": ""})

	stageLogIdStr := ctx.Query("stageLogId")

	if stageLogIdStr == "" {
		result, _ = json.Marshal(map[string]string{"errMsg": "request stage log id can't be empty"})
		return http.StatusBadRequest, result
	}

	stageLogIdStr = strings.TrimPrefix(stageLogIdStr, "s-")
	stageLogId, err := strconv.ParseInt(stageLogIdStr, 10, 64)

	if err != nil {
		result, _ = json.Marshal(map[string]string{"errMsg": "stage log id is illegal"})
		return http.StatusBadRequest, result
	}

	resultMap, err := module.GetStageHistoryInfo(stageLogId)
	if err != nil {
		result, _ = json.Marshal(map[string]string{"errMsg": err.Error()})
		return http.StatusBadRequest, result
	}

	result, _ = json.Marshal(map[string]interface{}{"result": resultMap})

	return http.StatusOK, result
}
