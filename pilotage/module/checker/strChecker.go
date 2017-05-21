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

package checker

import (
	"errors"
	"net/http"
	"strings"
)

type strChecker struct {
}

func (checker *strChecker) Support(eventMap map[string]string) bool {
	if eventMap["sourceType"] == "gitlab" || eventMap["sourceType"] == "github" {
		return true
	}

	return false
}

func (checker *strChecker) Check(eventMap map[string]string, expectedToken string, reqHeader http.Header, reqBody []byte) (bool, error) {
	passCheck := false

	if eventMap["sourceType"] != "gitlab" && eventMap["sourceType"] != "github" {
		passCheck = false
		return passCheck, errors.New("strChecker doesn't support type:" + eventMap["sourceType"])
	}

	switch eventMap["eventName"] {
	case "Merge Request Hook":
		if strings.Contains(string(reqBody), `"action":"open"`) {
			passCheck = true
		}
	case "pull_request":
		if strings.Contains(string(reqBody), `"action":"opened"`) {
			passCheck = true
		}
	}

	if passCheck {
		return passCheck, nil
	}

	return passCheck, errors.New("strChecker failed")
}
