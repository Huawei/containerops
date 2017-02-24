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
	"github.com/Huawei/containerops/component/module"
)

// CommonResp is common resp that each resp will have
type CommonResp struct {
	OK        bool    `json:"ok"`
	ErrorCode errCode `json:"error_code,omitempty"`
	Message   string  `json:"message",omitempty`
}

// ListComponentsResp is get component list's struct
type ListComponentsResp struct {
	CommonResp `json:"common"`
	Components []module.ComponentBaseData `json:"components"`
}

type CreateComponentResp struct {
	CommonResp    `json:"common"`
	ComponentInfo struct {
		ID int64 `json:"id"`
	} `json:"component"`
}

// ComponentDetailResp is get component detail's resp struct
type ComponentDetailResp struct {
	*module.ComponentData `json:"component,omitempty"`
	CommonResp            `json:"common"`
}

// Env is env define struct
type Env struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
