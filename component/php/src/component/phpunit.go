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
	// baseCommand string = "phpunit"
	baseCommand string = "/home/composer/.composer/vendor/bin/phpunit"
	composerCommand string = "/usr/local/bin/composer"
	// coveragePath string = "/tmp/coverage.xml"
	// coverageFormat string = "COVERAGE REPORT"
	reportPath string = "/tmp/report.xml"
	reportFormat string = "TEST REPORT"
)

func main() {
	data := os.Getenv("CO_DATA")
	keys := []string{
		"git-url",
		"composer",
		"bootstrap",
		"include-path",
		"configuration",
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

	if codata["composer"] == "true" {
		composer_command := fmt.Sprintf("%s %s", composerCommand, "install")

		if err := cmd.RunCommand(composer_command, basePath); err != nil {
			fmt.Fprintf(os.Stderr, "[COUT] Composer install dependence error: %s\n", err.Error())
			fmt.Fprintf(os.Stdout, "[COUT] CO_RESULT = %s\n", "false")
			os.Exit(1)
		}
	}

	command := baseCommand

	params := []string{
		"bootstrap",
		"include-path",
		"configuration",
	}

	for _, param := range params {
		if codata[param] != "" {
			command = fmt.Sprintf("%s --%s %s", command, param, codata[param])
		}
	}

	// command = fmt.Sprintf("%s --coverage-clover %s", command, coveragePath)
	command = fmt.Sprintf("%s --log-junit %s", command, reportPath)

	if err := cmd.RunCommand(command, basePath); err != nil {
		fmt.Fprintf(os.Stdout, "[COUT] CO_RESULT = %s\n", "false")
		os.Exit(1)
	}

	file.StdoutAll(reportPath, reportFormat)
	// file.StdoutAll(coveragePath, coverageFormat)

	fmt.Fprintf(os.Stdout, "[COUT] CO_RESULT = %s\n", "true")
	os.Exit(0)
}