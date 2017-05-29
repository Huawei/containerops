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
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
)

var (
	eventMap map[string]EventInfo
	dataChan chan string

	// POD_NAME is
	POD_NAME string
	// RUN_ID is
	RUN_ID string
	// SERVICE_ADDR is
	SERVICE_ADDR string
)

// EventInfo is
type EventInfo struct {
	id  string
	url string
}

func init() {
	eventMap = make(map[string]EventInfo)
	dataChan = make(chan string, 1)
}

func main() {
	initEvent()

	// when start a component ,send a COMPONENT_START notify and register self in workflow
	notifyEvent("COMPONENT_START", "application/json", strings.NewReader("component start..."))

	// on this func ,you can wait workflow send the data you will use
	for _, serviceInfo := range strings.Split(SERVICE_ADDR, ",") {
		if len(strings.Split(serviceInfo, ":")) >= 3 {
			go waitForData(strings.Split(serviceInfo, ":"))
		}
	}

	notifyEvent("TASK_STATE", "application/json", strings.NewReader("waitting data..."))
	codePath := <-dataChan

	// after get all data you need, send a TASH_START notify to workflow and start do what you need
	notifyEvent("TASK_START", "application/json", strings.NewReader("task start..."))

	result := generateGodoc(codePath)

	resultMap := make(map[string]interface{}, 0)
	resultMap["output"] = map[string]string{"doc": result}
	resultMap["status"] = true
	resultMap["result"] = map[string]interface{}{"status": resultMap["status"], "output": resultMap["output"]}
	resultByte, _ := json.Marshal(resultMap)
	// after task done send a TASK_RESULT notify to workflow
	notifyEvent("TASK_RESULT", "application/json", bytes.NewReader(resultByte))

	// on this func, you can do some after-task job,like recycle some resource
	// recycleResource()

	// after after-task job done, send a COMPONENT_STOP notify to workflow ,so that workflow will stop it as soon as possible
	// if you don't send a COMPONENT_STOP notify,pipelin will auto stop it after some time(default:60s)
	notifyEvent("COMPONENT_STOP", "application/json", strings.NewReader("component stoping..."))
}

func initEvent() {
	eventList := os.Getenv("EVENT_LIST")
	RUN_ID = os.Getenv("RUN_ID")
	POD_NAME = os.Getenv("POD_NAME")
	SERVICE_ADDR = os.Getenv("SERVICE_ADDR")

	for _, eventInfo := range strings.Split(eventList, ";") {
		if len(strings.Split(eventInfo, ",")) > 1 {
			eventKey := strings.Split(eventInfo, ",")[0]
			eventId := strings.Split(eventInfo, ",")[1]
			if os.Getenv(eventKey) != "" {
				event := new(EventInfo)
				event.id = eventId
				event.url = os.Getenv(eventKey)
				eventMap[eventKey] = *event
			}
		}
	}
}

func notifyEvent(eventKey, bodyType string, body io.Reader) {
	if event, ok := eventMap[eventKey]; ok {
		eventId := event.id
		notifyUrl := event.url
		if strings.Contains(notifyUrl, "?") {
			notifyUrl += "&runId=" + RUN_ID
		} else {
			notifyUrl += "?runId=" + RUN_ID
		}
		notifyUrl += "&event=" + eventKey + "&eventId=" + eventId

		client := &http.Client{}
		req, err := http.NewRequest(http.MethodPut, notifyUrl, body)
		if err != nil {
			fmt.Println("error when create a notify request:" + err.Error())
		}

		resp, err := client.Do(req)

		if err != nil {
			fmt.Println("error when notify event", eventKey, body)
		}

		respBody, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		fmt.Println("notify resp is", string(respBody))
	}
}

func waitForData(serviceInfo []string) {
	// send a message to pipelint to notify that current component is ready to receive data
	registerUrl := eventMap["REGISTER_URL"].url
	if strings.Contains(registerUrl, "?") {
		registerUrl += "&runId=" + RUN_ID
	} else {
		registerUrl += "?runId=" + RUN_ID
	}

	serviceAddr := ":" + serviceInfo[1]

	registerUrl += "&podName=" + POD_NAME
	registerUrl += "&receiveUrl=" + url.QueryEscape(serviceAddr+"/receivedata")

	client := &http.Client{}
	req, _ := http.NewRequest("PUT", registerUrl, nil)

	client.Do(req)

	http.HandleFunc("/receivedata", receiveDataHandler)
	http.ListenAndServe(":"+serviceInfo[2], nil)
}

func receiveDataHandler(w http.ResponseWriter, r *http.Request) {
	result, _ := json.Marshal(map[string]string{"message": "ok"})

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("error when get data body:" + err.Error())
	}

	codePathMap := make(map[string]string)
	json.Unmarshal([]byte(body), &codePathMap)
	dataChan <- codePathMap["path"]

	w.Write(result)
}

func generateGodoc(codePath string) string {
	cmd := exec.Command("godoc", codePath)
	var result bytes.Buffer
	cmd.Stdout = &result
	err := cmd.Run()
	if err != nil {
		return err.Error()
	}
	return result.String()
}
