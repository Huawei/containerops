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

package module

const (
	// AuthTypePipelineDefault is default auth type that pipeline doesn't spec a target
	AuthTypePipelineDefault = "default"
	// AuthTypePipelineStartAPI is auth that given by api request
	AuthTypePipelineStartAPI = "PipelineStartAPI"
	// AuthTypePipelineStartDone is auth that given by system when pipeline start auth is all done
	AuthTypePipelineStartDone = "PipelineStartDone"

	// AuthTypePipelineDefault is default auth type that pipeline doesn't spec a target
	AuthTypeStageDefault = "default"
	// AuthTypePreStageDone is auth that given by system when pre stage run success
	AuthTypePreStageDone = "PreStageDone"
	// AuthTyptStageStartDon is auth that given by system when a stage get all auth it request
	AuthTyptStageStartDone = "StageStartDone"

	// AuthAuthorizerDefault is default auth authorizer that pipeline doesn't spec a authorizer
	AuthTokenDefault = "default"
)
