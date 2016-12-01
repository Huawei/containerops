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

//PostWorkflowV1Handler is create workflow data with namespace/repository and basic workflow data.
func PostWorkflowV1Handler(ctx *macaron.Context) (int, []byte) {
	result, _ := json.Marshal(map[string]string{"message": ""})

	body := new(struct {
		Name    string `json:"name`
		Version string `json:"version`
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

	resultStr, err := module.CreateNewWorkflow(namespace, repository, body.Name, body.Version)
	if err != nil {
		result, _ = json.Marshal(map[string]string{"errMsg": "error when create workflow:" + err.Error()})
		return http.StatusBadRequest, result
	}

	result, _ = json.Marshal(map[string]string{"message": resultStr})

	return http.StatusOK, result
}

//PostWorkflowJSONV1Handler is create workflow with entire data.
func PostWorkflowJSONV1Handler(ctx *macaron.Context) (int, []byte) {
	result, _ := json.Marshal(map[string]string{"message": ""})
	return http.StatusOK, result
}

//GetWorkflowListV1Handler is get workflow list info
func GetWorkflowListV1Handler(ctx *macaron.Context) (int, []byte) {
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

	workflowList, err := module.GetWorkflowListByNamespaceAndRepository(namespace, repository)
	if err != nil {
		result, _ = json.Marshal(map[string]string{"errMsg": "error when get workflow list:" + err.Error()})
		return http.StatusBadRequest, result
	}

	result, _ = json.Marshal(map[string]interface{}{"list": workflowList})

	return http.StatusOK, result
}

//GetWorkflowV1Handler is get workflow data, json/yaml format.
func GetWorkflowV1Handler(ctx *macaron.Context) (int, []byte) {
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

	id := ctx.QueryInt64("id")
	if id == int64(0) {
		result, _ = json.Marshal(map[string]string{"errMsg": "workflow's id can't be zero"})
		return http.StatusBadRequest, result
	}

	resultMap, err := module.GetWorkflowInfo(namespace, repository, workflowName, id)
	if err != nil {
		result, _ = json.Marshal(map[string]string{"errMsg": "error when get workflow info:" + err.Error()})
		return http.StatusBadRequest, result
	}

	result, _ = json.Marshal(resultMap)

	return http.StatusOK, result
}

// GetWorkflowTokenV1Handler is
func GetWorkflowTokenV1Handler(ctx *macaron.Context) (int, []byte) {
	result, _ := json.Marshal(map[string]string{"message": ""})

	namespace := ctx.Params(":namespace")
	repository := ctx.Params(":repository")
	workflowName := ctx.Params(":workflow")
	id := ctx.QueryInt64("id")

	if namespace == "" {
		result, _ = json.Marshal(map[string]string{"errMsg": "namespace can't be empty"})
		return http.StatusBadRequest, result
	}

	if repository == "" {
		result, _ = json.Marshal(map[string]string{"errMsg": "repository can't be empty"})
		return http.StatusBadRequest, result
	}

	if workflowName == "" {
		result, _ = json.Marshal(map[string]string{"errMsg": "workflow can't be empty"})
		return http.StatusBadRequest, result
	}

	workflow, err := module.GetWorkflow(id)
	if err != nil {
		result, _ = json.Marshal(map[string]string{"errMsg": err.Error()})
		return http.StatusBadRequest, result
	}

	tokenInfo, err := workflow.GetWorkflowToken()
	if err != nil {
		result, _ = json.Marshal(map[string]string{"errMsg": err.Error()})
		return http.StatusBadRequest, result
	}

	result, _ = json.Marshal(tokenInfo)

	return http.StatusOK, result
}

// GetWorkflowHistoriesV1Handler is
func GetWorkflowHistoriesV1Handler(ctx *macaron.Context) (int, []byte) {
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

	// resultMap, err := module.GetWorkflowHistoriesList(namespace)
	resultMap, err := module.GetWorkflowRunHistoryList(namespace, repository)
	if err != nil {
		result, _ = json.Marshal(map[string]string{"errMsg": err.Error()})
		return http.StatusBadRequest, result
	}

	result, _ = json.Marshal(map[string]interface{}{"workflowList": resultMap})
	return http.StatusOK, result
}

func GetWorkflowHistoryDefineV1Handler(ctx *macaron.Context) (int, []byte) {
	result, _ := json.Marshal(map[string]string{"message": ""})

	sequence := ctx.QueryInt64("sequence")

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

	workflow := ctx.Params(":workflow")
	if workflow == "" {
		result, _ = json.Marshal(map[string]string{"errMsg": "workflow can't be empty"})
		return http.StatusBadRequest, result
	}

	version := ctx.Params(":version")
	if version == "" {
		result, _ = json.Marshal(map[string]string{"errMsg": "request workflow's version id can't be zero"})
		return http.StatusBadRequest, result
	}

	// resultMap, err := module.GetWorkflowDefineByRunSequence(versionId, sequenceId)
	workflowLog, err := module.GetWorkflowLog(namespace, repository, workflow, version, sequence)
	resultMap, err := workflowLog.GetDefineInfo()
	if err != nil {
		result, _ = json.Marshal(map[string]string{"errMsg": err.Error()})
		return http.StatusBadRequest, result
	}

	result, _ = json.Marshal(map[string]interface{}{"define": resultMap})
	return http.StatusOK, result
}

func GetSequenceLineHistoryV1Handler(ctx *macaron.Context) (int, []byte) {
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

	workflow := ctx.Params(":workflow")
	if workflow == "" {
		result, _ = json.Marshal(map[string]string{"errMsg": "workflow can't be empty"})
		return http.StatusBadRequest, result
	}

	version := ctx.Params(":version")
	if version == "" {
		result, _ = json.Marshal(map[string]string{"errMsg": "request workflow's version id can't be zero"})
		return http.StatusBadRequest, result
	}

	sequence := ctx.Params(":sequence")
	sequenceInt, err := strconv.ParseInt(sequence, 10, 64)
	if err != nil {
		result, _ = json.Marshal(map[string]string{"errMsg": "request workflow's sequence is illegal"})
		return http.StatusBadRequest, result
	}

	lineId := ctx.Params(":relation")
	if lineId == "" || len(strings.Split(lineId, "-")) != 4 {
		result, _ = json.Marshal(map[string]string{"errMsg": "request workflow's relation is illegal:" + lineId})
		return http.StatusBadRequest, result
	}

	lineInputData := make(map[string]interface{})
	if strings.Split(lineId, "-")[0] == "s" {
		workflowLog, err := module.GetWorkflowLog(namespace, repository, workflow, version, sequenceInt)
		if err != nil {
			result, _ = json.Marshal(map[string]string{"errMsg": "error when get request workflow info:" + err.Error()})
			return http.StatusBadRequest, result
		}

		lineInputData, err = workflowLog.GetStartStageData()
	} else {
		actionLogId, err := strconv.ParseInt(strings.Split(lineId, "-")[1], 10, 64)
		if err != nil {
			result, _ = json.Marshal(map[string]string{"errMsg": "error when get request action info:" + err.Error()})
			return http.StatusBadRequest, result
		}

		actionLog, err := module.GetActionLog(actionLogId)
		if err != nil {
			result, _ = json.Marshal(map[string]string{"errMsg": "error when get request action info:" + err.Error()})
			return http.StatusBadRequest, result
		}

		lineInputData, err = actionLog.GetOutputData()
	}

	if err != nil {
		result, _ = json.Marshal(map[string]string{"errMsg": "error when get request line input info:" + err.Error()})
		return http.StatusBadRequest, result
	}

	actionLogId, err := strconv.ParseInt(strings.Split(lineId, "-")[3], 10, 64)
	if err != nil {
		result, _ = json.Marshal(map[string]string{"errMsg": "error when get request action info:" + err.Error()})
		return http.StatusBadRequest, result
	}

	actionLog, err := module.GetActionLog(actionLogId)
	if err != nil {
		result, _ = json.Marshal(map[string]string{"errMsg": "error when get request action info:" + err.Error()})
		return http.StatusBadRequest, result
	}

	lineoutputData, err := actionLog.GetInputData()

	resultMap := make(map[string]interface{})
	resultMap["define"] = map[string]interface{}{
		"input":  lineInputData,
		"output": lineoutputData,
	}

	result, _ = json.Marshal(resultMap)
	return http.StatusOK, result
}

//PutWorkflowV1Handler is update the workflow data, json format data in the http body.
func PutWorkflowV1Handler(ctx *macaron.Context) (int, []byte) {
	result, _ := json.Marshal(map[string]string{"message": ""})

	body := new(struct {
		Id      int64                  `json:"id"`
		Version string                 `json:"version"`
		Define  map[string]interface{} `json:"define"`
	})

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

	// workflowInfo := new(models.Workflow)
	// err = workflowInfo.GetWorkflow().Where("id = ?", body.Id).Find(&workflowInfo).Error
	workflowInfo, err := module.GetWorkflow(body.Id)
	if err != nil {
		result, _ = json.Marshal(map[string]string{"errMsg": "error when get workflow info from db:" + err.Error()})
		return http.StatusBadRequest, result
	}

	if workflowInfo.ID == 0 {
		result, _ = json.Marshal(map[string]string{"errMsg": "workflow is not exist"})
		return http.StatusBadRequest, result
	}

	if workflowInfo.Namespace != ctx.Params(":namespace") || workflowInfo.Repository != ctx.Params(":repository") || workflowInfo.Workflow.Workflow != ctx.Params(":workflow") {
		result, _ = json.Marshal(map[string]string{"errMsg": "request workflow is not equal to the given one"})
		return http.StatusBadRequest, result
	}

	if err != nil {
		result, _ = json.Marshal(map[string]string{"errMsg": "error when unmarshal request workflow's define:" + err.Error()})
		return http.StatusBadRequest, result
	}

	if workflowInfo.Version == body.Version {
		err = workflowInfo.UpdateWorkflowInfo(body.Define)
	} else {
		err = workflowInfo.CreateNewVersion(body.Define, body.Version)
	}

	if err != nil {
		result, _ = json.Marshal(map[string]string{"errMsg": "error when save workflow info:" + err.Error()})
		return http.StatusBadRequest, result
	}

	result, _ = json.Marshal(map[string]string{"message": "success"})

	return http.StatusOK, result
}

// PutWorkflowEnvV1Handler is set a workflow's env
func PutWorkflowEnvV1Handler(ctx *macaron.Context) (int, []byte) {
	result, _ := json.Marshal(map[string]string{"message": ""})

	body := new(struct {
		Id  int64                  `json:"id"`
		Env map[string]interface{} `json:"env"`
	})

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

	workflowInfo := new(models.Workflow)
	err = workflowInfo.GetWorkflow().Where("id = ?", body.Id).Find(&workflowInfo).Error
	if err != nil {
		result, _ = json.Marshal(map[string]string{"errMsg": "error when get workflow info from db:" + err.Error()})
		return http.StatusBadRequest, result
	}

	if workflowInfo.ID == 0 {
		result, _ = json.Marshal(map[string]string{"errMsg": "workflow is not exist"})
		return http.StatusBadRequest, result
	}

	envByte, err := json.Marshal(body.Env)
	if err != nil {
		result, _ = json.Marshal(map[string]string{"errMsg": "error when unmarshal env info:" + err.Error()})
		return http.StatusBadRequest, result
	}

	workflowInfo.Env = string(envByte)
	err = workflowInfo.GetWorkflow().Save(workflowInfo).Error

	if err != nil {
		result, _ = json.Marshal(map[string]string{"errMsg": "error when save workflow info:" + err.Error()})
		return http.StatusBadRequest, result
	}

	result, _ = json.Marshal(map[string]string{"message": "success"})
	return http.StatusOK, result
}

// GetWorkflowEnvV1Handler
func GetWorkflowEnvV1Handler(ctx *macaron.Context) (int, []byte) {
	result, _ := json.Marshal(map[string]string{"message": ""})

	id := ctx.QueryInt64("id")
	if id == 0 {
		result, _ = json.Marshal(map[string]string{"errMsg": "workflow's id can't be zero"})
		return http.StatusBadRequest, result
	}

	workflowInfo := new(models.Workflow)
	err := workflowInfo.GetWorkflow().Where("id = ?", id).Find(&workflowInfo).Error
	if err != nil {
		result, _ = json.Marshal(map[string]string{"errMsg": "error when get workflow info from db:" + err.Error()})
		return http.StatusBadRequest, result
	}

	if workflowInfo.ID == 0 {
		result, _ = json.Marshal(map[string]string{"errMsg": "workflow is not exist"})
		return http.StatusBadRequest, result
	}

	envMap := make(map[string]interface{})
	json.Unmarshal([]byte(workflowInfo.Env), &envMap)

	result, err = json.Marshal(map[string]interface{}{"env": envMap})
	return http.StatusOK, result
}

// PutWorkflowVarV1Handler is set a workflow's Var
func PutWorkflowVarV1Handler(ctx *macaron.Context) (int, []byte) {
	result, _ := json.Marshal(map[string]string{"message": ""})

	body := new(struct {
		Id  int64                  `json:"id"`
		Var map[string]interface{} `json:"var"`
	})

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

	err = module.SetWorkflowVarInfo(body.Id, body.Var)
	if err != nil {
		result, _ = json.Marshal(map[string]string{"errMsg": "error when save var info:" + err.Error()})
		return http.StatusBadRequest, result
	}

	result, _ = json.Marshal(map[string]string{"message": "success"})
	return http.StatusOK, result
}

// GetWorkflowVarV1Handler
func GetWorkflowVarV1Handler(ctx *macaron.Context) (int, []byte) {
	result, _ := json.Marshal(map[string]string{"message": ""})

	id := ctx.QueryInt64("id")
	if id == 0 {
		result, _ = json.Marshal(map[string]string{"errMsg": "workflow's id can't be zero"})
		return http.StatusBadRequest, result
	}

	varMap, err := module.GetWorkflowVarInfo(id)
	if err != nil {
		result, _ = json.Marshal(map[string]string{"errMsg": "error when get workflow var info"})
		return http.StatusBadRequest, result
	}

	result, err = json.Marshal(map[string]interface{}{"var": varMap})
	return http.StatusOK, result
}

//DeleteWorkflowV1Handler is delete the workflow data.
func DeleteWorkflowV1Handler(ctx *macaron.Context) (int, []byte) {
	result, _ := json.Marshal(map[string]string{"message": ""})
	return http.StatusOK, result
}

//ExecuteWorkflowV1Handler is run a workflow
func ExecuteWorkflowV1Handler(ctx *macaron.Context) (int, []byte) {
	result := []byte("")

	version := ctx.Query("version")
	namespace := ctx.Params(":namespace")
	repository := ctx.Params(":repository")
	workflowName := ctx.Params(":workflow")

	reqHeader := ctx.Req.Request.Header
	reqBody, err := ctx.Req.Body().Bytes()
	if err != nil {
		result, _ = json.Marshal(map[string]string{"errMsg": "error when get request body:" + err.Error()})
		return http.StatusBadRequest, result
	}

	workflowInfo, err := module.GetLatestRunableWorkflow(namespace, repository, workflowName, version)
	if err != nil {
		result, _ = json.Marshal(map[string]string{"errMsg": "error when get workflow info:" + err.Error()})
		return http.StatusBadRequest, result
	}

	if ok, eventMap, err := workflowInfo.BeforeExecCheck(reqHeader, reqBody); !ok {
		result, _ = json.Marshal(map[string]string{"errMsg": "failed on before exec check" + err.Error()})
		return http.StatusBadRequest, result
	} else {
		authMap := make(map[string]interface{})
		authMap["type"] = module.AuthTypeWorkflowDefault
		authMap["token"] = module.AuthTokenDefault
		authMap["eventName"] = eventMap["eventName"]
		authMap["eventType"] = eventMap["sourceType"]
		authMap["time"] = time.Now().Format("2006-01-02 15:04:05")

		err := module.Run(workflowInfo.ID, authMap, string(reqBody))
		if err != nil {
			result, _ = json.Marshal(map[string]string{"result": "error when run workflow:" + err.Error()})
			return http.StatusBadRequest, result
		}

		result, _ = json.Marshal(map[string]string{"result": "success"})
	}

	return http.StatusOK, result
}

//GetOutcomeListV1Handler is
func GetOutcomeListV1Handler(ctx *macaron.Context) (int, []byte) {
	result, _ := json.Marshal(map[string]string{"message": ""})
	return http.StatusOK, result
}

//GetOutcomeV1Handler is
func GetOutcomeV1Handler(ctx *macaron.Context) (int, []byte) {
	result, _ := json.Marshal(map[string]string{"message": ""})
	return http.StatusOK, result
}

// PutWorkflowStateV1Handler is
func PutWorkflowStateV1Handler(ctx *macaron.Context) (int, []byte) {
	result, _ := json.Marshal(map[string]string{"message": ""})

	body := new(struct {
		Id    int64 `json:"id"`
		State int64 `json:"state"`
	})

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

	if body.State != models.WorkflowStateAble && body.State != models.WorkflowStateDisable {
		result, _ = json.Marshal(map[string]string{"errMsg": "state code is illegal"})
		return http.StatusBadRequest, result
	}

	workflowInfo := new(models.Workflow)
	err = workflowInfo.GetWorkflow().Where("id = ?", body.Id).Find(&workflowInfo).Error
	if err != nil {
		result, _ = json.Marshal(map[string]string{"errMsg": "error when get workflow info from db:" + err.Error()})
		return http.StatusBadRequest, result
	}

	if workflowInfo.ID == 0 {
		result, _ = json.Marshal(map[string]string{"errMsg": "workflow is not exist"})
		return http.StatusBadRequest, result
	}

	workflowInfo.State = body.State
	err = workflowInfo.GetWorkflow().Save(workflowInfo).Error

	if err != nil {
		result, _ = json.Marshal(map[string]string{"errMsg": "error when save workflow info:" + err.Error()})
		return http.StatusBadRequest, result
	}

	result, _ = json.Marshal(map[string]string{"message": "success"})
	return http.StatusOK, result
}
