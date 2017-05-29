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

package module

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Huawei/containerops/pilotage/models"
	"github.com/Huawei/containerops/pilotage/module/checker"

	log "github.com/Sirupsen/logrus"
	"github.com/containerops/configure"
)

const (
	// WorkflowStopReasonInstanceFull is
	WorkflowStopReasonInstanceFull = "NO_ROOM_FOR_RUN_INSTANCE"
	// WorkflowStopReasonTimeout is
	WorkflowStopReasonTimeout = "TIME_OUT"

	// WorkflowStopReasonRunSuccess is
	WorkflowStopReasonRunSuccess = "WORKFLOW_RUN_SUCCESS"
	// WorkflowStopReasonRunFailed is
	WorkflowStopReasonRunFailed = "WORKFLOW_RUN_FAILED"
)

var (
	startWorkflowChan  chan bool
	createWorkflowChan chan bool

	workflowlogAuthChan             chan bool
	workflowlogListenChan           chan bool
	workflowlogSequenceGenerateChan chan bool
)

func init() {
	startWorkflowChan = make(chan bool, 1)
	createWorkflowChan = make(chan bool, 1)
	workflowlogAuthChan = make(chan bool, 1)
	workflowlogListenChan = make(chan bool, 1)
	workflowlogSequenceGenerateChan = make(chan bool, 1)
}

// Workflow is
type Workflow struct {
	*models.Workflow
}

// WorkflowLog is
type WorkflowLog struct {
	*models.WorkflowLog
}

// CreateNewWorkflow is create a new workflow with given data
func CreateNewWorkflow(namespace, repository, workflowName, workflowVersion string) (string, error) {
	createWorkflowChan <- true
	defer func() {
		<-createWorkflowChan
	}()

	var count int64
	err := new(models.Workflow).GetWorkflow().Where("namespace = ?", namespace).Where("workflow = ?", workflowName).Order("-id").Count(&count).Error
	if err != nil {
		return "", errors.New("error when query workflow data in database:" + err.Error())
	}

	if count > 0 {
		return "", errors.New("workflow name is exist!")
	}

	workflowInfo := new(models.Workflow)
	workflowInfo.Namespace = strings.TrimSpace(namespace)
	workflowInfo.Repository = strings.TrimSpace(repository)
	workflowInfo.Workflow = strings.TrimSpace(workflowName)
	workflowInfo.Version = strings.TrimSpace(workflowVersion)
	workflowInfo.VersionCode = 1

	err = workflowInfo.GetWorkflow().Save(workflowInfo).Error
	if err != nil {
		return "", errors.New("error when save workflow info:" + err.Error())
	}

	return "create new workflow success", nil
}

// GetWorkflowListByNamespaceAndRepository is
func GetWorkflowListByNamespaceAndRepository(namespace, repository string) ([]map[string]interface{}, error) {
	resultMap := make([]map[string]interface{}, 0)
	workflowList := make([]models.Workflow, 0)
	workflowsMap := make(map[string]interface{}, 0)
	err := new(models.Workflow).GetWorkflow().Where("namespace = ?", namespace).Where("repository = ?", repository).Order("-updated_at").Find(&workflowList).Error
	if err != nil && !strings.Contains(err.Error(), "record not found") {
		log.Error("[workflow's GetWorkflowListByNamespaceAndRepository]error when get workflow list from db:" + err.Error())
		return nil, errors.New("error when get workflow list by namespace and repository from db:" + err.Error())
	}

	for _, workflowInfo := range workflowList {
		if _, ok := workflowsMap[workflowInfo.Workflow]; !ok {
			tempMap := make(map[string]interface{})
			tempMap["version"] = make(map[int64]interface{})
			workflowsMap[workflowInfo.Workflow] = tempMap
		}

		workflowMap := workflowsMap[workflowInfo.Workflow].(map[string]interface{})
		versionMap := workflowMap["version"].(map[int64]interface{})

		versionMap[workflowInfo.VersionCode] = workflowInfo
		workflowMap["id"] = workflowInfo.ID
		workflowMap["name"] = workflowInfo.Workflow
		workflowMap["version"] = versionMap
	}

	for _, workflow := range workflowList {

		workflowInfo := workflowsMap[workflow.Workflow].(map[string]interface{})
		if isSign, ok := workflowInfo["isSign"].(bool); ok && isSign {
			continue
		}

		workflowInfo["isSign"] = true
		workflowsMap[workflow.Workflow] = workflowInfo

		versionList := make([]map[string]interface{}, 0)
		for _, workflowVersion := range workflowList {
			if workflowVersion.Workflow == workflowInfo["name"].(string) {
				versionMap := make(map[string]interface{})
				versionMap["id"] = workflowVersion.ID
				versionMap["version"] = workflowVersion.Version
				versionMap["versionCode"] = workflowVersion.VersionCode

				latestWorkflowLog := new(models.WorkflowLog)
				err := latestWorkflowLog.GetWorkflowLog().Where("from_workflow = ?", workflowVersion.ID).Order("-id").First(latestWorkflowLog).Error
				if err != nil && err.Error() != "record not found" {
					log.Error("[workflow's GetWorkflowListByNamespaceAndRepository]:error when get workflow's latest run info:", err.Error())
				}

				if latestWorkflowLog.ID != 0 {
					statusMap := make(map[string]interface{})

					status := false
					if latestWorkflowLog.RunState != models.WorkflowLogStateRunFailed {
						status = true
					}

					statusMap["time"] = latestWorkflowLog.CreatedAt.Format("2006-01-02 15:04:05")
					statusMap["status"] = status

					versionMap["status"] = statusMap
				}

				versionList = append(versionList, versionMap)
			}
		}

		tempResult := make(map[string]interface{})
		tempResult["id"] = workflowInfo["id"]
		tempResult["name"] = workflowInfo["name"]
		tempResult["version"] = versionList

		resultMap = append(resultMap, tempResult)
	}

	return resultMap, nil
}

// GetWorkflowInfo is
func GetWorkflowInfo(namespace, repository, workflowName string, workflowId int64) (map[string]interface{}, error) {
	resultMap := make(map[string]interface{})
	workflowInfo := new(models.Workflow)
	err := workflowInfo.GetWorkflow().Where("id = ?", workflowId).First(&workflowInfo).Error
	if err != nil {
		return nil, errors.New("error when get workflow info from db:" + err.Error())
	}

	if workflowInfo.Namespace != namespace || workflowInfo.Repository != repository || workflowInfo.Workflow != workflowName {
		return nil, errors.New("workflow is not equal to target workflow")
	}

	// get workflow define json first, if has a define json,return it
	if workflowInfo.Manifest != "" {
		defineMap := make(map[string]interface{})
		json.Unmarshal([]byte(workflowInfo.Manifest), &defineMap)

		if defineInfo, ok := defineMap["define"]; ok {
			if defineInfoMap, ok := defineInfo.(map[string]interface{}); ok {
				defineInfoMap["status"] = workflowInfo.State == models.WorkflowStateAble
				return defineInfoMap, nil
			}
		}
	}

	// get all stage info of current workflow
	// if a workflow done have a define of itself
	// then the workflow is a new workflow ,so only get it's stage list is ok
	stageList, err := getDefaultStageListByWorkflow(*workflowInfo)
	if err != nil {
		return nil, err
	}
	resultMap["stageList"] = stageList
	// resultMap["stageList"] = make([]map[string]interface{}, 0)

	resultMap["lineList"] = make([]map[string]interface{}, 0)

	resultMap["setting"] = map[string]interface{}{
		"data": map[string]interface{}{
			"runningInstances": map[string]interface{}{
				"available": false,
				"number":    10},
			"timedTasks": map[string]interface{}{
				"available": false,
				"tasks":     make([]interface{}, 0),
			}}}

	resultMap["status"] = false

	return resultMap, nil
}

func getDefaultStageListByWorkflow(workflowInfo models.Workflow) ([]map[string]interface{}, error) {
	stageListMap := make([]map[string]interface{}, 0)

	startStage := make(map[string]interface{})
	startStage["id"] = "start-stage"
	startStage["type"] = "workflow-start"
	startStage["setupData"] = make(map[string]interface{})
	stageListMap = append(stageListMap, startStage)

	addStage := make(map[string]interface{})
	addStage["id"] = "add-stage"
	addStage["type"] = "workflow-add-stage"
	stageListMap = append(stageListMap, addStage)

	endStage := make(map[string]interface{})
	endStage["id"] = "end-stage"
	endStage["type"] = "workflow-end"
	endStage["setupData"] = make(map[string]interface{})
	stageListMap = append(stageListMap, endStage)

	return stageListMap, nil
}

// GetStageHistoryInfo is
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

// Run is
func Run(workflowId int64, authMap map[string]interface{}, startData string) (*WorkflowLog, error) {
	workflowInfo := new(models.Workflow)
	err := workflowInfo.GetWorkflow().Where("id = ?", workflowId).First(workflowInfo).Error
	if err != nil {
		log.Error("[workflow's Run]:error when get workflow's info from db:", err.Error())
		return nil, errors.New("error when get target workflow info:" + err.Error())
	}

	if workflowInfo.State == models.WorkflowStateDisable {
		return nil, errors.New("workflow is not runnable")
	}

	workflow := new(Workflow)
	workflow.Workflow = workflowInfo

	eventName, ok := authMap["eventName"].(string)
	if !ok {
		log.Error("[workflow's Run]:error when parse eventName,want a string, got:", authMap["eventName"])
		return nil, errors.New("error when get eventName")
	}

	eventType, ok := authMap["eventType"].(string)
	if !ok {
		log.Error("[workflow's Run]:error when parse eventName,want a string, got:", authMap["eventType"])
		return nil, errors.New("error when get eventType")
	}

	if eventType == "github" {
		for event, realName := range allEventMap["github"] {
			if realName == eventName {
				eventName = event
			}
		}
	}

	eventMap := make(map[string]string)
	eventMap["eventName"] = eventName
	eventMap["eventType"] = eventType

	// first generate a workflow log to record all current workflow's info which will be used in feature
	workflowLog, err := workflow.GenerateNewLog(eventMap)
	if err != nil {
		return nil, err
	}

	// let workflow log listen all auth, if all auth is ok, start run this workflow log
	err = workflowLog.Listen(startData)
	if err != nil {
		return nil, err
	}

	// auth this workflow log by given auth info
	err = workflowLog.Auth(authMap)
	if err != nil {
		return nil, err
	}

	return workflowLog, nil
}

// GetWorkflowList is
func GetWorkflowList(namespace, repository string, page, prePageCount int64, filter, filtertype string) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	if filtertype == "fuzzy" {
		filter = "%" + filter + "%"
	}

	workflowList := make([]models.WorkflowLog, 0)
	err := models.GetDB().Raw("SELECT * FROM (SELECT * FROM workflow_log WHERE `workflow_log`.deleted_at IS NULL AND workflow like ? AND workflow_log.namespace = ? AND workflow_log.repository = ? ORDER BY id DESC) i GROUP BY i.workflow LIMIT ? OFFSET ?", filter, namespace, repository, prePageCount, (page-1)*prePageCount).Scan(&workflowList).Error
	if err != nil && err.Error() != "record not found" {
		log.Error("[workflow's GetWorkflowList]:error when get workflow list from db:", err.Error())
		return nil, errors.New("error when get workflow list")
	}

	count := new(struct {
		Count int64
	})
	err = models.GetDB().Raw("select count(distinct(workflow)) as count from workflow_log where deleted_at IS NULL and namespace = ? and repository = ? and workflow like ?", namespace, repository, filter).Scan(count).Error
	if err != nil && err.Error() != "sql: no rows in result set" {
		log.Error("[workflow's GetWorkflowList]:error when count workfow num from db:", err.Error())
		return nil, errors.New("error when get workflow list")
	}

	workflows := make([]map[string]interface{}, 0)
	for _, workflowInfo := range workflowList {
		tempMap := make(map[string]interface{})
		tempMap["workflowName"] = workflowInfo.Workflow
		tempMap["workflowId"] = workflowInfo.ID

		workflows = append(workflows, tempMap)
	}

	result["totalWorkflows"] = count.Count
	result["workflows"] = workflows

	return result, nil
}

// GetWorkflowVersionList is
func GetWorkflowVersionList(namespace, repository, workflow string, workflowID int64) ([]map[string]interface{}, error) {
	result := make([]map[string]interface{}, 0)

	versionList := make([]models.WorkflowLog, 0)
	err := new(models.WorkflowLog).GetWorkflowLog().Where("namespace = ?", namespace).Where("repository = ?", repository).Where("workflow = ?", workflow).Group("version").Find(&versionList).Error
	if err != nil && err.Error() != "record not found" {
		log.Error("[workflow's GetWorkflowVersionList]:error when get workflow's version list from db:", err.Error())
		return nil, errors.New("error when get workflow version list")
	}

	for _, versionInfo := range versionList {
		tempMap := make(map[string]interface{})
		tempMap["versionName"] = versionInfo.Version
		tempMap["versionId"] = versionInfo.ID

		result = append(result, tempMap)
	}

	return result, nil
}

// GetWorkflowSequenceList is
func GetWorkflowSequenceList(namespace, repository, workflow, version string, versionID, sum int64) ([]map[string]interface{}, error) {
	result := make([]map[string]interface{}, 0)

	if sum < 1 || sum < 100 {
		sum = 10
	}

	workflows := make([]models.WorkflowLog, 0)
	err := new(models.WorkflowLog).GetWorkflowLog().Where("namespace = ?", namespace).Where("repository = ?", repository).Where("workflow = ?", workflow).Where("version = ?", version).Where("run_state > ?", 1).Order("-id").Limit(int(sum)).Find(&workflows).Error
	if err != nil {
		log.Error("[workflow's GetWorkflowSequenceList]:error when get workflow run log from db:", err.Error())
		return nil, errors.New("error when get sequence info")
	}

	for _, workflowInfo := range workflows {
		sequenceMap := make(map[string]interface{})
		sequenceMap["sequenceId"] = workflowInfo.ID
		sequenceMap["sequence"] = workflowInfo.Sequence
		sequenceMap["runTime"] = strconv.FormatFloat(workflowInfo.UpdatedAt.Sub(workflowInfo.CreatedAt).Seconds(), 'f', 0, 64)
		sequenceMap["runResult"] = workflowInfo.RunState
		sequenceMap["date"] = workflowInfo.CreatedAt.Format("2006-01-02")
		sequenceMap["time"] = workflowInfo.CreatedAt.Format("15:04")
		sequenceMap["error"] = workflowInfo.FailReason
		if workflowInfo.PreWorkflow != 0 {
			preWorkflow := new(models.WorkflowLog)
			err := preWorkflow.GetWorkflowLog().Where("id = ?", workflowInfo.PreWorkflow).First(&preWorkflow).Error
			if err != nil {
				log.Error("[workflow's GetWorkflowSequenceList]:error when get preworkflow info from db:", err.Error())
			} else {
				sequenceMap["startWorkflowName"] = preWorkflow.Workflow
			}
		}
		stageList, err := getSequenceStageInfo(namespace, repository, workflowInfo.ID, workflowInfo.Sequence)
		if err != nil {
			log.Error("[workflow's GetWorkflowSequenceList]:error when get sequence's stage list:", err.Error())
			return nil, errors.New("error when get sequence info")
		}

		sequenceMap["stages"] = stageList

		result = append(result, sequenceMap)
	}

	return result, nil
}

// GetActionLinkStartInfo is
func GetActionLinkStartInfo(namespace, repository, workflow, version, action string, sequence, workflowID, actionID int64) ([]map[string]interface{}, error) {
	result := make([]map[string]interface{}, 0)

	workflows := make([]models.WorkflowLog, 0)
	err := new(models.WorkflowLog).GetWorkflowLog().Where("namespace = ?", namespace).Where("repository = ?", repository).Where("pre_workflow = ?", workflowID).Where("pre_action = ?", actionID).Where("run_state > ?", 1).Find(&workflows).Error
	if err != nil {
		log.Error("[workflow's GetActionLinkStartInfo]:error when get workflow run log from db:", err.Error())
		return nil, errors.New("error when get sequence info")
	}

	for _, workflowInfo := range workflows {
		sequenceMap := make(map[string]interface{})
		sequenceMap["workflowName"] = workflowInfo.Workflow
		sequenceMap["workflowId"] = workflowInfo.ID
		sequenceMap["versionName"] = workflowInfo.Version
		sequenceMap["versionId"] = workflowInfo.ID
		sequenceMap["sequenceId"] = workflowInfo.ID
		sequenceMap["sequence"] = workflowInfo.Sequence
		sequenceMap["runTime"] = strconv.FormatFloat(workflowInfo.UpdatedAt.Sub(workflowInfo.CreatedAt).Seconds(), 'f', 0, 64)
		sequenceMap["runResult"] = workflowInfo.RunState
		sequenceMap["date"] = workflowInfo.CreatedAt.Format("2006-01-02")
		sequenceMap["time"] = workflowInfo.CreatedAt.Format("15:04")
		stageList, err := getSequenceStageInfo(namespace, repository, workflowInfo.ID, workflowInfo.Sequence)
		if err != nil {
			log.Error("[workflow's GetActionLinkStartInfo]:error when get sequence's stage list:", err.Error())
			return nil, errors.New("error when get sequence info")
		}

		sequenceMap["stages"] = stageList

		result = append(result, sequenceMap)
	}

	return result, nil
}

func getSequenceStageInfo(namespace, repository string, workflow, sequence int64) ([]map[string]interface{}, error) {
	result := make([]map[string]interface{}, 0)

	stages := make([]models.StageLog, 0)
	err := new(models.StageLog).GetStageLog().Where("namespace = ?", namespace).Where("repository = ?", repository).Where("workflow = ?", workflow).Where("sequence = ?", sequence).Find(&stages).Error
	if err != nil {
		log.Error("[workflow's getSequenceStageInfo]:error when get stage run log from db:", err.Error())
		return nil, errors.New("error when get sequence info")
	}

	for _, stageInfo := range stages {
		stageMap := make(map[string]interface{})
		stageMap["stageName"] = stageInfo.Stage
		stageMap["stageId"] = stageInfo.ID
		stageMap["runResult"] = stageInfo.RunState
		stageMap["isTimeout"] = stageInfo.FailReason == StageStopReasonTimeout
		stageMap["timeout"] = stageInfo.Timeout
		stageMap["runTime"] = strconv.FormatFloat(stageInfo.UpdatedAt.Sub(stageInfo.CreatedAt).Seconds(), 'f', 0, 64)
		stageMap["error"] = stageInfo.FailReason
		actionList, err := getSequenceActionInfo(namespace, repository, stageInfo.Workflow, stageInfo.Sequence, stageInfo.ID)
		if err != nil {
			log.Error("[workflow's getSequenceStageInfo]:error when get sequence's action list:", err.Error())
			return nil, errors.New("error when get sequence info")
		}

		stageMap["actions"] = actionList
		result = append(result, stageMap)
	}

	return result, nil
}

func getSequenceActionInfo(namespace, repository string, workflow, sequence, stageID int64) ([]map[string]interface{}, error) {
	result := make([]map[string]interface{}, 0)

	actions := make([]models.ActionLog, 0)
	err := new(models.ActionLog).GetActionLog().Where("namespace = ?", namespace).Where("repository = ?", repository).Where("workflow = ?", workflow).Where("sequence = ?", sequence).Where("stage = ?", stageID).Find(&actions).Error
	if err != nil {
		log.Error("[workflow's getSequenceActionInfo]:error when get action run log from db:", err.Error())
		return nil, errors.New("error when get sequence info")
	}

	for _, actionInfo := range actions {
		linkStartWorkflows := make([]models.WorkflowLog, 0)
		err := new(models.WorkflowLog).GetWorkflowLog().Where("namespace = ?", namespace).Where("repository = ?", repository).Where("pre_workflow = ?", actionInfo.Workflow).Where("pre_action = ?", actionInfo.ID).Find(&linkStartWorkflows).Error
		if err != nil && err.Error() != "record not found" {
			log.Error("[workflow's getSequenceActionInfo]:error when get link start workflow info from db:", err.Error())
			return nil, errors.New("error when get workflow info")
		}

		runResult := models.WorkflowLogStateRunSuccess
		for _, linkstartInfo := range linkStartWorkflows {
			if linkstartInfo.RunState != models.WorkflowLogStateRunSuccess {
				runResult = models.WorkflowLogStateRunFailed
				break
			}
		}

		actionMap := make(map[string]interface{})
		actionMap["actionName"] = actionInfo.Action
		actionMap["actionId"] = actionInfo.ID
		actionMap["runResult"] = actionInfo.RunState
		actionMap["isTimeout"] = actionInfo.FailReason == ActionStopReasonTimeout
		actionMap["timeout"] = actionInfo.Timeout
		actionMap["runTime"] = strconv.FormatFloat(actionInfo.UpdatedAt.Sub(actionInfo.CreatedAt).Seconds(), 'f', 0, 64)
		actionMap["isStartWorkflow"] = len(linkStartWorkflows) > 0
		actionMap["startWorkflowResult"] = runResult
		actionMap["error"] = actionInfo.FailReason

		result = append(result, actionMap)
	}

	return result, nil
}

// GetWorkflow is
func GetWorkflow(workflowId int64) (*Workflow, error) {
	if workflowId == int64(0) {
		return nil, errors.New("workflow's id is empty")
	}

	workflowInfo := new(models.Workflow)
	err := workflowInfo.GetWorkflow().Where("id = ?", workflowId).First(workflowInfo).Error
	if err != nil {
		log.Error("[workflow's GetWorkflow]:error when get workflow info from db:", err.Error())
		return nil, err
	}

	workflow := new(Workflow)
	workflow.Workflow = workflowInfo

	return workflow, nil
}

// GetLatestRunableWorkflow is
func GetLatestRunableWorkflow(namespace, repository, workflowName, version string) (*Workflow, error) {
	if namespace == "" || repository == "" {
		log.Error("[workflow's GetLatestRunableWorkflow]:given empty parms:namespace: ===>", namespace, "<===  repository:===>", repository, "<===")
		return nil, errors.New("parms is empty")
	}

	workflowInfo := new(models.Workflow)
	query := workflowInfo.GetWorkflow().Where("namespace = ?", namespace).Where("repository = ?", repository).Where("workflow = ?", workflowName)

	if version != "" {
		query = query.Where("version = ?", version)
	}

	err := query.Where("state = ?", models.WorkflowStateAble).Order("-id").First(&workflowInfo).Error
	if err != nil {
		log.Error("[workflow's GetLatestRunableWorkflow]:error when get workflow info from db:", err.Error())
		return nil, err
	}

	if workflowInfo.ID == 0 {
		return nil, errors.New("no runable workflow")
	}

	workflow := new(Workflow)
	workflow.Workflow = workflowInfo

	return workflow, nil
}

// GetWorkflowLog is
func GetWorkflowLog(namespace, repository, workflowName, versionName string, sequence int64) (*WorkflowLog, error) {
	var err error
	workflowLogInfo := new(models.WorkflowLog)

	query := workflowLogInfo.GetWorkflowLog().Where("namespace =? ", namespace).Where("repository = ?", repository).Where("workflow = ?", workflowName).Where("version = ?", versionName)
	if sequence == int64(0) {
		query = query.Order("-id")
	} else {
		query = query.Where("sequence = ?", sequence)
	}

	err = query.First(workflowLogInfo).Error
	if err != nil {
		log.Error("[workflowLog's GetWorkflowLog]:error when get workflowLog(version=", versionName, ", sequence=", sequence, ") info from db:", err.Error())
		return nil, err
	}

	workflowLog := new(WorkflowLog)
	workflowLog.WorkflowLog = workflowLogInfo

	return workflowLog, nil
}

func getWorkflowEnvList(workflowLogId int64) ([]map[string]interface{}, error) {
	resultList := make([]map[string]interface{}, 0)
	workflowLog := new(models.WorkflowLog)
	err := workflowLog.GetWorkflowLog().Where("id = ?", workflowLogId).First(workflowLog).Error
	if err != nil {
		log.Error("[workflowLog's getWorkflowEnvList]:error when get workflowlog info from db:", err.Error())
		return nil, errors.New("error when get workflow info from db:" + err.Error())
	}

	envMap := make(map[string]string)
	if workflowLog.Env != "" {
		err = json.Unmarshal([]byte(workflowLog.Env), &envMap)
		if err != nil {
			log.Error("[workflowLog's getWorkflowEnvList]:error when unmarshal workflow's env setting:", workflowLog.Env, " ===>error is:", err.Error())
			return nil, errors.New("error when unmarshal workflow env info" + err.Error())
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

// GetWorkflowRunHistoryList is
func GetWorkflowRunHistoryList(namespace, repository string) ([]map[string]interface{}, error) {
	resultList := make([]map[string]interface{}, 0)
	workflowLogIndexMap := make(map[string]int)
	workflowLogVersionIndexMap := make(map[string]interface{})
	workflowLogList := make([]models.WorkflowLog, 0)

	err := new(models.WorkflowLog).GetWorkflowLog().Where("namespace = ?", namespace).Where("repository = ?", repository).Order("-id").Find(&workflowLogList).Error
	if err != nil && !strings.Contains(err.Error(), "record not found") {
		log.Error("[workflow's GetWorkflowRunHistoryList]:error when get workflowLog list from db:", err.Error())
		return nil, err
	}

	for _, workflowlog := range workflowLogList {
		if _, ok := workflowLogIndexMap[workflowlog.Workflow]; !ok {
			workflowInfoMap := make(map[string]interface{})
			workflowVersionListInfoMap := make([]map[string]interface{}, 0)
			workflowInfoMap["id"] = workflowlog.FromWorkflow
			workflowInfoMap["name"] = workflowlog.Workflow
			workflowInfoMap["versionList"] = workflowVersionListInfoMap

			resultList = append(resultList, workflowInfoMap)
			workflowLogIndexMap[workflowlog.Workflow] = len(resultList) - 1

			versionIndexMap := make(map[string]int64)
			workflowLogVersionIndexMap[workflowlog.Workflow] = versionIndexMap
		}

		workflowInfoMap := resultList[workflowLogIndexMap[workflowlog.Workflow]]
		if _, ok := workflowLogVersionIndexMap[workflowlog.Workflow].(map[string]int64)[workflowlog.Version]; !ok {
			workflowVersionInfoMap := make(map[string]interface{})
			workflowVersionSequenceListInfoMap := make([]map[string]interface{}, 0)
			workflowVersionInfoMap["id"] = workflowlog.ID
			workflowVersionInfoMap["name"] = workflowlog.Version
			workflowVersionInfoMap["info"] = ""
			workflowVersionInfoMap["total"] = int64(0)
			workflowVersionInfoMap["success"] = int64(0)
			workflowVersionInfoMap["sequenceList"] = workflowVersionSequenceListInfoMap

			workflowInfoMap["versionList"] = append(workflowInfoMap["versionList"].([]map[string]interface{}), workflowVersionInfoMap)
			workflowLogVersionIndexMap[workflowlog.Workflow].(map[string]int64)[workflowlog.Version] = int64(len(workflowInfoMap["versionList"].([]map[string]interface{})) - 1)
		}

		workflowVersionInfoMap := workflowInfoMap["versionList"].([]map[string]interface{})[workflowLogVersionIndexMap[workflowlog.Workflow].(map[string](int64))[workflowlog.Version]]
		sequenceList := workflowVersionInfoMap["sequenceList"].([]map[string]interface{})

		sequenceInfoMap := make(map[string]interface{})
		sequenceInfoMap["workflowSequenceID"] = workflowlog.ID
		sequenceInfoMap["sequence"] = workflowlog.Sequence
		sequenceInfoMap["status"] = workflowlog.RunState
		sequenceInfoMap["time"] = workflowlog.CreatedAt.Format("2006-01-02 15:04:05")

		sequenceList = append(sequenceList, sequenceInfoMap)
		workflowVersionInfoMap["sequenceList"] = sequenceList
		workflowVersionInfoMap["total"] = workflowVersionInfoMap["total"].(int64) + 1

		if workflowlog.RunState == models.WorkflowLogStateRunSuccess {
			workflowVersionInfoMap["success"] = workflowVersionInfoMap["success"].(int64) + 1
		}
	}

	for _, workflowInfoMap := range resultList {
		for _, versionInfoMap := range workflowInfoMap["versionList"].([]map[string]interface{}) {
			success := versionInfoMap["success"].(int64)
			total := versionInfoMap["total"].(int64)

			versionInfoMap["info"] = "Success: " + strconv.FormatInt(success, 10) + " Total: " + strconv.FormatInt(total, 10)
		}
	}

	return resultList, nil
}

// CreateNewVersion is
func (workflowInfo *Workflow) CreateNewVersion(define map[string]interface{}, versionName string) error {
	var count int64
	err := new(models.Workflow).GetWorkflow().Where("namespace = ?", workflowInfo.Namespace).Where("repository = ?", workflowInfo.Repository).Where("workflow = ?", workflowInfo.Workflow.Workflow).Where("version = ?", versionName).Count(&count).Error
	if count > 0 {
		return errors.New("version code already exist!")
	}

	// get current least workflow's version
	leastWorkflow := new(models.Workflow)
	err = leastWorkflow.GetWorkflow().Where("namespace = ? ", workflowInfo.Namespace).Where("workflow = ?", workflowInfo.Workflow.Workflow).Order("-id").First(&leastWorkflow).Error
	if err != nil {
		return errors.New("error when get least workflow info :" + err.Error())
	}

	newWorkflowInfo := new(models.Workflow)
	newWorkflowInfo.Namespace = workflowInfo.Namespace
	newWorkflowInfo.Repository = workflowInfo.Repository
	newWorkflowInfo.Workflow = workflowInfo.Workflow.Workflow
	newWorkflowInfo.Event = workflowInfo.Event
	newWorkflowInfo.Version = strings.TrimSpace(versionName)
	newWorkflowInfo.VersionCode = leastWorkflow.VersionCode + 1
	newWorkflowInfo.State = models.WorkflowStateDisable
	newWorkflowInfo.Manifest = workflowInfo.Manifest
	newWorkflowInfo.Description = workflowInfo.Description
	newWorkflowInfo.SourceInfo = workflowInfo.SourceInfo
	newWorkflowInfo.Env = workflowInfo.Env
	newWorkflowInfo.Requires = workflowInfo.Requires

	err = newWorkflowInfo.GetWorkflow().Save(newWorkflowInfo).Error
	if err != nil {
		return err
	}

	newInfo := new(Workflow)
	newInfo.Workflow = newWorkflowInfo
	return newInfo.UpdateWorkflowInfo(define)
}

// GetWorkflowToken is
func (workflowInfo *Workflow) GetWorkflowToken() (map[string]interface{}, error) {
	result := make(map[string]interface{})

	if workflowInfo.ID == 0 {
		log.Error("[workflow's GetWorkflowToken]:got an empty pipelin:", workflowInfo)
		return nil, errors.New("workflow's info is empty")
	}

	token := ""
	tokenMap := make(map[string]interface{})
	if workflowInfo.SourceInfo == "" {
		// if sourceInfo is empty generate a token
		token = workflowInfo.Workflow.Workflow
	} else {
		json.Unmarshal([]byte(workflowInfo.SourceInfo), &tokenMap)

		if _, ok := tokenMap["token"].(string); !ok {
			token = workflowInfo.Workflow.Workflow
		} else {
			token = tokenMap["token"].(string)
		}
	}

	tokenMap["token"] = token
	sourceInfo, _ := json.Marshal(tokenMap)
	workflowInfo.SourceInfo = string(sourceInfo)
	err := workflowInfo.GetWorkflow().Save(workflowInfo).Error

	if err != nil {
		log.Error("[workflow's GetWorkflowToken]:error when save workflow's info to db:", err.Error())
		return nil, errors.New("error when get workflow info from db:" + err.Error())
	}

	result["token"] = token

	url := ""

	projectAddr := ""
	if configure.GetString("projectaddr") == "" {
		projectAddr = "current-workflow's-ip:port"
	} else {
		projectAddr = configure.GetString("projectaddr")
	}

	url += projectAddr
	url = strings.TrimSuffix(url, "/")
	url += "/v2" + "/" + workflowInfo.Namespace + "/" + workflowInfo.Repository + "/workflow/v1/exec/" + workflowInfo.Workflow.Workflow

	result["url"] = url

	return result, nil
}

// UpdateWorkflowInfo is
func (workflowInfo *Workflow) UpdateWorkflowInfo(define map[string]interface{}) error {
	db := models.GetDB().Begin()
	err := db.Error
	if err != nil {
		log.Error("[workflow's UpdateWorkflowInfo]:when db.Begin():", err.Error())
		return err
	}

	workflowOriginalManifestMap := make(map[string]interface{})
	if workflowInfo.Manifest != "" {
		err := json.Unmarshal([]byte(workflowInfo.Manifest), &workflowOriginalManifestMap)
		if err != nil {
			log.Error("[workflow's UpdateWorkflowInfo]:error unmarshal workflow's manifest info:", err.Error(), " set it to empty")
			workflowInfo.Manifest = ""
		}
	}

	workflowOriginalManifestMap["define"] = define
	workflowNewManifestBytes, err := json.Marshal(workflowOriginalManifestMap)
	if err != nil {
		log.Error("[workflow's UpdateWorkflowInfo]:error when marshal workflow's manifest info:", err.Error())
		rollbackErr := db.Rollback().Error
		if rollbackErr != nil {
			log.Error("[workflow's UpdateWorkflowInfo]:when rollback in save workflow's info:", rollbackErr.Error())
			return errors.New("errors occur:\nerror1:" + err.Error() + "\nerror2:" + rollbackErr.Error())
		}
		return errors.New("error when save workflow's define info:" + err.Error())
	}

	requestMap := make([]interface{}, 0)
	if request, ok := define["request"]; ok {
		if requestMap, ok = request.([]interface{}); !ok {
			log.Error("[workflow's UpdateWorkflowInfo]:error when get workflow's request info:want a json array,got:", request)
			return errors.New("error when get workflow's request info,want a json array")
		}
	} else {
		defaultRequestMap := make(map[string]interface{})
		defaultRequestMap["type"] = AuthTypeWorkflowDefault
		defaultRequestMap["token"] = AuthTokenDefault

		requestMap = append(requestMap, defaultRequestMap)
	}

	requestInfo, err := json.Marshal(requestMap)
	if err != nil {
		log.Error("[workflow's UpdateWorkflowInfo]:error when marshal workflow's request info:", requestMap, " ===>error is:", err.Error())
		return errors.New("error when save workflow's request info")
	}

	workflowInfo.State = models.WorkflowStateDisable
	workflowInfo.Manifest = string(workflowNewManifestBytes)
	workflowInfo.Requires = string(requestInfo)
	err = db.Save(workflowInfo).Error
	if err != nil {
		log.Error("[workflow's UpdateWorkflowInfo]:when save workflow's info:", workflowInfo, " ===>error is:", err.Error())
		rollbackErr := db.Rollback().Error
		if rollbackErr != nil {
			log.Error("[workflow's UpdateWorkflowInfo]:when rollback in save workflow's info:", rollbackErr.Error())
			return errors.New("errors occur:\nerror1:" + err.Error() + "\nerror2:" + rollbackErr.Error())
		}
		return err
	}

	relationMap, stageDefineList, settingMap, err := workflowInfo.getWorkflowDefineInfo(workflowInfo.Workflow)
	if err != nil {
		log.Error("[workflow's UpdateWorkflowInfo]:when get workflow's define info:", err.Error())
		rollbackErr := db.Rollback().Error
		if rollbackErr != nil {
			log.Error("[workflow's UpdateWorkflowInfo]:when rollback after get workflow define info:", rollbackErr.Error())
			return errors.New("errors occur:\nerror1:" + err.Error() + "\nerror2:" + rollbackErr.Error())
		}
		return err
	}

	// first delete old workflow define
	err = db.Model(&models.Action{}).Where("workflow = ?", workflowInfo.ID).Delete(&models.Action{}).Error
	if err != nil {
		log.Error("[workflow's UpdateWorkflowInfo]:when delete action's that belong workflow:", workflowInfo, " ===>error is:", err.Error())
		rollbackErr := db.Rollback().Error
		if rollbackErr != nil {
			log.Error("[workflow's UpdateWorkflowInfo]:when rollback in delete action info:", rollbackErr.Error())
			return errors.New("errors occur:\nerror1:" + err.Error() + "\nerror2:" + rollbackErr.Error())
		}
		return errors.New("error when remove old action info:" + err.Error())
	}

	err = db.Model(&models.Stage{}).Where("workflow = ?", workflowInfo.ID).Delete(&models.Stage{}).Error
	if err != nil {
		log.Error("[workflow's UpdateWorkflowInfo]:when delete stage's that belong workflow:", workflowInfo, " ===>error is:", err.Error())
		rollbackErr := db.Rollback().Error
		if rollbackErr != nil {
			log.Error("[workflow's UpdateWorkflowInfo]:when rollback in delete stage info:", rollbackErr.Error())
			return errors.New("errors occur:\nerror1:" + err.Error() + "\nerror2:" + rollbackErr.Error())
		}
		return errors.New("error when update stage info:" + err.Error())
	}

	// then create new workflow by define
	stageInfoMap := make(map[string]map[string]interface{})
	preStageId := int64(-1)
	allActionIdMap := make(map[string]int64)
	for _, stageDefine := range stageDefineList {
		stageId, stageTagId, actionMap, err := CreateNewStage(db, preStageId, workflowInfo.Workflow, stageDefine, relationMap)
		if err != nil {
			log.Error("[workflow's UpdateWorkflowInfo]:error when create new stage that workflow define:", stageDefine, " preStage is :", preStageId, " workflow is:", workflowInfo, " relation is:", relationMap)
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
					log.Error("[workflow's UpdateWorkflowInfo]:error when get action's relation info in map:", allActionIdMap, " want :", fromActionOriginId)
					rollbackErr := db.Rollback().Error
					if rollbackErr != nil {
						log.Error("[workflow's UpdateWorkflowInfo]:when rollback in get action relation info:", rollbackErr.Error())
						return errors.New("errors occur:\nerror1:" + err.Error() + "\nerror2:" + rollbackErr.Error())
					}
					return errors.New("action's relation is illegal")
				}

				tempRelation := make(map[string]interface{})
				tempRelation["toAction"] = actionID
				tempRelation["fromAction"] = fromActionId
				tempRelation["relation"] = realRelations

				actionRealtionList = append(actionRealtionList, tempRelation)
			}

			actionInfo := new(models.Action)
			err = db.Model(&models.Action{}).Where("id = ?", actionID).First(&actionInfo).Error
			if err != nil {
				log.Error("[workflow's UpdateWorkflowInfo]:error when get action info from db:", actionID, " ===>error is:", err.Error())
				rollbackErr := db.Rollback().Error
				if rollbackErr != nil {
					log.Error("[workflow's UpdateWorkflowInfo]:when rollback in get action info from db:", rollbackErr.Error())
					return errors.New("errors occur:\nerror1:" + err.Error() + "\nerror2:" + rollbackErr.Error())
				}
				return err
			}
			manifestMap := make(map[string]interface{})
			if actionInfo.Manifest != "" {
				json.Unmarshal([]byte(actionInfo.Manifest), &manifestMap)
			}

			manifestMap["relation"] = actionRealtionList
			relationBytes, _ := json.Marshal(manifestMap)
			actionInfo.Manifest = string(relationBytes)

			err = db.Model(&models.Action{}).Where("id = ?", actionID).UpdateColumn("manifest", actionInfo.Manifest).Error
			if err != nil {
				log.Error("[workflow's UpdateWorkflowInfo]:error when update action's column manifest:", actionInfo, " ===>error is:", err.Error())
				rollbackErr := db.Rollback().Error
				if rollbackErr != nil {
					log.Error("[workflow's UpdateWorkflowInfo]:when rollback in update action's column info from db:", rollbackErr.Error())
					return errors.New("errors occur:\nerror1:" + err.Error() + "\nerror2:" + rollbackErr.Error())
				}
				return err
			}
		}
	}

	if instanceSetting, ok := settingMap["runningInstances"]; ok {
		setting, ok := instanceSetting.(map[string]interface{})
		if !ok {
			log.Error("[workflow's UpdateWorkflowInfo]:error when parse instanceSetting map,want a json obj, got:", settingMap["runningInstances"])
			rollbackErr := db.Rollback().Error
			if rollbackErr != nil {
				log.Error("[workflow's UpdateWorkflowInfo]:when rollback in parse instance Setting map:", rollbackErr.Error())
				return errors.New("errors occur:\nerror1:" + err.Error() + "\nerror2:" + rollbackErr.Error())
			}
			return errors.New("error when save workflow's define info:" + err.Error())
		}

		available, ok := setting["available"].(bool)
		if !ok {
			available = false
		}

		instanceNum := int64(0)
		instanceNumF, ok := setting["number"].(float64)
		if ok {
			instanceNum = int64(instanceNumF)
		}

		workflowInfo.IsLimitInstance = available
		workflowInfo.LimitInstance = instanceNum

		err = db.Save(workflowInfo).Error
		if err != nil {
			log.Error("[workflow's UpdateWorkflowInfo]:when save workflow's info:", workflowInfo, " ===>error is:", err.Error())
			rollbackErr := db.Rollback().Error
			if rollbackErr != nil {
				log.Error("[workflow's UpdateWorkflowInfo]:when rollback in save workflow's info:", rollbackErr.Error())
				return errors.New("errors occur:\nerror1:" + err.Error() + "\nerror2:" + rollbackErr.Error())
			}
			return err
		}
	}

	if taskSetting, ok := settingMap["timedTasks"]; ok {
		taskMap, ok := taskSetting.(map[string]interface{})
		if !ok {
			log.Error("[workflow's UpdateWorkflowInfo]:error when get task map info,want a json obj, got:", taskMap)
			rollbackErr := db.Rollback().Error
			if rollbackErr != nil {
				log.Error("[workflow's UpdateWorkflowInfo]:when rollback in got a task setting map:", rollbackErr.Error())
				return errors.New("errors occur:\nerror1:" + err.Error() + "\nerror2:" + rollbackErr.Error())
			}
			return errors.New("error when save workflow's define info:" + err.Error())
		}
		err = workflowInfo.updateTimerTask(taskMap)
	}

	db.Commit()
	return nil
}

// DeleteWorkflow is
func (workflowInfo *Workflow) DeleteWorkflow() error {
	return workflowInfo.GetWorkflow().Delete(&workflowInfo).Error
}

func (workflowInfo *Workflow) updateTimerTask(taskMap map[string]interface{}) error {
	db := models.GetDB().Begin()
	available, ok := taskMap["available"].(bool)
	if !ok {
		available = false
	}

	db.Model(&models.Timer{}).Where("namespace = ?", workflowInfo.Namespace).Where("repository = ?", workflowInfo.Repository).Where("workflow = ?", workflowInfo.ID).Delete(&models.Timer{})
	if taskList, ok := taskMap["tasks"].([]interface{}); ok {
		for _, task := range taskList {
			if taskMap, ok := task.(map[string]interface{}); ok {
				cron, ok := taskMap["cronEntry"].(string)
				if !ok {
					log.Error("[workflow's updateTimerTask]:error when get cronEntry:want a string, got:", taskMap["cronEntry"])
					continue
				}

				if len(strings.Split(cron, " ")) == 5 {
					cron = "0 " + cron
				}

				eventName, ok := taskMap["eventName"].(string)
				if !ok {
					log.Error("[workflow's updateTimerTask]:error when get eventName:want a string, got:", taskMap["eventName"])
					continue
				}

				eventType, ok := taskMap["eventType"].(string)
				if !ok {
					log.Error("[workflow's updateTimerTask]:error when get eventType:want a string, got:", taskMap["eventType"])
					continue
				}

				startJson, ok := taskMap["startJson"].(map[string]interface{})
				if !ok {
					log.Error("[workflow's updateTimerTask]:error when get startJson:want a json obj, got:", taskMap["startJson"])
					continue
				}

				startJsonBytes, _ := json.Marshal(startJson)

				timer := new(models.Timer)
				timer.Namespace = workflowInfo.Namespace
				timer.Repository = workflowInfo.Repository
				timer.Workflow = workflowInfo.ID
				timer.Available = available
				timer.Cron = cron
				timer.EventType = eventType
				timer.EventName = eventName
				timer.StartJson = string(startJsonBytes)

				db.Save(timer)
			}
		}
	}
	db.Commit()
	UpdateWorkflowTimer(workflowInfo.Namespace, workflowInfo.Repository, workflowInfo.ID)
	return nil
}

func (workflowInfo *Workflow) getWorkflowDefineInfo(workflow *models.Workflow) (map[string]interface{}, []map[string]interface{}, map[string]interface{}, error) {
	lineList := make([]map[string]interface{}, 0)
	stageList := make([]map[string]interface{}, 0)

	manifestMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(workflow.Manifest), &manifestMap)
	if err != nil {
		log.Error("[workflow's getWorkflowDefineInfo]:error when unmarshal workflow's manifes info:", workflowInfo.Manifest, " ===>error is:", err.Error())
		return nil, nil, nil, errors.New("error when unmarshal workflow manifes info:" + err.Error())
	}

	defineMap, ok := manifestMap["define"].(map[string]interface{})
	if !ok {
		log.Error("[workflow's getWorkflowDefineInfo]:workflow's define is not a json obj:", manifestMap["define"])
		return nil, nil, nil, errors.New("workflow's define is not a json:" + err.Error())
	}

	realtionMap := make(map[string]interface{})
	if linesList, ok := defineMap["lineList"].([]interface{}); ok {
		if !ok {
			log.Error("[workflow's getWorkflowDefineInfo]:error in workflow's lineList define,want a array,got:", defineMap["lineList"])
			return nil, nil, nil, errors.New("workflow's lineList define is not an array")
		}

		for _, lineInfo := range linesList {
			lineInfoMap, ok := lineInfo.(map[string]interface{})
			if !ok {
				log.Error("[workflow's getWorkflowDefineInfo]:error in workflow's line define: want a json obj,got:", lineInfo)
				return nil, nil, nil, errors.New("workflow's line info is not a json")
			}

			lineList = append(lineList, lineInfoMap)
		}

		for _, lineInfo := range lineList {
			endData, ok := lineInfo["endData"].(map[string]interface{})
			if !ok {
				log.Error("[workflow's getWorkflowDefineInfo]:error in workflow's line define:line doesn't define any end point info:", lineInfo)
				return nil, nil, nil, errors.New("workflow's line define is illegal,don't have a end point info")
			}

			endPointId, ok := endData["id"].(string)
			if !ok {
				log.Error("[workflow's getWorkflowDefineInfo]:error in workflow's line define:end point's id is not a string:", endData)
				return nil, nil, nil, errors.New("workflow's line define is illegal,endPoint id is not a string")
			}

			if _, ok := realtionMap[endPointId]; !ok {
				realtionMap[endPointId] = make(map[string]interface{})
			}

			endPointMap := realtionMap[endPointId].(map[string]interface{})
			startData, ok := lineInfo["startData"].(map[string]interface{})
			if !ok {
				log.Error("[workflow's getWorkflowDefineInfo]:error in workflow's line define:line doesn't define any start point info:", lineInfo)
				return nil, nil, nil, errors.New("workflow's line define is illegal,don;t have a start point info")
			}

			startDataId, ok := startData["id"].(string)
			if !ok {
				log.Error("[workflow's getWorkflowDefineInfo]:error in workflow's line define:start point's id is not a string:", endData)
				return nil, nil, nil, errors.New("workflow's line define is illegal,startPoint id is not a string")
			}

			if startDataId == "start-stage" {
				if _, ok := endPointMap[startDataId]; !ok {
					endPointMap[startDataId] = make(map[string]interface{}, 0)
				}

				lineMap, ok := lineInfo["relation"].(map[string]interface{})
				if !ok {
					continue
				}

				lineOriginMap, _ := endPointMap[startDataId].(map[string]interface{})
				for key, value := range lineMap {
					lineOriginMap[key] = value
				}

				endPointMap[startDataId] = lineOriginMap
			} else {
				if _, ok := endPointMap[startDataId]; !ok {
					endPointMap[startDataId] = make([]interface{}, 0)
				}

				lineList, ok := lineInfo["relation"].([]interface{})
				if !ok {
					continue
				}

				endPointMap[startDataId] = append(endPointMap[startDataId].([]interface{}), lineList...)
			}
		}
	}

	stageListInfo, ok := defineMap["stageList"]
	if !ok {
		log.Error("[workflow's getWorkflowDefineInfo]:error in workflow's define:workflow doesn't define any stage info", defineMap)
		return nil, nil, nil, errors.New("workflow don't have a stage define")
	}

	stagesList, ok := stageListInfo.([]interface{})
	if !ok {
		log.Error("[workflow's getWorkflowDefineInfo]:error in stageList's define:want array,got:", stageListInfo)
		return nil, nil, nil, errors.New("workflow's stageList define is not an array")
	}

	for _, stageInfo := range stagesList {
		stageInfoMap, ok := stageInfo.(map[string]interface{})
		if !ok {
			log.Error("[workflow's getWorkflowDefineInfo]:error in stage's define,want a json obj,got:", stageInfo)
			return nil, nil, nil, errors.New("workflow's stage info is not a json obj")
		}

		stageList = append(stageList, stageInfoMap)
	}

	setting, ok := defineMap["setting"].(map[string]interface{})
	if !ok {
		log.Error("[workflow's getWorkflowDefineInfo]:error when get workflow setting: want a json obj, got:", defineMap["setting"])
		return nil, nil, nil, errors.New("workflow's setting info is not a json obj")
	}

	setting, ok = setting["data"].(map[string]interface{})
	if !ok {
		log.Error("[workflow's getWorkflowDefineInfo]:error when get workflow setting's data: want a json obj, got:", defineMap["setting"])
		return nil, nil, nil, errors.New("workflow's setting info is not a json obj")
	}

	return realtionMap, stageList, setting, nil
}

// BeforeExecCheck is
func (workflowInfo *Workflow) BeforeExecCheck(reqHeader http.Header, reqBody []byte) (bool, map[string]string, error) {
	if workflowInfo.SourceInfo == "" {
		return false, nil, errors.New("workflow's source info is empty")
	}

	sourceMap := make(map[string]interface{})
	sourceList := make([]interface{}, 0)
	err := json.Unmarshal([]byte(workflowInfo.SourceInfo), &sourceMap)
	if err != nil {
		log.Error("[workflow's BeforeExecCheck]:error when unmarshal workflow source info, want json obj, got:", workflowInfo.SourceInfo)
		return false, nil, errors.New("workflow's source define error")
	}

	expectedToken, ok := sourceMap["token"].(string)
	if !ok {
		log.Error("[workflow's BeforeExecCheck]:error when get source's expected token,want a string, got:", sourceMap["token"])
		return false, nil, errors.New("get token error")
	}

	sourceList, ok = sourceMap["sourceList"].([]interface{})
	if !ok {
		log.Error("[workflow's BeforeExecCheck]:error when get sourceList:want json array, got:", sourceMap["sourceList"])
		return false, nil, errors.New("workflow's sourceList define error")
	}

	eventInfoMap, err := getExecReqEventInfo(sourceList, reqHeader)
	if err != nil {
		log.Error("[workflow's BeforeExecCheck]:error when get exec request's event type and event info:", err.Error())
		return false, nil, errors.New("get req's event info failed:" + err.Error())
	}

	passCheck := true

	checkerList, err := checker.GetWorkflowExecCheckerList()
	if err != nil {
		log.Error("[workflow's BeforeExecCheck]:error when get checkerList:", err.Error())
		return false, nil, err
	}

	for _, c := range checkerList {
		if c.Support(eventInfoMap) {
			passCheck, err = c.Check(eventInfoMap, expectedToken, reqHeader, reqBody)
			if !passCheck {
				// log.Error("[workflow's BeforeExecCheck]:check failed:", c, "===>", err, "\neventInfoMap:", eventInfoMap, "\nreqHeader:", reqHeader, "\nreqBody:", string(reqBody))
				log.Error("[workflow's BeforeExecCheck]:check failed:", c, "===>", err)
				return false, nil, errors.New("failed when check exec req")
			}
		}
	}

	// check run instance number
	pass, err := checkInstanceNum(workflowInfo.ID)
	if !pass {
		return false, nil, err
	}

	return passCheck, eventInfoMap, nil
}

func checkInstanceNum(workflowInfoID int64) (bool, error) {
	// check run instance number
	db := models.GetDB().Begin()
	workflow := new(models.Workflow)
	err := db.Raw("select * from workflow where id = ? for update", workflowInfoID).Scan(workflow).Error
	if err != nil {
		log.Error("[workflow's BeforeExecCheck]:error when get workflow info  from db:", err.Error())
		db.Rollback()
		return false, errors.New("failed when check exec req")
	}

	if workflow.IsLimitInstance {
		if workflow.LimitInstance <= workflow.CurrentInstance {
			log.Error("[workflow's BeforeExecCheck]:workflow:", workflow.Workflow, " current run:", workflow.CurrentInstance, " max run num:", workflow.LimitInstance, " start failed ...")
			db.Rollback()
			return false, errors.New("no useable run instance")
		}
	}

	workflow.CurrentInstance += 1

	db.Save(workflow)
	db.Commit()

	return true, nil
}

func getExecReqEventInfo(sourceList []interface{}, reqHeader http.Header) (map[string]string, error) {
	result := make(map[string]string)
	for _, sourceConfigInfo := range sourceList {
		sourceConfig, ok := sourceConfigInfo.(map[string]interface{})
		if !ok {
			log.Error("[workflow's getExecReqEventInfo]:error when parse sourceConfig,want a json obj, got :", sourceConfigInfo)
			return nil, errors.New("source config is not a json obj")
		}

		tokenKey, ok := sourceConfig["headerKey"].(string)
		if !ok {
			log.Error("[workflow's getExecReqEventInfo]:error when get source's token key,want a string, got:", sourceConfig["headerKey"])
			return nil, errors.New("source's token key is not a string")
		}

		token := reqHeader.Get(tokenKey)
		if token != "" {
			supportEventList, ok := sourceConfig["eventList"].(string)
			if !ok {
				log.Error("[workflow's getExecReqEventInfo]:error when get source's support event list,want a string, got:", sourceConfig["eventList"])
				continue
			}

			sourceType, ok := sourceConfig["sourceType"].(string)
			if !ok {
				log.Error("[workflow's getExecReqEventInfo]:error when get source's sourceType,want a string, got:", sourceConfig["sourceType"])
				continue
			}

			eventName := getEventName(sourceType, reqHeader)
			if !strings.Contains(supportEventList, ","+eventName+",") {
				continue
			}

			result["sourceType"] = sourceType
			result["eventName"] = eventName
			result["token"] = token

			return result, nil
		}
	}
	return nil, errors.New("can't get event info from request header")
}

func getEventName(sourceType string, reqHeader http.Header) string {
	eventName := ""
	switch sourceType {
	case "github":
		eventName = reqHeader.Get("X-Github-Event")
	case "gitlab":
		eventName = reqHeader.Get("X-Gitlab-Event")
	case "customize":
		eventName = reqHeader.Get("X-Workflow-Event")
	}

	return eventName
}

// GenerateNewLog is
func (workflowInfo *Workflow) GenerateNewLog(eventMap map[string]string) (*WorkflowLog, error) {
	workflowlogSequenceGenerateChan <- true
	result := new(WorkflowLog)
	stageList := make([]models.Stage, 0)

	workflowSequence := new(models.WorkflowSequence)
	workflowSequence.Workflow = workflowInfo.ID
	err := workflowSequence.GetWorkflowSequence().Save(workflowSequence).Error
	if err != nil {
		<-workflowlogSequenceGenerateChan
		log.Error("[workflow's GenerateNewLog]:error when save workflow sequence info to db", workflowSequence, "===>error is :", err.Error())
		return nil, err
	}

	var count int64
	err = workflowSequence.GetWorkflowSequence().Where("id < ?", workflowSequence.ID).Where("workflow = ?", workflowInfo.ID).Count(&count).Error
	if err != nil && !strings.Contains(err.Error(), "record not found") {
		<-workflowlogSequenceGenerateChan
		log.Error("[workflow's GenerateNewLog]:error when get workflow sequence info to db:", err.Error())
		return nil, err
	}

	workflowSequence.Sequence = count + 1
	err = workflowSequence.GetWorkflowSequence().Save(workflowSequence).Error
	if err != nil {
		<-workflowlogSequenceGenerateChan
		log.Error("[workflow's GenerateNewLog]:error when save workflow sequence info to db", workflowSequence, "===>error is :", err.Error())
		return nil, err
	}

	err = new(models.Stage).GetStage().Where("workflow = ?", workflowInfo.ID).Find(&stageList).Error
	if err != nil {
		<-workflowlogSequenceGenerateChan
		log.Error("[workflow's GenerateNewLog]:error when get stage list by workflow info", workflowInfo, "===>error is :", err.Error())
		return nil, err
	}

	db := models.GetDB()
	db = db.Begin()

	eventInfoBytes, _ := json.Marshal(eventMap)

	// record workflow's info
	workflowLog := new(models.WorkflowLog)
	workflowLog.Namespace = workflowInfo.Namespace
	workflowLog.Repository = workflowInfo.Repository
	workflowLog.Workflow = workflowInfo.Workflow.Workflow
	workflowLog.FromWorkflow = workflowInfo.ID
	workflowLog.Version = workflowInfo.Version
	workflowLog.VersionCode = workflowInfo.VersionCode
	workflowLog.Sequence = workflowSequence.Sequence
	workflowLog.RunState = models.WorkflowLogStateCanListen
	workflowLog.Event = workflowInfo.Event
	workflowLog.Manifest = workflowInfo.Manifest
	workflowLog.Description = workflowInfo.Description
	workflowLog.SourceInfo = string(eventInfoBytes)
	workflowLog.Env = workflowInfo.Env
	workflowLog.Requires = workflowInfo.Requires
	workflowLog.AuthList = ""

	err = db.Save(workflowLog).Error
	if err != nil {
		<-workflowlogSequenceGenerateChan
		log.Error("[workflow's GenerateNewLog]:when save workflow log to db:", workflowLog, "===>error is :", err.Error())
		rollbackErr := db.Rollback().Error
		if rollbackErr != nil {
			log.Error("[workflow's GenerateNewLog]:when rollback in save workflow log:", rollbackErr.Error())
			return nil, errors.New("errors occur:\nerror1:" + err.Error() + "\nerror2:" + rollbackErr.Error())
		}
		return nil, err
	}

	preStageLogId := int64(-1)
	for _, stageInfo := range stageList {
		stage := new(Stage)
		stage.Stage = &stageInfo
		preStageLogId, err = stage.GenerateNewLog(db, workflowLog, preStageLogId)
		if err != nil {
			<-workflowlogSequenceGenerateChan
			log.Error("[workflow's GenerateNewLog]:when generate stage log:", err.Error())
			return nil, err
		}
	}

	allVarList := make([]models.WorkflowVar, 0)
	new(models.WorkflowVar).GetWorkflowVar().Where("workflow = ?", workflowInfo.ID).Find(&allVarList)

	for _, varInfo := range allVarList {
		tempVar := new(WorkflowVar)
		tempVar.WorkflowVar = &varInfo
		err := tempVar.GenerateNewLog(db, workflowLog)
		if err != nil {
			<-workflowlogSequenceGenerateChan
			log.Error("[workflow's GenerateNewLog]:when generate var log:", err.Error())
			return nil, err
		}
	}

	err = db.Commit().Error
	if err != nil {
		<-workflowlogSequenceGenerateChan
		log.Error("[workflow's GenerateNewLog]:when commit to db:", err.Error())
		return nil, errors.New("error when save workflow info to db:" + err.Error())
	}
	result.WorkflowLog = workflowLog
	<-workflowlogSequenceGenerateChan
	return result, nil
}

// GetDefineInfo is
func (workflowLog *WorkflowLog) GetDefineInfo() (map[string]interface{}, error) {
	defineMap := make(map[string]interface{})
	stageListMap := make([]map[string]interface{}, 0)
	lineList := make([]map[string]interface{}, 0)

	stageList := make([]*models.StageLog, 0)
	err := new(models.StageLog).GetStageLog().Where("workflow = ?", workflowLog.ID).Find(&stageList).Error
	if err != nil {
		log.Error("[StageLog's GetStageLogDefineListByWorkflowLogID]:error when get stage list from db:", err.Error())
		return nil, err
	}

	for _, stageInfo := range stageList {
		stage := new(StageLog)
		stage.StageLog = stageInfo
		stageDefineMap, err := stage.GetStageLogDefine()
		if err != nil {
			log.Error("[workflowLog's GetDefineInfo]:error when get stagelog define:", stage, " ===>error is:", err.Error())
			return nil, err
		}

		stageListMap = append(stageListMap, stageDefineMap)
	}

	for _, stageInfo := range stageList {
		if stageInfo.Type == models.StageTypeStart || stageInfo.Type == models.StageTypeEnd {
			continue
		}

		actionList := make([]*models.ActionLog, 0)
		err = new(models.ActionLog).GetActionLog().Where("stage = ?", stageInfo.ID).Find(&actionList).Error
		if err != nil {
			log.Error("[workflowLog's GetDefineInfo]:error when get actionlog list from db:", err.Error())
			continue
		}

		for _, actionInfo := range actionList {
			action := new(ActionLog)
			action.ActionLog = actionInfo
			actionLineInfo, err := action.GetActionLineInfo()
			if err != nil {
				log.Error("[workflowLog's GetDefineInfo]:error when get actionlog line info:", err.Error())
				continue
			}

			lineList = append(lineList, actionLineInfo...)
		}
	}

	defineMap["workflow"] = workflowLog.Workflow
	defineMap["version"] = workflowLog.Version
	defineMap["sequence"] = workflowLog.Sequence
	defineMap["status"] = workflowLog.RunState
	defineMap["lineList"] = lineList
	defineMap["stageList"] = stageListMap

	return defineMap, nil
}

// GetStartStageData is
func (workflowLog *WorkflowLog) GetStartStageData() (map[string]interface{}, error) {
	dataMap := make(map[string]interface{})
	outCome := new(models.Outcome)
	err := outCome.GetOutcome().Where("workflow = ?", workflowLog.ID).Where("sequence = ?", workflowLog.Sequence).Where("action = ?", models.OutcomeTypeStageStartActionID).First(outCome).Error
	if err != nil && !strings.Contains(err.Error(), "record not found") {
		log.Error("[workflowLog's GetStartStageData]:error when get start stage data from db:", err.Error())
		return nil, err
	}

	err = json.Unmarshal([]byte(outCome.Output), &dataMap)
	if err != nil {
		log.Error("[workflowLog's GetStartStageData]:error when unmarshal start stage's data:", outCome.Output, " ===>error is:", err.Error())
	}

	return dataMap, nil
}

// Listen is
func (workflowLog *WorkflowLog) Listen(startData string) error {
	workflowlogListenChan <- true
	defer func() { <-workflowlogListenChan }()

	err := workflowLog.GetWorkflowLog().Where("id = ?", workflowLog.ID).First(workflowLog).Error
	if err != nil {
		log.Error("[workflowLog's Listen]:error when get workflowLog info from db:", workflowLog, " ===>error is:", err.Error())
		return errors.New("error when get workflowlog's info from db:" + err.Error())
	}

	if workflowLog.RunState != models.WorkflowLogStateCanListen {
		log.Error("[workflowLog's Listen]:error workflowlog state:", workflowLog)
		return errors.New("can't listen curren workflowlog,current state is:" + strconv.FormatInt(workflowLog.RunState, 10))
	}

	workflowLog.RunState = models.WorkflowLogStateWaitToStart
	err = workflowLog.GetWorkflowLog().Save(workflowLog).Error
	if err != nil {
		log.Error("[workflowLog's Listen]:error when change workflowlog's run state to wait to start:", workflowLog, " ===>error is:", err.Error())
		return errors.New("can't listen target workflow,change workflow's state failed")
	}

	canStartChan := make(chan bool, 1)
	go func() {
		for true {
			time.Sleep(1 * time.Second)

			err := workflowLog.GetWorkflowLog().Where("id = ?", workflowLog.ID).First(workflowLog).Error
			if err != nil {
				log.Error("[workflowLog's Listen]:error when get workflowLog's info:", workflowLog, " ===>error is:", err.Error())
				canStartChan <- false
				break
			}
			if workflowLog.Requires == "" || workflowLog.Requires == "[]" {
				log.Info("[workflowLog's Listen]:workflowLog", workflowLog, "is ready and will start")
				canStartChan <- true
				break
			}
		}
	}()

	go func() {
		canStart := <-canStartChan
		if !canStart {
			log.Error("[workflowLog's Listen]:workflowLog can't start", workflowLog)
			workflowLog.Stop(WorkflowStopReasonRunFailed, models.WorkflowLogStateRunFailed)
		} else {
			go workflowLog.Start(startData)
		}

	}()

	return nil
}

// Auth is
func (workflowLog *WorkflowLog) Auth(authMap map[string]interface{}) error {
	workflowlogAuthChan <- true
	defer func() { <-workflowlogAuthChan }()

	authType, ok := authMap["type"].(string)
	if !ok {
		log.Error("[workflowLog's Auth]:error when get authType from given authMap:", authMap, " ===>to workflowlog:", workflowLog)
		return errors.New("authType is illegal")
	}

	token, ok := authMap["token"].(string)
	if !ok {
		log.Error("[workflowLog's Auth]:error when get token from given authMap:", authMap, " ===>to workflowlog:", workflowLog)
		return errors.New("token is illegal")
	}

	err := workflowLog.GetWorkflowLog().Where("id = ?", workflowLog.ID).First(workflowLog).Error
	if err != nil {
		log.Error("[workflowLog's Auth]:error when get workflowLog info from db:", workflowLog, " ===>error is:", err.Error())
		return errors.New("error when get workflowlog's info from db:" + err.Error())
	}

	if workflowLog.Requires == "" || workflowLog.Requires == "[]" {
		log.Error("[workflowLog's Auth]:error when set auth info,workflowlog's requires is empty", authMap, " ===>to workflowlog:", workflowLog)
		return errors.New("workflow don't need any more auth")
	}

	requireList := make([]interface{}, 0)
	remainRequireList := make([]interface{}, 0)
	err = json.Unmarshal([]byte(workflowLog.Requires), &requireList)
	if err != nil {
		log.Error("[workflowLog's Auth]:error when unmarshal workflowlog's require list:", workflowLog, " ===>error is:", err.Error())
		return errors.New("error when get workflow require auth info:" + err.Error())
	}

	hasAuthed := false
	for _, require := range requireList {
		requireMap, ok := require.(map[string]interface{})
		if !ok {
			log.Error("[workflowLog's Auth]:error when get workflowlog's require info:", workflowLog, " ===> require is:", require)
			return errors.New("error when get workflow require auth info,require is not a json object")
		}

		requireType, ok := requireMap["type"].(string)
		if !ok {
			log.Error("[workflowLog's Auth]:error when get workflowlog's require type:", workflowLog, " ===> require map is:", requireMap)
			return errors.New("error when get workflow require auth info,require don't have a type")
		}

		requireToken, ok := requireMap["token"].(string)
		if !ok {
			log.Error("[workflowLog's Auth]:error when get workflowlog's require token:", workflowLog, " ===> require map is:", requireMap)
			return errors.New("error when get workflow require auth info,require don't have a token")
		}

		if requireType == authType && requireToken == token {
			hasAuthed = true
			// record auth info to workflowlog's auth info list
			workflowLogAuthList := make([]interface{}, 0)
			if workflowLog.AuthList != "" {
				err = json.Unmarshal([]byte(workflowLog.AuthList), &workflowLogAuthList)
				if err != nil {
					log.Error("[workflowLog's Auth]:error when unmarshal workflowlog's auth list:", workflowLog, " ===>error is:", err.Error())
					return errors.New("error when set auth info to workflow")
				}
			}

			workflowLogAuthList = append(workflowLogAuthList, authMap)

			authListInfo, err := json.Marshal(workflowLogAuthList)
			if err != nil {
				log.Error("[workflowLog's Auth]:error when marshal workflowlog's auth list:", workflowLogAuthList, " ===>error is:", err.Error())
				return errors.New("error when save workflow auth info")
			}

			workflowLog.AuthList = string(authListInfo)
			err = workflowLog.GetWorkflowLog().Save(workflowLog).Error
			if err != nil {
				log.Error("[workflowLog's Auth]:error when save workflowlog's info to db:", workflowLog, " ===>error is:", err.Error())
				return errors.New("error when save workflow auth info")
			}
		} else {
			remainRequireList = append(remainRequireList, requireMap)
		}
	}

	if !hasAuthed {
		log.Error("[workflowLog's Auth]:error when auth a workflowlog to start, given auth:", authMap, " is not equal to any request one:", workflowLog.Requires)
		return errors.New("illegal auth info, auth failed")
	}

	remainRequireAuthInfo, err := json.Marshal(remainRequireList)
	if err != nil {
		log.Error("[workflowLog's Auth]:error when marshal workflowlog's remainRequireAuth list:", remainRequireList, " ===>error is:", err.Error())
		return errors.New("error when sync remain require auth info")
	}

	workflowLog.Requires = string(remainRequireAuthInfo)
	err = workflowLog.GetWorkflowLog().Save(workflowLog).Error
	if err != nil {
		log.Error("[workflowLog's Auth]:error when save workflowlog's remain require auth info:", workflowLog, " ===>error is:", err.Error())
		return errors.New("error when sync remain require auth info")
	}

	return nil
}

// Start is
func (workflowLog *WorkflowLog) Start(startData string) {
	// get current workflowlog's start stage
	startStageLog := new(models.StageLog)
	err := startStageLog.GetStageLog().Where("workflow = ?", workflowLog.ID).Where("pre_stage = ?", -1).Where("type = ?", models.StageTypeStart).First(startStageLog).Error
	if err != nil {
		log.Error("[workflowLog's Start]:error when get workflowlog's start stage info from db:", err.Error())
		workflowLog.Stop(WorkflowStopReasonRunFailed, models.WorkflowLogStateRunFailed)
		return
	}

	stage := new(StageLog)
	stage.StageLog = startStageLog
	err = stage.Listen()
	if err != nil {
		log.Error("[workflowLog's Start]:error when set workflow", workflowLog, " start stage:", startStageLog, "to listen:", err.Error())
		workflowLog.Stop(WorkflowStopReasonRunFailed, models.WorkflowLogStateRunFailed)
		return
	}

	authMap := make(map[string]interface{})
	authMap["type"] = AuthTypeWorkflowStartDone
	authMap["token"] = AuthTokenDefault
	authMap["authorizer"] = "system - " + workflowLog.Namespace + " - " + workflowLog.Repository + " - " +
		workflowLog.Workflow + "(" + strconv.FormatInt(workflowLog.FromWorkflow, 10) + ")"
	authMap["time"] = time.Now().Format("2006-01-02 15:04:05")

	err = stage.Auth(authMap)
	if err != nil {
		log.Error("[workflowLog's Start]:error when auth to start stage:", workflowLog, " start stage is ", startStageLog, " ===>error is:", err.Error())
		workflowLog.Stop(WorkflowStopReasonRunFailed, models.WorkflowLogStateRunFailed)
		return
	}

	err = workflowLog.recordWorkflowStartData(startData)
	if err != nil {
		log.Error("[workflowLog's Start]:error when record workflow's start data:", startData, " ===>error is:", err.Error())
		workflowLog.Stop(WorkflowStopReasonRunFailed, models.WorkflowLogStateRunFailed)
		return
	}
}

// Stop is
func (workflowLog *WorkflowLog) Stop(reason string, runState int64) {
	err := workflowLog.GetWorkflowLog().Where("id = ?", workflowLog.ID).First(workflowLog).Error
	if err != nil {
		log.Error("[workflowLog's Stop]:error when get workflowlog info from db:", err.Error())
		return
	}

	if runState != models.WorkflowLogStateRunSuccess {
		notEndStageLogList := make([]models.StageLog, 0)
		new(models.StageLog).GetStageLog().Where("workflow = ?", workflowLog.ID).Where("run_state != ?", models.StageLogStateRunSuccess).Where("run_state != ?", models.StageLogStateRunFailed).Find(&notEndStageLogList)

		for _, stageLogInfo := range notEndStageLogList {
			stage := new(StageLog)
			stage.StageLog = &stageLogInfo
			stage.Stop(StageStopScopeAll, StageStopReasonPreStageFailed, models.StageLogStateRunFailed)
		}
	}

	workflowLog.RunState = runState
	workflowLog.FailReason = reason
	err = workflowLog.GetWorkflowLog().Save(workflowLog).Error
	if err != nil {
		log.Error("[workflowLog's Stop]:error when change workflowlog's run state:", workflowLog, " ===>error is:", err.Error())
	}

	db := models.GetDB().Begin()
	workflowInfo := new(models.Workflow)
	err = db.Raw("select * from workflow where id = ?", workflowLog.FromWorkflow).Scan(workflowInfo).Error
	if err != nil {
		log.Error("[workflowLog's Stop]:error when get workflow info from db:", err.Error())
		db.Rollback()
	}

	workflowInfo.CurrentInstance -= 1
	err = db.Save(workflowInfo).Error
	if err != nil {
		log.Error("[workflowLog's Stop]:error when save workflow info from db:", err.Error())
		db.Rollback()
	}

	db.Commit()
}

func (workflowLog *WorkflowLog) recordWorkflowStartData(startData string) error {
	startStage := new(models.StageLog)
	err := startStage.GetStageLog().Where("workflow = ?", workflowLog.ID).Where("type = ?", models.StageTypeStart).First(startStage).Error
	if err != nil {
		log.Error("[workflowLog's recordWorkflowStartData]:error when get workflow startStage info:", startData, " ===>error is:", err.Error())
		return err
	}

	err = RecordOutcom(workflowLog.ID, workflowLog.FromWorkflow, startStage.ID, startStage.FromStage, models.OutcomeTypeStageStartActionID, models.OutcomeTypeStageStartActionID, workflowLog.Sequence, models.OutcomeTypeStageStartEventID, true, startData, startData)
	if err != nil {
		log.Error("[workflowLog's recordWorkflowStartData]:error when record workflow startData info:", " ===>error is:", err.Error())
		return err
	}

	return nil
}
