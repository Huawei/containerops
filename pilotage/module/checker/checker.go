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

package checker

import (
	"errors"
	"net/http"

	"github.com/containerops/configure"
)

// Checker is
type Checker interface {
	Support(eventMap map[string]string) bool
	Check(eventMap map[string]string, expectedToken string, reqHeader http.Header, reqBody []byte) (bool, error)
}

// GetWorkflowExecCheckerList is
func GetWorkflowExecCheckerList() ([]Checker, error) {
	checkers := make([]Checker, 0)
	checkerNameList := configure.GetStringSlice("auth.checker")

	for _, checkerName := range checkerNameList {
		checker, err := getChecker(checkerName)
		if err != nil {
			return nil, err
		}
		checkers = append(checkers, checker)
	}

	return checkers, nil
}

func getChecker(checkName string) (Checker, error) {
	switch checkName {
	case "tokenChecker":
		return new(tokenChecker), nil
	case "strChecker":
		return new(strChecker), nil
	}
	return nil, errors.New("error when getChecker by checkerName:" + checkName + ", checker not found")
}
