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
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Huawei/containerops/pilotage/models"
	"github.com/Huawei/containerops/pilotage/utils"

	log "github.com/Sirupsen/logrus"
	"github.com/containerops/configure"
)

var (
	startPipelineChan  chan bool
	createPipelineChan chan bool
)

const (
	PIPELINE_STAGE_TYPE_START = "pipeline-start"
	PIPELINE_STAGE_TYPE_RUN   = "pipeline-stage"
	PIPELINE_STAGE_TYPE_ADD   = "pipeline-add-stage"
	PIPELINE_STAGE_TYPE_END   = "pipeline-end"
)

type Relation struct {
	From string
	To   string
}

func init() {
	startPipelineChan = make(chan bool, 1)
	createPipelineChan = make(chan bool, 1)
}

func PipeExecRequestLegal(reqHeader http.Header, reqBody []byte, pipelineInfo models.Pipeline) bool {
	sourceList := make([]map[string]string, 0)

	err := json.Unmarshal([]byte(pipelineInfo.SourceInfo), &sourceList)
	if err != nil {
		log.Error("error when unmarshal pipelin source config" + err.Error())
		return false
	}

	legal := false
	for _, sourceConfig := range sourceList {
		token := reqHeader.Get(sourceConfig["headerKey"])
		if token != "" {
			eventType := getEventType(sourceConfig["sourceType"], reqHeader)
			if !strings.Contains(sourceConfig["eventList"], ","+eventType+",") {
				continue
			}

			legal = checkToken(sourceConfig["sourceType"], sourceConfig["secretKey"], reqHeader.Get(sourceConfig["headerKey"]), reqHeader, reqBody)

			if legal {
				break
			}
		}
	}

	// return legal
	return true
}

func DoPipelineLog(pipelineInfo models.Pipeline) (*models.PipelineLog, error) {
	pipelineLog := new(models.PipelineLog)
	// use chan to make sure pipeline sequence is unique
	startPipelineChan <- true

	tempSequence := new(struct {
		Sequence int64
	})

	err := new(models.Outcome).GetOutcome().Table("outcome").Where("real_pipeline = ?", pipelineInfo.ID).Order("-sequence").First(&tempSequence).Error
	if err != nil && err.Error() == "record not found" {
		tempSequence.Sequence = 0
	} else if err != nil {
		<-startPipelineChan
		return nil, errors.New("error when query outcome info by pipeline:" + err.Error())
	}

	pipelineSequence := tempSequence.Sequence + 1

	pipelineLog.Namespace = pipelineInfo.Namespace
	pipelineLog.Pipeline = pipelineInfo.Pipeline
	pipelineLog.FromPipeline = pipelineInfo.ID
	pipelineLog.Event = pipelineInfo.Event
	pipelineLog.Version = pipelineInfo.Version
	pipelineLog.VersionCode = pipelineInfo.VersionCode
	pipelineLog.State = pipelineInfo.State
	pipelineLog.Manifest = pipelineInfo.Manifest
	pipelineLog.Description = pipelineInfo.Description
	pipelineLog.SourceInfo = pipelineInfo.SourceInfo
	pipelineLog.Env = pipelineInfo.Env
	pipelineLog.Sequence = pipelineSequence

	err = pipelineLog.GetPipelineLog().Save(pipelineLog).Error
	if err != nil {
		<-startPipelineChan
		return nil, errors.New("error when create new pipeline log:" + err.Error())
	}

	<-startPipelineChan

	// copy all current pipeline's stage infos to stage log
	stageList := make([]models.Stage, 0)
	stageIdMap := make(map[int64]int64)
	err = new(models.Stage).GetStage().Where("pipeline = ?", pipelineInfo.ID).Find(&stageList).Error
	if err != nil {
		return nil, errors.New("error when get stage infos by pipeline id:" + strconv.FormatInt(pipelineInfo.ID, 10))
	}

	preStage := int64(-1)
	for _, stageInfo := range stageList {
		stageLog := new(models.StageLog)

		stageLog.Pipeline = pipelineLog.ID
		stageLog.Type = stageInfo.Type
		stageLog.PreStage = preStage
		stageLog.Stage = stageInfo.Stage
		stageLog.FromStage = stageInfo.ID
		stageLog.Title = stageInfo.Title
		stageLog.Description = stageInfo.Description
		stageLog.Event = stageInfo.Event
		stageLog.Manifest = stageInfo.Manifest
		stageLog.Env = stageInfo.Env
		stageLog.Timeout = stageInfo.Timeout

		err = stageLog.GetStageLog().Save(stageLog).Error
		if err != nil {
			return nil, errors.New("error when create new stage log" + err.Error())
		}

		preStage = stageLog.ID
		stageIdMap[stageInfo.ID] = stageLog.ID
	}

	// copy all action infos to action log
	for _, stageInfo := range stageList {
		actionList := make([]models.Action, 0)
		err = new(models.Action).GetAction().Where("stage = ?", stageInfo.ID).Find(&actionList).Error
		if err != nil {
			return nil, errors.New("error when get action infos by stage id:" + strconv.FormatInt(stageInfo.ID, 10))
		}

		for _, actionInfo := range actionList {
			actionLog := new(models.ActionLog)

			actionLog.Stage = stageIdMap[actionInfo.Stage]
			actionLog.Component = actionInfo.Component
			actionLog.Service = actionInfo.Service
			actionLog.Action = actionInfo.Action
			actionLog.FromAction = actionInfo.ID
			actionLog.Title = actionInfo.Title
			actionLog.Description = actionInfo.Description
			actionLog.Event = actionInfo.Event
			actionLog.Manifest = actionInfo.Manifest
			actionLog.Environment = actionInfo.Environment
			actionLog.Kubernetes = actionInfo.Kubernetes
			actionLog.Swarm = actionInfo.Swarm
			actionLog.Endpoint = actionInfo.Endpoint
			actionLog.Timeout = actionInfo.Timeout
			actionLog.Input = actionInfo.Input
			actionLog.Output = actionInfo.Output

			err = actionLog.GetActionLog().Save(actionLog).Error
			if err != nil {
				return nil, errors.New("error when create new action log:" + err.Error())
			}

			protocol := ""
			listenMode := configure.GetString("listenmode")
			switch listenMode {
			case "http":
				protocol = "http://"
				break
			case "https":
				protocol = "https://"
				break
			default:
				protocol = "https://"
				break
			}

			projectAddr := ""
			if configure.GetString("projectaddr") == "" {
				projectAddr = "localhost"
			} else {
				projectAddr = configure.GetString("projectaddr")
			}
			projectAddr = strings.TrimSuffix(projectAddr, "/")

			// add default event to actionlog
			eventList := []map[string]string{
				{"event": "COMPONENT_START", "value": protocol + projectAddr + "/pipeline/v1/" + pipelineInfo.Namespace + "/demo/" + pipelineInfo.Pipeline + "/event"},
				{"event": "COMPONENT_STOP", "value": protocol + projectAddr + "/pipeline/v1/" + pipelineInfo.Namespace + "/demo/" + pipelineInfo.Pipeline + "/event"},
				{"event": "TASK_START", "value": protocol + projectAddr + "/pipeline/v1/" + pipelineInfo.Namespace + "/demo/" + pipelineInfo.Pipeline + "/event"},
				{"event": "TASK_RESULT", "value": protocol + projectAddr + "/pipeline/v1/" + pipelineInfo.Namespace + "/demo/" + pipelineInfo.Pipeline + "/event"},
				{"event": "TASK_STATUS", "value": protocol + projectAddr + "/pipeline/v1/" + pipelineInfo.Namespace + "/demo/" + pipelineInfo.Pipeline + "/event"},
				{"event": "REGISTER_URL", "value": protocol + projectAddr + "/pipeline/v1/" + pipelineInfo.Namespace + "/demo/" + pipelineInfo.Pipeline + "/register"}}

			for _, event := range eventList {
				tempEvent := new(models.EventDefinition)
				tempEvent.Event = event["event"]
				tempEvent.Title = event["event"]
				tempEvent.Definition = event["value"]
				tempEvent.Action = actionLog.ID

				tempEvent.GetEventDefinition().Save(tempEvent)
			}

		}

	}

	return pipelineLog, nil
}

func getEventType(sourceType string, reqHeader http.Header) string {
	eventType := ""
	switch sourceType {
	case "Github":
		eventType = reqHeader.Get("X-Github-Event")
	}

	return eventType
}

func checkToken(sourceType, secretKey, token string, reqHeader http.Header, reqBody []byte) bool {
	legal := false

	switch sourceType {
	case "Github":
		mac := hmac.New(sha1.New, []byte(secretKey))
		mac.Write(reqBody)
		expectedMAC := mac.Sum(nil)
		expectedSig := "sha1=" + hex.EncodeToString(expectedMAC)

		if expectedSig == token {
			legal = true
		}
		break
	case "customize":
		if token == secretKey {
			legal = true
		}
		break
	}

	return legal
}

// start a pipeline
// user channel to make sure pipelineSequence is unique
func StartPipeline(pipelineInfo models.PipelineLog, reqBody string) string {

	// get start stage of current pipeline
	startStage := new(models.StageLog)
	err := startStage.GetStageLog().Where("pipeline = ?", pipelineInfo.ID).Where("pre_stage = ?", -1).First(&startStage).Error
	if err != nil {
		return "error when query pipeline start stage info:" + err.Error()
	}

	if startStage.ID == 0 {
		return "can't find start stage info for pipelineId:" + strconv.FormatInt(pipelineInfo.ID, 10)
	}

	pipelineSequence := pipelineInfo.Sequence
	// record pipeline start data
	startOutcome := new(models.Outcome)

	startOutcome.Sequence = pipelineSequence
	startOutcome.Pipeline = pipelineInfo.ID
	startOutcome.RealPipeline = pipelineInfo.FromPipeline
	startOutcome.Stage = startStage.ID
	startOutcome.RealStage = startStage.FromStage
	startOutcome.Status = false
	startOutcome.Result = "success"
	startOutcome.Output = reqBody

	err = startOutcome.GetOutcome().Save(startOutcome).Error
	if err != nil {
		return "error when save pipeline start data:" + err.Error()
	}

	envMap := make(map[string]string)
	if pipelineInfo.Env != "" {
		err = json.Unmarshal([]byte(pipelineInfo.Env), &envMap)
		if err != nil {
			return "error when unmarshal pipeline env info" + err.Error()
		}
	}

	go handleStage(pipelineInfo, *startStage, pipelineSequence, envMap)
	return "pipeline start ..."
}

// handler a stage to start a stage's all action and auto start next stage until all stage done or error
func handleStage(pipelineInfo models.PipelineLog, stageInfo models.StageLog, pipelineSequence int64, pipelineEnvMap map[string]string) {
	nextStage := new(models.StageLog)
	actionList := make([]models.ActionLog, 0)

	// if current stage is end stage, this pipeline run result is success, record success, stop run
	if stageInfo.Type == models.StageTypeEnd {
		log.Info("current stage is ", stageInfo, " now pipeline("+strconv.FormatInt(pipelineSequence, 10)+") is finish...")
		finalOutcome := new(models.Outcome)

		finalOutcome.Pipeline = pipelineInfo.ID
		finalOutcome.RealPipeline = pipelineInfo.FromPipeline
		finalOutcome.Stage = stageInfo.ID
		finalOutcome.RealStage = stageInfo.FromStage
		finalOutcome.Action = -1
		finalOutcome.RealAction = -1
		finalOutcome.Sequence = pipelineSequence
		finalOutcome.Status = true

		finalOutcome.GetOutcome().Save(finalOutcome)
	}

	// if current stage is nil , stop run
	if stageInfo.ID == 0 {
		log.Info("current stage is ", stageInfo, " now pipeline("+strconv.FormatInt(pipelineSequence, 10)+") is finish...")
		finalOutcome := new(models.Outcome)

		finalOutcome.Pipeline = pipelineInfo.ID
		finalOutcome.RealPipeline = pipelineInfo.FromPipeline
		finalOutcome.Stage = stageInfo.ID
		finalOutcome.RealStage = stageInfo.FromStage
		finalOutcome.Action = -1
		finalOutcome.RealAction = -1
		finalOutcome.Sequence = pipelineSequence
		finalOutcome.Status = false

		finalOutcome.GetOutcome().Save(finalOutcome)
		return
	}

	err := nextStage.GetStageLog().Where("pipeline = ?", pipelineInfo.ID).Where("pre_stage = ?", stageInfo.ID).First(&nextStage).Error
	if err != nil {
		log.Error("error when get nextStage info from db :" + err.Error())
		return
	}

	// set stage set env to stageEnvMap
	stageEnvMap := pipelineEnvMap
	if stageInfo.Env != "" {
		err := json.Unmarshal([]byte(stageInfo.Env), &stageEnvMap)
		if err != nil {
			log.Error("stage's env define is not a json obj:" + err.Error())
			return
		}
	}

	// get all action
	new(models.ActionLog).GetActionLog().Where("stage = ?", stageInfo.ID).Find(&actionList)

	// if current stage has action,start all action
	if stageInfo.PreStage != -1 && len(actionList) > 0 {
		// exec all action
		for _, action := range actionList {
			go execAction(pipelineInfo, stageInfo, action, pipelineSequence, stageEnvMap)
		}
	}

	if len(actionList) < 1 {
		// if stage don't have any action, start next stage
		handleStage(pipelineInfo, *nextStage, pipelineSequence, pipelineEnvMap)
	} else {
		stageTimeoutDuration, err := time.ParseDuration(strconv.FormatInt(stageInfo.Timeout, 10) + "s")
		if err != nil {
			log.Error("error when set stage" + stageInfo.Stage + "'s timeout:" + err.Error() + "set stage timer to default 36 hours")
			stageTimeoutDuration, _ = time.ParseDuration("36h")
		}
		actionResultChan := make(chan bool, 1)
		go waitAllActionFinish(actionList, pipelineSequence, actionResultChan)
		select {
		case <-time.After(stageTimeoutDuration):
			log.Error("stage " + stageInfo.Stage + " has a timeout ,stop ...")
			stopStage(actionList, pipelineInfo, stageInfo, pipelineSequence)
			return
		case isAllActionOk := <-actionResultChan:
			if isAllActionOk {
				log.Info("stage " + stageInfo.Stage + "'s all action run success, start next stage:" + nextStage.Stage)
				handleStage(pipelineInfo, *nextStage, pipelineSequence, pipelineEnvMap)
			} else {
				// if has a failer action ,then stop all other action's
				log.Info("stage " + stageInfo.Stage + " is stop with an action's error!")
				stopStage(actionList, pipelineInfo, stageInfo, pipelineSequence)
				return
			}
		}
	}
}

// stopStage is to stop all givent action in actionList
func stopStage(actionList []models.ActionLog, pipelineInfo models.PipelineLog, stageInfo models.StageLog, pipelineSequence int64) {
	for _, action := range actionList {
		if action.Component != 0 {
			go stopComponent(pipelineInfo, stageInfo, action, pipelineSequence)
		} else {
			go stopService(pipelineInfo, stageInfo, action, pipelineSequence)
		}
	}
}

// exec a action
func execAction(pipelineInfo models.PipelineLog, stageInfo models.StageLog, actionInfo models.ActionLog, pipelineSequence int64, envMap map[string]string) {
	fmt.Println("----------------------------------------------")
	fmt.Println(actionInfo.Environment)
	fmt.Println("=====>", envMap)
	if actionInfo.Environment != "" {
		err := json.Unmarshal([]byte(actionInfo.Environment), &envMap)
		if err != nil {
			log.Error("error when load action's env when try to start action" + actionInfo.Action + "(" + strconv.FormatInt(actionInfo.ID, 10) + ")")
		}
	}
	fmt.Println("=====>", envMap["CO_DATA"])

	if actionInfo.Component != 0 {
		startComponent(pipelineInfo, stageInfo, actionInfo, pipelineSequence, envMap)
	} else {
		startService(pipelineInfo, stageInfo, actionInfo, pipelineSequence, envMap)
	}
}

/////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////
func startComponent(pipelineInfo models.PipelineLog, stageInfo models.StageLog, actionInfo models.ActionLog, pipelineSequence int64, envMap map[string]string) {
	// get all event that bind this action
	eventList := make([]models.EventDefinition, 0)
	new(models.EventDefinition).GetEventDefinition().Where("action = ?", actionInfo.ID).Find(&eventList)

	// now component default run in k8s
	// TODO :component run in swarm etc.
	// componentId = pipelineId + stageId + actionId + pipelineSequence + componentId
	platformSetting, err := GetActionPlatformInfo(actionInfo)
	if err != nil {
		startErrOutcome := new(models.Outcome)
		startErrOutcome.Pipeline = pipelineInfo.ID
		startErrOutcome.RealPipeline = pipelineInfo.FromPipeline
		startErrOutcome.Stage = stageInfo.ID
		startErrOutcome.RealStage = stageInfo.FromStage
		startErrOutcome.Action = actionInfo.ID
		startErrOutcome.RealAction = actionInfo.FromAction
		startErrOutcome.Sequence = pipelineSequence
		startErrOutcome.Status = false
		startErrOutcome.Result = err.Error()
		startErrOutcome.Output = ""

		startErrOutcome.GetOutcome().Save(startErrOutcome)
		return
	}

	componentId := strconv.FormatInt(pipelineInfo.ID, 10) + "," + strconv.FormatInt(stageInfo.ID, 10) + "," + strconv.FormatInt(actionInfo.ID, 10) + "," + strconv.FormatInt(pipelineSequence, 10) + "," + strconv.FormatInt(actionInfo.Component, 10)

	c, err := InitComponet(actionInfo, platformSetting["platformType"], platformSetting["platformHost"], pipelineInfo.Namespace)
	if err != nil {
		// if has init error,stop this action and log it as start error
		startErrOutcome := new(models.Outcome)
		startErrOutcome.Pipeline = pipelineInfo.ID
		startErrOutcome.RealPipeline = pipelineInfo.FromPipeline
		startErrOutcome.Stage = stageInfo.ID
		startErrOutcome.RealStage = stageInfo.FromStage
		startErrOutcome.Action = actionInfo.ID
		startErrOutcome.RealAction = actionInfo.FromAction
		startErrOutcome.Sequence = pipelineSequence
		startErrOutcome.Status = false
		startErrOutcome.Result = "init error:" + err.Error()
		startErrOutcome.Output = ""

		startErrOutcome.GetOutcome().Save(startErrOutcome)
		return
	}
	err = c.Start(componentId, eventList, envMap)
	if err != nil {
		// if has start error,stop this action and log it as start error
		startErrOutcome := new(models.Outcome)
		startErrOutcome.Pipeline = pipelineInfo.ID
		startErrOutcome.RealPipeline = pipelineInfo.FromPipeline
		startErrOutcome.Stage = stageInfo.ID
		startErrOutcome.RealStage = stageInfo.FromStage
		startErrOutcome.Action = actionInfo.ID
		startErrOutcome.RealAction = actionInfo.FromAction
		startErrOutcome.Sequence = pipelineSequence
		startErrOutcome.Status = false
		startErrOutcome.Result = "start error:" + err.Error()
		startErrOutcome.Output = ""

		startErrOutcome.GetOutcome().Save(startErrOutcome)
	}
}

func stopComponent(pipelineInfo models.PipelineLog, stageInfo models.StageLog, actionInfo models.ActionLog, pipelineSequence int64) {
	platformSetting, err := GetActionPlatformInfo(actionInfo)
	if err != nil {
		initErrOutcome := new(models.Outcome)
		initErrOutcome.Pipeline = pipelineInfo.ID
		initErrOutcome.RealPipeline = pipelineInfo.FromPipeline
		initErrOutcome.Stage = stageInfo.ID
		initErrOutcome.RealStage = stageInfo.FromStage
		initErrOutcome.Action = actionInfo.ID
		initErrOutcome.RealAction = actionInfo.FromAction
		initErrOutcome.Sequence = pipelineSequence
		initErrOutcome.Status = false
		initErrOutcome.Result = err.Error()
		initErrOutcome.Output = ""

		initErrOutcome.GetOutcome().Save(initErrOutcome)
		return
	}

	c, err := InitComponet(actionInfo, platformSetting["platformType"], platformSetting["platformHost"], pipelineInfo.Namespace)
	if err != nil {
		// if has init error,stop this action and log it as start error
		initErrOutcome := new(models.Outcome)
		initErrOutcome.Pipeline = pipelineInfo.ID
		initErrOutcome.RealPipeline = pipelineInfo.FromPipeline
		initErrOutcome.Stage = stageInfo.ID
		initErrOutcome.RealStage = stageInfo.FromStage
		initErrOutcome.Action = actionInfo.ID
		initErrOutcome.RealAction = actionInfo.FromAction
		initErrOutcome.Sequence = pipelineSequence
		initErrOutcome.Status = false
		initErrOutcome.Result = "component init error:" + err.Error()
		initErrOutcome.Output = ""

		initErrOutcome.GetOutcome().Save(initErrOutcome)
		return
	}

	componentId := strconv.FormatInt(pipelineInfo.ID, 10) + "," + strconv.FormatInt(stageInfo.ID, 10) + "," + strconv.FormatInt(actionInfo.ID, 10) + "," + strconv.FormatInt(pipelineSequence, 10) + "," + strconv.FormatInt(actionInfo.Component, 10)
	c.Stop(componentId)
}

// TODO : start a service
func startService(pipelineInfo models.PipelineLog, stageInfo models.StageLog, actionInfo models.ActionLog, pipelineSequence int64, envMap map[string]string) {
	startErrOutcome := new(models.Outcome)
	startErrOutcome.Pipeline = pipelineInfo.ID
	startErrOutcome.RealPipeline = pipelineInfo.FromPipeline
	startErrOutcome.Stage = stageInfo.ID
	startErrOutcome.RealStage = stageInfo.FromStage
	startErrOutcome.Action = actionInfo.ID
	startErrOutcome.RealAction = actionInfo.FromAction
	startErrOutcome.Sequence = pipelineSequence
	startErrOutcome.Status = false
	startErrOutcome.Result = "start error: action component is 0"
	startErrOutcome.Output = ""

	startErrOutcome.GetOutcome().Save(startErrOutcome)
}

// TODO : stop a service
func stopService(pipelineInfo models.PipelineLog, stageInfo models.StageLog, actionInfo models.ActionLog, pipelineSequence int64) {

}

// TODO: need modify to ETCD or Redis to reduce DB IO
func waitAllActionFinish(actionList []models.ActionLog, Sequence int64, actionResultChan chan bool) {
	allActionIsOk := true

	actionIds := make([]int64, 0)
	for _, action := range actionList {
		actionIds = append(actionIds, action.ID)
	}

	for {
		time.Sleep(1 * time.Second)

		runResults := make([]struct {
			Status bool
		}, 0)

		new(models.Outcome).GetOutcome().Table("outcome").Where("sequence = ?", Sequence).Where("action in (?)", actionIds).Find(&runResults)

		for _, runResult := range runResults {
			if !runResult.Status {
				allActionIsOk = false
				break
			}
		}

		if len(runResults) == len(actionList) || !allActionIsOk {
			break
		}
	}

	actionResultChan <- allActionIsOk
	return
}

// send data to target url
func SendDataToAction(runId, targetUrl, podName string) {
	// runID = pipelineId + stageId + actionId + pipelineSequence + componentId
	pipelineId, err := strconv.ParseInt(strings.Split(runId, ",")[0], 10, 64)
	if err != nil {
		log.Error("error when get pipelineId:" + err.Error())
		return
	}
	stageId, err := strconv.ParseInt(strings.Split(runId, ",")[1], 10, 64)
	if err != nil {
		log.Error("error when get stageId:" + err.Error())
		return
	}
	actionId, err := strconv.ParseInt(strings.Split(runId, ",")[2], 10, 64)
	if err != nil {
		log.Error("error when get actionId:" + err.Error())
		return
	}
	pipelineSequence, err := strconv.ParseInt(strings.Split(runId, ",")[3], 10, 64)
	if err != nil {
		log.Error("error when get pipelineSequence:" + err.Error())
		return
	}
	componentId, err := strconv.ParseInt(strings.Split(runId, ",")[4], 10, 64)
	if err != nil {
		log.Error("error when get componentId:" + err.Error())
		return
	}

	// get action info
	actionInfo := new(models.ActionLog)
	err = actionInfo.GetActionLog().Where("id = ?", actionId).First(actionInfo).Error
	if err != nil {
		log.Error("error when get action info by given id:", actionId, ",err :", err.Error())
		return
	}

	// unmarshal action's manifestMap to get action data relation
	manifestMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(actionInfo.Manifest), &manifestMap)
	if err != nil {
		log.Error("error when get action manifest info:" + err.Error())
		return
	}

	stageInfo := new(models.StageLog)
	err = stageInfo.GetStageLog().Where("id = ?", stageId).First(&stageInfo).Error
	if err != nil {
		log.Error("error when get stage info from db:" + err.Error())
		return
	}

	fmt.Println("relation: ", manifestMap["relation"])

	dataMap := make(map[string]interface{})
	relations, ok := manifestMap["relation"]
	if ok {
		relationInfo, ok := relations.([]interface{})
		if !ok {
			log.Error("error when parse relations,relations is not an array")
			return
		}

		fmt.Println("start merage action data...")
		// get all data that current action is require
		dataMap, err = merageFromActionsOutputData(pipelineId, stageInfo.PreStage, actionId, pipelineSequence, componentId, relationInfo)
		fmt.Println("result data is :", dataMap)
		if err != nil {
			log.Error("error when get data map from action: " + err.Error())
			return
		}
	}

	if len(dataMap) == 0 {
		return
	}

	dataByte, err := json.Marshal(dataMap)
	if err != nil {
		log.Error("error when marshal dataMap from action:"+err.Error(), ", data map:", dataMap)
		return
	}

	// send data to component or service
	var resp *http.Response
	if actionInfo.Component != 0 {
		resp, err = sendDataToComponent(pipelineId, stageId, actionId, pipelineSequence, componentId, podName, targetUrl, (dataByte))
	} else {
		resp, err = sendDataToService((dataByte))
	}
	if err != nil {
		log.Error("error when send data to action:" + err.Error())
	}

	// record send info
	payload := make(map[string]interface{})
	payload["data"] = string(dataByte)
	if err != nil {
		payload["error"] = err.Error()
	} else {
		respBody, _ := ioutil.ReadAll(resp.Body)
		payload["resp"] = string(respBody)
	}

	payloadInfo, err := json.Marshal(payload)
	if err != nil {
		log.Error("error when marshal payload info:" + err.Error())
	}

	sendDataEvent := new(models.Event)
	sendDataEvent.Title = "SEND_DATA"
	sendDataEvent.Payload = string(payloadInfo)
	sendDataEvent.Type = models.TypeSystemEvent
	sendDataEvent.Pipeline = pipelineId
	sendDataEvent.Stage = stageId
	sendDataEvent.Action = actionId
	sendDataEvent.Sequence = pipelineSequence

	err = sendDataEvent.GetEvent().Save(sendDataEvent).Error
	if err != nil {
		log.Error("error when save send data info :" + err.Error())
	}
}

func merageFromActionsOutputData(pipelineId, stageId, actionId, pipelineSequence, componentId int64, relations []interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	for _, relation := range relations {
		relationMap, ok := relation.(map[string]interface{})
		if !ok {
			return nil, errors.New("error when parse relation info,relation is not a json!")
		}

		fromAction, ok := relationMap["fromAction"]
		if !ok {
			return nil, errors.New("error when read relation from data: relation don't have a fromAction data")
		}

		fromOutcome := new(models.Outcome)
		err := fromOutcome.GetOutcome().Where("pipeline = ? ", pipelineId).Where("real_action = ?", fromAction).Where("sequence = ?", pipelineSequence).First(&fromOutcome).Error
		if err != nil {
			return nil, errors.New("error when get from outcome, error:" + err.Error())
		}

		tempData := make(map[string]interface{})
		err = json.Unmarshal([]byte(fromOutcome.Output), &tempData)
		if err != nil {
			return nil, errors.New("error when parse from action data1:" + err.Error() + "\n" + fromOutcome.Output)
		}

		relationInfo, ok := relationMap["relation"]
		if !ok {
			return nil, errors.New("relation don't have a relation info")
		}

		relationArray, ok := relationInfo.([]interface{})
		if !ok {
			return nil, errors.New("error when get relation info ,relation is not a array")
		}

		relationList := make([]Relation, 0)
		for _, realationDefines := range relationArray {
			realationDefine := realationDefines.([]interface{})[0]
			fmt.Println("===================================================")
			fmt.Println("relationList", realationDefine)
			fmt.Println("===================================================")
			relationByte, err := json.Marshal(realationDefine.(map[string]interface{}))
			if err != nil {
				return nil, errors.New("error when marshal relation array:" + err.Error())
			}

			var r Relation
			err = json.Unmarshal(relationByte, &r)
			if err != nil {
				return nil, errors.New("error when parse relation info:" + err.Error())
			}

			relationList = append(relationList, r)
		}
		actionResult := make(map[string]interface{})
		err = getResultFromRelation(fromOutcome.Output, relationList, actionResult)

		if err != nil {
			return nil, errors.New("error when get from data:" + err.Error())
		}

		for key, value := range actionResult {
			result[key] = value
		}
	}

	return result, nil
}

// getResultFromRelation is get data from an action output
func getResultFromRelation(fromActionOutput string, relationList []Relation, result map[string]interface{}) error {
	fromActionData := make(map[string]interface{})

	err := json.Unmarshal([]byte(fromActionOutput), &fromActionData)
	if err != nil {
		return errors.New("error when parse from action data2:" + err.Error() + "\n" + fromActionOutput)
	}

	for _, relation := range relationList {
		// if len(relation.Child) == 0 {
		// get data from current relation path
		fromData, err := getJsonDataByPath(strings.TrimPrefix(relation.From, "."), fromActionData)
		if err != nil {
			return errors.New("error when get fromData :" + err.Error())
		}

		setDataToMapByPath(fromData, result, strings.TrimPrefix(relation.To, "."))
		// } else {
		// 	getResultFromRelation(fromActionOutput, relation.Child, result)
		// }
	}

	return nil
}

// getJsonDataByPath is get a value from a map by give path
func getJsonDataByPath(path string, data map[string]interface{}) (interface{}, error) {
	depth := len(strings.Split(path, "."))
	if depth == 1 {
		if info, ok := data[path]; !ok {
			// return nil, errors.New("key not exist:" + path)
			log.Error("error when get data from action,action's key not exist :" + path)
			return "", nil
		} else {
			return info, nil
		}
	}

	childDataInterface, ok := data[strings.Split(path, ".")[0]]
	if !ok {
		log.Error("error when get data from action,action's key not exist :" + path)
		return "", nil
		// return nil, errors.New("key not exist:" + path)
	}
	childData, ok := childDataInterface.(map[string]interface{})
	if !ok {
		return nil, errors.New("child data is not a json!")
	}

	childPath := strings.Join(strings.Split(path, ".")[1:], ".")
	return getJsonDataByPath(childPath, childData)
}

// setDataToMapByPath is set a data to a map by give path ,if parent path not exist,it will auto creat
func setDataToMapByPath(data interface{}, result map[string]interface{}, path string) {
	depth := len(strings.Split(path, "."))
	if depth == 1 {
		result[path] = data
		return
	}

	currentPath := strings.Split(path, ".")[0]
	currentMap := make(map[string]interface{})
	if _, ok := result[currentPath]; !ok {
		result[currentPath] = currentMap
	}

	var ok bool
	currentMap, ok = result[currentPath].(map[string]interface{})
	if !ok {
		return
	}

	childPath := strings.Join(strings.Split(path, ".")[1:], ".")
	setDataToMapByPath(data, currentMap, childPath)
	return
}

// sendDataToComponent is send data to component
func sendDataToComponent(pipelineId, stageId, actionId, pipelineSequence, componentId int64, podName, targetUrl string, data []byte) (*http.Response, error) {

	// componentInfo := new(models.ComponentLog)
	// componentInfo.GetComponentLog().Where("id = ?", componentId).First(componentInfo)

	actionInfo := new(models.ActionLog)
	actionInfo.GetActionLog().Where("id = ?", actionId).First(actionInfo)
	platformSetting, err := GetActionPlatformInfo(*actionInfo)
	if err != nil {
		return nil, err
	}

	pipelineInfo := new(models.PipelineLog)
	pipelineInfo.GetPipelineLog().Where("id = ? ", pipelineId).First(pipelineInfo)

	c, err := InitComponet(*actionInfo, platformSetting["platformType"], platformSetting["platformHost"], pipelineInfo.Namespace)

	// c, err := InitComponet(*componentInfo, RUNENV_KUBE)
	if err != nil {
		return nil, err
	}

	ip, err := c.GetIp(podName)
	if err != nil {
		return nil, err
	}

	fmt.Println("===========================================================")
	fmt.Println("===========================================================")
	fmt.Println("===========================================================")
	fmt.Println("===========================================================")
	fmt.Println("http://" + ip + targetUrl)
	fmt.Println(string(data))
	fmt.Println("===========================================================")
	fmt.Println("===========================================================")
	fmt.Println("===========================================================")
	fmt.Println("===========================================================")
	fmt.Println("===========================================================")
	fmt.Println("===========================================================")

	isSend := false
	count := 0
	var sendResp *http.Response
	for !isSend && count < 10 {
		sendResp, err = http.Post("http://"+ip+targetUrl, "application/json", bytes.NewReader(data))
		if err == nil {
			isSend = true
			break
		}
		count++
		// wait some time and send again
		time.Sleep(2 * time.Second)
	}

	if !isSend {
		return nil, errors.New("error when send data to component:" + err.Error())
	}
	return sendResp, nil
}

// sendDataToService is
// TODO : will sent it to service
func sendDataToService(data []byte) (*http.Response, error) {
	return nil, nil
}

// CreateNewPipeline is
func CreateNewPipeline(namespace, pipelineName, pipelineVersion string) (string, error) {
	createPipelineChan <- true
	defer func() {
		<-createPipelineChan
	}()

	var count int64
	err := new(models.Pipeline).GetPipeline().Where("namespace = ?", namespace).Where("pipeline = ?", pipelineName).Order("-id").Count(&count).Error
	if err != nil {
		return "", errors.New("error when query pipeline data in database:" + err.Error())
	}

	if count > 0 {
		return "", errors.New("pipelien name is exist!")
	}

	pipelineInfo := new(models.Pipeline)
	pipelineInfo.Namespace = namespace
	pipelineInfo.Pipeline = pipelineName
	pipelineInfo.Version = pipelineVersion
	pipelineInfo.VersionCode = 1

	err = pipelineInfo.GetPipeline().Save(pipelineInfo).Error
	if err != nil {
		return "", errors.New("error when save pipeline info:" + err.Error())
	}

	return "create new pipeline success", nil
}

func GetPipelineInfo(namespace, pipelineName string, pipelineId int64) (map[string]interface{}, error) {
	resultMap := make(map[string]interface{})
	pipelineInfo := new(models.Pipeline)
	err := pipelineInfo.GetPipeline().Where("id = ?", pipelineId).First(&pipelineInfo).Error
	if err != nil {
		return nil, errors.New("error when get pipeline info from db:" + err.Error())
	}

	if pipelineInfo.Pipeline != pipelineName {
		return nil, errors.New("pipeline's name is not equal to target pipeline")
	}

	// get pipeline define json first, if has a define json,return it
	if pipelineInfo.Manifest != "" {
		defineMap := make(map[string]interface{})
		json.Unmarshal([]byte(pipelineInfo.Manifest), &defineMap)

		if defineInfo, ok := defineMap["define"]; ok {
			if defineInfoMap, ok := defineInfo.(map[string]interface{}); ok {
				defineInfoMap["status"] = pipelineInfo.State == models.PipelineStateAble
				return defineInfoMap, nil
			}
		}
	}

	// get all stage info of current pipeline
	// if a pipeline done have a define of itself
	// then the pipeline is a new pipeline ,so only get it's stage list is ok
	// stageList, err := getStageListByPipeline(*pipelineInfo)
	// if err != nil {
	// 	return nil, err
	// }
	resultMap["stageList"] = make([]map[string]interface{}, 0)

	resultMap["lineList"] = make([]map[string]interface{}, 0)

	resultMap["status"] = false

	return resultMap, nil
}

func getStageListByPipeline(pipelineInfo models.Pipeline) ([]map[string]interface{}, error) {
	stageList := make([]models.Stage, 0)
	err := new(models.Stage).GetStage().Where("pipeline = ?", pipelineInfo.ID).Find(&stageList).Error
	if err != nil {
		return nil, err
	}

	stageListMap := make([]map[string]interface{}, 0)
	for i, stageInfo := range stageList {
		stageInfoMap := make(map[string]interface{})

		stageInfoMap["id"] = "stage-" + strconv.FormatInt(stageInfo.ID, 10)
		stageInfoMap["type"] = "pipeline-stage"
		stageInfoMap["setupData"] = make(map[string]interface{})
		if stageInfo.PreStage == -1 {
			stageInfoMap["type"] = "pipeline-start"
		}

		if stageInfo.PreStage != -1 {
			// if not a start stage,get all action in current stage
			allActionList := make([]models.Action, 0)
			err = new(models.Action).GetAction().Where("stage = ?", stageInfo.ID).Find(&allActionList).Error
			if err != nil {
				return nil, err
			}

			if len(allActionList) > 0 {
				allActionMap := make([]map[string]interface{}, 0)
				for _, actionInfo := range allActionList {
					actionMap := make(map[string]interface{})

					actionMap["id"] = "action-" + strconv.FormatInt(actionInfo.ID, 10)
					actionMap["type"] = "pipeline-action"

					allActionMap = append(allActionMap, actionMap)
				}

				stageInfoMap["actions"] = allActionMap
			}
		}

		if i == len(stageList)-1 {
			// if this is the least stage ,add a add-stage for display
			addStageInfo := make(map[string]interface{})
			addStageInfo["type"] = "pipeline-add-stage"
			addStageInfo["id"] = "pipeline-add-stage"
			stageListMap = append(stageListMap, addStageInfo)

			stageInfoMap["type"] = "pipeline-end"
		}

		stageListMap = append(stageListMap, stageInfoMap)
	}

	// if is a empty stage ,return init pipeline data
	if len(stageList) == 0 {
		startStage := make(map[string]interface{})
		startStage["id"] = "start-stage"
		startStage["type"] = "pipeline-start"
		startStage["setupData"] = make(map[string]interface{})
		stageListMap = append(stageListMap, startStage)

		addStage := make(map[string]interface{})
		addStage["id"] = "add-stage"
		addStage["type"] = "pipeline-add-stage"
		stageListMap = append(stageListMap, addStage)

		endStage := make(map[string]interface{})
		endStage["id"] = "end-stage"
		endStage["type"] = "pipeline-end"
		endStage["setupData"] = make(map[string]interface{})
		stageListMap = append(stageListMap, endStage)
	}

	return stageListMap, nil
}

func UpdatePipelineInfo(pipelineInfo models.Pipeline) error {
	// get pipeline define info and get pipeline's line list and stage list
	pipelineInfo.GetPipeline().Save(&pipelineInfo)
	relationMap, stageDefineList, err := getPipelineDefineInfo(pipelineInfo)
	if err != nil {
		return err
	}

	// first delete old pipeline define
	stageList := make([]models.Stage, 0)
	err = new(models.Stage).GetStage().Where("pipeline = ? ", pipelineInfo.ID).Find(&stageList).Error
	if err != nil {
		return errors.New("error when get stage list:" + err.Error())
	}

	stageIdList := make([]int64, 0)
	actionIdList := make([]int64, 0)
	for _, stage := range stageList {
		tempActionList := make([]models.Action, 0)
		err = new(models.Action).GetAction().Where("stage = ?", stage.ID).Find(&tempActionList).Error
		if err != nil {
			return errors.New("error when get action list:" + err.Error())
		}

		for _, action := range tempActionList {
			actionIdList = append(actionIdList, action.ID)
		}

		stageIdList = append(stageIdList, stage.ID)
	}

	err = new(models.Action).GetAction().Where("id in (?)", actionIdList).Delete(&models.Action{}).Error
	if err != nil {
		return errors.New("error when update action info:" + err.Error())
	}

	err = new(models.Stage).GetStage().Where("id in (?)", stageIdList).Delete(&models.Stage{}).Error
	if err != nil {
		return errors.New("error when update stage info:" + err.Error())
	}

	// then create new pipeline by define
	stageInfoMap := make(map[string]map[string]interface{})
	preStageId := int64(-1)
	allActionIdMap := make(map[string]int64)
	for _, stageDefine := range stageDefineList {
		stageId, stageTagId, actionMap, err := saveStageByStageDefine(stageDefine, pipelineInfo, preStageId, relationMap)
		if err != nil {
			return err
		}

		if stageId != 0 {
			preStageId = stageId
		}

		stageDefine["stageId"] = stageId
		stageInfoMap[stageTagId] = stageDefine
		for key, value := range actionMap {
			allActionIdMap[key] = value
		}
	}

	for actionOriginId, actionID := range allActionIdMap {
		if relations, ok := relationMap[actionOriginId].(map[string]interface{}); ok {
			actionRealtionList := make([]map[string]interface{}, 0)
			for fromActionOriginId, realRelations := range relations {
				fromActionId, ok := allActionIdMap[fromActionOriginId]
				if !ok {
					return errors.New("action's relation is illegal")
				}

				tempRelation := make(map[string]interface{})
				tempRelation["toAction"] = actionID
				tempRelation["fromAction"] = fromActionId
				tempRelation["relation"] = realRelations

				actionRealtionList = append(actionRealtionList, tempRelation)
			}

			actionInfo := new(models.Action)
			actionInfo.GetAction().Where("id = ?", actionID).First(&actionInfo)
			manifestMap := make(map[string]interface{})
			if actionInfo.Manifest != "" {
				json.Unmarshal([]byte(actionInfo.Manifest), &manifestMap)
			}

			manifestMap["relation"] = actionRealtionList
			relationBytes, _ := json.Marshal(manifestMap)
			actionInfo.Manifest = string(relationBytes)

			actionInfo.GetAction().Where("id = ?", actionID).UpdateColumn("manifest", actionInfo.Manifest)
		}
	}

	return nil
}

func getPipelineDefineInfo(pipelineInfo models.Pipeline) (map[string]interface{}, []map[string]interface{}, error) {
	lineList := make([]map[string]interface{}, 0)
	stageList := make([]map[string]interface{}, 0)

	manifestMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(pipelineInfo.Manifest), &manifestMap)
	if err != nil {
		return nil, nil, errors.New("error when unmarshal pipeline manifes info:" + err.Error())
	}

	defineMap, ok := manifestMap["define"].(map[string]interface{})
	if !ok {
		return nil, nil, errors.New("pipeline's define is not a json:" + err.Error())
	}

	// get all line data generate a map to record
	realtionMap := make(map[string]interface{})
	lineListInfo, ok := defineMap["lineList"]
	if ok {
		linesList, ok := lineListInfo.([]interface{})
		if !ok {
			return nil, nil, errors.New("pipeline's lineList define is not an array")
		}

		for _, lineInfo := range linesList {
			lineInfoMap, ok := lineInfo.(map[string]interface{})
			if !ok {
				return nil, nil, errors.New("pipeline's line info is not a json")
			}

			lineList = append(lineList, lineInfoMap)
		}

		for _, lineInfo := range lineList {
			endData, ok := lineInfo["endData"].(map[string]interface{})
			if !ok {
				return nil, nil, errors.New("pipeline's line define is illegal,don't have a end point info")
			}

			endPointId, ok := endData["id"].(string)
			if !ok {
				return nil, nil, errors.New("pipeline's line define is illegal,endPoint id is not a string")
			}

			if _, ok := realtionMap[endPointId]; !ok {
				realtionMap[endPointId] = make(map[string]interface{})
			}

			endPointMap := realtionMap[endPointId].(map[string]interface{})
			startData, ok := lineInfo["startData"].(map[string]interface{})
			if !ok {
				return nil, nil, errors.New("pipeline's line define is illegal,don;t have a start point info")
			}

			startDataId, ok := startData["id"].(string)
			if !ok {
				return nil, nil, errors.New("pipeline's line define is illegal,startPoint id is not a string")
			}

			if _, ok := endPointMap[startDataId]; !ok {
				endPointMap[startDataId] = make([]interface{}, 0)
			}

			lineList, ok := lineInfo["relation"].([]interface{})
			if !ok {
				continue
			}

			endPointMap[startDataId] = append(endPointMap[startDataId].([]interface{}), lineList)
		}
	}

	stageListInfo, ok := defineMap["stageList"]
	if !ok {
		return nil, nil, errors.New("pipeline don't have a stage define")
	}

	stagesList, ok := stageListInfo.([]interface{})
	if !ok {
		return nil, nil, errors.New("pipeline's stageList define is not an array")
	}

	for _, stageInfo := range stagesList {
		stageInfoMap, ok := stageInfo.(map[string]interface{})
		if !ok {
			return nil, nil, errors.New("pipeline's stage info is not a json")
		}

		stageList = append(stageList, stageInfoMap)
	}

	return realtionMap, stageList, nil
}

func saveStageByStageDefine(stageDefine map[string]interface{}, pipelineInfo models.Pipeline, preStageId int64, relationMap map[string]interface{}) (int64, string, map[string]int64, error) {
	stageType := models.StageTypeRun
	actionIdMap := make(map[string]int64)
	stageName := ""
	timeout := int64(60 * 60 * 24 * 36)
	manifestMap := make(map[string]interface{})

	idStr, ok := stageDefine["id"].(string)
	if !ok {
		return 0, "", nil, errors.New("stage define does not have a string id")
	}

	stageDefineType, ok := stageDefine["type"].(string)
	if !ok {
		return 0, "", nil, errors.New("stage type define is not a string")
	}

	if stageDefineType == PIPELINE_STAGE_TYPE_ADD {
		return 0, "", nil, nil
	} else if stageDefineType == PIPELINE_STAGE_TYPE_START {
		stageType = models.StageTypeStart
		stageName = pipelineInfo.Pipeline + "-start-stage"
		timeout = 0

		if stageSetupDataMap, ok := stageDefine["setupData"].(map[string]interface{}); ok {
			updatePipelineSourceInfo(stageSetupDataMap, pipelineInfo)
		}
	} else if stageDefineType == PIPELINE_STAGE_TYPE_END {
		stageType = models.StageTypeEnd
		stageName = pipelineInfo.Pipeline + "-end-stage"
		timeout = 0
	} else if stageDefineType == PIPELINE_STAGE_TYPE_RUN {
		setupData, ok := stageDefine["setupData"]
		if ok {
			setupDataMap, ok := setupData.(map[string]interface{})
			if !ok {
				return 0, "", nil, errors.New("stage's setupData is not a json")
			}
			defineName, ok := setupDataMap["name"]
			if ok {
				defineNameStr, ok := defineName.(string)
				if !ok {
					return 0, "", nil, errors.New("stage's name is not a string")
				}

				stageName = defineNameStr
			}

			defineTimeoutStr, ok := setupDataMap["timeout"].(string)
			if ok {
				var err error
				timeout, err = strconv.ParseInt(defineTimeoutStr, 10, 64)
				if err != nil {
					return 0, "", nil, errors.New("stage's timeout is not a string")
				}
			}
		}
	} else {
		return 0, "", nil, nil
	}

	manifestByte, _ := json.Marshal(manifestMap)

	stage := new(models.Stage)
	stage.Pipeline = pipelineInfo.ID
	stage.Type = int64(stageType)
	stage.PreStage = preStageId
	stage.Stage = stageName
	stage.Title = stageName
	stage.Description = stageName
	stage.Manifest = string(manifestByte)
	stage.Timeout = timeout

	err := stage.GetStage().Save(stage).Error
	if err != nil {
		return 0, "", nil, err
	}

	if stageDefineType == PIPELINE_STAGE_TYPE_START {
		actionIdMap[idStr] = 0
	}

	if actionDefine, ok := stageDefine["actions"]; ok {
		actionList, ok := actionDefine.([]interface{})
		if !ok {
			return 0, "", nil, errors.New("action list is not an array")
		}

		actionDefineList := make([]map[string]interface{}, 0)
		for _, action := range actionList {
			actionDefineMap, ok := action.(map[string]interface{})
			if !ok {
				return 0, "", nil, errors.New("action's define is not a json")
			}
			actionDefineList = append(actionDefineList, actionDefineMap)
		}

		actionIdMap, err = createActionByDefine(actionDefineList, stage.ID)
		if err != nil {
			return 0, "", nil, err
		}

	}

	return stage.ID, idStr, actionIdMap, err
}

func updatePipelineSourceInfo(stageSetupDataMap map[string]interface{}, pipelineInfo models.Pipeline) {
	sourceMap := make(map[string]interface{})
	json.Unmarshal([]byte(pipelineInfo.SourceInfo), &sourceMap)
	sourceType, ok := stageSetupDataMap["type"].(string)
	if !ok {
		return
	}

	eventType, ok := stageSetupDataMap["event"].(string)
	if !ok {
		return
	}

	headerKey := ""
	switch sourceType {
	case "github":
		headerKey = "X-Hub-Signature"
	case "customize":
		headerKey = "X-Pipeline-Signature"
	}

	tempSourceMap := make(map[string]string)
	tempSourceMap["sourceType"] = sourceType
	tempSourceMap["eventList"] = "," + eventType + ","
	tempSourceMap["headerKey"] = headerKey

	sourceList := make([]interface{}, 0)
	// if _, ok := sourceMap["sourceList"].([]interface{}); ok {
	// 	sourceList = sourceMap["sourceList"].([]interface{})
	// }

	sourceList = append(sourceList, tempSourceMap)

	sourceMap["sourceList"] = sourceList

	sourceInfoBytes, _ := json.Marshal(sourceMap)
	pipelineInfo.SourceInfo = string(sourceInfoBytes)
	pipelineInfo.GetPipeline().Save(&pipelineInfo)
}

func createActionByDefine(actionDefineList []map[string]interface{}, stageId int64) (map[string]int64, error) {
	actionIdMap := make(map[string]int64)
	for _, actionDefine := range actionDefineList {
		actionName := ""
		actionImage := ""
		kubernetesSetting := ""
		inputStr := ""
		outputStr := ""
		actionTimeout := int64(60 * 60 * 24 * 36)
		componentId := int64(0)
		serviceId := int64(0)
		platformMap := make(map[string]string)

		// get component info
		component, ok := actionDefine["component"]
		if ok {
			componentMap, ok := component.(map[string]interface{})
			if !ok {
				return nil, errors.New("action's component is not a json")
			}

			componentVersion, ok := componentMap["versionid"].(float64)
			if !ok {
				return nil, errors.New("action's component info error !")
			}

			componentId = int64(componentVersion)
		}

		// get action setup data info map
		if setupDataMap, ok := actionDefine["setupData"].(map[string]interface{}); ok {
			if actionSetupDataMap, ok := setupDataMap["action"].(map[string]interface{}); ok {
				if name, ok := actionSetupDataMap["name"].(string); ok {
					actionName = name
				}

				if image, ok := actionSetupDataMap["image"].(map[string]interface{}); ok {
					actionImage = ""
					if name, ok := image["name"]; ok {
						actionImage = name.(string) + ":"
						if tag, ok := image["tag"]; ok {
							actionImage += tag.(string)
						} else {
							actionImage += "latest"
						}
					}
				}

				if timeoutStr, ok := actionSetupDataMap["timeout"].(string); ok {
					var err error
					actionTimeout, err = strconv.ParseInt(timeoutStr, 10, 64)
					if err != nil {
						return nil, errors.New("action's timeout is not string")
					}
				}

				configMap := make(map[string]interface{})
				// record platform info
				if platFormType, ok := actionSetupDataMap["type"].(string); ok {
					platformMap["platformType"] = strings.ToUpper(platFormType)
				}

				if platformHost, ok := actionSetupDataMap["apiserver"].(string); ok {
					platformMap["platformHost"] = strings.ToUpper(platformHost)
				}

				if ip, ok := actionSetupDataMap["ip"].(string); ok {
					configMap["reachableIPs"] = []string{ip}
				}

				// unmarshal k8s info
				if useAdvanced, ok := actionSetupDataMap["useAdvanced"].(bool); ok {
					podConfigKey := "pod"
					serviceConfigKey := "service"
					if useAdvanced {
						podConfigKey = "pod_advanced"
						serviceConfigKey = "service_advanced"
					}

					podConfig, ok := setupDataMap[podConfigKey].(map[string]interface{})
					if !ok {
						configMap["podConfig"] = make(map[string]interface{})
					} else {
						configMap["podConfig"] = podConfig
					}

					serviceConfig, ok := setupDataMap[serviceConfigKey].(map[string]interface{})
					if !ok {
						configMap["serviceConfig"] = make(map[string]interface{})
					} else {
						configMap["serviceConfig"] = serviceConfig
					}

					kuberSettingBytes, _ := json.Marshal(configMap)
					kubernetesSetting = string(kuberSettingBytes)
				}
			}
		}

		inputMap, ok := actionDefine["inputJson"].(map[string]interface{})
		if ok {
			inputDescribe, err := describeJSON(inputMap, "")
			if err != nil {
				return nil, errors.New("error in component output json define:" + err.Error())
			}

			inputDescBytes, _ := json.Marshal(inputDescribe)
			inputStr = string(inputDescBytes)
		}

		outputMap, ok := actionDefine["outputJson"].(map[string]interface{})
		if ok {
			outputDescribe, err := describeJSON(outputMap, "")
			if err != nil {
				return nil, errors.New("error in component output json define:" + err.Error())
			}

			outputDescBytes, _ := json.Marshal(outputDescribe)
			outputStr = string(outputDescBytes)
		}

		allEnvMap := make(map[string]string)
		if envMap, ok := actionDefine["env"].([]interface{}); ok {
			for _, envInfo := range envMap {
				envInfoMap, ok := envInfo.(map[string]interface{})
				if !ok {
					fmt.Println(envInfo)
					fmt.Println(envInfo.(int64))
					return nil, errors.New("action's env set is not a json")
				}

				key, ok := envInfoMap["key"].(string)
				if !ok {
					return nil, errors.New("action's key is not a string")
				}

				value, ok := envInfoMap["value"].(string)
				if !ok {
					return nil, errors.New("action's value is not a string")
				}
				allEnvMap[key] = value
			}
		}
		envBytes, _ := json.Marshal(allEnvMap)

		// get aciont line info
		actionId, ok := actionDefine["id"].(string)
		if !ok {
			return nil, errors.New("action's id is not a string")
		}

		manifestMap := make(map[string]interface{})
		manifestMap["platform"] = platformMap
		manifestBytes, _ := json.Marshal(manifestMap)

		action := new(models.Action)
		action.Stage = stageId
		action.Component = componentId
		action.Service = serviceId
		action.Action = actionName
		action.Title = actionName
		action.Description = actionName
		action.Manifest = string(manifestBytes)
		action.Endpoint = actionImage
		action.Kubernetes = kubernetesSetting
		action.Input = inputStr
		action.Output = outputStr
		action.Timeout = actionTimeout
		action.Environment = string(envBytes)

		err := action.GetAction().Save(action).Error
		if err != nil {
			return nil, errors.New("error when save action info to db:" + err.Error())
		}
		actionIdMap[actionId] = action.ID
	}

	return actionIdMap, nil
}

func CreateNewPipelineVersion(pipelineInfo models.Pipeline, versionName string) error {
	var count int64
	new(models.Pipeline).GetPipeline().Where("namespace = ?", pipelineInfo.Namespace).Where("pipeline = ?", pipelineInfo.Pipeline).Where("version = ?", versionName).Count(&count)
	if count > 0 {
		return errors.New("version code already exist!")
	}

	// get current least pipeline's version
	leastPipeline := new(models.Pipeline)
	err := leastPipeline.GetPipeline().Where("namespace = ? ", pipelineInfo.Namespace).Where("pipeline = ?", pipelineInfo.Pipeline).Order("-id").First(&leastPipeline).Error
	if err != nil {
		return errors.New("error when get least pipeline info :" + err.Error())
	}

	newPipelineInfo := new(models.Pipeline)
	newPipelineInfo.Namespace = pipelineInfo.Namespace
	newPipelineInfo.Pipeline = pipelineInfo.Pipeline
	newPipelineInfo.Event = pipelineInfo.Event
	newPipelineInfo.Version = versionName
	newPipelineInfo.VersionCode = leastPipeline.VersionCode + 1
	newPipelineInfo.Manifest = pipelineInfo.Manifest
	newPipelineInfo.Description = pipelineInfo.Description
	newPipelineInfo.SourceInfo = pipelineInfo.SourceInfo
	newPipelineInfo.Env = pipelineInfo.Env

	err = newPipelineInfo.GetPipeline().Save(newPipelineInfo).Error
	if err != nil {
		return err
	}

	return UpdatePipelineInfo(*newPipelineInfo)
}

func GetActionPlatformInfo(actionInfo models.ActionLog) (map[string]string, error) {
	manifestMap := make(map[string]interface{})
	json.Unmarshal([]byte(actionInfo.Manifest), &manifestMap)

	platformSetting, ok := manifestMap["platform"].(map[string]interface{})
	if !ok {
		return nil, errors.New("action " + actionInfo.Action + "'s platform setting is illegal")
	}

	platformType, ok := platformSetting["platformType"].(string)
	if !ok {
		log.Error("action " + actionInfo.Action + "'s platform type is illegal")
		return nil, errors.New("action " + actionInfo.Action + "'s platform type is illegal")
	}

	platformHost, ok := platformSetting["platformHost"].(string)
	if !ok {
		fmt.Println(platformSetting["platformHost"].(bool))
		log.Error("action " + actionInfo.Action + "'s platform host is illegal")
		return nil, errors.New("action " + actionInfo.Action + "'s platform host is illegal")
	}

	result := make(map[string]string)
	result["platformType"] = platformType
	result["platformHost"] = platformHost

	return result, nil
}

func GetPipelineToken(namespace, pipelineName string, pipelineId int64) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	pipelineInfo := new(models.Pipeline)
	pipelineInfo.GetPipeline().Where("id = ?", pipelineId).First(&pipelineInfo)

	if pipelineInfo.ID == 0 {
		return nil, errors.New("pipeline's info is empty")
	}

	token := ""
	tokenMap := make(map[string]interface{})
	if pipelineInfo.SourceInfo == "" {
		// if sourceInfo is empty generate a token
		token = utils.MD5(pipelineInfo.Pipeline)
	} else {
		json.Unmarshal([]byte(pipelineInfo.SourceInfo), &tokenMap)

		if _, ok := tokenMap["token"].(string); !ok {
			token = utils.MD5(pipelineInfo.Pipeline)
		} else {
			token = tokenMap["token"].(string)
		}
	}

	tokenMap["token"] = token
	sourceInfo, _ := json.Marshal(tokenMap)
	pipelineInfo.SourceInfo = string(sourceInfo)
	pipelineInfo.GetPipeline().Save(pipelineInfo)

	result["token"] = token

	url := ""

	listenMode := configure.GetString("listenmode")
	switch listenMode {
	case "http":
		url = "http://"
		break
	case "https":
		url = "https://"
		break
	default:
		url = "https://"
		break
	}

	projectAddr := ""
	if configure.GetString("projectaddr") == "" {
		projectAddr = "current-pipeline's-ip:port"
	} else {
		projectAddr = configure.GetString("projectaddr")
	}

	url += projectAddr
	url = strings.TrimSuffix(url, "/")
	url += "/" + pipelineInfo.Namespace + "/" + "demo" + "/" + pipelineInfo.Pipeline

	result["url"] = url

	return result, nil
}

func GetPipelineHistoriesList(namespace string) ([]map[string]interface{}, error) {
	resultList := make([]map[string]interface{}, 0)
	pipelinesMap := make(map[int64]interface{})
	pipelineList := make([]models.Pipeline, 0)
	new(models.Pipeline).GetPipeline().Where("namespace = ?", namespace).Find(&pipelineList)

	for _, pipelineInfo := range pipelineList {
		if _, ok := pipelinesMap[pipelineInfo.ID]; !ok {
			tempMap := make(map[string]interface{})
			tempMap["versionsMap"] = make(map[int64]interface{})
			pipelinesMap[pipelineInfo.ID] = tempMap
		}

		pipelineMap := pipelinesMap[pipelineInfo.ID].(map[string]interface{})
		versionMap := pipelineMap["versionsMap"].(map[int64]interface{})

		versionMap[pipelineInfo.VersionCode] = pipelineInfo
		pipelineMap["id"] = pipelineInfo.ID
		pipelineMap["name"] = pipelineInfo.Pipeline
		pipelineMap["versionsMap"] = versionMap
	}

	for _, pipeline := range pipelinesMap {
		pipelineMap := pipeline.(map[string]interface{})
		tempPipelineMap := make(map[string]interface{})
		versionList := make([]map[string]interface{}, 0)

		versionsMap, ok := pipelineMap["versionsMap"].(map[int64]interface{})
		if ok {
			for _, version := range versionsMap {
				pipelineVersionInfo := version.(models.Pipeline)
				// 获取version所有运行信息
				startSequenceList := make([]models.Outcome, 0)
				new(models.Outcome).GetOutcome().Where("real_pipeline = ?", pipelineVersionInfo.ID).Where("action = ?", 0).Find(&startSequenceList)

				sequencesMap := make([]map[string]interface{}, 0)
				successNum := int64(0)
				sumNum := int64(0)
				for _, outcome := range startSequenceList {
					sumNum++
					resultOutcome := new(models.Outcome)
					resultOutcome.GetOutcome().Where("pipeline = ?", outcome.Pipeline).Where("sequence = ?", outcome.Sequence).Where("real_action = ?", -1).First(resultOutcome)

					tempMap := make(map[string]interface{})
					tempMap["pipelineSequenceID"] = outcome.ID
					tempMap["sequence"] = outcome.Sequence
					tempMap["status"] = false
					tempMap["time"] = outcome.CreatedAt.Format("2006-01-02 15:04:05")

					if resultOutcome.ID != 0 && resultOutcome.Status {
						successNum++
						tempMap["status"] = resultOutcome.Status
					}

					sequencesMap = append(sequencesMap, tempMap)
				}

				versionMap := make(map[string]interface{})
				versionMap["id"] = pipelineVersionInfo.ID
				versionMap["name"] = pipelineVersionInfo.Version
				versionMap["info"] = "Success :" + strconv.FormatInt(successNum, 10) + " Total :" + strconv.FormatInt(sumNum, 10)
				versionMap["sequenceList"] = sequencesMap

				versionList = append(versionList, versionMap)
			}
		}

		tempPipelineMap["id"] = pipelineMap["id"]
		tempPipelineMap["name"] = pipelineMap["name"]
		tempPipelineMap["versionList"] = versionList
		resultList = append(resultList, tempPipelineMap)
	}

	return resultList, nil
}

func GetPipelineDefineByRunSequence(sequenceId int64) (map[string]interface{}, error) {
	defineMap := make(map[string]interface{})

	pipelineStatus := true

	startOutcome := new(models.Outcome)
	startOutcome.GetOutcome().Where("id = ?", sequenceId).First(startOutcome)

	pipelineInfo := new(models.PipelineLog)
	new(models.PipelineLog).GetPipelineLog().Where("id = ?", startOutcome.Pipeline).First(pipelineInfo)

	if startOutcome.ID == 0 {
		return nil, errors.New("error sequence, can't get target sequence")
	}

	stageList := make([]models.StageLog, 0)
	new(models.StageLog).GetStageLog().Where("pipeline = ?", startOutcome.Pipeline).Order("id").Find(&stageList)

	// init pipeline id map
	idMap := make(map[string]string)
	for _, stageInfo := range stageList {
		idMap["s-"+strconv.FormatInt(stageInfo.FromStage, 10)] = "s-" + strconv.FormatInt(stageInfo.ID, 10)

		actionList := make([]models.ActionLog, 0)
		new(models.ActionLog).GetActionLog().Where("stage = ?", stageInfo.ID).Find(&actionList)

		for _, actionInfo := range actionList {
			idMap["a-"+strconv.FormatInt(actionInfo.FromAction, 10)] = "a-" + strconv.FormatInt(actionInfo.ID, 10)
		}
	}

	lineListMap := make([]map[string]interface{}, 0)

	stageListMap := make([]map[string]interface{}, 0)
	new(models.StageLog).GetStageLog().Where("pipeline = ?", startOutcome.Pipeline).Order("id").Find(&stageList)

	for _, stage := range stageList {
		stageInfoMap := make(map[string]interface{})

		stageSetupData := make(map[string]interface{})
		stageSetupData["name"] = stage.Stage

		stageType := "pipeline-stage"
		if stage.Type == models.StageTypeStart {
			stageType = "pipeline-start"
		} else if stage.Type == models.StageTypeEnd {
			stageType = "pipeline-end"
		}

		stagetStatus := true
		actionList := make([]models.ActionLog, 0)
		actionListMap := make([]map[string]interface{}, 0)
		new(models.ActionLog).GetActionLog().Where("stage = ?", stage.ID).Find(&actionList)

		for _, action := range actionList {
			actionInfoMap := make(map[string]interface{})

			actionStatus := true

			actionSetupData := make(map[string]interface{})
			actionSetupData["name"] = action.Action

			actionType := "pipeline-action"

			actionInfoMap["stetupData"] = actionSetupData
			actionInfoMap["id"] = "a-" + strconv.FormatInt(action.ID, 10)
			actionInfoMap["type"] = actionType

			actionOutput := new(models.Outcome)
			actionOutput.GetOutcome().Where("action = ?", action.ID).First(actionOutput)

			if actionOutput.ID == 0 && !actionOutput.Status {
				actionStatus = false
				stagetStatus = false
			}

			actionInfoMap["status"] = actionStatus

			actionListMap = append(actionListMap, actionInfoMap)

			// if current action has relation ,add line info to lineList
			if action.Manifest != "" {
				manifestMap := make(map[string]interface{})

				json.Unmarshal([]byte(action.Manifest), &manifestMap)
				if relationMap, ok := manifestMap["relation"].([]interface{}); ok {
					for _, relation := range relationMap {
						if _, ok := relation.(map[string]interface{}); !ok {
							continue
						}

						relationInfo := relation.(map[string]interface{})

						fromActionIdF := relationInfo["fromAction"].(float64)
						toActionIdF := relationInfo["toAction"].(float64)

						fromActionId := int64(fromActionIdF)
						toActionId := int64(toActionIdF)

						startDataMap := make(map[string]interface{})
						endDataMap := make(map[string]interface{})

						if fromActionId == 0 {
							// get start stage info
							startStage := new(models.StageLog)
							startStage.GetStageLog().Where("pipeline = ?", pipelineInfo.ID).Where("type = ? ", models.StageTypeStart).First(startStage)
							startDataMap["id"] = "s-" + strconv.FormatInt(startStage.ID, 10)
							startDataMap["type"] = "pipeline-start"
						} else {
							actionInfo := new(models.ActionLog)
							actionInfo.GetActionLog().Where("id = ?", strings.TrimPrefix(idMap["a-"+strconv.FormatInt(fromActionId, 10)], "a-")).First(actionInfo)

							startDataMap["id"] = "a-" + strconv.FormatInt(actionInfo.ID, 10)
							startDataMap["type"] = "pipeline-action"
							startDataMap["setupData"] = map[string]interface{}{"action": map[string]interface{}{"name": actionInfo.Action}}
						}

						endActionInfo := new(models.ActionLog)
						endActionInfo.GetActionLog().Where("id = ?", strings.TrimPrefix(idMap["a-"+strconv.FormatInt(toActionId, 10)], "a-")).First(endActionInfo)

						endDataMap["id"] = "a-" + strconv.FormatInt(endActionInfo.ID, 10)
						endDataMap["type"] = "pipeline-action"
						endDataMap["setupData"] = map[string]interface{}{"action": map[string]interface{}{"name": endActionInfo.Action}}

						lineInfoMap := make(map[string]interface{})
						lineInfoMap["pipelineLineViewId"] = "pipeline-line-view"
						lineInfoMap["startData"] = startDataMap
						lineInfoMap["endData"] = endDataMap
						lineInfoMap["id"] = startDataMap["id"].(string) + "-" + endDataMap["id"].(string)

						lineListMap = append(lineListMap, lineInfoMap)
					}
				}
			}

		}

		if len(actionListMap) > 0 {
			stageInfoMap["actions"] = actionListMap
		}
		stageInfoMap["id"] = "s-" + strconv.FormatInt(stage.ID, 10)
		stageInfoMap["setupData"] = stageSetupData
		stageInfoMap["type"] = stageType
		stageInfoMap["status"] = stagetStatus

		if !stagetStatus {
			pipelineStatus = stagetStatus
		}

		stageListMap = append(stageListMap, stageInfoMap)
	}

	defineMap["stageList"] = stageListMap
	defineMap["status"] = pipelineStatus
	defineMap["lineList"] = lineListMap

	return defineMap, nil
}

func GetStageHistoryInfo(stageLogId int64) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	// get all actions that belong to current stage ,
	// return stage info and action{id:actionId, name:actionName, status:true/false}
	actionList := make([]models.ActionLog, 0)
	actionListMap := make([]map[string]interface{}, 0)
	stageInfo := new(models.StageLog)
	stageStatus := true

	stageInfo.GetStageLog().Where("id = ?", stageLogId).First(stageInfo)
	new(models.ActionLog).GetActionLog().Where("stage = ?", stageLogId).Find(&actionList)

	for _, action := range actionList {
		actionOutcome := new(models.Outcome)
		new(models.Outcome).GetOutcome().Where("action = ?", action.ID).First(actionOutcome)

		actionMap := make(map[string]interface{})
		actionMap["id"] = "a-" + strconv.FormatInt(action.ID, 10)
		actionMap["name"] = action.Action
		if actionOutcome.Status {
			actionMap["status"] = true
		} else {
			actionMap["status"] = false
			stageStatus = false
		}

		actionListMap = append(actionListMap, actionMap)
	}

	firstActionStartEvent := new(models.Event)
	leastActionStopEvent := new(models.Event)

	firstActionStartEvent.GetEvent().Where("stage = ?", stageInfo.ID).Where("title = ?", "COMPONENT_START").Order("created_at").First(firstActionStartEvent)
	leastActionStopEvent.GetEvent().Where("stage = ?", stageInfo.ID).Where("title = ?", "COMPONENT_STOP").Order("-created_at").First(leastActionStopEvent)

	stageRunTime := ""
	if firstActionStartEvent.ID != 0 {
		stageRunTime = firstActionStartEvent.CreatedAt.Format("2006-01-02 15:04:05")
		stageRunTime += " - "
		if leastActionStopEvent.ID != 0 {
			stageRunTime += leastActionStopEvent.CreatedAt.Format("2006-01-02 15:04:05")
		}
	}

	result["name"] = stageInfo.Stage
	result["status"] = stageStatus
	result["actions"] = actionListMap
	result["runTime"] = stageRunTime

	return result, nil
}

func GetActionHistoryInfo(actionLogId int64) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	inputInfo := new(models.Event)
	inputInfo.GetEvent().Where("action = ?", actionLogId).Where("title = ?", "SEND_DATA").First(inputInfo)

	outputInfo := new(models.Event)
	outputInfo.GetEvent().Where("action = ?", actionLogId).Where("title = ?", "TASK_RESULT").First(outputInfo)

	dataMap := make(map[string]interface{})
	dataMap["input"] = inputInfo.Payload
	dataMap["output"] = outputInfo.Payload

	logList := make([]models.Event, 0)
	new(models.Event).GetEvent().Where("action = ?", actionLogId).Order("id").Find(&logList)

	logListStr := make([]string, 0)
	for _, log := range logList {
		logStr := log.CreatedAt.Format("2006-01-02 15:04:05") + " -> " + log.Payload

		logListStr = append(logListStr, logStr)
	}

	result["data"] = dataMap
	result["logList"] = logListStr
	return result, nil
}
