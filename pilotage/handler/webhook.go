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
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Huawei/containerops/pilotage/config"
	"github.com/Huawei/containerops/pilotage/module"
	log "github.com/Sirupsen/logrus"
	macaron "gopkg.in/macaron.v1"
)

func WebHook(ctx *macaron.Context) (int, []byte) {
	// Run the flow:
	// 1. Check if the singular changes
	// 2. (if changed) Build new singular and push the binary to dockyard
	// 3. ssh into the target server, update the singular service.
	// TODO Parameter validation
	namespace := ctx.Params("namespace")
	repository := ctx.Params("repository")
	tag := ctx.Params("tag")
	flowName := ctx.Params("flow")

	url := fmt.Sprintf("%s/flow/v1/%s/%s/%s/%s/%s", config.WebHook.Host, namespace, repository, flowName, tag, "yaml")

	flowYamlPath := fmt.Sprintf("%s/%s/%s/%s/flow.yaml", config.WebHook.FlowBaseDir, namespace, repository, tag)
	f := &module.Flow{}
	if err := f.ParseFlowFromFile(flowYamlPath, module.DaemonRun, false, true); err != nil {
		log.Error(err)
		return http.StatusInternalServerError, []byte("Failed to parse flow yaml file")
	}

	ns, repo, name, err := f.URIs()
	if err != nil {
		log.Error(err)
		return http.StatusInternalServerError, []byte("Failed to validate flow URI")
	}
	if ns != namespace || repo != repository || name != flowName {
		log.Error(err)
		return http.StatusInternalServerError, []byte("Invalid flow URI")
	}

	yamlBytes, err := f.YAML()
	if err != nil {
		log.Error(err)
		return http.StatusInternalServerError, []byte("Failed to yamlify flow")
	}
	bytesReader := bytes.NewReader(yamlBytes)

	client := http.Client{}
	req, _ := http.NewRequest(http.MethodPost, url, bytesReader)
	res, err := client.Do(req)
	if err != nil {
		log.Error(err)
		return http.StatusInternalServerError, []byte("Failed to request flow creation")
	}
	defer res.Body.Close()

	bs, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error(err)
		return http.StatusInternalServerError, []byte("Failed to get response from pilotage flow creation API")
	}

	return res.StatusCode, bs
}
