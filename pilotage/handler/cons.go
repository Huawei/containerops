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

type errCode uint64

const (
	ComponentError  errCode = 00010000
	workflowErrCode errCode = 00100000
	stageErrCode    errCode = 01000000
	actionErrCode   errCode = 10000000
)

const (
	_                       = iota
	ComponentReqBodyError   //errCode = 0001
	ComponentUnmarshalError //errCode = 0002
	ComponentCreateError    //errCode = 0003
	ComponentParseIDError   //errCode = 0004
	ComponentGetError       //errCode = 0005
	ComponentUpdateError    //errCode = 0006
	ComponentDeleteError    //errCode = 0007
	ComponentDebugError     //errCode = 0008
	ComponentmarshalError
	ComponentListError
)
