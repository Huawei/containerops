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

package main

import (
	"github.com/Huawei/containerops/singular/cmd"
	"github.com/Huawei/containerops/singular/deploy"
	"github.com/Huawei/containerops/singular/download"
	"github.com/Huawei/containerops/singular/vm"
)

func main() {

	cmd.Execute()
	// create vmlist
	vm.CreateVMs()
	//vm.CreateNewVM("lidian-unbantu-wk-master")
	//get iplist while
	download.Download_main()
	deploy.DeployNodes()

}
