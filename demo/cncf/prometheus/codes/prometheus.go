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

const (
	//FailuerExit exit code is -1.
	FailuerExit = -1
	//MissingParameter exit code is -2.
	MissingParameter = -2
	//ParseEnvFailure exit code is -3.
	ParseEnvFailure = -3
	//CloneError exit code is -4
	CloneError = -4
	//UnknownAction exit code is -5
	UnknownAction = -5
)

//Parse CO_DATA value, and return Kubernetes repository URI and action (build/test/release).
func parseEnv(env string) (uri string, action string, err error) {
	files := strings.Fields(env)
	if len(files) == 0 {
		return "", "", fmt.Errorf("CO_DATA value is null\n")
	}

	for _, v := range files {
		s := strings.Split(v, "=")
		key, value := s[0], s[1]

		switch key {
		case "prometheus":
			uri = value
		case "action":
			action = value
		default:
			fmt.Fprintf(os.Stdout, "[COUT] Unknown Parameter: [%s]\n", s)
		}
	}

	return uri, action, nil
}

//Git clone the prometheus repository, and process will redirect to system stdout.
func gitClone(repo, dest string) error {
	cmd := exec.Command("git", "clone", repo, dest)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "[COUT] Git clone error: %s\n", err.Error())
		fmt.Fprintf(os.Stdout, "[COUT] CO_RESULT = %s\n", "false")
		os.Exit(FailuerExit)
	}

	return nil
}

//Execute `make build` in the prometheus folder.
func prometheusBuild() error {
	cmd := exec.Command("make", "build")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "[COUT] Make test error: %s\n", err.Error())
		fmt.Fprintf(os.Stdout, "[COUT] CO_RESULT = %s\n", "false")

		return err
	}

	return nil
}

//Execute `make test` in the prometheus folder.
func prometheusTest() error {
	cmd := exec.Command("make", "test")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "[COUT] Make test error: %s\n", err.Error())
		fmt.Fprintf(os.Stdout, "[COUT] CO_RESULT = %s\n", "false")

		return err
	}

	return nil
}

//TODO Build the Prometheus all binary files, and release to containerops repository.
func release(dest string) error {
	fmt.Fprintf(os.Stderr, "[COUT] %s", "No release function in the demo component yet.")
	fmt.Fprintf(os.Stdout, "[COUT] CO_RESULT = %s\n", "false")

	return fmt.Errorf("No release function yet.")
}

func main() {
	//Get the CO_DATA from environment parameter "CO_DATA"
	data := os.Getenv("CO_DATA")
	if len(data) == 0 {
		fmt.Fprintf(os.Stderr, "[COUT] %s\n", "The CO_DATA value is null.")
		fmt.Fprintf(os.Stdout, "[COUT] CO_RESULT = %s\n", "The CO_DATA value is null.")
		os.Exit(MissingParameter)
	}

	if prometheusRepo, action, err := parseEnv(data); err != nil {
		fmt.Fprintf(os.Stderr, "[COUT] Parse the CO_DATA error: %s\n", err.Error())
		fmt.Fprintf(os.Stdout, "[COUT] CO_RESULT = %s\n", "false")
		os.Exit(ParseEnvFailure)
	} else {
		//Create the base path within GOPATH.
		basePath := path.Join(os.Getenv("GOPATH"), "src", "github.com", "prometheus", "prometheus")
		os.MkdirAll(basePath, os.ModePerm)

		//Clone the git repository
		if err := gitClone(prometheusRepo, basePath); err != nil {
			fmt.Fprintf(os.Stderr, "[COUT] Clone the prometheus repository error: %s\n", err.Error())
			fmt.Fprintf(os.Stdout, "[COUT] CO_RESULT = %s\n", "false")
			os.Exit(CloneError)
		}

		//Execute action
		switch action {
		case "build":

			if err := prometheusBuild(); err != nil {
				os.Exit(FailuerExit)
			}

		case "test":

			if err := prometheusTest(); err != nil {
				os.Exit(FailuerExit)
			}

		case "release":

			if err := release(basePath); err != nil {
				os.Exit(FailuerExit)
			}

		default:
			fmt.Fprintf(os.Stderr, "[COUT] %s\n", "Unknown action, the component only support build, test and publish action.")
			fmt.Fprintf(os.Stdout, "[COUT] CO_RESULT = %s\n", "false")
			os.Exit(UnknownAction)
		}

	}

	//Print result
	fmt.Fprintf(os.Stdout, "[COUT] CO_RESULT = %s\n", "true")
	os.Exit(0)
}
