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

package module

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Run
func (f *Flow) Run(verbose bool) error {

	return nil
}

// ExecuteFlowFromFile
func (f *Flow) ExecuteFlowFromFile(flowFile string, verbose bool) error {
	f.Model = RunModelCli

	if data, err := ioutil.ReadFile(flowFile); err != nil {
		fmt.Println("Read ", flowFile, " error: ", err.Error())
		return err
	} else {
		if err := yaml.Unmarshal(data, &f); err != nil {
			fmt.Println("Unmarshal the flow file error:", err.Error())
			return err
		} else {
			f.Run(verbose)
		}
	}

	return nil
}
