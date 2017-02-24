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

// errCode is
type errCode uint64

const (
	// ComponentError is error code for Component
	ComponentError errCode = 00010000

	// EventError is error code for Event
	EventError errCode = 00020000

	// WorkflowError is error code for Workflow
	WorkflowError errCode = 00100000

	// StageError is error code for Stage
	StageError errCode = 01000000

	// ActionError is error code for Action
	ActionError errCode = 10000000
)

const (
	_ = iota
	// ComponentReqBodyError is
	ComponentReqBodyError
	// ComponentUnmarshalError is
	ComponentUnmarshalError
	// ComponentCreateError is
	ComponentCreateError
	// ComponentParseIDError is
	ComponentParseIDError
	// ComponentGetError is
	ComponentGetError
	// ComponentUpdateError is
	ComponentUpdateError
	// ComponentDeleteError is
	ComponentDeleteError
	// ComponentDebugError is
	ComponentDebugError
	// ComponentmarshalError is
	ComponentmarshalError
	// ComponentListError is
	ComponentListError
	// ComponentEmptyNamespaceError is
	ComponentEmptyNamespaceError
)

const (
	_ = iota

	// EventReqBodyError is
	EventReqBodyError

	// EventUnmarshalError is
	EventUnmarshalError

	// EventIllegalDataError is
	EventIllegalDataError

	// EventGetActionError is
	EventGetActionError
)
