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
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Huawei/containerops/pilotage/models"

	log "github.com/Sirupsen/logrus"
)

type Relation struct {
	From  string
	To    string
	Child []Relation
}

var (
	startPipelineChan  chan bool
	createPipelineChan chan bool
)

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

	return legal
}

func DoPipelineLog(pipelineInfo models.Pipeline) (*models.PipelineLog, error) {
	pipelineLog := new(models.PipelineLog)
	// use chan to make sure pipeline sequence is unique
	startPipelineChan <- true

	tempSequence := new(struct {
		Sequence int64
	})

	err := new(models.Outcome).GetOutcome().Table("outcome").Where("pipeline = ?", pipelineInfo.ID).Order("-sequence").First(&tempSequence).Error
	if err != nil {
		<-startPipelineChan
		return nil, errors.New("error when query outcome info by pipeline:" + err.Error())
	}

	pipelineSequence := tempSequence.Sequence + 1

	pipelineLog.Namespace = pipelineInfo.Namespace
	pipelineLog.Pipeline = pipelineInfo.Pipeline
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
	err = new(models.Stage).GetStage().Where("pipeline = ?", pipelineInfo.ID).Find(&stageList).Error
	if err != nil {
		return nil, errors.New("error when get stage infos by pipeline id:" + strconv.FormatInt(pipelineInfo.ID, 10))
	}

	for _, stageInfo := range stageList {
		stageLog := new(models.StageLog)

		stageLog.Pipeline = stageInfo.Pipeline
		stageLog.Type = stageInfo.Type
		stageLog.PreStage = stageInfo.PreStage
		stageLog.Stage = stageInfo.Stage
		stageLog.Title = stageInfo.Title
		stageLog.Description = stageInfo.Description
		stageLog.Event = stageInfo.Event
		stageLog.Manifest = stageInfo.Manifest

		err = stageLog.GetStageLog().Save(stageLog).Error
		if err != nil {
			return nil, errors.New("error when create new stage log" + err.Error())
		}
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

			actionLog.Stage = actionInfo.Stage
			actionLog.Component = actionInfo.Component
			actionLog.Service = actionInfo.Service
			actionLog.Action = actionInfo.Action
			actionLog.Title = actionInfo.Title
			actionLog.Description = actionInfo.Description
			actionLog.Event = actionInfo.Event
			actionLog.Manifest = actionInfo.Manifest
			actionLog.Input = actionInfo.Input
			actionLog.Output = actionInfo.Output

			err = actionLog.GetActionLog().Save(actionLog).Error
			if err != nil {
				return nil, errors.New("error when create new action log:" + err.Error())
			}

			if actionInfo.Component != 0 {
				componentInfo := new(models.Component)
				err = componentInfo.GetComponent().Where("namespace = ?", pipelineInfo.Namespace).Where("id = ?", actionInfo.Component).First(&componentInfo).Error
				if err != nil {
					return nil, errors.New("error when get component info by id:" + strconv.FormatInt(actionInfo.Component, 10))
				}

				componentLog := new(models.ComponentLog)

				componentLog.Namespace = componentInfo.Namespace
				componentLog.Version = componentInfo.Version
				componentLog.VersionCode = componentInfo.VersionCode
				componentLog.Component = componentInfo.Component
				componentLog.Type = componentInfo.Type
				componentLog.Title = componentInfo.Title
				componentLog.Gravatar = componentInfo.Gravatar
				componentLog.Description = componentInfo.Description
				componentLog.Endpoint = componentInfo.Endpoint
				componentLog.Source = componentInfo.Source
				componentLog.Environment = componentInfo.Environment
				componentLog.Tag = componentInfo.Tag
				componentLog.VolumeLocation = componentInfo.VolumeLocation
				componentLog.VolumeData = componentInfo.VolumeData
				componentLog.Makefile = componentInfo.Makefile
				componentLog.Kubernetes = componentInfo.Kubernetes
				componentLog.Swarm = componentInfo.Swarm
				componentLog.Input = componentInfo.Input
				componentLog.Output = componentInfo.Output

				err = componentLog.GetComponentLog().Save(componentLog).Error
				if err != nil {
					return nil, errors.New("error when create new component log:" + err.Error())
				}
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
	case "Github", "Manual":
		mac := hmac.New(sha1.New, []byte(secretKey))
		mac.Write(reqBody)
		expectedMAC := mac.Sum(nil)
		expectedSig := "sha1=" + hex.EncodeToString(expectedMAC)

		if expectedSig == token {
			legal = true
		}
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
	startOutcome.Stage = startStage.ID
	startOutcome.Action = 0
	startOutcome.Event = 0
	startOutcome.Status = false
	startOutcome.Result = "success"
	startOutcome.Output = reqBody

	err = startOutcome.GetOutcome().Save(startOutcome).Error
	if err != nil {
		return "error when save pipeline start data:" + err.Error()
	}

	go handleStage(pipelineInfo, *startStage, pipelineSequence)
	return "pipeline start ..."
}

// handler a stage
func handleStage(pipelineInfo models.PipelineLog, stageInfo models.StageLog, pipelineSequence int64) {
	nextStage := new(models.StageLog)
	actionList := make([]models.ActionLog, 0)

	// if current stage is nil or is end stage, stop run
	if stageInfo.ID == 0 || stageInfo.Type == models.StageTypeEnd {
		log.Error("current stage is ", stageInfo, " now pipeline("+strconv.FormatInt(pipelineSequence, 10)+") is finish...")
		return
	}

	new(models.ActionLog).GetActionLog().Where("stage = ?", stageInfo.ID).Find(&actionList)

	if stageInfo.PreStage != -1 && len(actionList) > 0 {
		// exec all action
		for _, action := range actionList {
			go execAction(pipelineInfo, stageInfo, action, pipelineSequence)
		}
	}

	isAllActionOk := waitAllActionFinish(actionList, pipelineSequence)
	if isAllActionOk {
		nextStage.GetStageLog().Where("pipeline = ?", pipelineInfo.ID).Where("pre_stage = ?", stageInfo.ID).First(&nextStage)
		handleStage(pipelineInfo, *nextStage, pipelineSequence)
	} else {
		// if has a failer action ,then stop all other action's
		for _, action := range actionList {
			if action.Component != 0 {
				go stopComponent(pipelineInfo, stageInfo, action, pipelineSequence)
			} else {
				go stopService(pipelineInfo, stageInfo, action, pipelineSequence)
			}
		}
	}
}

// exec a action
func execAction(pipelineInfo models.PipelineLog, stageInfo models.StageLog, actionInfo models.ActionLog, pipelineSequence int64) {
	if actionInfo.Component != 0 {
		startComponent(pipelineInfo, stageInfo, actionInfo, pipelineSequence)
	} else {
		startService(pipelineInfo, stageInfo, actionInfo, pipelineSequence)
	}
}

func startComponent(pipelineInfo models.PipelineLog, stageInfo models.StageLog, actionInfo models.ActionLog, pipelineSequence int64) {
	componentInfo := new(models.ComponentLog)
	componentInfo.GetComponentLog().Where("id = ?", actionInfo.Component).First(componentInfo)

	// get all event that bind this action
	eventList := make([]models.EventDefinition, 0)
	new(models.EventDefinition).GetEventDefinition().Where("action = ?", actionInfo.ID).Find(&eventList)

	// now component default run in k8s
	// TODO :component run in swarm etc.
	// componentId = pipelineId + stageId + actionId + pipelineSequence + componentId
	componentId := strconv.FormatInt(pipelineInfo.ID, 10) + "," + strconv.FormatInt(stageInfo.ID, 10) + "," + strconv.FormatInt(actionInfo.ID, 10) + "," + strconv.FormatInt(pipelineSequence, 10) + "," + strconv.FormatInt(componentInfo.ID, 10)
	c, err := InitComponet(*componentInfo, RUNENV_KUBE)
	if err != nil {
		// if has init error,stop this action and log it as start error
		startErrOutcome := new(models.Outcome)
		startErrOutcome.Pipeline = pipelineInfo.ID
		startErrOutcome.Stage = stageInfo.ID
		startErrOutcome.Action = actionInfo.ID
		startErrOutcome.Event = 0
		startErrOutcome.Sequence = pipelineSequence
		startErrOutcome.Status = false
		startErrOutcome.Result = "init error:" + err.Error()
		startErrOutcome.Output = ""

		startErrOutcome.GetOutcome().Save(startErrOutcome)
		return
	}
	err = c.Start(componentId, eventList)
	if err != nil {
		// if has start error,stop this action and log it as start error
		startErrOutcome := new(models.Outcome)
		startErrOutcome.Pipeline = pipelineInfo.ID
		startErrOutcome.Stage = stageInfo.ID
		startErrOutcome.Action = actionInfo.ID
		startErrOutcome.Event = 0
		startErrOutcome.Sequence = pipelineSequence
		startErrOutcome.Status = false
		startErrOutcome.Result = "start error:" + err.Error()
		startErrOutcome.Output = ""

		startErrOutcome.GetOutcome().Save(startErrOutcome)
	}
}

func stopComponent(pipelineInfo models.PipelineLog, stageInfo models.StageLog, actionInfo models.ActionLog, pipelineSequence int64) {
	componentInfo := new(models.ComponentLog)
	componentInfo.GetComponentLog().Where("id = ?", actionInfo.Component).First(componentInfo)

	c, err := InitComponet(*componentInfo, RUNENV_KUBE)
	if err != nil {
		// if has init error,stop this action and log it as start error
		initErrOutcome := new(models.Outcome)
		initErrOutcome.Pipeline = pipelineInfo.ID
		initErrOutcome.Stage = stageInfo.ID
		initErrOutcome.Action = actionInfo.ID
		initErrOutcome.Event = 0
		initErrOutcome.Sequence = pipelineSequence
		initErrOutcome.Status = false
		initErrOutcome.Result = "component init error:" + err.Error()
		initErrOutcome.Output = ""

		initErrOutcome.GetOutcome().Save(initErrOutcome)
		return
	}

	componentId := strconv.FormatInt(pipelineInfo.ID, 10) + "," + strconv.FormatInt(stageInfo.ID, 10) + "," + strconv.FormatInt(actionInfo.ID, 10) + "," + strconv.FormatInt(pipelineSequence, 10) + "," + strconv.FormatInt(componentInfo.ID, 10)
	c.Stop(componentId)
}

// TODO : start a service
func startService(pipelineInfo models.PipelineLog, stageInfo models.StageLog, actionInfo models.ActionLog, pipelineSequence int64) {
	startErrOutcome := new(models.Outcome)
	startErrOutcome.Pipeline = pipelineInfo.ID
	startErrOutcome.Stage = stageInfo.ID
	startErrOutcome.Action = actionInfo.ID
	startErrOutcome.Event = 0
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
func waitAllActionFinish(actionList []models.ActionLog, Sequence int64) bool {
	allActionIsOk := true

	actionIds := make([]int64, 0)
	for _, action := range actionList {
		actionIds = append(actionIds, action.ID)
	}

	for {
		time.Sleep(1 * time.Second)

		runResults := make([]struct {
			Result bool
		}, 0)

		new(models.Outcome).GetOutcome().Table("outcome").Where("sequence = ?", Sequence).Where("action in (?)", actionIds).Find(&runResults)

		for _, runResult := range runResults {
			if !runResult.Result {
				allActionIsOk = false
				break
			}
		}

		if len(runResults) == len(actionList) || !allActionIsOk {
			break
		}
	}

	return allActionIsOk
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

	dataMap := make(map[string]interface{})
	relations, ok := manifestMap["relations"]
	if ok {
		relationInfo, ok := relations.([]interface{})
		if !ok {
			log.Error("error when parse relations,relations is not an array")
			return
		}

		// get all data that current action is require
		dataMap, err = merageFromActionsOutputData(pipelineId, stageInfo.PreStage, actionId, pipelineSequence, componentId, relationInfo)
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
		payload["resp"] = resp
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
		err := fromOutcome.GetOutcome().Where("pipeline = ? ", pipelineId).Where("stage = ?", stageId).Where("action = ?", fromAction).Where("sequence = ?", pipelineSequence).First(&fromOutcome).Error
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

		relationByte, err := json.Marshal(relationArray)
		if err != nil {
			return nil, errors.New("error when marshal relation array:" + err.Error())
		}

		relationList := make([]Relation, 0)

		err = json.Unmarshal(relationByte, &relationList)
		if err != nil {
			return nil, errors.New("error when parse relation info" + err.Error())
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
		if len(relation.Child) == 0 {
			// get data from current relation path
			fromData, err := getJsonDataByPath(strings.TrimPrefix(relation.From, "."), fromActionData)
			if err != nil {
				return errors.New("error when get fromData :" + err.Error())
			}

			setDataToMapByPath(fromData, result, strings.TrimPrefix(relation.To, "."))
		} else {
			getResultFromRelation(fromActionOutput, relation.Child, result)
		}
	}

	return nil
}

// getJsonDataByPath is get a value from a map by give path
func getJsonDataByPath(path string, data map[string]interface{}) (interface{}, error) {
	depth := len(strings.Split(path, "."))
	if depth == 1 {
		if info, ok := data[path]; !ok {
			return nil, errors.New("key not exist:" + path)
		} else {
			return info, nil
		}
	}

	childDataInterface, ok := data[strings.Split(path, ".")[0]]
	if !ok {
		return nil, errors.New("key not exist:" + path)
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
	componentInfo := new(models.ComponentLog)
	componentInfo.GetComponentLog().Where("id = ?", componentId).First(componentInfo)
	c, err := InitComponet(*componentInfo, RUNENV_KUBE)
	if err != nil {
		return nil, err
	}

	ip, err := c.GetIp(podName)
	if err != nil {
		return nil, err
	}

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

func GetPipelineList(namespace, pipelineName string, pipelineId int64) (map[string]interface{}, error) {
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
				return defineInfoMap, nil
			}
		}
	}

	// get all stage info of current pipeline
	stageList, err := getStageListByPipeline(*pipelineInfo)
	if err != nil {
		return nil, err
	}
	resultMap["stageList"] = stageList

	lineList, err := getLineListByPipeline(*pipelineInfo)
	if err != nil {
		return nil, err
	}
	resultMap["lineList"] = lineList

	// resultMap["stageList"] = stageListMap

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

func getLineListByPipeline(pipelineInfo models.Pipeline) (interface{}, error) {
	return make([]map[string]interface{}, 0), nil
}
