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
	"html/template"
	"net/http"
	"strconv"

	"github.com/Huawei/containerops/singular/controller"
	"github.com/cloudflare/cfssl/log"

	"gopkg.in/macaron.v1"
)

func GetIndexPageV1Handler(ctx *macaron.Context) {
	funcs := template.FuncMap{
		// "component_names": controller.StringifyComponentsNames,
		"component_names": func(args ...interface{}) (string, error) {
			return "hello", nil
		},
		"inc": func(i int) int {
			return i + 1
		},
	}

	deployments, err := controller.GetHtmlDeploymentList()
	if err != nil {
		log.Errorf("Failed to get deployment list: %s", err.Error())
		ctx.Resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	// The template's name should be the same with file name
	listTmpl, err := template.New("list.html").Funcs(funcs).ParseFiles("./templates/list.html")
	if err != nil {
		log.Error(err)
		ctx.Resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	renderData := map[string]interface{}{
		"infra_titles": controller.InfraTitles,
		"deployments":  deployments,
	}
	err = listTmpl.Execute(ctx.Resp, renderData)
	if err != nil {
		log.Error(err)
		return
	}
}

func GetDetailPageV1Handler(ctx *macaron.Context) {
	// Get the deployment information
	namespace := ctx.Params("namespace")
	repository := ctx.Params("repository")
	name := ctx.Params("name")
	tag := ctx.Params("tag")
	versionStr := ctx.Params("version")
	version, _ := strconv.Atoi(versionStr)

	deployment := controller.GetHtmlDeploymentDetail(namespace, repository, name, tag, int64(version))
	if deployment == nil {
		ctx.Resp.WriteHeader(http.StatusNotFound)
		ctx.Resp.Write([]byte("Deployment not found"))
		return
	}

	// The template's name should be the same with file name
	listTmpl, err := template.New("detail.html").ParseFiles("./templates/detail.html")
	if err != nil {
		log.Error(err)
		ctx.Resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = listTmpl.Execute(ctx.Resp, deployment)
	if err != nil {
		log.Error(err)
		return
	}
}
