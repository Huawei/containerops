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

package middleware

import (
	"gopkg.in/macaron.v1"
)

//setRespHeaders is set Resp header value.
//TODO: Add a config option for provide Docker Registry V1.
func setRespHeaders() macaron.Handler {
	return func(ctx *macaron.Context) {
		//Docker Registry V2
		ctx.Resp.Header().Set("Docker-Distribution-Api-Version", "registry/2.0")
	}
}
