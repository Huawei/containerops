/*
Copyright 2014 - 2017 Huawei Technologies Co., Ltd. All rights reserved.

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
	// AuthTypeWorkflowDefault is default auth type that workflow doesn't spec a target
	AuthTypeWorkflowDefault = "default"
	// AuthTypeWorkflowStartAPI is auth that given by api request
	AuthTypeWorkflowStartAPI = "WorkflowStartAPI"
	// AuthTypeWorkflowStartDone is auth that given by system when workflow start auth is all done
	AuthTypeWorkflowStartDone = "WorkflowStartDone"

	// AuthTypeStageDefault is default auth type that workflow doesn't spec a target
	AuthTypeStageDefault = "default"
	// AuthTypePreStageDone is auth that given by system when pre stage run success
	AuthTypePreStageDone = "PreStageDone"
	// AuthTyptStageStartDone is auth that given by system when a stage get all auth it request
	AuthTyptStageStartDone = "StageStartDone"

	// AuthTokenDefault is default auth authorizer that workflow doesn't spec a authorizer
	AuthTokenDefault = "default"
)
