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

	"github.com/containerops/pilotage/models"
	"github.com/containerops/pilotage/module"

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

	resultStr, err := module.CreateNewPipeline(namespace, body.Name, body.Version)
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

//GetPipelineV1Handler is get pipeline data, json/yaml format.
func GetPipelineV1Handler(ctx *macaron.Context) (int, []byte) {
	result, _ := json.Marshal(map[string]string{"message": ""})

	namespace := ctx.Params(":namespace")
	pipelineName := ctx.Params(":pipeline")
	id := ctx.QueryInt64("id")

	if namespace == "" {
		result, _ = json.Marshal(map[string]string{"errMsg": "namespace can't be empty"})
		return http.StatusBadRequest, result
	}

	if pipelineName == "" {
		result, _ = json.Marshal(map[string]string{"errMsg": "pipeline can't be empty"})
		return http.StatusBadRequest, result
	}

	resultMap, err := module.GetPipelineList(namespace, pipelineName, id)
	if err != nil {
		result, _ = json.Marshal(map[string]string{"errMsg": "error when get pipeline info:" + err.Error()})
		return http.StatusBadRequest, result
	}

	result, _ = json.Marshal(resultMap)

	return http.StatusOK, result
}

//GetPipelineListV1Handler is get pipeline list info
func GetPipelineListV1Handler(ctx *macaron.Context) (int, []byte) {
	result, _ := json.Marshal(map[string]string{"message": ""})
	namespace := ctx.Params(":namespace")

	if namespace == "" {
		result, _ = json.Marshal(map[string]string{"errMsg": "namespace can't be null"})
		return http.StatusBadRequest, result
	}

	resultMap := make([]map[string]interface{}, 0)
	pipelineList := make([]models.Pipeline, 0)
	new(models.Pipeline).GetPipeline().Where("namespace = ?", namespace).Order("pipeline").Order("-version_code").Find(&pipelineList)

	for _, pipelineInfo := range pipelineList {

		shouldAppend := false
		pipelineMap := make(map[string]interface{}, 0)

		if len(resultMap) > 0 {
			if resultMap[len(resultMap)-1]["name"] == pipelineInfo.Pipeline {
				pipelineMap = resultMap[len(resultMap)-1]
			}
		}

		if pipelineMap["id"] == nil {
			pipelineMap["id"] = pipelineInfo.ID
			pipelineMap["name"] = pipelineInfo.Pipeline
			pipelineMap["version"] = make([]map[string]interface{}, 0)
			shouldAppend = true
		}

		versionMap := make(map[string]interface{})
		versionMap["id"] = pipelineInfo.ID
		versionMap["version"] = pipelineInfo.Version
		versionMap["versionCode"] = pipelineInfo.VersionCode

		versionList := pipelineMap["version"].([]map[string]interface{})
		pipelineMap["version"] = append(versionList, versionMap)

		if shouldAppend {
			resultMap = append(resultMap, pipelineMap)
		}
	}

	result, _ = json.Marshal(map[string]interface{}{"list": resultMap})

	return http.StatusOK, result
}

//PutPipelineV1Handler is update the pipeline data, json format data in the http body.
func PutPipelineV1Handler(ctx *macaron.Context) (int, []byte) {
	result, _ := json.Marshal(map[string]string{"message": ""})

	body := new(struct {
		Id     int64                  `json:"id"`
		Define map[string]interface{} `json:"define"`
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
	err = pipelineInfo.GetPipeline().Where("id = ?", body.Id).Error
	if err != nil {
		result, _ = json.Marshal(map[string]string{"errMsg": "error when get pipeline info from db:" + err.Error()})
		return http.StatusBadRequest, result
	}

	if pipelineInfo.ID == 0 {
		result, _ = json.Marshal(map[string]string{"errMsg": "pipeline is not exist"})
		return http.StatusBadRequest, result
	}

	defineMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(pipelineInfo.Manifest), &defineMap)
	if err != nil {
		result, _ = json.Marshal(map[string]string{"errMsg": "error when save pipeline info:" + err.Error()})
		return http.StatusBadRequest, result
	}

	defineMap["define"] = body.Define
	defineByte, err := json.Marshal(defineMap)
	if err != nil {
		result, _ = json.Marshal(map[string]string{"errMsg": "error when save pipeline info:" + err.Error()})
		return http.StatusBadRequest, result
	}

	pipelineInfo.Manifest = string(defineByte)

	err = pipelineInfo.GetPipeline().Save(pipelineInfo).Error
	if err != nil {
		result, _ = json.Marshal(map[string]string{"errMsg": "error when save pipeline info:" + err.Error()})
		return http.StatusBadRequest, result
	}

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

	// get pipeline info
	pipelineInfo := new(models.Pipeline)
	version := ctx.Query("version")
	namespace := ctx.Params(":namespace")
	pipeline := ctx.Params(":pipeline")
	//TODO: is version,namespace,pipeline illegal,if not ,return error

	// if version is nil, select a least and not disabled version
	if version == "" {
		tempPipelineInfo := new(models.Pipeline)
		err := tempPipelineInfo.GetPipeline().Where("namespace = ?", namespace).Where("pipeline = ?", pipeline).Where("state = ?", models.PipelineStateAble).Order("-version_code").First(&tempPipelineInfo).Error
		if err != nil || tempPipelineInfo.ID == 0 {
			result, _ := json.Marshal(map[string]string{"result": "error when get least useable pipeline version:" + err.Error()})
			return http.StatusOK, result
		}

		version = tempPipelineInfo.Version
	}

	var reqBody []byte
	err := pipelineInfo.GetPipeline().Where("namespace = ?", namespace).Where("pipeline = ?", pipeline).Where("version = ?", version).First(&pipelineInfo).Error
	if err != nil {
		result, _ = json.Marshal(map[string]string{"result": "error when get pipeline info:" + err.Error()})
	} else if pipelineInfo.ID == 0 || pipelineInfo.Version == "" {
		result, _ = json.Marshal(map[string]string{"result": "error when get pipeline info from namespace(" + ctx.Params(":namespace") + ") and pipeline(" + ctx.Params(":pipeline") + ")"})
	} else if pipelineInfo.SourceInfo == "" {
		result, _ = json.Marshal(map[string]string{"result": "pipeline does not config source info"})
	} else if pipelineInfo.State == models.PipelineStateDisable {
		result, _ = json.Marshal(map[string]string{"result": "pipeline is disabled!"})
	} else if reqBody, err = ctx.Req.Body().Bytes(); err != nil {
		result, _ = json.Marshal(map[string]string{"result": "error when get request body:" + err.Error()})
	} else if !module.PipeExecRequestLegal(ctx.Req.Request.Header, reqBody, *pipelineInfo) {
		result, _ = json.Marshal(map[string]string{"result": "request token is illegal!"})
	} else if pipelineLog, err := module.DoPipelineLog(*pipelineInfo); err != nil {
		// pipeline is ready , copy current pipelin info and all remain action will use the copy data
		result, _ = json.Marshal(map[string]string{"result": "error when do pipeline log:" + err.Error()})
	} else {
		resultStr := module.StartPipeline(*pipelineLog, string(reqBody))
		result, _ = json.Marshal(map[string]string{"result": resultStr})
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
