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

//PostPipelineV1Handler is create pipeline data with namespace/repository and basic pipeline data.
func PostPipelineV1Handler(ctx *macaron.Context) (int, []byte) {
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

	resultStr, err := module.CreateNewPipeline(namespace, repository, body.Name, body.Version)
	if err != nil {
		result, _ = json.Marshal(map[string]string{"errMsg": "error when create pipeline:" + err.Error()})
		return http.StatusBadRequest, result
	}

	result, _ = json.Marshal(map[string]string{"message": resultStr})

	return http.StatusOK, result
}

//PostPipelineJSONV1Handler is create pipeline with entire data.
func PostPipelineJSONV1Handler(ctx *macaron.Context) (int, []byte) {
	result, _ := json.Marshal(map[string]string{"message": ""})
	return http.StatusOK, result
}

//GetPipelineListV1Handler is get pipeline list info
func GetPipelineListV1Handler(ctx *macaron.Context) (int, []byte) {
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

	pipelineList, err := module.GetPipelineListByNamespaceAndRepository(namespace, repository)
	if err != nil {
		result, _ = json.Marshal(map[string]string{"errMsg": "error when get pipeline list:" + err.Error()})
		return http.StatusBadRequest, result
	}

	result, _ = json.Marshal(map[string]interface{}{"list": pipelineList})

	return http.StatusOK, result
}

//GetPipelineV1Handler is get pipeline data, json/yaml format.
func GetPipelineV1Handler(ctx *macaron.Context) (int, []byte) {
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

	id := ctx.QueryInt64("id")
	if id == int64(0) {
		result, _ = json.Marshal(map[string]string{"errMsg": "pipeline's id can't be zero"})
		return http.StatusBadRequest, result
	}

	resultMap, err := module.GetPipelineInfo(namespace, repository, pipelineName, id)
	if err != nil {
		result, _ = json.Marshal(map[string]string{"errMsg": "error when get pipeline info:" + err.Error()})
		return http.StatusBadRequest, result
	}

	result, _ = json.Marshal(resultMap)

	return http.StatusOK, result
}

// GetPipelineTokenV1Handler is
func GetPipelineTokenV1Handler(ctx *macaron.Context) (int, []byte) {
	result, _ := json.Marshal(map[string]string{"message": ""})

	namespace := ctx.Params(":namespace")
	repository := ctx.Params(":repository")
	pipelineName := ctx.Params(":workflow")
	id := ctx.QueryInt64("id")

	if namespace == "" {
		result, _ = json.Marshal(map[string]string{"errMsg": "namespace can't be empty"})
		return http.StatusBadRequest, result
	}

	if repository == "" {
		result, _ = json.Marshal(map[string]string{"errMsg": "repository can't be empty"})
		return http.StatusBadRequest, result
	}

	if pipelineName == "" {
		result, _ = json.Marshal(map[string]string{"errMsg": "pipeline can't be empty"})
		return http.StatusBadRequest, result
	}

	pipeline, err := module.GetPipeline(id)
	if err != nil {
		result, _ = json.Marshal(map[string]string{"errMsg": err.Error()})
		return http.StatusBadRequest, result
	}

	tokenInfo, err := pipeline.GetPipelineToken()
	if err != nil {
		result, _ = json.Marshal(map[string]string{"errMsg": err.Error()})
		return http.StatusBadRequest, result
	}

	result, _ = json.Marshal(tokenInfo)

	return http.StatusOK, result
}

// GetPipelineHistoriesV1Handler is
func GetPipelineHistoriesV1Handler(ctx *macaron.Context) (int, []byte) {
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

	// resultMap, err := module.GetPipelineHistoriesList(namespace)
	resultMap, err := module.GetPipelineRunHistoryList(namespace, repository)
	if err != nil {
		result, _ = json.Marshal(map[string]string{"errMsg": err.Error()})
		return http.StatusBadRequest, result
	}

	result, _ = json.Marshal(map[string]interface{}{"pipelineList": resultMap})
	return http.StatusOK, result
}

func GetPipelineHistoryDefineV1Handler(ctx *macaron.Context) (int, []byte) {
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
		result, _ = json.Marshal(map[string]string{"errMsg": "request pipeline's version id can't be zero"})
		return http.StatusBadRequest, result
	}

	// resultMap, err := module.GetPipelineDefineByRunSequence(versionId, sequenceId)
	pipelineLog, err := module.GetPipelineLog(namespace, repository, workflow, version, sequence)
	resultMap, err := pipelineLog.GetDefineInfo()
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
		result, _ = json.Marshal(map[string]string{"errMsg": "request pipeline's version id can't be zero"})
		return http.StatusBadRequest, result
	}

	sequence := ctx.Params(":sequence")
	sequenceInt, err := strconv.ParseInt(sequence, 10, 64)
	if err != nil {
		result, _ = json.Marshal(map[string]string{"errMsg": "request pipeline's sequence is illegal"})
		return http.StatusBadRequest, result
	}

	lineId := ctx.Params(":relation")
	if lineId == "" || len(strings.Split(lineId, "-")) != 4 {
		result, _ = json.Marshal(map[string]string{"errMsg": "request pipeline's relation is illegal:" + lineId})
		return http.StatusBadRequest, result
	}

	lineInputData := make(map[string]interface{})
	if strings.Split(lineId, "-")[0] == "s" {
		pipelineLog, err := module.GetPipelineLog(namespace, repository, workflow, version, sequenceInt)
		if err != nil {
			result, _ = json.Marshal(map[string]string{"errMsg": "error when get request pipeline info:" + err.Error()})
			return http.StatusBadRequest, result
		}

		lineInputData, err = pipelineLog.GetStartStageData()
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

//PutPipelineV1Handler is update the pipeline data, json format data in the http body.
func PutPipelineV1Handler(ctx *macaron.Context) (int, []byte) {
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

	// pipelineInfo := new(models.Pipeline)
	// err = pipelineInfo.GetPipeline().Where("id = ?", body.Id).Find(&pipelineInfo).Error
	pipelineInfo, err := module.GetPipeline(body.Id)
	if err != nil {
		result, _ = json.Marshal(map[string]string{"errMsg": "error when get pipeline info from db:" + err.Error()})
		return http.StatusBadRequest, result
	}

	if pipelineInfo.ID == 0 {
		result, _ = json.Marshal(map[string]string{"errMsg": "pipeline is not exist"})
		return http.StatusBadRequest, result
	}

	if pipelineInfo.Namespace != ctx.Params(":namespace") || pipelineInfo.Repository != ctx.Params(":repository") || pipelineInfo.Pipeline.Pipeline != ctx.Params(":workflow") {
		result, _ = json.Marshal(map[string]string{"errMsg": "request pipeline is not equal to the given one"})
		return http.StatusBadRequest, result
	}

	if err != nil {
		result, _ = json.Marshal(map[string]string{"errMsg": "error when unmarshal request pipeline's define:" + err.Error()})
		return http.StatusBadRequest, result
	}

	if pipelineInfo.Version == body.Version {
		err = pipelineInfo.UpdatePipelineInfo(body.Define)
	} else {
		err = pipelineInfo.CreateNewVersion(body.Define, body.Version)
	}

	if err != nil {
		result, _ = json.Marshal(map[string]string{"errMsg": "error when save pipeline info:" + err.Error()})
		return http.StatusBadRequest, result
	}

	result, _ = json.Marshal(map[string]string{"message": "success"})

	return http.StatusOK, result
}

// PutPipelineEnvV1Handler is set a pipeline's env
func PutPipelineEnvV1Handler(ctx *macaron.Context) (int, []byte) {
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

	pipelineInfo := new(models.Pipeline)
	err = pipelineInfo.GetPipeline().Where("id = ?", body.Id).Find(&pipelineInfo).Error
	if err != nil {
		result, _ = json.Marshal(map[string]string{"errMsg": "error when get pipeline info from db:" + err.Error()})
		return http.StatusBadRequest, result
	}

	if pipelineInfo.ID == 0 {
		result, _ = json.Marshal(map[string]string{"errMsg": "pipeline is not exist"})
		return http.StatusBadRequest, result
	}

	envByte, err := json.Marshal(body.Env)
	if err != nil {
		result, _ = json.Marshal(map[string]string{"errMsg": "error when unmarshal env info:" + err.Error()})
		return http.StatusBadRequest, result
	}

	pipelineInfo.Env = string(envByte)
	err = pipelineInfo.GetPipeline().Save(pipelineInfo).Error

	if err != nil {
		result, _ = json.Marshal(map[string]string{"errMsg": "error when save pipeline info:" + err.Error()})
		return http.StatusBadRequest, result
	}

	result, _ = json.Marshal(map[string]string{"message": "success"})
	return http.StatusOK, result
}

// GetPipelineEnvV1Handler
func GetPipelineEnvV1Handler(ctx *macaron.Context) (int, []byte) {
	result, _ := json.Marshal(map[string]string{"message": ""})

	id := ctx.QueryInt64("id")
	if id == 0 {
		result, _ = json.Marshal(map[string]string{"errMsg": "pipeline's id can't be zero"})
		return http.StatusBadRequest, result
	}

	pipelineInfo := new(models.Pipeline)
	err := pipelineInfo.GetPipeline().Where("id = ?", id).Find(&pipelineInfo).Error
	if err != nil {
		result, _ = json.Marshal(map[string]string{"errMsg": "error when get pipeline info from db:" + err.Error()})
		return http.StatusBadRequest, result
	}

	if pipelineInfo.ID == 0 {
		result, _ = json.Marshal(map[string]string{"errMsg": "pipeline is not exist"})
		return http.StatusBadRequest, result
	}

	envMap := make(map[string]interface{})
	json.Unmarshal([]byte(pipelineInfo.Env), &envMap)

	result, err = json.Marshal(map[string]interface{}{"env": envMap})
	return http.StatusOK, result
}

//DeletePipelineV1Handler is delete the pipeline data.
func DeletePipelineV1Handler(ctx *macaron.Context) (int, []byte) {
	result, _ := json.Marshal(map[string]string{"message": ""})
	return http.StatusOK, result
}

//ExecutePipelineV1Handler is run a pipeline
func ExecutePipelineV1Handler(ctx *macaron.Context) (int, []byte) {
	result := []byte("")

	version := ctx.Query("version")
	namespace := ctx.Params(":namespace")
	repository := ctx.Params(":repository")
	pipelineName := ctx.Params(":workflow")

	reqHeader := ctx.Req.Request.Header
	reqBody, err := ctx.Req.Body().Bytes()
	if err != nil {
		result, _ = json.Marshal(map[string]string{"errMsg": "error when get request body:" + err.Error()})
		return http.StatusBadRequest, result
	}

	pipelineInfo, err := module.GetLatestRunablePipeline(namespace, repository, pipelineName, version)
	if err != nil {
		result, _ = json.Marshal(map[string]string{"errMsg": "error when get pipeline info:" + err.Error()})
		return http.StatusBadRequest, result
	}

	if ok, err := pipelineInfo.BeforeExecCheck(reqHeader, reqBody); !ok {
		result, _ = json.Marshal(map[string]string{"errMsg": "failed on before exec check" + err.Error()})
		return http.StatusBadRequest, result
	} else {
		authMap := make(map[string]interface{})
		authMap["type"] = module.AuthTypePipelineDefault
		authMap["token"] = module.AuthTokenDefault
		authMap["time"] = time.Now().Format("2006-01-02 15:04:05")

		err := module.Run(pipelineInfo.ID, authMap, string(reqBody))
		if err != nil {
			result, _ = json.Marshal(map[string]string{"result": "error when run pipeline:" + err.Error()})
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

// PutPipelineStateV1Handler is
func PutPipelineStateV1Handler(ctx *macaron.Context) (int, []byte) {
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

	if body.State != models.PipelineStateAble && body.State != models.PipelineStateDisable {
		result, _ = json.Marshal(map[string]string{"errMsg": "state code is illegal"})
		return http.StatusBadRequest, result
	}

	pipelineInfo := new(models.Pipeline)
	err = pipelineInfo.GetPipeline().Where("id = ?", body.Id).Find(&pipelineInfo).Error
	if err != nil {
		result, _ = json.Marshal(map[string]string{"errMsg": "error when get pipeline info from db:" + err.Error()})
		return http.StatusBadRequest, result
	}

	if pipelineInfo.ID == 0 {
		result, _ = json.Marshal(map[string]string{"errMsg": "pipeline is not exist"})
		return http.StatusBadRequest, result
	}

	pipelineInfo.State = body.State
	err = pipelineInfo.GetPipeline().Save(pipelineInfo).Error

	if err != nil {
		result, _ = json.Marshal(map[string]string{"errMsg": "error when save pipeline info:" + err.Error()})
		return http.StatusBadRequest, result
	}

	result, _ = json.Marshal(map[string]string{"message": "success"})
	return http.StatusOK, result
}
