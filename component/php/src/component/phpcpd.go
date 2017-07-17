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
	"util/git"
	"util/input"
	"util/file"
	"util/cmd"
)

const (
	basePath string = "./workspace"
	// baseCommand string = "phpcpd"
	baseCommand string = "/home/composer/.composer/vendor/bin/phpcpd"
	reportPath string = "/tmp/PMD-CPD.xml"
	reportFormat string = "XML_REPORT"
)

func main() {
	data := os.Getenv("CO_DATA")
	keys := []string{
		"git-url",
		"path",
		"exclude",
		"names",
		"names-exclude",
		"count-tests",
		"regexps-exclude",
		"min-lines",
		"min-tokens",
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

	command := baseCommand

	if codata["path"] == "" {
		codata["path"] = "."
	}
	command = fmt.Sprintf("%s %s", command, codata["path"])

	params := []string{
		"names",
		"names-exclude",
		"regexps-exclude",
		"exclude",
		"min-lines",
		"min-tokens",
	}

	for _, param := range params {
		if codata[param] != "" {
			command = fmt.Sprintf("%s --%s=%s", command, param, codata[param])
		}
	}

	command = fmt.Sprintf("%s --log-pmd=%s", command, reportPath)

	if err := cmd.RunCommand(command, basePath); err != nil {
		fmt.Fprintf(os.Stdout, "[COUT] CO_RESULT = %s\n", "false")
		os.Exit(1)
	}

	file.StdoutAll(reportPath, reportFormat)

	fmt.Fprintf(os.Stdout, "[COUT] CO_RESULT = %s\n", "true")
	os.Exit(0)
}