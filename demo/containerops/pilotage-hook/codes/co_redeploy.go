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
	"fmt"
	"os"
	"strings"

	common "github.com/Huawei/containerops/common/utils"
)

func main() {

	data := os.Getenv("CO_DATA")
	if len(data) == 0 {
		fmt.Fprintf(os.Stderr, "[COUT] %s\n", "The CO_DATA value is null.")
		fmt.Fprintf(os.Stdout, "[COUT] CO_RESULT = %s\n", "false")
		os.Exit(1)
	}

	target, url, key, err := parseEnv(data)

	err = update(target, url, key)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[COUT] Failed to update service %s: %s\n", target, err.Error())
		fmt.Fprintf(os.Stdout, "[COUT] CO_RESULT = false\n")
		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, "[COUT] CO_RESULT=true\n")
}

func parseEnv(env string) (target, url, sshKey string, err error) {
	files := strings.Fields(env)
	if len(files) == 0 {
		err = fmt.Errorf("CO_DATA value is null\n")
		return
	}

	for _, v := range files {
		s := strings.Split(v, "=")
		key, value := s[0], s[1]

		switch key {
		case "target":
			target = value
		case "url":
			url = value
		case "key":
			sshKey = value
		default:
			fmt.Fprintf(os.Stdout, "[COUT] Unknown Parameter: [%s]\n", s)
		}
	}
	return
}

func update(target, url, key string) error {
	cmd := fmt.Sprintf("/var/containerops/scripts/%s/deploy.sh '%s'", target, url)
	return common.SSHCommand("root", key, "hub.opshub.sh", 22, cmd, os.Stdout, os.Stderr)
}
