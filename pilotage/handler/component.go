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
	"github.com/prometheus/common/log"
	"strconv"
)

func ListComponents(ctx *macaron.Context) (int, []byte) {
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

func CreateComponent(ctx *macaron.Context) (int, []byte) {
	var httpStatus int64
	var resp ComponentResp
	body, err := ctx.Req.Body().Bytes()
	if err != nil {
		httpStatus = http.StatusBadRequest
		resp.OK = false
		resp.ErrorCode = componentErrCode + 1
		resp.Message = "Get requrest body error: " + err.Error()
	}
	if id, err := module.CreateComponent(body); err != nil {
		httpStatus = http.StatusBadRequest
		resp.OK = false
		resp.ErrorCode = componentErrCode + 2
		resp.Message = "Create component error: " + err.Error()
	} else {
		httpStatus = http.StatusOK
		resp.ID = id
		resp.OK = true
		resp.Message = "Component Created"
	}

	result, err := json.Marshal(resp)
	if err != nil {
		log.Errorln("Create component marshal data error: " + err.Error())
	}
	return httpStatus, result
}

func GetComponent(ctx *macaron.Context) (int, []byte) {
	var httpStatus int64
	var resp ComponentResp
	componentID := ctx.Params(":component_id")
	id, err := strconv.Atoi(componentID)
	if err != nil {
		httpStatus = http.StatusBadRequest
		resp.OK = false
		resp.ErrorCode = componentErrCode + 3
		resp.Message = "Parse component id error: " + err.Error()
	}
	if component, err := module.GetComponentByID(id); err != nil {
		httpStatus = http.StatusBadRequest
		resp.OK = false
		resp.ErrorCode = componentErrCode + 4
		resp.Message = "get component by id error: " + err.Error()
	} else {
		httpStatus = http.StatusOK
		resp.OK = true
		resp.Component = component
	}
	result, err := json.Marshal(resp)
	if err != nil {
		log.Errorln("Get component marshal data error: " + err.Error())
	}
	return httpStatus, result
}

func UpdateComponent(ctx *macaron.Context) (int, []byte) {
	var httpStatus int64
	var resp ComponentResp
	body, err := ctx.Req.Body().Bytes()
	if err != nil {
		httpStatus = http.StatusBadRequest
		resp.OK = false
		resp.ErrorCode = componentErrCode + 5
		resp.Message = "Get requrest body error: " + err.Error()
	}

	componentID := ctx.Params(":component_id")
	id, err := strconv.Atoi(componentID)
	if err != nil {
		httpStatus = http.StatusBadRequest
		resp.OK = false
		resp.ErrorCode = componentErrCode + 6
		resp.Message = "Parse component id error: " + err.Error()
	}

	if err := module.UpdateComponent(id, body); err != nil {
		httpStatus = http.StatusBadRequest
		resp.OK = false
		resp.ErrorCode = componentErrCode + 7
		resp.Message = "update component error: " + err.Error()
	} else {
		httpStatus = http.StatusOK
		resp.OK = true
	}

	result, err := json.Marshal(resp)
	if err != nil {
		log.Errorln("Update component marshal data error: " + err.Error())
	}
	return httpStatus, result
	//
	//err = json.Unmarshal(reqBody, &body)
	//if err != nil {
	//	result, _ = json.Marshal(map[string]string{"errMsg": "error when unmarshal request body:" + err.Error()})
	//	return http.StatusBadRequest, result
	//}
	//
	//componentInfo := new(models.Component)
	//err = componentInfo.GetComponent().Where("id = ?", body.Id).Find(&componentInfo).Error
	//if err != nil {
	//	result, _ = json.Marshal(map[string]string{"errMsg": "error when get component info from db:" + err.Error()})
	//	return http.StatusBadRequest, result
	//}
	//
	//if componentInfo.ID == 0 {
	//	result, _ = json.Marshal(map[string]string{"errMsg": "component is not exist"})
	//	return http.StatusBadRequest, result
	//}
	//
	//defineMap := make(map[string]interface{})
	//if componentInfo.Manifest != "" {
	//	err = json.Unmarshal([]byte(componentInfo.Manifest), &defineMap)
	//	if err != nil {
	//		result, _ = json.Marshal(map[string]string{"errMsg": "error when save component info:" + err.Error()})
	//		return http.StatusBadRequest, result
	//	}
	//}
	//
	//defineMap["define"] = body.Define
	//defineByte, err := json.Marshal(defineMap)
	//if err != nil {
	//	result, _ = json.Marshal(map[string]string{"errMsg": "error when save component info:" + err.Error()})
	//	return http.StatusBadRequest, result
	//}
	//
	//componentInfo.Manifest = string(defineByte)
	//if componentInfo.Version == body.Version {
	//	// err = componentInfo.GetComponent().Save(componentInfo).Error
	//	err = module.UpdateComponent(*componentInfo)
	//} else {
	//	err = module.CreateNewComponentVersion(*componentInfo, body.Version)
	//}
	//
	//if err != nil {
	//	result, _ = json.Marshal(map[string]string{"errMsg": "error when save component info:" + err.Error()})
	//	return http.StatusBadRequest, result
	//}
	//
	//result, _ = json.Marshal(map[string]string{"message": "success"})
	//
	//return http.StatusOK, result
}

func DeleteComponent(ctx *macaron.Context) (int, []byte) {
	result, _ := json.Marshal(map[string]string{"message": ""})
	return http.StatusOK, result
}

func DebugComponent(ctx *macaron.Context) (int, []byte) {
	return nil, nil
}
