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

package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Huawei/containerops/pilotage/handler"

	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/macaron.v1"
)

func TestListWorkflowsV1(t *testing.T) {

	Convey("Get workflow list", t, func() {

		m := macaron.New()

		m.Get("/:namespace/:repository", handler.ListWorkflowsV1)

		Convey("With namespace is nil", func() {

			resp := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "/ /test", nil)
			So(err, ShouldBeNil)

			m.ServeHTTP(resp, req)

			result, _ := json.Marshal(map[string]string{"errMsg": "namespace or repository can't be empty"})

			So(resp.Code, ShouldEqual, http.StatusBadRequest)
			So(resp.Body.String(), ShouldEqual, string(result))

		})

		Convey("With repository is nil", func() {

			resp := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "/test/ ", nil)
			So(err, ShouldBeNil)

			m.ServeHTTP(resp, req)

			result, _ := json.Marshal(map[string]string{"errMsg": "namespace or repository can't be empty"})

			So(resp.Code, ShouldEqual, http.StatusBadRequest)
			So(resp.Body.String(), ShouldEqual, string(result))
		})
	})
}
