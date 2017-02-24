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

	"github.com/Huawei/containerops/component/module"

	logs "github.com/Huawei/containerops/component/log"
	"gopkg.in/macaron.v1"
)

//  log is log pkg instance
var log *logs.Logger

// init is init log pkg instance
func init() {
	log = logs.New()
}

// ListComponents is return component list with selector
func ListComponents(ctx *macaron.Context) (httpStatus int, result []byte) {
	resp := new(ListComponentsResp)
	resp.Components = make([]module.ComponentBaseData, 0)

	namespace := ctx.ParamsEscape("namespace")
	if namespace == "" {
		httpStatus = http.StatusBadRequest
		resp.OK = false
		resp.ErrorCode = ComponentError + ComponentEmptyNamespaceError
		resp.Message = "got empty namespace"

		result, _ = json.Marshal(resp)
		return
	}

	name := ctx.QueryTrim("name")
	fuzzy := ctx.QueryBool("fuzzy")

	pageNum := ctx.QueryInt("pageNum")
	if pageNum <= 0 {
		pageNum = 10
	}

	versionNum := ctx.QueryInt("versionNum")
	if versionNum <= 0 {
		versionNum = 3
	}

	offset := ctx.QueryInt("offset")
	if offset < 0 {
		offset = 0
	}

	components, err := module.GetComponents(namespace, name, fuzzy, pageNum, versionNum, offset)
	if err != nil {
		httpStatus = http.StatusBadRequest
		resp.OK = false
		resp.ErrorCode = ComponentError + ComponentListError
		resp.Message = "List components error: " + err.Error()

		result, _ = json.Marshal(resp)
		return
	}

	httpStatus = http.StatusOK
	resp.Components = components
	resp.OK = true

	result, _ = json.Marshal(resp)
	return
}

// CreateComponent is create new component with given info
func CreateComponent(ctx *macaron.Context) (httpStatus int, result []byte) {
	var resp CreateComponentResp

	namespace := ctx.ParamsEscape(":namespace")
	if namespace == "" {
		httpStatus = http.StatusBadRequest
		resp.OK = false
		resp.ErrorCode = ComponentError + ComponentEmptyNamespaceError
		resp.Message = "got empty namespace"

		result, _ = json.Marshal(resp)
		return
	}

	body, err := ctx.Req.Body().Bytes()
	if err != nil {
		httpStatus = http.StatusBadRequest
		resp.OK = false
		resp.ErrorCode = ComponentError + ComponentReqBodyError
		resp.Message = "Get requrest body error: " + err.Error()

		result, _ = json.Marshal(resp)
		return
	}

	var id int64
	if id, err = module.CreateComponent(namespace, body); err != nil {
		httpStatus = http.StatusBadRequest
		resp.OK = false
		resp.ErrorCode = ComponentError + ComponentCreateError
		resp.Message = "Create component error: " + err.Error()

		result, err = json.Marshal(resp)
		if err != nil {
			log.Errorln("Create component marshal data error: " + err.Error())
		}
		return
	}

	httpStatus = http.StatusCreated
	resp.ComponentInfo.ID = id
	resp.OK = true
	resp.Message = "Component created"

	result, _ = json.Marshal(resp)
	return
}

// GetComponent is get an component detail with given info
func GetComponent(ctx *macaron.Context) (httpStatus int, result []byte) {
	var resp ComponentDetailResp

	namespace := ctx.ParamsEscape(":namespace")
	if namespace == "" {
		httpStatus = http.StatusBadRequest
		resp.OK = false
		resp.ErrorCode = ComponentError + ComponentEmptyNamespaceError
		resp.Message = "got empty namespace"

		result, _ = json.Marshal(resp)
		return
	}

	componentID := ctx.ParamsInt64(":component")

	component, err := module.GetComponentByID(namespace, componentID)
	if err != nil {
		httpStatus = http.StatusBadRequest
		resp.OK = false
		resp.ErrorCode = ComponentError + ComponentGetError
		resp.Message = "get component detail error: " + err.Error()

		result, _ = json.Marshal(resp)
		return
	}

	if component == nil {
		httpStatus = http.StatusNotFound
		resp.OK = false
		resp.ErrorCode = ComponentError + ComponentGetError
		resp.Message = "component not found"

		result, _ = json.Marshal(resp)
		return
	}

	httpStatus = http.StatusOK
	resp.OK = true
	resp.ComponentData = component
	result, _ = json.Marshal(resp)
	return
}

// UpdateComponent is update an component's info with given info
func UpdateComponent(ctx *macaron.Context) (httpStatus int, result []byte) {
	var resp CommonResp

	namespace := ctx.ParamsEscape(":namespace")
	if namespace == "" {
		httpStatus = http.StatusBadRequest
		resp.OK = false
		resp.ErrorCode = ComponentError + ComponentEmptyNamespaceError
		resp.Message = "got empty namespace"

		result, _ = json.Marshal(resp)
		return
	}

	componentID := ctx.ParamsInt64(":component")

	body, err := ctx.Req.Body().Bytes()
	if err != nil {
		httpStatus = http.StatusBadRequest
		resp.OK = false
		resp.ErrorCode = ComponentError + ComponentReqBodyError
		resp.Message = "Get requrest body error: " + err.Error()

		result, _ = json.Marshal(resp)
		return
	}

	if err = module.UpdateComponent(namespace, componentID, body); err != nil {
		httpStatus = http.StatusBadRequest
		resp.OK = false
		resp.ErrorCode = ComponentError + ComponentUpdateError
		resp.Message = "update component error: " + err.Error()

		result, err = json.Marshal(resp)
		if err != nil {
			log.Errorln("update component marshal data error: " + err.Error())
		}
		return
	}

	httpStatus = http.StatusOK
	resp.OK = true
	resp.Message = "update success"

	result, _ = json.Marshal(resp)
	return
}

// DeleteComponent is is delete an selected component
func DeleteComponent(ctx *macaron.Context) (httpStatus int, result []byte) {
	var resp ComponentDetailResp

	namespace := ctx.ParamsEscape(":namespace")
	if namespace == "" {
		httpStatus = http.StatusBadRequest
		resp.OK = false
		resp.ErrorCode = ComponentError + ComponentEmptyNamespaceError
		resp.Message = "got empty namespace"

		result, _ = json.Marshal(resp)
		return
	}

	componentID := ctx.ParamsInt64(":component")

	err := module.DeleteComponentByID(namespace, componentID)
	if err != nil {
		httpStatus = http.StatusBadRequest
		resp.OK = false
		resp.ErrorCode = ComponentError + ComponentGetError
		resp.Message = "get component detail error: " + err.Error()

		result, _ = json.Marshal(resp)
		return
	}

	httpStatus = http.StatusOK
	resp.OK = true
	result, _ = json.Marshal(resp)
	return
}
