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

	"github.com/Huawei/containerops/pilotage/models"
	"github.com/Huawei/containerops/pilotage/module"

	"gopkg.in/macaron.v1"
)

// GetComponentListV1Handler is get all component list
func GetComponentListV1Handler(ctx *macaron.Context) (int, []byte) {
	result, _ := json.Marshal(map[string]string{"message": ""})

	namespace := ctx.Params(":namespace")

	if namespace == "" {
		result, _ = json.Marshal(map[string]string{"errMsg": "namespace can't be null"})
		return http.StatusBadRequest, result
	}

	resultMap := make([]map[string]interface{}, 0)
	componentList := make([]models.Component, 0)
	componentsMap := make(map[int64]interface{})
	new(models.Component).GetComponent().Where("namespace = ?", namespace).Order("-id").Find(&componentList)

	for _, componentInfo := range componentList {
		if _, ok := componentsMap[componentInfo.ID]; !ok {
			tempMap := make(map[string]interface{})
			tempMap["version"] = make(map[int64]interface{})
			componentsMap[componentInfo.ID] = tempMap
		}

		componentMap := componentsMap[componentInfo.ID].(map[string]interface{})
		versionMap := componentMap["version"].(map[int64]interface{})

		versionMap[componentInfo.VersionCode] = componentInfo
		componentMap["id"] = componentInfo.ID
		componentMap["name"] = componentInfo.Component
		componentMap["version"] = versionMap
	}

	for _, component := range componentList {
		componentInfo := componentsMap[component.ID].(map[string]interface{})

		versionList := make([]map[string]interface{}, 0)
		for _, componentVersion := range componentList {
			if componentVersion.Component == componentInfo["name"].(string) {
				versionMap := make(map[string]interface{})
				versionMap["id"] = componentVersion.ID
				versionMap["version"] = componentVersion.Version
				versionMap["versionCode"] = componentVersion.VersionCode

				versionList = append(versionList, versionMap)
			}
		}

		tempResult := make(map[string]interface{})
		tempResult["id"] = componentInfo["id"]
		tempResult["name"] = componentInfo["name"]
		tempResult["version"] = versionList

		resultMap = append(resultMap, tempResult)
	}

	result, _ = json.Marshal(map[string]interface{}{"list": resultMap})

	return http.StatusOK, result
}

//PostComponentV1Handler is
func PostComponentV1Handler(ctx *macaron.Context) (int, []byte) {
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

	resultStr, err := module.CreateNewComponent(namespace, body.Name, body.Version)
	if err != nil {
		result, _ = json.Marshal(map[string]string{"errMsg": "error when create pipeline:" + err.Error()})
		return http.StatusBadRequest, result
	}

	result, _ = json.Marshal(map[string]string{"message": resultStr})

	return http.StatusOK, result
}

//PutComponentV1Handler is
func PutComponentV1Handler(ctx *macaron.Context) (int, []byte) {
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

	componentInfo := new(models.Component)
	err = componentInfo.GetComponent().Where("id = ?", body.Id).Find(&componentInfo).Error
	if err != nil {
		result, _ = json.Marshal(map[string]string{"errMsg": "error when get component info from db:" + err.Error()})
		return http.StatusBadRequest, result
	}

	if componentInfo.ID == 0 {
		result, _ = json.Marshal(map[string]string{"errMsg": "component is not exist"})
		return http.StatusBadRequest, result
	}

	defineMap := make(map[string]interface{})
	if componentInfo.Manifest != "" {
		err = json.Unmarshal([]byte(componentInfo.Manifest), &defineMap)
		if err != nil {
			result, _ = json.Marshal(map[string]string{"errMsg": "error when save component info:" + err.Error()})
			return http.StatusBadRequest, result
		}
	}

	defineMap["define"] = body.Define
	defineByte, err := json.Marshal(defineMap)
	if err != nil {
		result, _ = json.Marshal(map[string]string{"errMsg": "error when save component info:" + err.Error()})
		return http.StatusBadRequest, result
	}

	componentInfo.Manifest = string(defineByte)
	if componentInfo.Version == body.Version {
		// err = componentInfo.GetComponent().Save(componentInfo).Error
		err = module.UpdateComponentInfo(*componentInfo)
	} else {
		err = module.CreateNewComponentVersion(*componentInfo, body.Version)
	}

	if err != nil {
		result, _ = json.Marshal(map[string]string{"errMsg": "error when save component info:" + err.Error()})
		return http.StatusBadRequest, result
	}

	result, _ = json.Marshal(map[string]string{"message": "success"})

	return http.StatusOK, result
}

//GetComponentV1Handler is
func GetComponentV1Handler(ctx *macaron.Context) (int, []byte) {
	result, _ := json.Marshal(map[string]string{"message": ""})

	namespace := ctx.Params(":namespace")
	componentName := ctx.Params(":component")
	id := ctx.QueryInt64("id")

	if namespace == "" {
		result, _ = json.Marshal(map[string]string{"errMsg": "namespace can't be empty"})
		return http.StatusBadRequest, result
	}

	if componentName == "" {
		result, _ = json.Marshal(map[string]string{"errMsg": "component can't be empty"})
		return http.StatusBadRequest, result
	}

	resultMap, err := module.GetComponentInfo(namespace, componentName, id)
	if err != nil {
		result, _ = json.Marshal(map[string]string{"errMsg": "error when get component info:" + err.Error()})
		return http.StatusBadRequest, result
	}

	result, _ = json.Marshal(resultMap)
	return http.StatusOK, result
}

//DeleteComponentv1Handler is
func DeleteComponentv1Handler(ctx *macaron.Context) (int, []byte) {
	result, _ := json.Marshal(map[string]string{"message": ""})
	return http.StatusOK, result
}
