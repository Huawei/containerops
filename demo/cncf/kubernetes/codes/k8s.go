/*
Copyright 2014 - 2017 Huawei Technologies Co., Ltd. All rights reserved.

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
	FailuerExit      = -1
	MissingParamater = -2
	ParseEnvFailure  = -3
	CLONE_ERROR      = -4
	UNKNOWN_ACTION   = -5
)

//Parse CO_DATA value, and return Kubernetes repository URI and action (build/test/publish).
func parse_env(env string) (uri string, action string, err error) {
	files := strings.Fields(env)
	if len(files) == 0 {
		return "", "", fmt.Errorf("CO_DATA value is null\n")
	}

	for _, v := range files {
		s := strings.Split(v, "=")
		key, value := s[0], s[1]

		switch key {
		case "kubernetes":
			uri = value
		case "action":
			action = value
		default:
			fmt.Fprintf(os.Stdout, "[COUT] Unknown Parameter: [%s]\n", s)
		}
	}

	return uri, action, nil
}

//Git clone the kubernetes repository, and process will redirect to system stdout.
func git_clone(repo, dest string) error {
	cmd := exec.Command("git", "clone", repo, dest)
	cmd.Path = dest
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "[COUT] Git clone error: %s\n", err.Error())
		fmt.Fprintf(os.Stdout, "[COUT] CO_RESULT = false\n")
		os.Exit(FailuerExit)
	}

	/*
		if _, err := git.PlainClone(dest, false, &git.CloneOptions{
			URL:               repo,
			RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
			Progress:          os.Stdout,
		}); err != nil {
			return err
		}
	*/
	return nil
}

//make bazel-test
func bazel_test(dest string) {
	cmd := exec.Command("make", "bazel-test")
	cmd.Path = dest
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "[COUT] Bazel test error: %s\n", err.Error())
		fmt.Fprintf(os.Stdout, "[COUT] CO_RESULT = false\n")
		os.Exit(FailuerExit)
	}

}

//`make bazel-build`
func bazel_build(dest string) {
	cmd := exec.Command("make", "bazel-build")
	cmd.Path = dest
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "[COUT] Bazel build error: %s\n", err.Error())
		fmt.Fprintf(os.Stdout, "[COUT] CO_RESULT = false\n")
		os.Exit(FailuerExit)
	}
}

//TODO Build the kubernetes all binrary files, and publish to containerops repository. And not execute the `make bazel-publish` command.
func publish(dest string) {

}

func main() {
	//Get the CO_DATA from environment parameter "CO_DATA"
	co_data := os.Getenv("CO_DATA")
	if len(co_data) == 0 {
		fmt.Fprintf(os.Stderr, "[COUT] The CO_DATA value is null.\n")
		fmt.Fprintf(os.Stdout, "[COUT] CO_RESULT = false\n")
		os.Exit(MissingParamater)
	}

	//Parse the CO_DATA, get the kubernetes repository URI and action
	if k8s_repo, action, err := parse_env(co_data); err != nil {
		fmt.Fprintf(os.Stderr, "[COUT] Parse the CO_DATA error: %s\n", err.Error())
		fmt.Fprintf(os.Stdout, "[COUT] CO_RESULT = false\n")
		os.Exit(ParseEnvFailure)
	} else {
		//Create the base path within GOPATH.
		base_path := path.Join(os.Getenv("GOPATH"), "src", "github.com", "kubernetes")
		os.MkdirAll(base_path, os.ModePerm)

		//Clone the git repository
		if err := git_clone(k8s_repo, base_path); err != nil {
			fmt.Fprintf(os.Stderr, "[COUT] Clone the kubernetes repository error: %s\n", err.Error())
			fmt.Fprintf(os.Stdout, "[COUT] CO_RESULT = false\n")
			os.Exit(CLONE_ERROR)
		}

		//Execute action
		switch action {
		case "build":
			bazel_build(path.Join(base_path, "kubernetes"))
		case "test":
			bazel_test(path.Join(base_path, "kubernetes"))
		case "publish":
			publish(base_path)
		default:
			fmt.Fprintf(os.Stderr, "[COUT] Unknown action, the component only support build, test and publish action.\n")
			fmt.Fprintf(os.Stdout, "[COUT] CO_RESULT = false\n")
			os.Exit(UNKNOWN_ACTION)
		}

	}

	//Print result
	fmt.Fprintf(os.Stdout, "[COUT] CO_RESULT = true\n")
	os.Exit(0)
}
