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
	// baseCommand string = "phploc"
	baseCommand string = "/home/composer/.composer/vendor/bin/phploc"
	reportPath string = "/tmp/report.xml"
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

	if codata["exclude"] != "" {
		command = fmt.Sprintf("%s --exclude=%s", command, codata["exclude"])
	}

	if codata["names"] != "" {
		command = fmt.Sprintf("%s --names=%v", command, codata["names"])
	}

	if codata["names-exclude"] != "" {
		command = fmt.Sprintf("%s --names-exclude=%v", command, codata["names-exclude"])
	}
	
	if codata["count_tests"] != "true" {
		command = fmt.Sprintf("%s --count-tests", command)
	}

	command = fmt.Sprintf("%s --log-xml=%s", command, reportPath)

	if err := cmd.RunCommand(command, basePath); err != nil {
		fmt.Fprintf(os.Stdout, "[COUT] CO_RESULT = %s\n", "false")
		os.Exit(1)
	}

	file.StdoutAll(reportPath)

	fmt.Fprintf(os.Stdout, "[COUT] CO_RESULT = %s\n", "true")
	os.Exit(0)
}