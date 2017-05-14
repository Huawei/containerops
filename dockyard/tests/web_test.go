/*
Copyright 2014 Huawei Technologies Co., Ltd. All rights reserved.

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

package tests

import (
	"fmt"
	"net/http"
	"os"
	"testing"
)

var (
	testServer = ""
)

func init() {
	// Start a dockyard web server and set the enviornment like this:
	//     $ export US_TEST_SERVER=https://containerops.me
	testServer = os.Getenv("US_TEST_SERVER")
}

func Test_IndexHandler(t *testing.T) {
	if testServer == "" {
		fmt.Println("Skip index handler testing since 'US_TEST_SERVER' is not set")
		return
	}
	endpoint := testServer

	resp, err := http.Get(endpoint)
	if err != nil {
		t.Errorf("Test REST API \"/\" Error: %s .", err.Error())
		return
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Test REST API \"/\" Error StatusCode: %d .", resp.StatusCode)
	} else {
		t.Log("Test REST API \"/\" Successfully.")
	}
}
