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
	"time"
)

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

	actionName := ctx.Params(":action")
	if actionName == "" {
		result, _ = json.Marshal(map[string]string{"errMsg": "action can't be empty"})
		return http.StatusBadRequest, result
	}

	actionLogInfo, err := module.GetActionLogByName(namespace, repository, workflowName, sequenceInt, stageName, actionName)
	if err != nil {
		result, _ = json.Marshal(map[string]string{"errMsg": err.Error()})
		return http.StatusBadRequest, result
	}

	resultMap, err := actionLogInfo.GetActionHistoryInfo()
	if err != nil {
		result, _ = json.Marshal(map[string]string{"errMsg": err.Error()})
		return http.StatusBadRequest, result
	}

	result, _ = json.Marshal(map[string]interface{}{"result": resultMap})

	return http.StatusOK, result
}

func GetActionConsoleLogV1Handler(ctx *macaron.Context) (int, []byte) {
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

	actionName := ctx.Params(":action")
	if actionName == "" {
		result, _ = json.Marshal(map[string]string{"errMsg": "action can't be empty"})
		return http.StatusBadRequest, result
	}

	key := ctx.Query("key")
	size := ctx.QueryInt64("size")
	if size == 0 {
		size = 10
	}

	actionLogInfo, err := module.GetActionLogByName(namespace, repository, workflowName, sequenceInt, stageName, actionName)
	if err != nil {
		result, _ = json.Marshal(map[string]string{"errMsg": err.Error()})
		return http.StatusBadRequest, result
	}

	logResult, err := actionLogInfo.GetActionConsoleLog(key, size)
	if err != nil {
		result, _ = json.Marshal(map[string]string{"errMsg": err.Error()})
		return http.StatusBadRequest, result
	}

	result, _ = json.Marshal(logResult)

	return http.StatusOK, result
}

func CreateEvent(ctx *macaron.Context) (httpStatus int, result []byte) {
	var resp CommonResp
	body, err := ctx.Req.Body().Bytes()
	if err != nil {
		httpStatus = http.StatusBadRequest
		resp.OK = false
		resp.ErrorCode = EventError + EventReqBodyError
		resp.Message = "Get requrest body error: " + err.Error()

		result, err = json.Marshal(resp)
		if err != nil {
			log.Errorln("Create event marshal data error: " + err.Error())
		}
		return
	}

	var req EventReq
	err = json.Unmarshal(body, &req)
	if err != nil {
		log.Errorln("CreateEvent unmarshal data error: ", err.Error())
		httpStatus = http.StatusBadRequest
		resp.OK = false
		resp.ErrorCode = EventError + EventUnmarshalError
		resp.Message = "unmarshal data error: " + err.Error()

		result, err = json.Marshal(resp)
		if err != nil {
			log.Errorln("Create event marshal data error: " + err.Error())
		}
		return
	}

	//err = json.Unmarshal(body, &reqBody)
	//if err != nil {
	//	log.Error("[action's PostActionEventV1Handler]:error when unmarshal reqBody:", string(body), " ===>error is:", err.Error())
	//	result, _ := json.Marshal(map[string]string{"message": "illegal request body,want a json obj,got:" + string(body)})
	//	return http.StatusBadRequest, result
	//}

	//eventKey, ok := reqBody["EVENT"].(string)
	//if !ok {
	//	log.Error("[action's PostActionEventV1Handler]:error when get eventKey from request, want a string, got:", reqBody["EVENT"])
	//	result, _ := json.Marshal(map[string]string{"message": "eventKey is not a string"})
	//	return http.StatusBadRequest, result
	//}

	//eventIdF, ok := reqBody["EVENT_ID"].(float64)
	//if !ok {
	//	log.Error("[action's PostActionEventV1Handler]:error when get event_ID from request, want a number, got:", reqBody["EVENT_ID"])
	//	result, _ := json.Marshal(map[string]string{"message": "eventId is not a number"})
	//	return http.StatusBadRequest, result
	//}

	//eventId := int64(eventIdF)
	//runId, ok := reqBody["RUN_ID"].(string)
	//if !ok {
	//	log.Error("[action's PostActionEventV1Handler]:error when get runID from request, want a string, got:", reqBody["RUN_ID"])
	//	result, _ := json.Marshal(map[string]string{"message": "runId is not a string"})
	//	return http.StatusBadRequest, result
	//}

	slice := strings.Split(req.RunID, "-")
	if len(slice) != 3 {
		httpStatus = http.StatusBadRequest
		resp.OK = false
		resp.ErrorCode = EventError + EventIllegalDataError
		resp.Message = "Illegal run_id: " + req.RunID

		result, err = json.Marshal(resp)
		if err != nil {
			log.Errorln("Create event marshal data error: " + err.Error())
		}
		return
	}

	actionLogId, err := strconv.ParseInt(slice[2], 10, 64)
	if err != nil {
		httpStatus = http.StatusBadRequest
		resp.OK = false
		resp.ErrorCode = EventError + EventIllegalDataError
		resp.Message = "Parse run_id error: " + err.Error()

		result, err = json.Marshal(resp)
		if err != nil {
			log.Errorln("Create event marshal data error: " + err.Error())
		}
		return
	}

	actionLog, err := module.GetActionLog(actionLogId)
	if err != nil {
		httpStatus = http.StatusBadRequest
		resp.OK = false
		resp.ErrorCode = EventError + EventGetActionError
		resp.Message = "get action log by id error: " + err.Error()

		result, err = json.Marshal(resp)
		if err != nil {
			log.Errorln("Create event marshal data error: " + err.Error())
		}
		return
	}

	go func(){
		value, ok := cache.Get(actionLogId)
		if !ok {
			log.Warnf("Component message channel key %d not exist\n", actionLogId)
			return
		}

		c, ok := value.(chan DebugEvent)
		if !ok {
			log.Errorf("Can't convert type %T to message channel\n", value)
			return
		}

		c <- DebugEvent{
			Type: req.EventType,
			Content: time.Now().Format("2006-01-02 15:04:05") + " -> " + string(body),
		}
	}()

	err = actionLog.RecordEvent(req.EventID, req.EventType, req.Info, body, ctx.Req.Header)
	if err != nil {
		httpStatus = http.StatusBadRequest
		resp.OK = false
		resp.ErrorCode = EventError + EventGetActionError
		resp.Message = "get action log by id error: " + err.Error()

		result, err = json.Marshal(resp)
		if err != nil {
			log.Errorln("Create event marshal data error: " + err.Error())
		}
		return
	}

	httpStatus = http.StatusOK
	resp.OK = true
	resp.Message = "Event received"

	result, err = json.Marshal(resp)
	if err != nil {
		log.Errorln("Create event marshal data error: " + err.Error())
	}
	return
}

func PostActionSetVarV1Handler(ctx *macaron.Context) (int, []byte) {
	result, _ := json.Marshal(map[string]string{"message": ""})

	bodyByte, _ := ctx.Req.Body().Bytes()

	reqBody := make(map[string]interface{})
	err := json.Unmarshal(bodyByte, &reqBody)
	if err != nil {
		log.Error("[action's PostActionSetVarV1Handler]:error when unmarshal reqBody:", string(bodyByte), " ===>error is:", err.Error())
		result, _ := json.Marshal(map[string]string{"message": "illegal request body,want a json obj,got:" + string(bodyByte)})
		return http.StatusBadRequest, result
	}

	runId, ok := reqBody["RUN_ID"].(string)
	if !ok {
		log.Error("[action's PostActionSetVarV1Handler]:error when get runID from request, want a string, got:", reqBody["RUN_ID"])
		result, _ := json.Marshal(map[string]string{"message": "runId is not a string"})
		return http.StatusBadRequest, result
	}

	if len(strings.Split(runId, "-")) < 3 {
		log.Error("[action's PostActionSetVarV1Handler]:runID illegal,want XX-XX-XX, got:", runId)
		result, _ := json.Marshal(map[string]string{"message": "illegal runID"})
		return http.StatusBadRequest, result
	}

	actionLogId, err := strconv.ParseInt(strings.Split(runId, "-")[2], 10, 64)
	if err != nil {
		log.Error("[action's PostActionSetVarV1Handler]:error when get actionLogId from runID, want number, got:", runId)
		result, _ := json.Marshal(map[string]string{"message": "illegal actionLogId id"})
		return http.StatusBadRequest, result
	}

	varMap, ok := reqBody["varMap"].(map[string]interface{})
	if !ok {
		log.Error("[action's PostActionSetVarV1Handler]:error when get varMap from request, want a obj, got:", reqBody["varMap"])
		result, _ := json.Marshal(map[string]string{"message": "runId is not a string"})
		return http.StatusBadRequest, result
	}

	varKey, ok := varMap["KEY"].(string)
	if !ok {
		log.Error("[action's PostActionSetVarV1Handler]:error when get varKey from request's varMap, want a string, got:", varMap["KEY"])
		result, _ := json.Marshal(map[string]string{"message": "varKey is not a string"})
		return http.StatusBadRequest, result
	}

	varValue, ok := varMap["VALUE"].(string)
	if !ok {
		log.Error("[action's PostActionSetVarV1Handler]:error when get varValue from request's varMap, want a string, got:", varMap["VALUE"])
		result, _ := json.Marshal(map[string]string{"message": "varValue is not a string"})
		return http.StatusBadRequest, result
	}

	actionLog, err := module.GetActionLog(actionLogId)
	if err != nil {
		log.Error("[action's PostActionEventV1Handler]:error when get action's info:", err.Error())
		result, _ := json.Marshal(map[string]string{"message": "error when get target action"})
		return http.StatusBadRequest, result
	}

	err = actionLog.ChangeWorkflowRuntimeVar(runId, varKey, varValue)
	if err != nil {
		log.Error("[action's PostActionSetVarV1Handler]:error when change action's var:", err.Error())
		result, _ := json.Marshal(map[string]string{"message": "error when change action's var"})
		return http.StatusBadRequest, result
	}

	result, _ = json.Marshal(map[string]string{"message": "ok"})
	return http.StatusOK, result
}

func PostActionLinkStartV1Handler(ctx *macaron.Context) (int, []byte) {
	result, _ := json.Marshal(map[string]string{"message": "ok"})

	bodyByte, _ := ctx.Req.Body().Bytes()

	reqBody := make(map[string]interface{})
	err := json.Unmarshal(bodyByte, &reqBody)
	if err != nil {
		log.Error("[action's PostActionEventV1Handler]:error when unmarshal reqBody:", string(bodyByte), " ===>error is:", err.Error())
		result, _ := json.Marshal(map[string]string{"message": "illegal request body,want a json obj,got:" + string(bodyByte)})
		return http.StatusBadRequest, result
	}

	runId, ok := reqBody["RUN_ID"].(string)
	if !ok {
		log.Error("[action's PostActionEventV1Handler]:error when get runID from request, want a string, got:", reqBody["RUN_ID"])
		result, _ := json.Marshal(map[string]string{"message": "runId is not a string"})
		return http.StatusBadRequest, result
	}

	if len(strings.Split(runId, "-")) < 3 {
		log.Error("[action's PostActionEventV1Handler]:runID illegal,want XX-XX-XX, got:", runId)
		result, _ := json.Marshal(map[string]string{"message": "illegal runID"})
		return http.StatusBadRequest, result
	}

	actionLogId, err := strconv.ParseInt(strings.Split(runId, "-")[2], 10, 64)
	if err != nil {
		log.Error("[action's PostActionEventV1Handler]:error when get actionLogId from runID, want number, got:", runId)
		result, _ := json.Marshal(map[string]string{"message": "illegal actionLogId id"})
		return http.StatusBadRequest, result
	}

	linkInfoMap, ok := reqBody["linkInfoMap"].(map[string]interface{})
	if !ok {
		log.Error("[action's PostActionSetVarV1Handler]:error when get linkInfoMap from request, want a obj, got:", reqBody["linkInfoMap"])
		result, _ := json.Marshal(map[string]string{"message": "linkInfoMap is illegal"})
		return http.StatusBadRequest, result
	}

	token, ok := linkInfoMap["token"].(string)
	if !ok {
		log.Error("[action's PostActionSetVarV1Handler]:error when get token from request's linkInfoMap, want a string, got:", linkInfoMap["token"])
		result, _ := json.Marshal(map[string]string{"message": "token is not a string"})
		return http.StatusBadRequest, result
	}

	workflowName, ok := linkInfoMap["workflowName"].(string)
	if !ok {
		log.Error("[action's PostActionSetVarV1Handler]:error when get workflowName from request's linkInfoMap, want a string, got:", linkInfoMap["workflowName"])
		result, _ := json.Marshal(map[string]string{"message": "workflowName is not a string"})
		return http.StatusBadRequest, result
	}

	workflowVersion, ok := linkInfoMap["workflowVersion"].(string)
	if !ok {
		log.Error("[action's PostActionSetVarV1Handler]:error when get workflowVersion from request's linkInfoMap, want a string, got:", linkInfoMap["workflowVersion"])
		result, _ := json.Marshal(map[string]string{"message": "workflowVersion is not a string"})
		return http.StatusBadRequest, result
	}

	eventName, ok := linkInfoMap["eventName"].(string)
	if !ok {
		log.Error("[action's PostActionSetVarV1Handler]:error when get eventName from request's linkInfoMap, want a string, got:", linkInfoMap["eventName"])
		result, _ := json.Marshal(map[string]string{"message": "eventName is not a string"})
		return http.StatusBadRequest, result
	}

	eventType, ok := linkInfoMap["eventType"].(string)
	if !ok {
		log.Error("[action's PostActionSetVarV1Handler]:error when get eventType from request's linkInfoMap, want a string, got:", linkInfoMap["eventType"])
		result, _ := json.Marshal(map[string]string{"message": "eventType is not a string"})
		return http.StatusBadRequest, result
	}

	startJsonStr, ok := linkInfoMap["startJson"].(string)
	if !ok {
		log.Error("[action's PostActionSetVarV1Handler]:error when get startJson from request's linkInfoMap, want a string, got:", linkInfoMap["startJson"])
		result, _ := json.Marshal(map[string]string{"message": "startJson is not a string"})
		return http.StatusBadRequest, result
	}

	startJson := make(map[string]interface{})
	err = json.Unmarshal([]byte(startJsonStr), &startJson)

	actionLog, err := module.GetActionLog(actionLogId)
	if err != nil {
		log.Error("[action's PostActionEventV1Handler]:error when get action's info:", err.Error())
		result, _ := json.Marshal(map[string]string{"message": "error when get target action"})
		return http.StatusBadRequest, result
	}

	err = actionLog.LinkStartWorkflow(runId, token, workflowName, workflowVersion, eventName, eventType, startJson)
	if err != nil {
		log.Error("[action's PostActionEventV1Handler]:error when record action's event:", err.Error())
		result, _ := json.Marshal(map[string]string{"message": "error when record action's event"})
		return http.StatusBadRequest, result
	}

	result, _ = json.Marshal(map[string]string{"message": "ok"})
	return http.StatusOK, result
}
