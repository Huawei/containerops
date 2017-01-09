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

package module

import (
	"testing"

	"github.com/Huawei/containerops/pilotage/module"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetWorkflows(t *testing.T) {
	Convey("Get workflow list", t, func() {
		namespace := "demo"
		repository := "demo"
		name := "testName"
		nameFuzzy := "false"
		version := "testVersion"
		versionFuzzy := "false"
		offset := int64(0)
		pageNum := int64(10)
		versionNum := int64(3)

		Convey("empty param", func() {
			Convey("With empty namespace", func() {
				namespace = ""
				result, err := module.GetWorkflows(namespace, repository, name, nameFuzzy, version, versionFuzzy, offset, pageNum, versionNum)

				So(err.Error(), ShouldEqual, "namespace or repository can't be empty")
				So(result, ShouldBeEmpty)
			})

			Convey("With empty repository", func() {
				repository = ""
				result, err := module.GetWorkflows(namespace, repository, name, nameFuzzy, version, versionFuzzy, offset, pageNum, versionNum)

				So(err.Error(), ShouldEqual, "namespace or repository can't be empty")
				So(result, ShouldBeEmpty)
			})
		})

		Convey("normal", func() {
			Convey("get workflow list", func() {
				name = ""
				nameFuzzy = "true"
				version = ""
				versionFuzzy = "true"

				result, err := module.GetWorkflows(namespace, repository, name, nameFuzzy, version, versionFuzzy, offset, pageNum, versionNum)

				So(err, ShouldBeNil)
				So(len(result), ShouldBeGreaterThanOrEqualTo, 0)
			})

			Convey("get more workflow", func() {
				name = ""
				nameFuzzy = "true"
				version = ""
				versionFuzzy = "true"
				offset = int64(10)

				result, err := module.GetWorkflows(namespace, repository, name, nameFuzzy, version, versionFuzzy, offset, pageNum, versionNum)

				So(err, ShouldBeNil)
				So(len(result), ShouldEqual, 0)
			})

			Convey("get more version", func() {
				name = ""
				nameFuzzy = "false"
				version = ""
				versionFuzzy = "true"
				offset = int64(3)

				result, err := module.GetWorkflows(namespace, repository, name, nameFuzzy, version, versionFuzzy, offset, pageNum, versionNum)

				So(err, ShouldBeNil)
				So(len(result), ShouldEqual, 0)
			})
		})

		Convey("search", func() {
			Convey("search by workflow name", func() {
				name = ""
				nameFuzzy = "false"
				version = ""
				versionFuzzy = "true"

				result, err := module.GetWorkflows(namespace, repository, name, nameFuzzy, version, versionFuzzy, offset, pageNum, versionNum)

				So(err, ShouldBeNil)
				So(len(result), ShouldEqual, 0)
			})

			Convey("search by version", func() {
				name = ""
				nameFuzzy = "false"
				version = ""
				versionFuzzy = "true"

				result, err := module.GetWorkflows(namespace, repository, name, nameFuzzy, version, versionFuzzy, offset, pageNum, versionNum)

				So(err, ShouldBeNil)
				So(len(result), ShouldEqual, 0)
			})

			Convey("get more workflow after search", func() {
				name = ""
				nameFuzzy = "false"
				version = ""
				versionFuzzy = "true"

				result, err := module.GetWorkflows(namespace, repository, name, nameFuzzy, version, versionFuzzy, offset, pageNum, versionNum)

				So(err, ShouldBeNil)
				So(len(result), ShouldEqual, 0)
			})

			Convey("get more version after search", func() {
				name = ""
				nameFuzzy = "false"
				version = ""
				versionFuzzy = "true"

				result, err := module.GetWorkflows(namespace, repository, name, nameFuzzy, version, versionFuzzy, offset, pageNum, versionNum)

				So(err, ShouldBeNil)
				So(len(result), ShouldEqual, 0)
			})
		})

		Reset(func() {
			namespace = "demo"
			repository = "demo"
			name = "testName"
			nameFuzzy = "false"
			version = "testVersion"
			versionFuzzy = "false"
			offset = int64(0)
			pageNum = int64(10)
			versionNum = int64(3)
		})
	})
}
