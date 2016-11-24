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

	"github.com/Huawei/containerops/pilotage/module"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/macaron.v1"
)

//PostActionV1Handler is
func PostActionV1Handler(ctx *macaron.Context) (int, []byte) {
	result, _ := json.Marshal(map[string]string{"message": ""})
	return http.StatusOK, result
}

//GetActionV1Handler is
func GetActionV1Handler(ctx *macaron.Context) (int, []byte) {
	result, _ := json.Marshal(map[string]string{"message": ""})
	return http.StatusOK, result
}

func GetActionHistoryInfoV1Handler(ctx *macaron.Context) (int, []byte) {
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

	pipelineName := ctx.Params(":workflow")
	if pipelineName == "" {
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

	actionName := ctx.Params(":action")
	if actionName == "" {
		result, _ = json.Marshal(map[string]string{"errMsg": "action can't be empty"})
		return http.StatusBadRequest, result
	}

	actionLogInfo, err := module.GetActionLogByName(namespace, repository, pipelineName, sequenceInt, stageName, actionName)
	if err != nil {
		result, _ = json.Marshal(map[string]string{"errMsg": err.Error()})
		return http.StatusBadRequest, result
	}

	// resultMap, err := module.GetActionHistoryInfo(actionLogId)
	resultMap, err := actionLogInfo.GetActionHistoryInfo()
	if err != nil {
		result, _ = json.Marshal(map[string]string{"errMsg": err.Error()})
		return http.StatusBadRequest, result
	}

	result, _ = json.Marshal(map[string]interface{}{"result": resultMap})

	return http.StatusOK, result
}

//PutActionV1Handler is
func PutActionV1Handler(ctx *macaron.Context) (int, []byte) {
	result, _ := json.Marshal(map[string]string{"message": ""})
	return http.StatusOK, result
}

// PutActionEventV1Handler is all action callback handler
func PutActionEventV1Handler(ctx *macaron.Context) (int, []byte) {
	result, _ := json.Marshal(map[string]string{"message": "ok"})

	bodyByte, _ := ctx.Req.Body().Bytes()

	reqBody := make(map[string]interface{})
	err := json.Unmarshal(bodyByte, &reqBody)
	if err != nil {
		log.Error("[action's PutActionEventV1Handler]:error when unmarshal reqBody:", string(bodyByte), " ===>error is:", err.Error())
		result, _ := json.Marshal(map[string]string{"message": "illegal request body,want a json obj,got:" + string(bodyByte)})
		return http.StatusBadRequest, result
	}

	eventKey, ok := reqBody["EVENT"].(string)
	if !ok {
		log.Error("[action's PutActionEventV1Handler]:error when get eventKey from request, want a string, got:", reqBody["EVENT"])
		result, _ := json.Marshal(map[string]string{"message": "eventKey is not a string"})
		return http.StatusBadRequest, result
	}

	eventIdF, ok := reqBody["EVENTID"].(float64)
	if !ok {
		log.Error("[action's PutActionEventV1Handler]:error when get eventID from request, want a number, got:", reqBody["EVENTID"])
		result, _ := json.Marshal(map[string]string{"message": "eventId is not a number"})
		return http.StatusBadRequest, result
	}

	eventId := int64(eventIdF)
	runId, ok := reqBody["RUN_ID"].(string)
	if !ok {
		log.Error("[action's PutActionEventV1Handler]:error when get runID from request, want a string, got:", reqBody["RUN_ID"])
		result, _ := json.Marshal(map[string]string{"message": "runId is not a string"})
		return http.StatusBadRequest, result
	}

	if len(strings.Split(runId, "-")) < 3 {
		log.Error("[action's PutActionEventV1Handler]:runID illegal,want XX-XX-XX, got:", runId)
		result, _ := json.Marshal(map[string]string{"message": "illegal runID"})
		return http.StatusBadRequest, result
	}

	actionLogId, err := strconv.ParseInt(strings.Split(runId, "-")[2], 10, 64)
	if err != nil {
		log.Error("[action's PutActionEventV1Handler]:error when get actionLogId from runID, want number, got:", runId)
		result, _ := json.Marshal(map[string]string{"message": "illegal actionLogId id"})
		return http.StatusBadRequest, result
	}

	actionLog, err := module.GetActionLog(actionLogId)
	if err != nil {
		log.Error("[action's PutActionEventV1Handler]:error when get action's info:", err.Error())
		result, _ := json.Marshal(map[string]string{"message": "error when get target action"})
		return http.StatusBadRequest, result
	}

	err = actionLog.RecordEvent(eventId, eventKey, reqBody, ctx.Req.Header)
	if err != nil {
		log.Error("[action's PutActionEventV1Handler]:error when record action's event:", err.Error())
		result, _ := json.Marshal(map[string]string{"message": "error when record action's event"})
		return http.StatusBadRequest, result
	}

	result, _ = json.Marshal(map[string]string{"message": "ok"})
	return http.StatusOK, result
}

// PutActionRegisterV1Handler is all action register here
func PutActionRegisterV1Handler(ctx *macaron.Context) (int, []byte) {
	result, _ := json.Marshal(map[string]string{"message": "ok"})

	bodyByte, err := ctx.Req.Body().Bytes()
	if err != nil {
		log.Error("[action's PutActionRegisterV1Handler]:error when get request body:", err.Error())
		result, _ := json.Marshal(map[string]string{"message": "error when getrequest body:" + err.Error()})
		return http.StatusBadRequest, result
	}

	if string(bodyByte) == "" {
		log.Error("[action's PutActionRegisterV1Handler]:got an empty reqBody")
		result, _ := json.Marshal(map[string]string{"message": "illegal request body: empty body"})
		return http.StatusBadRequest, result
	}

	reqBody := make(map[string]interface{})
	err = json.Unmarshal(bodyByte, &reqBody)
	if err != nil {
		log.Error("[action's PutActionRegisterV1Handler]:error when unmarshal reqBody:", string(bodyByte), " ===>error is:", err.Error())
		result, _ := json.Marshal(map[string]string{"message": "illegal request body: want a json obj,got:" + string(bodyByte)})
		return http.StatusBadRequest, result
	}

	runId, ok := reqBody["RUN_ID"].(string)
	if !ok {
		result, _ := json.Marshal(map[string]string{"message": "illegal runId, runId is not a string"})
		return http.StatusBadRequest, result
	}

	receiveUrl, ok := reqBody["RECEIVE_URL"].(string)
	if !ok {
		result, _ := json.Marshal(map[string]string{"message": "illegal receiveUrl, receiveUrl is not a string"})
		return http.StatusBadRequest, result
	}

	if len(strings.Split(runId, "-")) != 3 {
		// illegal runId return
		result, _ := json.Marshal(map[string]string{"message": "illegal id:" + runId})
		return http.StatusBadRequest, result
	}

	actionLogId, err := strconv.ParseInt(strings.Split(runId, "-")[2], 10, 64)
	if err != nil {
		log.Error("[action's PutActionEventV1Handler]:error when get actionLogId from runID, want number, got:", runId)
		result, _ := json.Marshal(map[string]string{"message": "illegal actionLogId id"})
		return http.StatusBadRequest, result
	}

	actionLog, err := module.GetActionLog(actionLogId)
	if err != nil {
		log.Error("[action's PutActionEventV1Handler]:error when get action's info:", err.Error())
		result, _ := json.Marshal(map[string]string{"message": "error when get target action"})
		return http.StatusBadRequest, result
	}

	go actionLog.SendDataToAction(receiveUrl)

	return http.StatusOK, result
}

//PutStartActionV1Handler is
func PutStartActionV1Handler(ctx *macaron.Context) (int, []byte) {
	result, _ := json.Marshal(map[string]string{"message": ""})
	return http.StatusOK, result
}

//PutExecuteActionV1Handler is
func PutExecuteActionV1Handler(ctx *macaron.Context) (int, []byte) {
	result, _ := json.Marshal(map[string]string{"message": ""})
	return http.StatusOK, result
}

//PutStatusActionV1Handler is
func PutStatusActionV1Handler(ctx *macaron.Context) (int, []byte) {
	result, _ := json.Marshal(map[string]string{"message": ""})
	return http.StatusOK, result
}

//PutResultActionV1Handler is
func PutResultActionV1Handler(ctx *macaron.Context) (int, []byte) {
	result, _ := json.Marshal(map[string]string{"message": ""})
	return http.StatusOK, result
}

//PutDeleteActionV1Handle is
func PutDeleteActionV1Handle(ctx *macaron.Context) (int, []byte) {
	result, _ := json.Marshal(map[string]string{"message": ""})
	return http.StatusOK, result
}

//DeleteActionV1Handler is
func DeleteActionV1Handler(ctx *macaron.Context) (int, []byte) {
	result, _ := json.Marshal(map[string]string{"message": ""})
	return http.StatusOK, result
}
