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
	"os/exec"
	"path"
	"strings"
)

//Parse CO_DATA value, and return lintcode repository URI ).
func parseEnv(env string) (uri string, err error) {
	files := strings.Fields(env)
	if len(files) == 0 {
		return "", fmt.Errorf("CO_DATA value is null\n")
	}

	for _, v := range files {
		s := strings.Split(v, "=")
		key, value := s[0], s[1]

		switch key {
		case "coderepo":
			uri = value
		default:
			fmt.Fprintf(os.Stdout, "[COUT] Unknown Parameter: [%s]\n", s)
		}
	}

	return uri, nil
}

//Git clone the code repository, and process will redirect to system stdout.
func gitClone(repo, dest string) error {
	cmd := exec.Command("git", "clone", repo, dest)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "[COUT] Git clone error: %s\n", err.Error())
		fmt.Fprintf(os.Stdout, "[COUT] CO_RESULT = %s\n", "false")
		os.Exit(1)
	}

	return nil
}


func main() {
	//Get the CO_DATA from environment parameter "CO_DATA"
	data := os.Getenv("CO_DATA")
	if len(data) == 0 {
		fmt.Fprintf(os.Stderr, "[COUT] %s\n", "The CO_DATA value is null.")
		fmt.Fprintf(os.Stdout, "[COUT] CO_RESULT = %s\n", "false")
		os.Exit(1)
	}

	//Parse the CO_DATA, get the lintcode repository URI and action
	if codeRepo, err := parseEnv(data); err != nil {
		fmt.Fprintf(os.Stderr, "[COUT] Parse the CO_DATA error: %s\n", err.Error())
		fmt.Fprintf(os.Stdout, "[COUT] CO_RESULT = %s\n", "false")
		os.Exit(1)
	} else {
		//Create the base path within GOPATH.
		basePath := path.Join(os.Getenv("GOPATH"), "src","tmp")

		//Clone the git repository
		if err := gitClone(codeRepo, basePath); err != nil {
			fmt.Fprintf(os.Stderr, "[COUT] Clone the code repository error: %s\n", err.Error())
			fmt.Fprintf(os.Stdout, "[COUT] CO_RESULT = %s\n", "false")
			os.Exit(1)
		}

		//lint code path
		lintPath := path.Join(os.Getenv("GOPATH"), "src","tmp","...")
		fmt.Fprintf(os.Stdout, "[COUT] Lint codeRepo(s%) and report：\n", codeRepo)

		//Execute golint
		cmd := exec.Command("golint",lintPath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "[COUT] golint error: %s\n", err.Error())
			fmt.Fprintf(os.Stdout, "[COUT] CO_RESULT = %s\n", "false")
			os.Exit(1)
		}
	}

	//Print result
	fmt.Fprintf(os.Stdout, "[COUT] CO_RESULT = %s\n", "true")
	os.Exit(0)
}

// Author: tanhaijun@huawei.com
// 2017-7-11 Create
