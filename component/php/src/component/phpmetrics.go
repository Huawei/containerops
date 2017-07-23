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
	"util/git"
	"util/input"
	"util/file"
)

const (
	basePath string = "./workspace"
	// baseCommand string = "phpmetrics"
	baseCommand string = "/home/composer/.composer/vendor/bin/phpmetrics"
	reportPath string = "/tmp/phpmetrics.xml"
	reportFormat string = "REPORT"
)

func main() {
	// Get the CO_DATA from environment parameter "CO_DATA"
	data := os.Getenv("CO_DATA")
	keys := []string{
		"git-url",
		"path",
		"exclude",
	}
	codata := map[string]string{}

	if err := input.HandleInput(data, keys, codata); err != nil {
		fmt.Fprintf(os.Stderr, "[COUT] Handle input error: %s\n", err.Error())
		fmt.Fprintf(os.Stdout, "[COUT] CO_RESULT = %s\n", "false")
		os.Exit(1)
	}

	if err := git.Clone(codata["git-url"], basePath); err != nil {
		fmt.Fprintf(os.Stderr, "[COUT] Clone the repository error: %s\n", err.Error())
		fmt.Fprintf(os.Stdout, "[COUT] CO_RESULT = %s\n", "false")
		os.Exit(1)
	}

	if codata["path"] == "" {
		codata["path"] = "."
	}
	exclude := fmt.Sprintf("--exclude=%v", codata["exclude"])

	cmd := exec.Command(baseCommand, codata["path"], exclude,"--report-violations=/tmp/report.xml")
	cmd.Dir = basePath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "[COUT] Create report error: %s\n", err.Error())
		fmt.Fprintf(os.Stdout, "[COUT] CO_RESULT = %s\n", "false")
		os.Exit(1)
	}

	file.StdoutAll(reportPath, reportFormat)

	fmt.Fprintf(os.Stdout, "[COUT] CO_RESULT = %s\n", "true")
	os.Exit(0)
}