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

package download

import (
	"fmt"

	cmd "github.com/Huawei/containerops/singular/cmd"
	"github.com/Huawei/containerops/singular/init_config"
)

func Master_Download(ip string) {

	var fileslist = init_config.Get_files()
	for key, url := range fileslist {
		cmd.ExecCMDparams("wget", []string{"-P", "/tmp/", "-c", url})
		fmt.Printf("%s\n\n", key)
	}
	//scp -r ./config/. root@138.68.14.193:/tmp/
	cmd.LocalExecCMDparams("scp", []string{"-r", "./config/.", init_config.User + "@" + ip + ":" + "/tmp/config/"})

}
