
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
	"os/exec"
)

type CO_DATA struct {
    gitUrl string
	action string
}

func main() {
	// Get the CO_DATA from environment parameter "CO_DATA"
	data := os.Getenv("CO_DATA")
	if len(data) == 0 {
		fmt.Fprintf(os.Stderr, "[COUT] %s\n", "The CO_DATA value is null.")
		fmt.Fprintf(os.Stdout, "[COUT] CO_RESULT = %s\n", "false")
		os.Exit(1)
	}

	// Handle with CO_DATA
	codata, err := handleCO_DATA(data)
	if err != nil {
		os.Exit(1)
	}

	basePath := "./projects"

	if err := gitClone(codata.gitUrl, basePath); err != nil {
		fmt.Fprintf(os.Stderr, "[COUT] Clone the kubernetes repository error: %s\n", err.Error())
		fmt.Fprintf(os.Stdout, "[COUT] CO_RESULT = %s\n", "false")
		os.Exit(1)
	}

	switch codata.action {
		case "install": 
			cmd := exec.Command("composer", "install")
			cmd.Dir = basePath
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			if err := cmd.Run(); err != nil {
				fmt.Fprintf(os.Stderr, "[COUT] Composer install error: %s\n", err.Error())
				fmt.Fprintf(os.Stdout, "[COUT] CO_RESULT = %s\n", "false")
				os.Exit(1)
			}
		default:
			fmt.Fprintf(os.Stderr, "[COUT] No such action: %s\n", codata.action)
			fmt.Fprintf(os.Stdout, "[COUT] CO_RESULT = %s\n", "false")
			os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, "[COUT] CO_RESULT = %s\n", "true")
	os.Exit(0)
}

func handleCO_DATA(data string) (codata CO_DATA, err error) {
	files := strings.Fields(data)
	if len(files) == 0 {
		return codata, fmt.Errorf("CO_DATA value is null\n")
	}

	for _, v := range files {
		s := strings.Split(v, "=")
		key, value := s[0], s[1]

		switch key {
		case "git-url": 
			codata.gitUrl = value
		case "action":
			codata.action = value
		default:
			fmt.Fprintf(os.Stdout, "[COUT] Unknown Parameter: [%s]\n", key)
		}
	}

	return codata, nil;
}

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