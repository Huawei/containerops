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
	"github.com/Huawei/containerops/pilotage/module"

	"gopkg.in/macaron.v1"
)

// GetComponentListV1Handler is get all component list
func GetComponentListV1Handler(ctx *macaron.Context) (int, []byte) {
	result, _ := json.Marshal(map[string]string{"message": ""})

	namespace := ctx.Params(":namespace")

	if namespace == "" {
		result, _ = json.Marshal(map[string]string{"errMsg": "namespace can't be empty"})
		return http.StatusBadRequest, result
	}

	componentList, err := module.GetComponentListByNamespace(namespace)

	if err != nil {
		result, _ = json.Marshal(map[string]string{"errMsg": "error when get component list:" + err.Error()})
		return http.StatusBadRequest, result
	}

	result, _ = json.Marshal(map[string]interface{}{"list": componentList})

	return http.StatusOK, result
}

//PostComponentV1Handler is create a new component
func PostComponentV1Handler(ctx *macaron.Context) (int, []byte) {
	result, _ := json.Marshal(map[string]string{"message": ""})

	body := new(struct {
		Name    string `json:"name"`
		Version string `json:"version"`
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
		result, _ = json.Marshal(map[string]string{"errMsg": "error when create workflow:" + err.Error()})
		return http.StatusBadRequest, result
	}

	result, _ = json.Marshal(map[string]string{"message": resultStr})

	return http.StatusOK, result
}

//GetComponentV1Handler is get a specified component info by component id
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

//PutComponentV1Handler is update a component info or create a new component's version
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

//DeleteComponentv1Handler is
func DeleteComponentv1Handler(ctx *macaron.Context) (int, []byte) {
	result := []byte("")

	componentID := ctx.ParamsInt64(":component")

	err := module.DeleteComponentInfo(componentID)
	if err != nil {
		result, _ = json.Marshal(map[string]string{"errMsg": "error when delete component info from db:" + err.Error()})
		return http.StatusBadRequest, result
	}

	result, _ = json.Marshal(map[string]string{"message": "success"})
	return http.StatusOK, result
}
