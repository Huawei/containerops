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
	"time"

	"github.com/Huawei/containerops/pilotage/models"
	"github.com/Huawei/containerops/pilotage/module"

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

//PutActionV1Handler is
func PutActionV1Handler(ctx *macaron.Context) (int, []byte) {
	result, _ := json.Marshal(map[string]string{"message": ""})
	return http.StatusOK, result
}

//DeleteActionV1Handler is
func DeleteActionV1Handler(ctx *macaron.Context) (int, []byte) {
	result, _ := json.Marshal(map[string]string{"message": ""})
	return http.StatusOK, result
}

// PutActionEventV1Handler is all action callback handler
func PutActionEventV1Handler(ctx *macaron.Context) (int, []byte) {
	result, _ := json.Marshal(map[string]string{"message": "ok"})

	// eventKey := ctx.Query("event")
	// eventId := ctx.QueryInt64("eventId")
	// runId := ctx.Query("runId")
	bodyByte, _ := ctx.Req.Body().Bytes()

	reqBody := make(map[string]interface{})
	json.Unmarshal(bodyByte, &reqBody)

	eventKey, ok := reqBody["EVENT"].(string)
	if !ok {
		result, _ := json.Marshal(map[string]string{"message": "illegal eventKey, is not a string"})
		return http.StatusOK, result
	}

	eventIdF, ok := reqBody["EVENTID"].(float64)
	if !ok {
		result, _ := json.Marshal(map[string]string{"message": "illegal eventId, eventId is not a number"})
		return http.StatusOK, result
	}
	eventId := int64(eventIdF)

	runId, ok := reqBody["RUN_ID"].(string)
	if !ok {
		result, _ := json.Marshal(map[string]string{"message": "illegal runId, runId is not a string"})
		return http.StatusOK, result
	}

	if len(strings.Split(runId, ",")) != 5 {
		// illegal runId return
		result, _ := json.Marshal(map[string]string{"message": "illegal id"})
		return http.StatusOK, result
	}

	pipelineId, err := strconv.ParseInt(strings.Split(runId, ",")[0], 10, 64)
	if err != nil {
		result, _ := json.Marshal(map[string]string{"message": "illegal pipeline id"})
		return http.StatusOK, result
	}

	stageId, err := strconv.ParseInt(strings.Split(runId, ",")[1], 10, 64)
	if err != nil {
		result, _ := json.Marshal(map[string]string{"message": "illegal stageId id"})
		return http.StatusOK, result
	}

	actionId, err := strconv.ParseInt(strings.Split(runId, ",")[2], 10, 64)
	if err != nil {
		result, _ := json.Marshal(map[string]string{"message": "illegal actionId id"})
		return http.StatusOK, result
	}

	pipelineSequence, err := strconv.ParseInt(strings.Split(runId, ",")[3], 10, 64)
	if err != nil {
		result, _ := json.Marshal(map[string]string{"message": "illegal pipelineSequence id"})
		return http.StatusOK, result
	}

	componentId, err := strconv.ParseInt(strings.Split(runId, ",")[4], 10, 64)
	if err != nil {
		result, _ := json.Marshal(map[string]string{"message": "illegal componentId id"})
		return http.StatusOK, result
	}

	pipelineInfo := new(models.PipelineLog)
	pipelineInfo.GetPipelineLog().Where("id = ?", pipelineId).First(pipelineInfo)

	stageInfo := new(models.StageLog)
	stageInfo.GetStageLog().Where("id = ?", stageId).First(stageInfo)

	actionInfo := new(models.ActionLog)
	actionInfo.GetActionLog().Where("id = ?", actionId).First(actionInfo)

	platformSetting, err := module.GetActionPlatformInfo(*actionInfo)
	if err != nil {
		// do an extra log for outcome to record action result
		componentInitErrOutcome := new(models.Outcome)

		componentInitErrOutcome.Pipeline = pipelineId
		componentInitErrOutcome.RealPipeline = pipelineInfo.FromPipeline
		componentInitErrOutcome.Stage = stageId
		componentInitErrOutcome.RealStage = stageInfo.FromStage
		componentInitErrOutcome.Action = actionId
		componentInitErrOutcome.RealAction = actionInfo.FromAction
		componentInitErrOutcome.Event = eventId
		componentInitErrOutcome.Sequence = pipelineSequence
		componentInitErrOutcome.Status = false
		componentInitErrOutcome.Result = "component init error:" + err.Error()
		componentInitErrOutcome.Output = ""

		componentInitErrOutcome.GetOutcome().Save(componentInitErrOutcome)

		result, _ := json.Marshal(map[string]string{"message": "component init error:" + err.Error()})
		return http.StatusOK, result
	}

	componentInfo := new(models.ComponentLog)
	componentInfo.GetComponentLog().Where("id = ?", componentId).First(&componentInfo)
	c, err := module.InitComponet(*actionInfo, platformSetting["platformType"], platformSetting["platformHost"], pipelineInfo.Namespace)
	if err != nil {
		// do an extra log for outcome to record action result
		componentInitErrOutcome := new(models.Outcome)

		componentInitErrOutcome.Pipeline = pipelineId
		componentInitErrOutcome.RealPipeline = pipelineInfo.FromPipeline
		componentInitErrOutcome.Stage = stageId
		componentInitErrOutcome.RealStage = stageInfo.FromStage
		componentInitErrOutcome.Action = actionId
		componentInitErrOutcome.RealAction = actionInfo.FromAction
		componentInitErrOutcome.Event = eventId
		componentInitErrOutcome.Sequence = pipelineSequence
		componentInitErrOutcome.Status = false
		componentInitErrOutcome.Result = "component init error:" + err.Error()
		componentInitErrOutcome.Output = ""

		componentInitErrOutcome.GetOutcome().Save(componentInitErrOutcome)

		result, _ := json.Marshal(map[string]string{"message": "component init error:" + err.Error()})
		return http.StatusOK, result
	}

	if eventKey == models.EVENT_TASK_RESULT {
		reqBody, ok = reqBody["INFO"].(map[string]interface{})
		if !ok {
			result, _ := json.Marshal(map[string]string{"message": "illegal info, info is not a json"})
			return http.StatusOK, result
		}

		status, ok := reqBody["status"].(bool)
		if !ok {
			status = false
		}

		result, ok := reqBody["result"].(string)
		if !ok {
			result = ""
		}

		output, ok := reqBody["output"].(map[string]interface{})
		outputStr := ""
		if !ok {
			outputStr = ""
		} else {
			outputBytes, _ := json.Marshal(output)
			outputStr = string(outputBytes)
		}

		// do an extra log for outcome to record action result
		actionOutcome := new(models.Outcome)

		actionOutcome.Pipeline = pipelineId
		actionOutcome.RealPipeline = pipelineInfo.FromPipeline
		actionOutcome.Stage = stageId
		actionOutcome.RealStage = stageInfo.FromStage
		actionOutcome.Action = actionId
		actionOutcome.RealAction = actionInfo.FromAction
		actionOutcome.Event = eventId
		actionOutcome.Sequence = pipelineSequence
		actionOutcome.Status = status
		actionOutcome.Result = result
		actionOutcome.Output = outputStr

		actionOutcome.GetOutcome().Save(actionOutcome)

		// set a timer to auto stop component or service if it not close by itself
		go func() {
			time.Sleep(60 * time.Second)
			// check is this component already closed
			var isClose int64
			new(models.Event).GetEvent().Where("pipeline = ?", pipelineId).Where("stage = ?", stageId).Where("action = ?", actionId).Where("sequence = ?", pipelineSequence).Where("title = ?", models.EVENT_COMPONENT_STOP).Count(&isClose)

			if isClose < 1 {
				c.Stop(runId)
			}

		}()
	}

	if eventKey == models.EVENT_COMPONENT_STOP {
		// need send single to stop component or service
		c.Stop(runId)
	}

	headerMap := make(map[string]interface{})
	for key, value := range ctx.Resp.Header() {
		headerMap[key] = value
	}
	headerBytes, _ := json.Marshal(headerMap)

	eventDefine := new(models.EventDefinition)
	eventDefine.GetEventDefinition().Where("id = ?", eventId).First(&eventDefine)

	authStr := ""
	auths, ok := headerMap["Authorization"].([]string)
	if !ok {
		authStr = ""
	} else {
		authStr = strings.Join(auths, ";")
	}

	// log evnet
	eventLog := new(models.Event)
	eventLog.Pipeline = pipelineId
	eventLog.Stage = stageId
	eventLog.Action = actionId
	eventLog.Sequence = pipelineSequence
	eventLog.Definition = eventId
	eventLog.Title = eventKey
	eventLog.Header = string(headerBytes)
	eventLog.Payload = string(bodyByte)
	eventLog.Authorization = authStr
	eventLog.Type = eventDefine.Type
	eventLog.Source = eventDefine.Source

	eventLog.GetEvent().Save(eventLog)

	return http.StatusOK, result
}

// PutActionRegisterV1Handler is all action register here
func PutActionRegisterV1Handler(ctx *macaron.Context) (int, []byte) {
	result, _ := json.Marshal(map[string]string{"message": "ok"})

	// runId := ctx.Query("runId")
	// podName := ctx.Query("podName")
	// receiveUrl := ctx.Query("receiveUrl")

	bodyByte, _ := ctx.Req.Body().Bytes()

	reqBody := make(map[string]interface{})
	json.Unmarshal(bodyByte, &reqBody)
	// reqBody, ok = reqBody["INFO"].(map[string]interface{})
	// if !ok {
	// 	result, _ := json.Marshal(map[string]string{"message": "illegal reqData, reqData.INFO is not a json"})
	// 	return http.StatusOK, result
	// }

	runId, ok := reqBody["RUN_ID"].(string)
	if !ok {
		result, _ := json.Marshal(map[string]string{"message": "illegal runId, runId is not a string"})
		return http.StatusOK, result
	}

	podName, ok := reqBody["POD_NAME"].(string)
	if !ok {
		result, _ := json.Marshal(map[string]string{"message": "illegal podName, podName is not a string"})
		return http.StatusOK, result
	}

	receiveUrl, ok := reqBody["RECEIVE_URL"].(string)
	if !ok {
		result, _ := json.Marshal(map[string]string{"message": "illegal receiveUrl, receiveUrl is not a string"})
		return http.StatusOK, result
	}

	if len(strings.Split(runId, ",")) != 5 {
		// illegal runId return
		result, _ := json.Marshal(map[string]string{"message": "illegal id:" + runId})
		return http.StatusOK, result
	}

	go module.SendDataToAction(runId, receiveUrl, podName)
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
