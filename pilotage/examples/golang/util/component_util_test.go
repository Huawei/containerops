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

package util

import (
	"io/ioutil"
	"log"
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_EventNotify(t *testing.T) {
	go startSimpaleServer()

	Convey("test notify event", t, func() {
		Convey("test notify custom event", func() {
			eventURLMap["custom"] = "http://localhost:9999/"
			eventIDMap["custom"] = 100

			err := NotifyEvent("custom", true, "this is a result", "this is a output")
			So(err, ShouldBeNil)
		})
	})

	Convey("test notify event", t, func() {

		Convey("test notify event with all param is set", func() {
			eventURLMap["custom"] = "http://localhost:9999/"
			eventURLMap[CO_COMPONENT_START] = "http://localhost:9999/"
			eventURLMap[CO_COMPONENT_STOP] = "http://localhost:9999/"
			eventURLMap[CO_TASK_START] = "http://localhost:9999/"
			eventURLMap[CO_TASK_RESULT] = "http://localhost:9999/"
			eventURLMap[CO_TASK_STATUS] = "http://localhost:9999/"

			eventIDMap["custom"] = 100
			eventIDMap[CO_COMPONENT_START] = 101
			eventIDMap[CO_COMPONENT_STOP] = 102
			eventIDMap[CO_TASK_START] = 103
			eventIDMap[CO_TASK_RESULT] = 104
			eventIDMap[CO_TASK_STATUS] = 105

			Convey("test notify custom event", func() {
				err := NotifyEvent("custom", true, "this is a result", "this is a output")
				So(err, ShouldBeNil)
			})

			Convey("test notify component start", func() {
				err := ComponentStart("this is a info")
				So(err, ShouldBeNil)
			})

			Convey("test notify component stop", func() {
				err := ComponentStop("this is a info")
				So(err, ShouldBeNil)
			})

			Convey("test notify task start", func() {
				err := TaskStart("this is a info")
				So(err, ShouldBeNil)
			})

			Convey("test notify task result", func() {
				err := TaskResult("this is a info")
				So(err, ShouldBeNil)
			})

			Convey("test notify task status", func() {
				err := TaskStatus(true, "this is a info", "this is a output")
				So(err, ShouldBeNil)
			})

		})

		Convey("test notify event without eventURL", func() {

			Convey("test notify custom event", func() {
				err := NotifyEvent("custom", true, "this is a result", "this is a output")
				So(err.Error(), ShouldEqual, "event is illegal")
			})

			Convey("test notify component start", func() {
				err := ComponentStart("this is a info")
				So(err.Error(), ShouldEqual, "event is illegal")
			})

			Convey("test notify component stop", func() {
				err := ComponentStop("this is a info")
				So(err.Error(), ShouldEqual, "event is illegal")
			})

			Convey("test notify task start", func() {
				err := TaskStart("this is a info")
				So(err.Error(), ShouldEqual, "event is illegal")
			})

			Convey("test notify task result", func() {
				err := TaskResult("this is a info")
				So(err.Error(), ShouldEqual, "event is illegal")
			})

			Convey("test notify task status", func() {
				err := TaskStatus(true, "this is a info", "this is a output")
				So(err.Error(), ShouldEqual, "event is illegal")
			})

		})

		Convey("test notify event without eventID", func() {

			Convey("test notify custom event", func() {
				err := NotifyEvent("custom", true, "this is a result", "this is a output")
				So(err.Error(), ShouldEqual, "event is illegal")
			})

			Convey("test notify component start", func() {
				err := ComponentStart("this is a info")
				So(err.Error(), ShouldEqual, "event is illegal")
			})

			Convey("test notify component stop", func() {
				err := ComponentStop("this is a info")
				So(err.Error(), ShouldEqual, "event is illegal")
			})

			Convey("test notify task start", func() {
				err := TaskStart("this is a info")
				So(err.Error(), ShouldEqual, "event is illegal")
			})

			Convey("test notify task result", func() {
				err := TaskResult("this is a info")
				So(err.Error(), ShouldEqual, "event is illegal")
			})

			Convey("test notify task status", func() {
				err := TaskStatus(true, "this is a info", "this is a output")
				So(err.Error(), ShouldEqual, "event is illegal")
			})

		})

		Reset(func() {
			delete(eventURLMap, "custom")
			delete(eventURLMap, CO_COMPONENT_START)
			delete(eventURLMap, CO_COMPONENT_STOP)
			delete(eventURLMap, CO_TASK_START)
			delete(eventURLMap, CO_TASK_RESULT)
			delete(eventURLMap, CO_TASK_STATUS)

			delete(eventIDMap, "custom")
			delete(eventIDMap, CO_COMPONENT_START)
			delete(eventIDMap, CO_COMPONENT_STOP)
			delete(eventIDMap, CO_TASK_START)
			delete(eventIDMap, CO_TASK_RESULT)
			delete(eventIDMap, CO_TASK_STATUS)
		})
	})

	Convey("test get data", t, func() {

		Convey("test get data from env", func() {
			Convey("get data from env with env is set", func() {
				eventINFOMap[CO_DATA] = `{"data":"this is a data json"}`

				dataChan := make(chan map[string]interface{}, 1)
				GetData(9090, false, dataChan)
			})

			Convey("get data from env with env is not set", func() {

			})

			Reset(func() {
				delete(eventINFOMap, CO_DATA)
			})
		})

		Convey("test get data from workflow", func() {

		})

	})

}

func startSimpaleServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("error when get data body:" + err.Error())
		}

		log.Println("coming request's body:", string(body))

		w.Write(body)
	})

	http.ListenAndServe(":9999", nil)
}

func sendDATA() {

}
