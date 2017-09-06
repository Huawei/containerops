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
	"fmt"
	"net/http"

	"gopkg.in/macaron.v1"
	"gopkg.in/yaml.v2"

	"time"

	"github.com/Huawei/containerops/pilotage/module"
)

// GetFlowRuntime is return flow runtime status and logs.
func GetFlowRuntime(ctx *macaron.Context) (int, []byte) {
	switch ctx.Data["mode"] {
	case module.DaemonRun:
		if f := ctx.Data["flow"].(*module.Flow); f != nil {
			t := ctx.Params("type")

			switch t {
			case "json":
				if data, err := f.JSON(); err != nil {
					result, _ := json.Marshal(map[string]string{
						"message": fmt.Sprintf("Get flow JSON definition error: %s", err.Error())})
					return http.StatusBadRequest, result
				} else {
					return http.StatusOK, data
				}
			case "yaml":
				if data, err := f.YAML(); err != nil {
					result, _ := json.Marshal(map[string]string{
						"message": fmt.Sprintf("Get flow YAML definition error: %s", err.Error())})
					return http.StatusBadRequest, result
				} else {
					return http.StatusOK, data
				}
			default:
				result, _ := json.Marshal(map[string]string{
					"message": fmt.Sprintf("Unsupport definition type: %s", ctx.Data["mode"])})
				return http.StatusBadRequest, result
			}
		} else {
			result, _ := json.Marshal(map[string]string{"message": "No flow data"})
			return http.StatusBadRequest, result
		}
	case module.DaemonStart:
		// TODO
	default:
		result, _ := json.Marshal(map[string]string{
			"message": fmt.Sprintf("Unsupport engine run mode: %s", ctx.Data["mode"])})
		return http.StatusBadRequest, result
	}

	result, _ := json.Marshal(map[string]string{})
	return http.StatusOK, result
}

type PostFlowResponse struct {
	ID         string `json:"id"`
	Namespace  string `json:"namespace"`
	Repository string `json:"repository"`
	Name       string `json:"name"`
	Tag        string `json:"tag"`
	Title      string `json:"title"`
	Version    int64  `json:"version"`
	Status     string `json:"status"`
}

func PostFlowRuntime(ctx *macaron.Context) (int, []byte) {
	namespace := ctx.Params("namespace")
	repository := ctx.Params("repository")
	flowName := ctx.Params("flow")
	data, _ := ctx.Req.Body().Bytes()


	f := module.Flow{Number: 1, Status: module.Pending}
	switch ctx.Params("type") {
	case "json":
		if err := json.Unmarshal(data, &f); err != nil {
			info:=fmt.Sprintf("Unmarshal the flow file error: %s", err.Error())
			f.Log(info, true, true)
			result, _ := json.Marshal(map[string]string{"message": info})
			return http.StatusBadRequest, result
		}

	case "yaml":
		if err := yaml.Unmarshal(data, &f); err != nil {
			info:=fmt.Sprintf("Unmarshal the flow file error: %s", err.Error())
			f.Log(info, true, true)
			result, _ := json.Marshal(map[string]string{"message": info})
			return http.StatusBadRequest, result
		}
	default:
		result, _ := json.Marshal(map[string]string{
			"message": fmt.Sprintf("Unsupport type: %s", ctx.Params("type"))})
		return http.StatusBadRequest, result
	}

	go func() {
		f.LocalRun(true, true)
	}()
	// Sleep one second to wait the init status change of flow
	time.Sleep(1 * time.Second)
	resp := PostFlowResponse{Namespace: namespace, Repository: repository, Name: flowName, Tag: f.Tag,
		Version: f.Version, Title: f.Title, Status: f.Status}
	result, _ := json.Marshal(resp)
	return http.StatusCreated, result
}

// GetFlowJobLog is return log of a Job
func GetFlowJobLog(ctx *macaron.Context) (int, []byte) {
	result, _ := json.Marshal(map[string]string{})
	return http.StatusOK, result
}
