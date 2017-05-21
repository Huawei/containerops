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

package module

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Huawei/containerops/pilotage/models"
	"github.com/Huawei/containerops/pilotage/utils"

	log "github.com/Sirupsen/logrus"
	"github.com/jinzhu/gorm"
)

const (
	// ActionStopReasonTimeout is
	ActionStopReasonTimeout = "TIME_OUT"

	// ActionStopReasonSendDataFailed is
	ActionStopReasonSendDataFailed = "SEND_DATA_FAILED"

	// ActionStopReasonRunSuccess is
	ActionStopReasonRunSuccess = "ACTION_RUN_SUCCESS"
	// ActionStopReasonRunFailed is
	ActionStopReasonRunFailed = "ACTION_RUN_FAILED"
)

var (
	actionlogAuthChan         chan bool
	actionlogListenChan       chan bool
	actionlogSetGlobalVarChan chan bool
)

// Action is
type Action struct {
	*models.Action
}

// ActionLog is
type ActionLog struct {
	*models.ActionLog
}

// Relation is
type Relation struct {
	From string
	To   string
}

func init() {
	actionlogAuthChan = make(chan bool, 1)
	actionlogListenChan = make(chan bool, 1)
	actionlogSetGlobalVarChan = make(chan bool, 1)
}

func getActionEnvList(actionLogId int64) ([]map[string]interface{}, error) {
	resultList := make([]map[string]interface{}, 0)
	actionLog := new(models.ActionLog)
	err := actionLog.GetActionLog().Where("id = ?", actionLogId).First(actionLog).Error
	if err != nil {
		log.Error("[actionLog's getActionEnvList]:error when get actionLog info from db:", err.Error())
		return nil, errors.New("error when get action info from db:" + err.Error())
	}

	envMap := make(map[string]string)
	if actionLog.Environment != "" {
		err = json.Unmarshal([]byte(actionLog.Environment), &envMap)
		if err != nil {
			log.Error("[actionLog's getActionEnvList]:error when unmarshal action's env setting:", actionLog.Environment, " ===>error is:", err.Error())
			return nil, errors.New("error when unmarshal action's env info" + err.Error())
		}
	}

	for key, value := range envMap {
		tempEnvMap := make(map[string]interface{})
		tempEnvMap["name"] = key
		tempEnvMap["value"] = value

		resultList = append(resultList, tempEnvMap)
	}

	return resultList, nil
}

// CreateNewActions is
func CreateNewActions(db *gorm.DB, workflowInfo *models.Workflow, stageInfo *models.Stage, defineList []map[string]interface{}) (map[string]int64, error) {
	if db == nil {
		db = models.GetDB()
		db = db.Begin()
	}

	actionIdMap := make(map[string]int64)
	for _, actionDefine := range defineList {
		actionName := ""
		imageName := ""
		imageTag := ""
		kubernetesSetting := ""
		inputStr := ""
		outputStr := ""
		actionTimeout := strconv.FormatInt(int64(60*60*24*36), 10)
		componentId := int64(0)
		serviceId := int64(0)
		platformMap := make(map[string]string)
		requestMapList := make([]interface{}, 0)

		// get component info
		component, ok := actionDefine["component"]
		if ok {
			componentMap, ok := component.(map[string]interface{})
			if !ok {
				log.Error("[action's CreateNewActions]:error when get action's component info, want a json obj, got:", component)
				return nil, errors.New("action's component is not a json")
			}

			componentVersion, ok := componentMap["versionid"].(float64)
			if !ok {
				log.Error("[action's CreateNewActions]:error when get action's component info,compoent doesn't has a versionid,component define is:", componentMap)
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
					if name, ok := image["name"].(string); ok {
						imageName = name
					}
					if tag, ok := image["tag"].(string); ok {
						imageTag = tag
					}
				}

				defineTimeoutStr, ok := actionSetupDataMap["timeout"].(string)
				if ok && !strings.Contains(defineTimeoutStr, "@") && !strings.Contains(defineTimeoutStr, "@") {
					timeoutInt, err := strconv.ParseInt(defineTimeoutStr, 10, 64)
					if err != nil {
						actionTimeout = "0"
					} else {
						actionTimeout = strconv.FormatInt(timeoutInt, 10)
					}
				} else if ok {
					actionTimeout = actionSetupDataMap["timeout"].(string)
				}

				configMap := make(map[string]interface{})
				// record platform info
				if platFormType, ok := actionSetupDataMap["type"].(string); ok {
					platformMap["platformType"] = strings.ToUpper(platFormType)
				}

				if platformHost, ok := actionSetupDataMap["apiserver"].(string); ok {
					platformHost = strings.TrimSuffix(platformHost, "/")
					platformMap["platformHost"] = platformHost
				}

				if ip, ok := actionSetupDataMap["ip"].(string); ok {
					configMap["nodeIP"] = ip
				}

				// unmarshal k8s info
				if useAdvanced, ok := actionSetupDataMap["useAdvanced"].(bool); ok {
					podConfigKey := "pod"
					serviceConfigKey := "service"
					if useAdvanced {
						configMap["useAdvanced"] = true
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
				log.Error("[action's CreateNewActions]:error when describe action's input json define:", inputMap, " ===>error is:", err.Error())
				return nil, errors.New("error in component output json define:" + err.Error())
			}

			inputDescBytes, _ := json.Marshal(inputDescribe)
			inputStr = string(inputDescBytes)
		}

		outputMap, ok := actionDefine["outputJson"].(map[string]interface{})
		if ok {
			outputDescribe, err := describeJSON(outputMap, "")
			if err != nil {
				log.Error("[action's CreateNewActions]:error when describe action's output json define:", inputMap, " ===>error is:", err.Error())
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
					log.Error("[action's CreateNewActions]:error when get action's env setting, want a json obj,got:", envInfo)
					return nil, errors.New("action's env set is not a json")
				}

				key, ok := envInfoMap["key"].(string)
				if !ok {
					log.Error("[action's CreateNewActions]:error when get action's env setting, want string key,got:", envInfoMap)
					return nil, errors.New("action's key is not a string")
				}

				value, ok := envInfoMap["value"].(string)
				if !ok {
					log.Error("[action's CreateNewActions]:error when get action's env setting, want string value,got:", envInfoMap)
					return nil, errors.New("action's value is not a string")
				}
				allEnvMap[key] = value
			}
		}
		envBytes, _ := json.Marshal(allEnvMap)

		// get aciont line info
		actionId, ok := actionDefine["id"].(string)
		if !ok {
			log.Error("[action's CreateNewActions]:error when action's id from action define, want string, got:", actionDefine)
			return nil, errors.New("action's id is not a string")
		}

		manifestMap := make(map[string]interface{})
		manifestMap["platform"] = platformMap
		manifestBytes, _ := json.Marshal(manifestMap)

		stageRequest, ok := actionDefine["request"].([]interface{})
		if !ok {
			defaultRequestMap := make(map[string]interface{})
			defaultRequestMap["type"] = AuthTyptStageStartDone
			defaultRequestMap["token"] = AuthTokenDefault

			requestMapList = append(requestMapList, defaultRequestMap)
		} else {
			requestMapList = stageRequest
		}
		requestInfos, _ := json.Marshal(requestMapList)

		action := new(models.Action)
		action.Namespace = workflowInfo.Namespace
		action.Repository = workflowInfo.Repository
		action.Workflow = stageInfo.Workflow
		action.Stage = stageInfo.ID
		action.Component = componentId
		action.Service = serviceId
		action.Action = actionName
		action.Title = actionName
		action.Description = actionName
		action.Manifest = string(manifestBytes)
		action.Environment = string(envBytes)
		action.Kubernetes = kubernetesSetting
		action.Input = inputStr
		action.Output = outputStr
		action.ImageName = imageName
		action.ImageTag = imageTag
		action.Timeout = actionTimeout
		action.Requires = string(requestInfos)

		err := db.Model(&models.Action{}).Save(action).Error
		if err != nil {
			log.Error("[action's CreateNewActions]:error when save action info to db:", err.Error())
			return nil, errors.New("error when save action info to db:" + err.Error())
		}
		actionIdMap[actionId] = action.ID
	}

	return actionIdMap, nil
}

// GetActionLog is
func GetActionLog(actionLogId int64) (*ActionLog, error) {
	action := new(ActionLog)
	actionLog := new(models.ActionLog)
	err := actionLog.GetActionLog().Where("id = ?", actionLogId).First(actionLog).Error
	if err != nil {
		log.Error("[actionLog's GetActionLog]:error when get action log info from db:", err.Error())
		return nil, err
	}

	action.ActionLog = actionLog
	return action, nil
}

// GetActionLogByName is
func GetActionLogByName(namespace, repository, workflowName string, sequence int64, stageName, actionName string) (*ActionLog, error) {
	action := new(ActionLog)
	workflowLog := new(models.WorkflowLog)
	stageLog := new(models.StageLog)
	actionLog := new(models.ActionLog)

	err := workflowLog.GetWorkflowLog().Where("namespace = ?", namespace).Where("repository = ?", repository).Where("workflow = ?", workflowName).Where("sequence = ?", sequence).First(workflowLog).Error
	if err != nil {
		if err != nil {
			log.Error("[actionLog's GetActionLog]:error when get workflowLog info from db:", err.Error())
			return nil, err
		}
	}

	err = stageLog.GetStageLog().Where("namespace = ?", namespace).Where("repository = ?", repository).Where("workflow = ?", workflowLog.ID).Where("sequence = ?", sequence).Where("stage = ?", stageName).First(stageLog).Error
	if err != nil {
		if err != nil {
			log.Error("[actionLog's GetActionLog]:error when get stageLog info from db:", err.Error())
			return nil, err
		}
	}

	err = actionLog.GetActionLog().Where("namespace = ?", namespace).Where("repository = ?", repository).Where("workflow = ?", workflowLog.ID).Where("sequence = ?", sequence).Where("stage = ?", stageLog.ID).Where("action = ?", actionName).First(actionLog).Error
	if err != nil {
		log.Error("[actionLog's GetActionLog]:error when get action log info from db:", err.Error())
		return nil, err
	}

	action.ActionLog = actionLog
	return action, nil
}

// GenerateNewLog is
func (actionInfo *Action) GenerateNewLog(db *gorm.DB, workflowLog *models.WorkflowLog, stageLog *models.StageLog) error {
	if db == nil {
		db = models.GetDB().Begin()
		err := db.Error
		if err != nil {
			log.Error("[action's GenerateNewLog]:when db.Begin():", err.Error())
			return err
		}
	}

	// record action's info
	actionLog := new(models.ActionLog)
	actionLog.Namespace = actionInfo.Namespace
	actionLog.Repository = actionInfo.Repository
	actionLog.Workflow = workflowLog.ID
	actionLog.FromWorkflow = workflowLog.FromWorkflow
	actionLog.Sequence = workflowLog.Sequence
	actionLog.Stage = stageLog.ID
	actionLog.FromStage = stageLog.FromStage
	actionLog.FromAction = actionInfo.ID
	actionLog.RunState = models.ActionLogStateCanListen
	actionLog.Component = actionInfo.Component
	actionLog.Service = actionInfo.Service
	actionLog.Action = actionInfo.Action.Action
	actionLog.Title = actionInfo.Title
	actionLog.Description = actionInfo.Description
	actionLog.Event = actionInfo.Event
	actionLog.Manifest = actionInfo.Manifest
	actionLog.Environment = actionInfo.Environment
	actionLog.Kubernetes = actionInfo.Kubernetes
	actionLog.Swarm = actionInfo.Swarm
	actionLog.Input = actionInfo.Input
	actionLog.Output = actionInfo.Output
	actionLog.ImageName = actionInfo.ImageName
	actionLog.ImageTag = actionInfo.ImageTag
	actionLog.Timeout = actionInfo.Timeout
	actionLog.Requires = actionInfo.Requires
	actionLog.AuthList = ""

	err := db.Save(actionLog).Error
	if err != nil {
		log.Error("[action's GenerateNewLog]:when save action log to db:", actionLog, " ===>error is:", err.Error())
		rollbackErr := db.Rollback().Error
		if rollbackErr != nil {
			log.Error("[action's GenerateNewLog]:when rollback in save action log:", rollbackErr.Error())
			return errors.New("errors occur:\nerror1:" + err.Error() + "\nerror2:" + rollbackErr.Error())
		}
		return err
	}

	err = setSystemEvent(db, actionLog)
	if err != nil {
		log.Error("[action's GenerateNewLog]:when save action log to db:", err.Error())
		return err
	}

	return nil
}

// GetActionLineInfo is
func (actionLog *ActionLog) GetActionLineInfo() ([]map[string]interface{}, error) {
	lineList := make([]map[string]interface{}, 0)

	manifestMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(actionLog.Manifest), &manifestMap)
	if err != nil {
		log.Error("[actionLog's GetActionLineInfo]:error when unmarshal action's manifest, want a json obj,got:", actionLog.Manifest)
		return nil, err
	}

	relationList, ok := manifestMap["relation"].([]interface{})
	if !ok {
		log.Error("[actionLog's GetActionLineInfo]:error when get action's relation from action's manifestMap:", manifestMap)
		return lineList, nil
	}

	for _, relation := range relationList {
		relationInfo, ok := relation.(map[string]interface{})
		if err != nil {
			log.Error("[actionLog's GetActionLineInfo]:error when get action's relation info, want a json obj,got:", relation)
			continue
		}

		fromRealActionIDF, ok := relationInfo["fromAction"].(float64)
		if !ok {
			log.Error("[actionLog's GetActionLineInfo]:error when get fromRealActionID from action's relation,want a number,got:", relationInfo["fromAction"])
			continue
		}

		fromRealActionID := int64(fromRealActionIDF)
		fromActionInfoMap := make(map[string]string)
		if fromRealActionID == int64(0) {
			// if action's id == 0 ,this is a relation from workflow's start stage
			startStage := new(models.StageLog)
			err := startStage.GetStageLog().Where("namespace = ?", actionLog.Namespace).Where("repository = ?", actionLog.Repository).Where("workflow = ?", actionLog.Workflow).Where("type = ?", models.StageTypeStart).First(startStage).Error
			if err != nil {
				log.Error("[actionLog's GetActionLineInfo]:error when get pipline's start stage from db:", err.Error())
				continue
			}

			fromActionInfoMap["id"] = "s-" + strconv.FormatInt(startStage.ID, 10)
			fromActionInfoMap["type"] = models.StageTypeForWeb[startStage.Type]
		} else {
			fromActionInfo := new(models.ActionLog)
			err = fromActionInfo.GetActionLog().Where("namespace = ?", actionLog.Namespace).Where("repository = ?", actionLog.Repository).Where("workflow = ?", actionLog.Workflow).Where("sequence = ?", actionLog.Sequence).Where("from_action = ?", fromRealActionID).First(fromActionInfo).Error
			if err != nil {
				log.Error("[actionLog's GetActionLineInfo]:error when get preActionlog info from db:", err.Error())
				continue
			}

			fromActionInfoMap["id"] = "a-" + strconv.FormatInt(fromActionInfo.ID, 10)
			fromActionInfoMap["type"] = "workflow-action"
		}

		toActionInfoMap := make(map[string]string)
		toActionInfoMap["id"] = "a-" + strconv.FormatInt(actionLog.ID, 10)
		toActionInfoMap["type"] = "workflow-action"
		toActionInfoMap["name"] = actionLog.Action

		lineMap := make(map[string]interface{})
		lineMap["id"] = fromActionInfoMap["id"] + "-" + toActionInfoMap["id"]
		lineMap["workflowLineViewId"] = "workflow-line-view"

		lineMap["startData"] = map[string]string{
			"id":   fromActionInfoMap["id"],
			"type": fromActionInfoMap["type"],
		}

		lineMap["endData"] = map[string]interface{}{
			"id": toActionInfoMap["id"],
			"setupData": map[string]interface{}{
				"action": map[string]string{
					"name": toActionInfoMap["name"],
				},
			},
		}

		lineList = append(lineList, lineMap)
	}

	return lineList, nil
}

// GetActionHistoryInfo is
func (actionLog *ActionLog) GetActionHistoryInfo() (map[string]interface{}, error) {
	result := make(map[string]interface{})

	inputMap, err := actionLog.GetInputData()
	if err != nil {
		log.Error("[actionLog's GetActionHistoryInfo]:error when get action's input data:", err.Error())
		return nil, err
	}

	outputMap, err := actionLog.GetOutputData()
	if err != nil {
		log.Error("[actionLog's GetActionHistoryInfo]:error when get action's output data:", err.Error())
		return nil, err
	}

	dataMap := make(map[string]interface{})
	dataMap["input"] = inputMap
	dataMap["output"] = outputMap

	logList := make([]models.Event, 0)
	err = new(models.Event).GetEvent().Where("namespace = ?", actionLog.Namespace).Where("repository = ?", actionLog.Repository).Where("workflow = ?", actionLog.Workflow).Where("stage = ?", actionLog.Stage).Where("action = ?", actionLog.ID).Order("id").Find(&logList).Error
	if err != nil && !strings.Contains(err.Error(), "record not found") {
		log.Error("[actionLog's GetActionHistoryInfo]:error when get actionlog's log from db:", err.Error())
		return nil, err
	}

	logListStr := make([]string, 0)
	for _, log := range logList {
		logStr := log.CreatedAt.Format("2006-01-02 15:04:05") + " -> " + log.Payload

		logListStr = append(logListStr, logStr)
	}

	result["data"] = dataMap
	result["logList"] = logListStr
	return result, nil
}

// GetActionConsoleLog is
func (actionLog *ActionLog) GetActionConsoleLog(key string, size int64) (map[string]interface{}, error) {
	resultList := make([]map[string]interface{}, 0)
	resultKey := ""
	result := make(map[string]interface{})

	var kube component
	if actionLog.Component != 0 {
		c, err := InitComponetNew(actionLog)
		if err != nil {
			log.Error("[actionLog's GetActionConsoleLog]:error when init component:", err.Error())
			RecordOutcom(actionLog.Workflow, actionLog.FromWorkflow, actionLog.Stage, actionLog.FromStage, actionLog.ID, actionLog.FromAction, actionLog.Sequence, 0, false, "start action error", "error when init component:"+err.Error())
			actionLog.Stop(ActionStopReasonRunFailed, models.ActionLogStateRunFailed)
			return nil, nil
		} else {
			kube = c
		}
	}

	if actionLog.ContainerId == "" {
		temp := make(map[string]interface{})

		temp["log"] = "can't get action's container log"
		temp["stream"] = "workflow"
		temp["time"] = time.Now().Format("2006-01-02 15:04:05")

		resultList = append(resultList, temp)
		result["list"] = resultList

		kube.Update()

		return result, nil
	}

	if kube != nil {
		platformSetting, err := actionLog.GetActionPlatformInfo()
		if err != nil {
			log.Error("[actionLog's GetActionConsoleLog]:error when get given actionLog's platformSetting:", actionLog, " ===>error is:", err.Error())
			return nil, errors.New("error when get action's info")
		}

		logServerUrl := platformSetting["platformHost"] + "/api/v1/proxy/namespaces/kube-system/services/elasticsearch-logging/_search"
		logReqBody := []byte("")
		if key == "" {
			logServerUrl += "?scroll=10m"
			bodyMap := map[string]interface{}{
				"query": map[string]interface{}{
					"match": map[string]interface{}{
						"tag": map[string]string{
							"query": "kubernetes.var.log.containers." + actionLog.ContainerId,
							"type":  "phrase"}}},
				"size": size,
				"sort": "@timestamp"}

			logReqBody, _ = json.Marshal(bodyMap)
		} else {
			logServerUrl += "/scroll"
			bodyMap := map[string]string{
				"scroll":    "10m",
				"scroll_id": key}
			logReqBody, _ = json.Marshal(bodyMap)
		}

		log.Info("[actionLog's GetActionConsoleLog]:send req to ", logServerUrl, " req body is :", string(logReqBody))
		resp, err := http.Post(logServerUrl, "application/json", bytes.NewReader(logReqBody))
		if err != nil {
			go kube.Update()
			log.Error("[actionLog's GetActionConsoleLog]:error when get action's log:", err.Error())
			return nil, errors.New("error when get log info from server")
		}

		respBody, _ := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		respMap := make(map[string]interface{})
		json.Unmarshal(respBody, &respMap)
		isEnd := true
		if hits, ok := respMap["hits"].(map[string]interface{}); ok {
			if hitsList, ok := hits["hits"].([]interface{}); ok {
				for _, info := range hitsList {
					if hitMap, ok := info.(map[string]interface{}); ok {
						if sourceMap, ok := hitMap["_source"].(map[string]interface{}); ok {
							tempMap := make(map[string]interface{})
							tempMap["stream"] = sourceMap["stream"]
							tempMap["time"] = sourceMap["@timestamp"]
							tempMap["log"] = sourceMap["log"]

							resultList = append(resultList, tempMap)
						}
					}
				}

				if len(hitsList) == 0 && key != "" {
					reqBodyMap := map[string]interface{}{"scroll_id": []string{key}}
					client := &http.Client{}
					reqBody, _ := json.Marshal(reqBodyMap)
					req, _ := http.NewRequest(http.MethodDelete, logServerUrl, bytes.NewReader(reqBody))
					resp, _ := client.Do(req)
					defer resp.Body.Close()
				} else {
					isEnd = false
				}
			}
		}

		if scrollID, ok := respMap["_scroll_id"].(string); ok && !isEnd {
			resultKey = scrollID
		}
	}

	result["list"] = resultList
	result["key"] = resultKey
	return result, nil
}

// GetInputData is
func (actionLog *ActionLog) GetInputData() (map[string]interface{}, error) {
	inputMap := make(map[string]interface{})

	inputInfo := new(models.Event)
	err := inputInfo.GetEvent().Where("namespace = ?", actionLog.Namespace).Where("repository = ?", actionLog.Repository).Where("workflow = ?", actionLog.Workflow).Where("sequence = ?", actionLog.Sequence).Where("action = ?", actionLog.ID).Where("title = ?", "SEND_DATA").First(inputInfo).Error
	if err != nil && !strings.Contains(err.Error(), "record not found") {
		log.Error("[actionLog's GetInputData]:error when get actionlog's input info from db:", err.Error())
		return nil, err
	}

	err = json.Unmarshal([]byte(inputInfo.Payload), &inputMap)
	if err != nil {
		log.Error("[actionLog's GetInputData]:error when unmarshal input info:", inputInfo.Payload, " ===>error is:"+err.Error())
	}

	inputStr, ok := inputMap["data"].(string)
	if !ok {
		log.Error("[actionLog's GetInputData]:error when get inputMap str, want a string, got:", inputMap["data"])
		inputStr = ""
	}

	realinputMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(inputStr), &realinputMap)
	if err != nil {
		log.Error("[actionLog's GetInputData]:error when unmarshal real input info:", inputStr, "===>error is:", err.Error())
	}

	return inputMap, nil
}

// GetOutputData is
func (actionLog *ActionLog) GetOutputData() (map[string]interface{}, error) {
	outputMap := make(map[string]interface{})

	outputInfo := new(models.Outcome)
	err := outputInfo.GetOutcome().Where("workflow = ?", actionLog.Workflow).Where("sequence = ?", actionLog.Sequence).Where("stage = ?", actionLog.Stage).Where("action = ?", actionLog.ID).First(outputInfo).Error
	if err != nil && !strings.Contains(err.Error(), "record not found") {
		log.Error("[actionLog's GetOutputData]:error when get actionlog's output info from db:", err.Error())
		return nil, err
	}

	err = json.Unmarshal([]byte(outputInfo.Output), &outputMap)
	if err != nil {
		log.Error("[actionLog's GetOutputData]:error when unmarshal output info:", outputInfo.Output, "===>error is:", err.Error())
	}

	return outputMap, nil
}

// Listen is
func (actionLog *ActionLog) Listen() error {
	actionlogListenChan <- true
	defer func() { <-actionlogListenChan }()

	err := actionLog.GetActionLog().Where("id = ?", actionLog.ID).First(actionLog).Error
	if err != nil {
		log.Error("[actionLog's Listen]:error when get action info from db:", actionLog, " ===>error is:", err.Error())
		return errors.New("error when get actionlog's info from db:" + err.Error())
	}

	if actionLog.RunState != models.ActionLogStateCanListen {
		log.Error("[actionLog's Listen]:error actionlog state:", actionLog)
		return errors.New("can't listen curren actionlog,current state is:" + strconv.FormatInt(actionLog.RunState, 10))
	}

	actionLog.RunState = models.ActionLogStateWaitToStart
	err = actionLog.GetActionLog().Save(actionLog).Error
	if err != nil {
		log.Error("[actionLog's Listen]:error when change actionLog's run state to wait to start:", actionLog, " ===>error is:", err.Error())
		return errors.New("can't listen target action,change action's state failed")
	}

	canStartChan := make(chan models.ActionLog, 1)
	go func() {
		aLog := *actionLog.ActionLog
		for true {
			time.Sleep(1 * time.Second)

			err := aLog.GetActionLog().Where("id = ?", aLog.ID).First(&aLog).Error
			if err != nil {
				log.Error("[actionLog's Listen]:error when get actionLog's info:", aLog, " ===>error is:", err.Error())
				canStartChan <- *new(models.ActionLog)
				break
			}
			if aLog.Requires == "" || aLog.Requires == "[]" {
				canStartChan <- aLog
				break
			}
		}
	}()

	go func() {
		aLog := <-canStartChan

		Log := new(ActionLog)
		Log.ActionLog = &aLog
		if aLog.ID == 0 {
			log.Error("[actionLog's Listen]:actionLog can't start", aLog)
			Log.Stop(StageStopReasonRunFailed, models.ActionLogStateRunFailed)
			return
		}
		go Log.Start()
	}()

	return nil
}

// Auth is
func (actionLog *ActionLog) Auth(authMap map[string]interface{}) error {
	actionlogAuthChan <- true
	defer func() { <-actionlogAuthChan }()

	authType, ok := authMap["type"].(string)
	if !ok {
		log.Error("[actionLog's Auth]:error when get authType from given authMap:", authMap, " ===>to actionLog:", actionLog)
		return errors.New("authType is illegal")
	}

	token, ok := authMap["token"].(string)
	if !ok {
		log.Error("[actionLog's Auth]:error when get token from given authMap:", authMap, " ===>to actionLog:", actionLog)
		return errors.New("token is illegal")
	}

	err := actionLog.GetActionLog().Where("id = ?", actionLog.ID).First(actionLog).Error
	if err != nil {
		log.Error("[actionLog's Auth]:error when get actionLog info from db:", actionLog, " ===>error is:", err.Error())
		return errors.New("error when get stagelog's info from db:" + err.Error())
	}

	if actionLog.Requires == "" || actionLog.Requires == "[]" {
		log.Error("[actionLog's Auth]:error when set auth info,actionLog's requires is empty", authMap, " ===>to actionLog:", actionLog)
		return errors.New("action don't need any more auth")
	}

	requireList := make([]interface{}, 0)
	remainRequireList := make([]interface{}, 0)
	err = json.Unmarshal([]byte(actionLog.Requires), &requireList)
	if err != nil {
		log.Error("[actionLog's Auth]:error when unmarshal actionLog's require list:", actionLog, " ===>error is:", err.Error())
		return errors.New("error when get action require auth info:" + err.Error())
	}

	hasAuthed := false
	for _, require := range requireList {
		requireMap, ok := require.(map[string]interface{})
		if !ok {
			log.Error("[actionLog's Auth]:error when get actionLog's require info:", actionLog, " ===> require is:", require)
			return errors.New("error when get actionLog require auth info,require is not a json object")
		}

		requireType, ok := requireMap["type"].(string)
		if !ok {
			log.Error("[actionLog's Auth]:error when get actionLog's require type:", actionLog, " ===> require map is:", requireMap)
			return errors.New("error when get action require auth info,require don't have a type")
		}

		requireToken, ok := requireMap["token"].(string)
		if !ok {
			log.Error("[actionLog's Auth]:error when get actionLog's require token:", actionLog, " ===> require map is:", requireMap)
			return errors.New("error when get action require auth info,require don't have a token")
		}

		if requireType == authType && requireToken == token {
			hasAuthed = true
			// record auth info to actionLog's Auth info list
			actionLogAuthList := make([]interface{}, 0)
			if actionLog.AuthList != "" {
				err = json.Unmarshal([]byte(actionLog.AuthList), &actionLogAuthList)
				if err != nil {
					log.Error("[actionLog's Auth]:error when unmarshal actionLog's Auth list:", actionLog, " ===>error is:", err.Error())
					return errors.New("error when set auth info to action")
				}
			}

			actionLogAuthList = append(actionLogAuthList, authMap)

			authListInfo, err := json.Marshal(actionLogAuthList)
			if err != nil {
				log.Error("[actionLog's Auth]:error when marshal actionLog's Auth list:", actionLogAuthList, " ===>error is:", err.Error())
				return errors.New("error when save action auth info")
			}

			actionLog.AuthList = string(authListInfo)
			err = actionLog.GetActionLog().Save(actionLog).Error
			if err != nil {
				log.Error("[actionLog's Auth]:error when save actionLog's info to db:", actionLog, " ===>error is:", err.Error())
				return errors.New("error when save action auth info")
			}
		} else {
			remainRequireList = append(remainRequireList, requireMap)
		}
	}

	if !hasAuthed {
		log.Error("[actionLog's Auth]:error when auth a actionLog to start, given auth:", authMap, " is not equal to any request one:", actionLog.Requires)
		return errors.New("illegal auth info, auth failed")
	}

	remainRequireAuthInfo, err := json.Marshal(remainRequireList)
	if err != nil {
		log.Error("[actionLog's Auth]:error when marshal actionLog's remainRequireAuth list:", remainRequireList, " ===>error is:", err.Error())
		return errors.New("error when sync remain require auth info")
	}

	actionLog.Requires = string(remainRequireAuthInfo)
	err = actionLog.GetActionLog().Save(actionLog).Error
	if err != nil {
		log.Error("[actionLog's Auth]:error when save actionLog's remain require auth info:", actionLog, " ===>error is:", err.Error())
		return errors.New("error when sync remain require auth info")
	}

	return nil
}

// Start is
func (actionLog *ActionLog) Start() {
	err := actionLog.changeGlobalVar()
	if err != nil {
		log.Error("[actionLog's Start]:error when change action's global var:", err.Error())
		RecordOutcom(actionLog.Workflow, actionLog.FromWorkflow, actionLog.Stage, actionLog.FromStage, actionLog.ID, actionLog.FromAction, actionLog.Sequence, 0, false, "start action error", "error when replace action's global var:"+err.Error())
		actionLog.Stop(ActionStopReasonRunFailed, models.ActionLogStateRunFailed)
		return
	}

	if actionLog.Timeout != "" {
		go actionLog.WaitActionDone()
	}

	if actionLog.Component != 0 {
		c, err := InitComponetNew(actionLog)
		if err != nil {
			log.Error("[actionLog's Start]:error when init component:", err.Error())
			RecordOutcom(actionLog.Workflow, actionLog.FromWorkflow, actionLog.Stage, actionLog.FromStage, actionLog.ID, actionLog.FromAction, actionLog.Sequence, 0, false, "start action error", "error when init component:"+err.Error())
			actionLog.Stop(ActionStopReasonRunFailed, models.ActionLogStateRunFailed)
			return
		}

		err = c.Start()
		if err != nil {
			log.Error("[actionLog's Start]:error when start component:", err.Error())
			RecordOutcom(actionLog.Workflow, actionLog.FromWorkflow, actionLog.Stage, actionLog.FromStage, actionLog.ID, actionLog.FromAction, actionLog.Sequence, 0, false, "start action error", err.Error())
			actionLog.Stop(ActionStopReasonRunFailed, models.ActionLogStateRunFailed)
		}
	} else if actionLog.Service != 0 {
		log.Info("[actionLog's Start]:start an action that use service:", actionLog)
		RecordOutcom(actionLog.Workflow, actionLog.FromWorkflow, actionLog.Stage, actionLog.FromStage, actionLog.ID, actionLog.FromAction, actionLog.Sequence, 0, false, "start action error", "use service but not support")
		actionLog.Stop(ActionStopReasonRunSuccess, models.ActionLogStateRunSuccess)
	} else {
		log.Error("[actionLog's Start]:error when start action,action doesn't spec a type", actionLog)
		RecordOutcom(actionLog.Workflow, actionLog.FromWorkflow, actionLog.Stage, actionLog.FromStage, actionLog.ID, actionLog.FromAction, actionLog.Sequence, 0, false, "start action error", "action doesn't spec a component or a service")
		actionLog.Stop(ActionStopReasonRunFailed, models.ActionLogStateRunFailed)
	}
}

// Stop is
func (actionLog *ActionLog) Stop(reason string, runState int64) {
	err := actionLog.GetActionLog().Where("id = ?", actionLog.ID).First(actionLog).Error
	if err != nil {
		log.Error("[actionLog's Stop]:error when get actionLog's info from db:", err.Error())
		return
	}

	if actionLog.RunState == models.ActionLogStateRunFailed || actionLog.RunState == models.ActionLogStateRunSuccess {
		return
	}

	actionLog.RunState = runState
	actionLog.FailReason = reason
	err = actionLog.GetActionLog().Save(actionLog).Error
	if err != nil {
		log.Error("[actionLog's Stop]:error when change action state:", actionLog, " ===>error is:", err.Error())
		return
	}

	if actionLog.Component != 0 {
		c, err := InitComponetNew(actionLog)
		if err != nil {
			log.Error("[actionLog's Stop]:error when init component:", err.Error())
			return
		}

		err = c.Stop()
		if err != nil {
			log.Error("[actionLog's Stop]:error when stop component:", err.Error())
		}
	} else if actionLog.Service != 0 {
		log.Info("[actionLog's Stop]:stop an action that use service:", actionLog)
	} else {
		log.Error("[actionLog's Stop]:error when stop action,action doesn't spec a type", actionLog)
	}
}

// RecordEvent is
func (actionLog *ActionLog) RecordEvent(eventId int64, eventKey string, reqBody map[string]interface{}, headerInfo http.Header) error {
	c, err := InitComponetNew(actionLog)
	if err != nil {
		recordErr := RecordOutcom(actionLog.Workflow, actionLog.FromWorkflow, actionLog.Stage, actionLog.FromStage, actionLog.ID, actionLog.FromAction, actionLog.Sequence, eventId, false, "component init error:"+err.Error(), "")
		if recordErr != nil {
			log.Error("[actionLog's RecordEvent]:error when record outcome info:", recordErr.Error())
			return recordErr
		}

		log.Error("[actionLog's RecordEvent]:error when get action's platformInfo:", actionLog, " ===>error is:", err.Error())
		return err
	}

	if eventKey == models.EventTaskStatus {
		resultReqBody, ok := reqBody["INFO"].(map[string]interface{})
		if !ok {
			log.Error("[actionLog's RecordEvent]:error when get request's info body, want a json obj, got:", reqBody["INFO"])
			return errors.New("request body's info is not a json obj")
		}

		status, ok := resultReqBody["status"].(bool)
		if !ok {
			status = false
		}

		result, ok := resultReqBody["result"].(string)
		if !ok {
			result = ""
		}

		outputStr, ok := resultReqBody["output"].(string)
		if !ok {
			outputStr = ""
		}

		recordErr := RecordOutcom(actionLog.Workflow, actionLog.FromWorkflow, actionLog.Stage, actionLog.FromStage, actionLog.ID, actionLog.FromAction, actionLog.Sequence, eventId, status, result, outputStr)
		if recordErr != nil {
			log.Error("[actionLog's RecordEvent]:error when record outcome info:", recordErr.Error())
			return recordErr
		}

		stopStatus := models.ActionLogStateRunFailed
		stopReason := ActionStopReasonRunFailed
		if status {
			stopStatus = models.ActionLogStateRunSuccess
			stopReason = ActionStopReasonRunSuccess
		}
		actionLog.Stop(stopReason, int64(stopStatus))
	}

	if eventKey == models.EventComponentStop {
		c.Stop()
	}

	headerMap := make(map[string]interface{})
	for key, value := range headerInfo {
		headerMap[key] = value
	}
	headerBytes, _ := json.Marshal(headerMap)

	eventDefine := new(models.EventDefinition)
	err = eventDefine.GetEventDefinition().Where("id = ?", eventId).First(&eventDefine).Error
	if err != nil {
		log.Error("[actionLog's RecordEvent]:error when get eventDefine from db:", err.Error())
		return err
	}

	authStr := ""
	auths, ok := headerMap["Authorization"].([]string)
	if !ok {
		authStr = ""
	} else {
		authStr = strings.Join(auths, ";")
	}

	bodyBytes, _ := json.Marshal(reqBody)

	// log evnet
	err = RecordEventInfo(eventId, actionLog.Sequence, string(headerBytes), string(bodyBytes), authStr)
	if err != nil {
		log.Error("[actionLog's RecordEvent]:error when save event to db:", err.Error())
		return err
	}

	return nil
}

// SendDataToAction is
func (actionLog *ActionLog) SendDataToAction(targetUrl string) {
	manifestMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(actionLog.Manifest), &manifestMap)
	if err != nil {
		log.Error("[actionLog's SendDataToAction]:error when get action manifest info:" + err.Error())
		return
	}

	dataMap := make(map[string]interface{})
	relations, ok := manifestMap["relation"]
	if ok {
		relationInfo, ok := relations.([]interface{})
		if !ok {
			log.Error("[actionLog's SendDataToAction]:error when parse relations,want an array,got:", relations)
			return
		}

		dataMap, err = actionLog.merageFromActionsOutputData(relationInfo)
		if err != nil {
			log.Error("[actionLog's SendDataToAction]:error when get data map from action: " + err.Error())
		}
	}

	log.Info("[actionLog's SendDataToAction]:action", actionLog, " got data:", dataMap)

	var dataByte []byte

	if len(dataMap) == 0 {
		dataByte = make([]byte, 0)
	} else {
		dataByte, err = json.Marshal(dataMap)
		if err != nil {
			log.Error("[actionLog's SendDataToAction]:error when marshal dataMap:", dataMap, " ===>error is:", err.Error())
			return
		}
	}

	character := int(0)
	// send data to component or service
	resps := make([]*http.Response, 0)
	if actionLog.Component != 0 {
		character = models.CharacterComponentEvent
		resps, err = actionLog.sendDataToComponent(targetUrl, dataByte)
	} else {
		character = models.CharacterServiceEvent
		resps, err = actionLog.sendDataToService(dataByte)
	}

	resultStr := ""
	status := false
	payload := make(map[string]interface{})
	if err != nil {
		resultStr = err.Error()
		status = false
		go actionLog.Stop(ActionStopReasonSendDataFailed, models.ActionLogStateRunFailed)
	} else {
		respMap := make(map[int64]string, len(resps))
		for count, resp := range resps {
			if resp != nil {
				respBody, _ := ioutil.ReadAll(resp.Body)
				respStr := string(respBody)

				respMap[int64(count)] = respStr
			}
		}

		result, _ := json.Marshal(respMap)
		resultStr = string(result)
		status = true
	}

	payload["EVENT"] = "SEND_DATA"
	payload["EVENTID"] = "SEND_DATA"
	payload["INFO"] = map[string]interface{}{"output": string(dataByte), "result": resultStr, "status": status}
	payload["RUN_ID"] = strconv.FormatInt(actionLog.Workflow, 10) + "-" + strconv.FormatInt(actionLog.Stage, 10) + "-" + strconv.FormatInt(actionLog.ID, 10)

	payloadInfo, marshalErr := json.Marshal(payload)
	if marshalErr != nil {
		go actionLog.Stop(ActionStopReasonSendDataFailed, models.ActionLogStateRunFailed)
		log.Error("[actionLog's SendDataToAction]:error when marshal payload info:" + marshalErr.Error())
	}

	if err != nil {
		go actionLog.Stop(ActionStopReasonSendDataFailed, models.ActionLogStateRunFailed)
		log.Error("[actionLog's SendDataToAction]:error when send data to action:" + err.Error())
	}

	err = RecordEventInfo(models.EventDefineIDSendDataToAction, actionLog.Sequence, "", string(payloadInfo), "", "SEND_DATA", strconv.FormatInt(int64(character), 10), actionLog.Namespace, actionLog.Repository, strconv.FormatInt(actionLog.Workflow, 10), strconv.FormatInt(actionLog.Stage, 10), strconv.FormatInt(actionLog.ID, 10))
	if err != nil {
		go actionLog.Stop(ActionStopReasonSendDataFailed, models.ActionLogStateRunFailed)
		log.Error("[actionLog's SendDataToAction]:error when save send data info :" + err.Error())
	}
}

func (actionLog *ActionLog) merageFromActionsOutputData(relationInfo []interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	for _, relation := range relationInfo {
		relationMap, ok := relation.(map[string]interface{})
		if !ok {
			log.Println("[actionLog's merageFromActionsOutputData]:error when get relation info:want json obj,got:", relation)
			return nil, errors.New("error when parse relation info,relation is not a json!")
		}

		fromOutcome := new(models.Outcome)
		err := fromOutcome.GetOutcome().Where("real_workflow = ?", actionLog.FromWorkflow).Where("workflow = ? ", actionLog.Workflow).Where("real_action = ?", relationMap["fromAction"]).First(&fromOutcome).Order("-id").Error
		if err != nil {
			log.Error("[actionLog's merageFromActionsOutputData]:error when get request action's output from db:want get action(", actionLog.ID, ")'s output, ===>error is:", err.Error())
			return nil, errors.New("error when get from outcome, error:" + err.Error())
		}

		tempData := make(map[string]interface{})
		err = json.Unmarshal([]byte(fromOutcome.Output), &tempData)
		if err != nil {
			log.Error("[actionLog's merageFromActionsOutputData]:error when unmarshal action(", actionLog.ID, ")'s output info:", fromOutcome.Output, " ===>error is:", err.Error())
			return nil, errors.New("error when parse from action data1:" + err.Error() + "\n" + fromOutcome.Output)
		}

		relationArray := make([]interface{}, 0)

		fromRealActionIDF, ok := relationMap["fromAction"].(float64)
		if !ok {
			log.Error("[actionLog's merageFromActionsOutputData]:error when get fromRealActionID from action's relation,want a number,got:", relationMap["fromAction"])
			return nil, errors.New("error when parse from action id:relation's fromAction's id is not a number")
		}

		fromRealActionID := int64(fromRealActionIDF)

		if fromRealActionID == 0 {
			// get workflow source info, then get relation array by eventType and eventName
			workflowLog := new(models.WorkflowLog)
			err := workflowLog.GetWorkflowLog().Where("id = ?", actionLog.Workflow).First(workflowLog).Error
			if err != nil {
				log.Error("[actionLog's merageFromActionsOutputData]:error when get workflow log info from db:", err.Error())
				return nil, err
			}
			sourceInfo := workflowLog.SourceInfo
			sourceMap := make(map[string]string)
			err = json.Unmarshal([]byte(sourceInfo), &sourceMap)
			if err != nil {
				log.Error("[actionLog's merageFromActionsOutputData]:error when unmarshal workflowLog's sourceInfo:", err.Error())
				return nil, err
			}

			realRelation, ok := relationMap["relation"].(map[string]interface{})
			if !ok {
				log.Error("[actionLog's merageFromActionsOutputData]:error in workflowLog's relation define,want a obj,got: ", relationMap["relation"])
				return nil, errors.New("error in workflow's relation define")
			}

			relationArray, ok = realRelation[sourceMap["eventName"]+"_"+sourceMap["eventType"]].([]interface{})
			if !ok {
				log.Error("[actionLog's merageFromActionsOutputData]:error when get real relation from relationMap:", relationMap)
				return nil, errors.New("relation doesn't have a relation info")
			}
		} else {
			relationArray, ok = relationMap["relation"].([]interface{})
			if !ok {
				log.Error("[actionLog's merageFromActionsOutputData]:error when get relation from relationMap:", relationMap)
				return nil, errors.New("relation doesn't have a relation info")
			}
		}

		relationList := make([]Relation, 0)
		if len(relationArray) > 0 {
			for _, realationDefines := range relationArray {
				relationByte, err := json.Marshal(realationDefines)
				if err != nil {
					log.Error("[actionLog's merageFromActionsOutputData]:error went marshal relation array:", realationDefines, " ===>error is:", err.Error())
					return nil, errors.New("error when marshal relation array:" + err.Error())
				}

				var r Relation
				err = json.Unmarshal(relationByte, &r)
				if err != nil {
					log.Error("[actionLog's merageFromActionsOutputData]:error when parse relation info:", string(relationByte), " ===>error is:", err.Error())
					return nil, errors.New("error when parse relation info:" + err.Error())
				}

				relationList = append(relationList, r)
			}
		}

		actionResult := make(map[string]interface{})
		err = getResultFromRelation(fromOutcome.Output, relationList, actionResult)
		if err != nil {
			log.Error("[actionLog's merageFromActionsOutputData]:error when get result from action's relation:", err.Error())
			return nil, errors.New("error when get from data:" + err.Error())
		}

		for key, value := range actionResult {
			result[key] = value
		}
	}

	return result, nil
}

func (actionLog *ActionLog) sendDataToComponent(targetUrl string, data []byte) ([]*http.Response, error) {
	c, err := InitComponetNew(actionLog)
	if err != nil {
		log.Error("[actionLog's sendDataToComponent]:error when init component info:", err.Error())
		return nil, errors.New("error when init component info:" + err.Error())
	}

	log.Info("start send data to component...")
	return c.SendData(targetUrl, data)
}

func (actionLog *ActionLog) sendDataToService(data []byte) ([]*http.Response, error) {
	return nil, nil
}

// WaitActionDone is
func (actionLog *ActionLog) WaitActionDone() {
	timeout, err := strconv.ParseInt(actionLog.Timeout, 10, 64)
	if err != nil || timeout < 0 {
		log.Error("[actionLog's WaitActionDone]:error when parse action's timeout, want a number, got:", actionLog.Timeout)
		actionLog.Stop(ActionStopReasonRunFailed, models.ActionLogStateRunFailed)
	}

	canStop := false
	actionRunResultChan := make(chan bool, 1)
	go func() {
		for !canStop {
			actionLogInfo := new(models.ActionLog)
			err := actionLogInfo.GetActionLog().Where("id = ?", actionLog.ID).First(actionLogInfo).Error
			if err != nil {
				log.Error("[actionLog's WaitActionDone]:error when get actionLog's info from db:", err.Error())
				actionRunResultChan <- false
				return
			}

			if actionLogInfo.RunState == models.ActionLogStateRunFailed {
				actionRunResultChan <- false
				return
			} else if actionLogInfo.RunState == models.ActionLogStateRunSuccess {
				actionRunResultChan <- true
				return
			}

			time.Sleep(1 * time.Second)
		}
	}()

	duration, _ := time.ParseDuration(actionLog.Timeout + "s")
	select {
	case <-time.After(duration):
		canStop = true
		actionLog.Stop(ActionStopReasonTimeout, models.ActionLogStateRunFailed)
	case runSuccess := <-actionRunResultChan:
		if runSuccess {
			actionLog.Stop(ActionStopReasonRunSuccess, models.ActionLogStateRunSuccess)
		} else {
			actionLog.Stop(ActionStopReasonRunFailed, models.ActionLogStateRunFailed)
		}
	}
}

// GetActionPlatformInfo is
func (actionLog *ActionLog) GetActionPlatformInfo() (map[string]string, error) {
	manifestMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(actionLog.Manifest), &manifestMap)
	if err != nil {
		log.Error("[actionLog's GetActionPlatformInfo]:error when unmarshal action's manifest:", actionLog.Manifest, " ===>error is:", err.Error())
		return nil, errors.New("action " + actionLog.ActionLog.Action + "'s manifest is illegal")
	}

	platformSetting, ok := manifestMap["platform"].(map[string]interface{})
	if !ok {
		log.Error("[actionLog's GetActionPlatformInfo]:error when unmarshal action's platform info:", manifestMap, " platform setting is not a map[string]interface{}")
		return nil, errors.New("action " + actionLog.ActionLog.Action + "'s platform setting is illegal")
	}

	platformType, ok := platformSetting["platformType"].(string)
	if !ok {
		log.Error("[actionLog's GetActionPlatformInfo]:error when get action's platformType:", platformSetting, " platformType is not a string")
		return nil, errors.New("action " + actionLog.ActionLog.Action + "'s platform type is illegal")
	}

	platformHost, ok := platformSetting["platformHost"].(string)
	if !ok {
		log.Error("[actionLog's GetActionPlatformInfo]:error when get action's platformHost:", platformSetting, " platformHost is not a string")
		return nil, errors.New("action " + actionLog.ActionLog.Action + "'s platform host is illegal")
	}

	result := make(map[string]string)
	result["platformType"] = platformType
	result["platformHost"] = platformHost

	return result, nil
}

// ChangeWorkflowRuntimeVar is
func (actionLog *ActionLog) ChangeWorkflowRuntimeVar(runId string, varKey, varValue string) error {
	actionlogSetGlobalVarChan <- true
	defer func() {
		<-actionlogSetGlobalVarChan
	}()

	db := models.GetDB()
	db = db.Begin()

	varInfo := new(models.WorkflowVarLog)
	err := db.Model(&models.WorkflowVarLog{}).Where("workflow = ?", actionLog.Workflow).Where("sequence = ?", actionLog.Sequence).Where("`key` = ?", varKey).First(varInfo).Error

	if err != nil && err.Error() != "result not found" {
		log.Error("[actionLog's ChangeWorkflowRuntimeVar]:error when get workflow var info from db:", err.Error())
		rollbackErr := db.Rollback().Error
		if rollbackErr != nil {
			log.Error("[actionLog's ChangeWorkflowRuntimeVar]:when rollback in get workflow var info:", rollbackErr.Error())
			return errors.New("errors occur:\nerror1:" + err.Error() + "\nerror2:" + rollbackErr.Error())
		}
		return err
	}

	changeLogMap := make(map[string]interface{})
	changeLogMap["user"] = runId
	changeLogMap["time"] = time.Now().Format("2006-01-02 15:04:05")
	changeLogMap["action"] = "set key:" + varKey + " 's value to:" + varValue

	changeLogList := make([]interface{}, 0)
	if varInfo.ChangeLog != "" {
		err := json.Unmarshal([]byte(varInfo.ChangeLog), &changeLogList)
		if err != nil {
			log.Error("[actionLog's ChangeWorkflowRuntimeVar]:error when unmarshal var's changelog to a list,got:", varInfo.ChangeLog)
			return errors.New("error when save var info")
		}
	}
	changeLogList = append(changeLogList, changeLogMap)

	changeInfoBytes, _ := json.Marshal(changeLogList)

	varInfo.ChangeLog = string(changeInfoBytes)
	varInfo.Vaule = varValue

	err = db.Model(&models.WorkflowVarLog{}).Save(varInfo).Error
	if err != nil {
		log.Error("[actionLog's ChangeWorkflowRuntimeVar]:error when save workflow var info to db:", err.Error())
		rollbackErr := db.Rollback().Error
		if rollbackErr != nil {
			log.Error("[actionLog's ChangeWorkflowRuntimeVar]:when rollback in save workflow var info:", rollbackErr.Error())
			return errors.New("errors occur:\nerror1:" + err.Error() + "\nerror2:" + rollbackErr.Error())
		}
		return err
	}

	db.Commit()
	return nil
}

func (actionLog *ActionLog) changeGlobalVar() error {
	// change action's name
	if strings.HasPrefix(actionLog.Action, "@") && strings.HasSuffix(actionLog.Action, "@") {
		varKey := actionLog.Action[1 : len(actionLog.Action)-1]

		varValue, err := getWorkflowVarLogInfo(actionLog.Workflow, actionLog.Sequence, varKey)
		if err != nil {
			log.Error("[actionLog's changeGlobalVar]:action:", actionLog.Action, " use name both start and end with '@',but not a global value")
		} else {
			actionLog.Action = varValue
		}
	}

	// change action's imageName
	if strings.HasPrefix(actionLog.ImageName, "@") && strings.HasSuffix(actionLog.ImageName, "@") {
		varKey := actionLog.ImageName[1 : len(actionLog.ImageName)-1]

		varValue, err := getWorkflowVarLogInfo(actionLog.Workflow, actionLog.Sequence, varKey)
		if err != nil {
			log.Error("[actionLog's changeGlobalVar]:action:", actionLog.Action, " use ImageName both start and end with '@',but not a global value")
			return errors.New("error action's image name is illegal")
		}

		actionLog.ImageName = varValue
	}

	// change action's imageTag
	if strings.HasPrefix(actionLog.ImageTag, "@") && strings.HasSuffix(actionLog.ImageTag, "@") {
		varKey := actionLog.ImageTag[1 : len(actionLog.ImageTag)-1]

		varValue, err := getWorkflowVarLogInfo(actionLog.Workflow, actionLog.Sequence, varKey)
		if err != nil {
			log.Error("[actionLog's changeGlobalVar]:action:", actionLog.Action, " use a name both start and end with '@',but not a global value")
			return errors.New("error action's image tag is illegal")
		}

		actionLog.ImageTag = varValue
	}

	// change action's timeout
	if strings.HasPrefix(actionLog.Timeout, "@") && strings.HasSuffix(actionLog.Timeout, "@") {
		varKey := actionLog.Timeout[1 : len(actionLog.Timeout)-1]

		varValue, err := getWorkflowVarLogInfo(actionLog.Workflow, actionLog.Sequence, varKey)
		if err != nil {
			log.Error("[actionLog's changeGlobalVar]:action:", actionLog.Action, " got an error when get:", varKey, " from db:", err.Error())
			return errors.New("error when get workflow var info")
		}

		timeoutInt, err := strconv.ParseInt(varValue, 10, 64)
		if err != nil {
			log.Error("[actionLog's changeGlobalVar]:action:", actionLog.Action, " set time as:", actionLog.Timeout, " but when parse var's value(", varValue, ") to int got a error:", err.Error())
			return errors.New("use a NaN value to action's timeout")
		}

		actionLog.Timeout = strconv.FormatInt(timeoutInt, 10)
	}

	// change action's title
	if strings.HasPrefix(actionLog.Title, "@") && strings.HasSuffix(actionLog.Title, "@") {
		varKey := actionLog.Title[1 : len(actionLog.Title)-1]

		varValue, err := getWorkflowVarLogInfo(actionLog.Workflow, actionLog.Sequence, varKey)
		if err != nil {
			log.Error("[actionLog's changeGlobalVar]:action:", actionLog.Action, " use a title both start and end with '@',but not a global value")
		} else {
			actionLog.Title = varValue
		}
	}

	// change action's description
	if strings.HasPrefix(actionLog.Description, "@") && strings.HasSuffix(actionLog.Description, "@") {
		varKey := actionLog.Description[1 : len(actionLog.Description)-1]

		varValue, err := getWorkflowVarLogInfo(actionLog.Workflow, actionLog.Sequence, varKey)
		if err != nil {
			log.Error("[actionLog's changeGlobalVar]:action:", actionLog.Action, " use a title both start and end with '@',but not a global value")
		} else {
			actionLog.Description = varValue
		}
	}

	// change action's manifest
	// only replace platform info
	manifestMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(actionLog.Manifest), &manifestMap)
	if err != nil {
		log.Error("[actionLog's changeGlobalVar]:action:", actionLog.Action, "error when get platform Map:", err.Error())
		return errors.New("action's platform info error")
	}

	if platformMap, ok := manifestMap["platform"].(map[string]interface{}); ok {
		host, ok := platformMap["platformHost"].(string)
		if ok && strings.HasPrefix(host, "@") && strings.HasSuffix(host, "@") {
			varKey := host[1 : len(host)-1]

			varValue, err := getWorkflowVarLogInfo(actionLog.Workflow, actionLog.Sequence, varKey)
			if err != nil {
				log.Error("[actionLog's changeGlobalVar]:action:", host, " use a name both start and end with '@',but not a global value")
				return errors.New("action's platform setting is illegal")
			}

			platformMap["platformHost"] = varValue
			manifestMap["platform"] = platformMap

			manifestBytes, _ := json.Marshal(manifestMap)

			actionLog.Manifest = string(manifestBytes)
		}
	}

	// change action's env
	envMap := make(map[string]string)
	err = json.Unmarshal([]byte(actionLog.Environment), &envMap)
	if err != nil {
		log.Error("[actionLog's changeGlobalVar]:action:", actionLog.Action, " error when unmarshal action's env info,want a json obj, got:", actionLog.Environment)
		return errors.New("error in action's env define")
	}

	afterChangeEnvMap := make(map[string]string)
	for key, value := range envMap {
		if strings.HasPrefix(key, "@") && strings.HasSuffix(key, "@") {
			varKey := key[1 : len(key)-1]

			varValue, err := getWorkflowVarLogInfo(actionLog.Workflow, actionLog.Sequence, varKey)
			if err != nil {
				log.Error("[actionLog's changeGlobalVar]:action:", actionLog.Action, " use env key both start and end with '@',but not a global value")
			} else {
				key = varValue
			}
		}

		if strings.HasPrefix(value, "@") && strings.HasSuffix(value, "@") {
			varKey := value[1 : len(value)-1]

			varValue, err := getWorkflowVarLogInfo(actionLog.Workflow, actionLog.Sequence, varKey)
			if err != nil {
				log.Error("[actionLog's changeGlobalVar]:action:", actionLog.Action, " use env value both start and end with '@',but not a global value")
			} else {
				value = varValue
			}
		}

		afterChangeEnvMap[key] = value
	}

	envBytes, _ := json.Marshal(afterChangeEnvMap)

	actionLog.Environment = string(envBytes)

	// change action's kube info
	kubeSettingMap := make(map[string]interface{})
	json.Unmarshal([]byte(actionLog.Kubernetes), &kubeSettingMap)

	nodeIPStr, ok := kubeSettingMap["nodeIP"].(string)
	if !ok {
		log.Error("[actionLog's changeGlobalVar]:error when get kube's nodeip, want a string, got:", kubeSettingMap["nodeIP"])
		return errors.New("error when get kube's nodeip")
	}

	if strings.HasPrefix(nodeIPStr, "@") && strings.HasSuffix(nodeIPStr, "@") {
		varKey := nodeIPStr[1 : len(nodeIPStr)-1]

		varValue, err := getWorkflowVarLogInfo(actionLog.Workflow, actionLog.Sequence, varKey)
		if err != nil {
			log.Error("[actionLog's changeGlobalVar]:action:", actionLog.Action, " use nodeIP both start and end with '@',but not a global value")
			return errors.New("nodeIP is not a addressable address")
		}

		nodeIPStr = varValue

		kubeSettingMap["nodeIP"] = nodeIPStr
	}

	if useAdvanced, ok := kubeSettingMap["useAdvanced"].(bool); !ok || !useAdvanced {
		podConfigMap, ok := kubeSettingMap["podConfig"].(map[string]interface{})
		if ok {
			podConfig, err := changeMapInfoWithGlobalVar(actionLog.Workflow, actionLog.Sequence, podConfigMap)
			if err != nil {
				log.Error("[actionLog's changeGlobalVar]:action:", actionLog.Action, " error when change pod config:", err.Error())
				return errors.New("action's pod config is illegal")
			}

			kubeSettingMap["podConfig"] = podConfig
		}

		serviceConfigMap, ok := kubeSettingMap["serviceConfig"].(map[string]interface{})
		if ok {
			serviceConfig, err := changeMapInfoWithGlobalVar(actionLog.Workflow, actionLog.Sequence, serviceConfigMap)
			if err != nil {
				log.Error("[actionLog's changeGlobalVar]:action:", actionLog.Action, " error when change service config:", err.Error())
				return errors.New("action's service config is illegal")
			}

			kubeSettingMap["serviceConfig"] = serviceConfig
		}

		kubeSettingBytes, _ := json.Marshal(kubeSettingMap)
		actionLog.Kubernetes = string(kubeSettingBytes)
	}

	err = new(models.ActionLog).GetActionLog().Save(actionLog.ActionLog).Error
	if err != nil {
		log.Error("[actionLog's changeGlobalVar]:error when save action's change to db:", err.Error())
		return errors.New("error when save change to db")
	}

	return nil
}

// LinkStartWorkflow is
func (actionLog *ActionLog) LinkStartWorkflow(runId, token, workflowName, workflowVersion, eventName, eventType string, startJson map[string]interface{}) error {
	expectToken := utils.MD5(actionLog.Action + runId)
	if expectToken != token {
		log.Info("[actionLog's LinkStartWorkflow]:action(", actionLog.ID, ") runid is:(", runId, ") error when check token: want:", expectToken, " got:", token)
		return errors.New("token is illegal")
	}

	workflowInfo, err := GetLatestRunableWorkflow(actionLog.Namespace, actionLog.Repository, workflowName, workflowVersion)
	if err != nil {
		log.Error("[actionLog's LinkStartWorkflow]:error when get workflow info:", err.Error())
		return errors.New("get workflow error")
	}

	startDataBytes, _ := json.Marshal(startJson)

	authMap := make(map[string]interface{})
	authMap["type"] = AuthTypeWorkflowDefault
	authMap["token"] = AuthTokenDefault
	authMap["runID"] = runId
	authMap["eventName"] = eventName
	authMap["eventType"] = eventType
	authMap["time"] = time.Now().Format("2006-01-02 15:04:05")
	pass, err := checkInstanceNum(workflowInfo.ID)
	if !pass {
		log.Error("[actionLog's LinkStartWorkflow]:error when checkworkflow instance num:", err.Error())
		return err
	}

	workflowLog, err := Run(workflowInfo.ID, authMap, string(startDataBytes))
	if err != nil {
		log.Error("[actionLog's LinkStartWorkflow]:error when run workflow:", err.Error())
		return errors.New("workflow run error")
	}

	preInfoMap := make(map[string]interface{})
	preInfoMap["workflowID"] = actionLog.Workflow
	preInfoMap["stageID"] = actionLog.Stage
	preInfoMap["actionID"] = actionLog.ID
	preInfoMap["token"] = token

	preInfoBytes, _ := json.Marshal(preInfoMap)

	workflowLog.PreWorkflow = actionLog.Workflow
	workflowLog.PreStage = actionLog.Stage
	workflowLog.PreAction = actionLog.ID
	workflowLog.PreWorkflowInfo = string(preInfoBytes)

	err = workflowLog.GetWorkflowLog().Save(workflowLog).Error
	if err != nil {
		log.Error("[actionLog's LinkStartWorkflow]:error when update workflow info to db:", err.Error())
		workflowLog.Stop(WorkflowStopReasonRunFailed, models.WorkflowLogStateRunFailed)
		return errors.New("error when update workflow info")
	}

	return nil
}

func getResultFromRelation(outputJson string, relationList []Relation, result map[string]interface{}) error {
	fromActionData := make(map[string]interface{})

	err := json.Unmarshal([]byte(outputJson), &fromActionData)
	if err != nil {
		log.Error("[actionLog's getResultFromRelation]:error when unmarshal action's output json:", outputJson, " ===>error is:", err.Error())
		return errors.New("error when parse from action data2:" + err.Error() + "\n" + outputJson)
	}

	for _, relation := range relationList {
		fromData, err := getJsonDataByPath(strings.TrimPrefix(relation.From, "."), fromActionData)
		if err != nil {
			return errors.New("error when get fromData :" + err.Error())
		}

		setDataToMapByPath(fromData, result, strings.TrimPrefix(relation.To, "."))
	}

	return nil
}

func changeMapInfoWithGlobalVar(workflow, sequence int64, sourceMap map[string]interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	for key, value := range sourceMap {
		if strings.HasPrefix(key, "@") && strings.HasSuffix(key, "@") {
			varKey := key[1 : len(key)-1]

			varValue, err := getWorkflowVarLogInfo(workflow, sequence, varKey)
			if err != nil {
				log.Error("[actionLog's changeMapInfoWithGlobalVar]:map use key both start and end with '@',but not a global value")
			} else {
				key = varValue
			}
		}

		switch value.(type) {
		case string:
			valueStr := value.(string)

			var afterReplace string
			if strings.HasPrefix(valueStr, "@") && strings.HasSuffix(valueStr, "@") {
				varKey := valueStr[1 : len(valueStr)-1]

				varValue, err := getWorkflowVarLogInfo(workflow, sequence, varKey)
				if err != nil {
					log.Error("[actionLog's changeMapInfoWithGlobalVar]:map use value value both start and end with '@',but not a global value")
				} else {
					afterReplace = varValue
				}
			} else {
				afterReplace = value.(string)
			}

			if key == "memory" {
				afterReplace += "Mi"
			}

			valueInt, err := strconv.ParseInt(afterReplace, 10, 64)
			if err != nil {
				value = afterReplace
			} else {
				value = valueInt
			}

		case map[string]interface{}:
			childValue, err := changeMapInfoWithGlobalVar(workflow, sequence, value.(map[string]interface{}))
			if err != nil {
				return nil, err
			}

			value = childValue
		case []interface{}:
			resultArray := make([]interface{}, 0)
			for _, info := range value.([]interface{}) {
				switch info.(type) {
				case string:
					valueStr := value.(string)

					var afterReplace string
					if strings.HasPrefix(valueStr, "@") && strings.HasSuffix(valueStr, "@") {
						varKey := valueStr[1 : len(valueStr)-1]

						varValue, err := getWorkflowVarLogInfo(workflow, sequence, varKey)
						if err != nil {
							log.Error("[actionLog's changeMapInfoWithGlobalVar]:map use value value both start and end with '@',but not a global value")
						} else {
							valueStr = varValue
							afterReplace = varValue
						}
					} else {
						afterReplace = value.(string)
					}

					if key == "memory" {
						valueStr += "Mi"
						value = valueStr
					}

					if key == "nodePort" || key == "port" || key == "targetPort" {
						valueInt, err := strconv.ParseInt(afterReplace, 10, 64)
						if err != nil {
							value = afterReplace
						} else {
							value = valueInt
						}
					}

					resultArray = append(resultArray, value)
				case map[string]interface{}:
					tempResult, err := changeMapInfoWithGlobalVar(workflow, sequence, info.(map[string]interface{}))
					if err != nil {
						return nil, err
					}

					resultArray = append(resultArray, tempResult)
				}
			}

			value = resultArray
		default:
			log.Error("[actionLog's changeMapInfoWithGlobalVar]:map's value is illegal, want a string, got:", value)
			return nil, errors.New("got unknow type")
		}

		result[key] = value
	}

	return result, nil
}
