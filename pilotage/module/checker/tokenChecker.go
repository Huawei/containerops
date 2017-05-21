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
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"net/http"
)

type tokenChecker struct {
}

func (checker *tokenChecker) Support(eventMap map[string]string) bool {
	return true
}

func (checker *tokenChecker) Check(eventMap map[string]string, expectedToken string, reqHeader http.Header, reqBody []byte) (bool, error) {
	passCheck := false

	switch eventMap["sourceType"] {
	case "customize":
		if expectedToken == eventMap["token"] {
			passCheck = true
		}
	case "gitlab":
		if expectedToken == eventMap["token"] {
			passCheck = true
		}
	case "github":
		mac := hmac.New(sha1.New, []byte(expectedToken))
		mac.Write(reqBody)
		expectedMAC := mac.Sum(nil)
		expectedSig := "sha1=" + hex.EncodeToString(expectedMAC)

		if expectedSig == eventMap["token"] {
			passCheck = true
		}
	}

	if passCheck {
		return passCheck, nil
	}

	return passCheck, errors.New("token checker failed")
}
