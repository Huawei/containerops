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

// GetFlowJobLog is return log of a Job
func GetFlowJobLog(ctx *macaron.Context) (int, []byte) {
	result, _ := json.Marshal(map[string]string{})
	return http.StatusOK, result
}
