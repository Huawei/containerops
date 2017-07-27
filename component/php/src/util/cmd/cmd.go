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

package cmd

import (
	"os"
	"os/exec"
	"fmt"
	"strings"
)

func RunCommand(command string, basePath string) error {
	list := strings.Split(command, " ")
	cmd := exec.Command(list[0], list[1:]...)

	cmd.Dir = basePath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "[COUT] Run command error: %s\n[COUT] Error msg: %s\n", command, err.Error())
		
		return err
	}

	return nil;
}