/*
Copyright 2014 - 2017 Huawei Technologies Co., Ltd. All rights reserved.

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

// GetStageHistoryInfoV1Handler is
func GetStageHistoryInfoV1Handler(ctx *macaron.Context) (int, []byte) {
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

	workflowName := ctx.Params(":workflow")
	if workflowName == "" {
		result, _ = json.Marshal(map[string]string{"errMsg": "workflow can't be empty"})
		return http.StatusBadRequest, result
	}

	sequence := ctx.Params(":sequence")
	sequenceInt, err := strconv.ParseInt(sequence, 10, 64)
	if err != nil || sequenceInt == int64(0) {
		result, _ = json.Marshal(map[string]string{"errMsg": "sequence error"})
		return http.StatusBadRequest, result
	}

	stageName := ctx.Params(":stage")
	if stageName == "" {
		result, _ = json.Marshal(map[string]string{"errMsg": "stage can't be empty"})
		return http.StatusBadRequest, result
	}

	stageLogInfo, err := module.GetStageLogByName(namespace, repository, workflowName, sequenceInt, stageName)
	if err != nil {
		result, _ = json.Marshal(map[string]string{"errMsg": err.Error()})
		return http.StatusBadRequest, result
	}

	resultMap, err := stageLogInfo.GetStageLogDefine()
	if err != nil {
		result, _ = json.Marshal(map[string]string{"errMsg": err.Error()})
		return http.StatusBadRequest, result
	}

	result, _ = json.Marshal(map[string]interface{}{"result": resultMap})

	return http.StatusOK, result
}
