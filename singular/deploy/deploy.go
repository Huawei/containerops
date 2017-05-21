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

package deploy

import (
	"fmt"

	"github.com/Huawei/containerops/singular/init_config"
)

func DeployNodes() {
	var nodelist = init_config.Get_nodes()
	for k, ip := range nodelist {
		init_config.TargetIP = ip
		fmt.Printf("k=%v, v=%v\n", k, ip)
		if k == init_config.Master_name {
			Deploymaster(nodelist, ip)
		}
		if k == init_config.Minion_name {
			Deploynode(nodelist, ip)
		}
	}
}
