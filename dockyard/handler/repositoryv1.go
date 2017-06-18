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

	log "github.com/Sirupsen/logrus"
	"gopkg.in/macaron.v1"

	"fmt"
	"github.com/Huawei/containerops/dockyard/model"
	"github.com/Huawei/containerops/dockyard/module"
)

// PostRepositoryV1Handler is creating a repository in Dockyard.
// The repository has different types:
//   1. docker -> Docker V2 repository -> model: DockerV2
//   2. binary -> Binary V1 repository -> model: BinaryV1
// Notes:
//   1. When developer push a docker image, it will be create a Docker V2 repository in the database automatically.
//   2. When developer push a binary file, he will be crate a Binary repository first using API/UI/CLI
// API Specification:
//   POST /v1/:namespace/:repository/:type
//   Create a repository.
//   Parameters:
//     :namespace -> username or organization name
//     :repository -> repository name
//     :type ->
//       1. docker -> Docker V2 repository
//       2. binary -> Binary V1 repository
//   Return:
//      201 -> Creation successful
//          	-> {
//               	"namespace" : "genedna",
//               	"repository" : "dockyard",
//               	"type" : "docker"
//             	}
//      400 -> Creation failed
//          	-> {
//								"errors": [{
//                   "code": "",
//                   "message": "",
//                   "description": ""
//                 }]
//               }
//      401 -> Authentication Failed
//            -> {
//								"errors": [{
//                   "code": "AUTHENTICATION_FAILED",
//                   "message": "authentication failed",
//                   "description": ""
//                 }]
//               }
func PostRepositoryV1Handler(ctx *macaron.Context) (int, []byte) {
	namespace := ctx.Params(":namespace")
	repository := ctx.Params(":repository")
	repoType := ctx.Params(":type")

	switch repoType {
	case "docker":
		r := new(model.DockerV2)

		if err := r.Put(namespace, repository); err != nil {
			log.Errorf("Put Docker V2 repository error: %s", err.Error())

			result, _ := module.EncodingError(module.REPOSITORY_CREATE_FAILED, fmt.Sprintf("%s/%s", namespace, repository))
			return http.StatusBadRequest, result
		}
	case "binary":
		b := new(model.BinaryV1)

		if err := b.Put(namespace, repository); err != nil {
			log.Errorf("Put Binary V2 repository error: %s", err.Error())

			result, _ := module.EncodingError(module.REPOSITORY_CREATE_FAILED, fmt.Sprintf("%s/%s", namespace, repository))
			return http.StatusBadRequest, result
		}

	default:
		log.Errorf("Unknown repository type: %s", repoType)
		result, _ := module.EncodingError(module.REPOSITORY_CREATE_FAILED, fmt.Sprintf("Unknown repository type: %s", repoType))
		return http.StatusBadRequest, result
	}

	result, _ := json.Marshal(map[string]string{})
	return http.StatusCreated, result
}
