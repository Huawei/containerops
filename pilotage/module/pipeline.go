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
	PipelineStopReasonTimeout = "TIME_OUT"

	PipelineStopReasonRunSuccess = "RunSuccess"
	PipelineStopReasonRunFailed  = "RunFailed"
)

var (
	startPipelineChan  chan bool
	createPipelineChan chan bool

	pipelinelogAuthChan             chan bool
	pipelinelogListenChan           chan bool
	pipelinelogSequenceGenerateChan chan bool
)

func init() {
	startPipelineChan = make(chan bool, 1)
	createPipelineChan = make(chan bool, 1)
	pipelinelogAuthChan = make(chan bool, 1)
	pipelinelogListenChan = make(chan bool, 1)
	pipelinelogSequenceGenerateChan = make(chan bool, 1)
}

type Pipeline struct {
	*models.Pipeline
}

type PipelineLog struct {
	*models.PipelineLog
}

// CreateNewPipeline is create a new pipeline with given data
func CreateNewPipeline(namespace, repository, pipelineName, pipelineVersion string) (string, error) {
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
	pipelineInfo.Repository = repository
	pipelineInfo.Pipeline = pipelineName
	pipelineInfo.Version = pipelineVersion
	pipelineInfo.VersionCode = 1

	err = pipelineInfo.GetPipeline().Save(pipelineInfo).Error
	if err != nil {
		return "", errors.New("error when save pipeline info:" + err.Error())
	}

	return "create new pipeline success", nil
}

func GetPipelineListByNamespaceAndRepository(namespace, repository string) ([]map[string]interface{}, error) {
	resultMap := make([]map[string]interface{}, 0)
	pipelineList := make([]models.Pipeline, 0)
	pipelinesMap := make(map[string]interface{}, 0)
	err := new(models.Pipeline).GetPipeline().Where("namespace = ?", namespace).Where("repository = ?", repository).Order("-updated_at").Find(&pipelineList).Error
	if err != nil && !strings.Contains(err.Error(), "record not found") {
		log.Error("[pipeline's GetPipelineListByNamespaceAndRepository]error when get pipeline list from db:" + err.Error())
		return nil, errors.New("error when get pipeline list by namespace and repository from db:" + err.Error())
	}

	for _, pipelineInfo := range pipelineList {
		if _, ok := pipelinesMap[pipelineInfo.Pipeline]; !ok {
			tempMap := make(map[string]interface{})
			tempMap["version"] = make(map[int64]interface{})
			pipelinesMap[pipelineInfo.Pipeline] = tempMap
		}

		pipelineMap := pipelinesMap[pipelineInfo.Pipeline].(map[string]interface{})
		versionMap := pipelineMap["version"].(map[int64]interface{})

		versionMap[pipelineInfo.VersionCode] = pipelineInfo
		pipelineMap["id"] = pipelineInfo.ID
		pipelineMap["name"] = pipelineInfo.Pipeline
		pipelineMap["version"] = versionMap
	}

	for _, pipeline := range pipelineList {

		pipelineInfo := pipelinesMap[pipeline.Pipeline].(map[string]interface{})
		if isSign, ok := pipelineInfo["isSign"].(bool); ok && isSign {
			continue
		}

		pipelineInfo["isSign"] = true
		pipelinesMap[pipeline.Pipeline] = pipelineInfo

		versionList := make([]map[string]interface{}, 0)
		for _, pipelineVersion := range pipelineList {
			if pipelineVersion.Pipeline == pipelineInfo["name"].(string) {
				versionMap := make(map[string]interface{})
				versionMap["id"] = pipelineVersion.ID
				versionMap["version"] = pipelineVersion.Version
				versionMap["versionCode"] = pipelineVersion.VersionCode

				latestPipelineLog := new(models.PipelineLog)
				err := latestPipelineLog.GetPipelineLog().Where("from_pipeline = ?", pipelineVersion.ID).Order("-id").First(latestPipelineLog).Error
				if err != nil {
					log.Error("[pipeline's GetPipelineListByNamespaceAndRepository]:error when get pipeline's latest run info:", err.Error())
				}

				if latestPipelineLog.ID != 0 {
					statusMap := make(map[string]interface{})

					status := false
					if latestPipelineLog.RunState != models.PipelineLogStateRunFailed {
						status = true
					}

					statusMap["time"] = latestPipelineLog.CreatedAt.Format("2006-01-02 15:04:05")
					statusMap["status"] = status

					versionMap["status"] = statusMap
				}

				versionList = append(versionList, versionMap)
			}
		}

		tempResult := make(map[string]interface{})
		tempResult["id"] = pipelineInfo["id"]
		tempResult["name"] = pipelineInfo["name"]
		tempResult["version"] = versionList

		resultMap = append(resultMap, tempResult)
	}

	return resultMap, nil
}

func GetPipelineInfo(namespace, repository, pipelineName string, pipelineId int64) (map[string]interface{}, error) {
	resultMap := make(map[string]interface{})
	pipelineInfo := new(models.Pipeline)
	err := pipelineInfo.GetPipeline().Where("id = ?", pipelineId).First(&pipelineInfo).Error
	if err != nil {
		return nil, errors.New("error when get pipeline info from db:" + err.Error())
	}

	if pipelineInfo.Namespace != namespace || pipelineInfo.Repository != repository || pipelineInfo.Pipeline != pipelineName {
		return nil, errors.New("pipeline is not equal to target pipeline")
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
	stageList, err := getDefaultStageListByPipeline(*pipelineInfo)
	if err != nil {
		return nil, err
	}
	resultMap["stageList"] = stageList
	// resultMap["stageList"] = make([]map[string]interface{}, 0)

	resultMap["lineList"] = make([]map[string]interface{}, 0)

	resultMap["status"] = false

	return resultMap, nil
}

func getDefaultStageListByPipeline(pipelineInfo models.Pipeline) ([]map[string]interface{}, error) {
	stageListMap := make([]map[string]interface{}, 0)

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

	return stageListMap, nil
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

func Run(pipelineId int64, authMap map[string]interface{}, startData string) error {
	pipelineInfo := new(models.Pipeline)
	err := pipelineInfo.GetPipeline().Where("id = ?", pipelineId).First(pipelineInfo).Error
	if err != nil {
		log.Error("[pipeline's Run]:error when get pipeline's info from db:", err.Error())
		return errors.New("error when get target pipeline info:" + err.Error())
	}
	pipeline := new(Pipeline)
	pipeline.Pipeline = pipelineInfo

	eventName, ok := authMap["eventName"].(string)
	if !ok {
		log.Error("[pipeline's Run]:error when parse eventName,want a string, got:", authMap["eventName"])
		return errors.New("error when get eventName")
	}

	eventType, ok := authMap["eventType"].(string)
	if !ok {
		log.Error("[pipeline's Run]:error when parse eventName,want a string, got:", authMap["eventType"])
		return errors.New("error when get eventType")
	}

	eventMap := make(map[string]string)
	eventMap["eventName"] = eventName
	eventMap["eventType"] = eventType

	// first generate a pipeline log to record all current pipeline's info which will be used in feature
	pipelineLog, err := pipeline.GenerateNewLog(eventMap)
	if err != nil {
		return err
	}

	// let pipeline log listen all auth, if all auth is ok, start run this pipeline log
	err = pipelineLog.Listen(startData)
	if err != nil {
		return err
	}

	// auth this pipeline log by given auth info
	err = pipelineLog.Auth(authMap)
	if err != nil {
		return err
	}

	return nil
}

func GetPipeline(pipelineId int64) (*Pipeline, error) {
	if pipelineId == int64(0) {
		return nil, errors.New("pipeline's id is empty")
	}

	pipelineInfo := new(models.Pipeline)
	err := pipelineInfo.GetPipeline().Where("id = ?", pipelineId).First(pipelineInfo).Error
	if err != nil {
		log.Error("[pipeline's GetPipeline]:error when get pipeline info from db:", err.Error())
		return nil, err
	}

	pipeline := new(Pipeline)
	pipeline.Pipeline = pipelineInfo

	return pipeline, nil
}

func GetLatestRunablePipeline(namespace, repository, pipelineName, version string) (*Pipeline, error) {
	if namespace == "" || repository == "" {
		log.Error("[pipeline's GetLatestRunablePipeline]:given empty parms:namespace: ===>", namespace, "<===  repository:===>", repository, "<===")
		return nil, errors.New("parms is empty")
	}

	pipelineInfo := new(models.Pipeline)
	query := pipelineInfo.GetPipeline().Where("namespace = ?", namespace).Where("repository = ?", repository).Where("pipeline = ?", pipelineName)

	if version != "" {
		query = query.Where("version = ?", version)
	}

	err := query.Where("state = ?", models.PipelineStateAble).Order("-id").First(&pipelineInfo).Error
	if err != nil {
		log.Error("[pipeline's GetLatestRunablePipeline]:error when get pipeline info from db:", err.Error())
		return nil, err
	}

	if pipelineInfo.ID == 0 {
		return nil, errors.New("no runable workflow")
	}

	pipeline := new(Pipeline)
	pipeline.Pipeline = pipelineInfo

	return pipeline, nil
}

func GetPipelineLog(namespace, repository, workflowName, versionName string, sequence int64) (*PipelineLog, error) {
	var err error
	pipelineLogInfo := new(models.PipelineLog)

	query := pipelineLogInfo.GetPipelineLog().Where("namespace =? ", namespace).Where("repository = ?", repository).Where("pipeline = ?", workflowName).Where("version = ?", versionName)
	if sequence == int64(0) {
		query = query.Order("-id")
	} else {
		query = query.Where("sequence = ?", sequence)
	}

	err = query.First(pipelineLogInfo).Error
	if err != nil {
		log.Error("[pipelineLog's GetPipelineLog]:error when get pipelineLog(version=", versionName, ", sequence=", sequence, ") info from db:", err.Error())
		return nil, err
	}

	pipelineLog := new(PipelineLog)
	pipelineLog.PipelineLog = pipelineLogInfo

	return pipelineLog, nil
}

func getPipelineEnvList(pipelineLogId int64) ([]map[string]interface{}, error) {
	resultList := make([]map[string]interface{}, 0)
	pipelineLog := new(models.PipelineLog)
	err := pipelineLog.GetPipelineLog().Where("id = ?", pipelineLogId).First(pipelineLog).Error
	if err != nil {
		log.Error("[pipelineLog's getPipelineEnvList]:error when get pipelinelog info from db:", err.Error())
		return nil, errors.New("error when get pipeline info from db:" + err.Error())
	}

	envMap := make(map[string]string)
	if pipelineLog.Env != "" {
		err = json.Unmarshal([]byte(pipelineLog.Env), &envMap)
		if err != nil {
			log.Error("[pipelineLog's getPipelineEnvList]:error when unmarshal pipeline's env setting:", pipelineLog.Env, " ===>error is:", err.Error())
			return nil, errors.New("error when unmarshal pipeline env info" + err.Error())
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

func GetPipelineRunHistoryList(namespace, repository string) ([]map[string]interface{}, error) {
	resultList := make([]map[string]interface{}, 0)
	pipelineLogIndexMap := make(map[string]int)
	pipelineLogVersionIndexMap := make(map[string]interface{})
	pipelineLogList := make([]models.PipelineLog, 0)

	err := new(models.PipelineLog).GetPipelineLog().Where("namespace = ?", namespace).Where("repository = ?", repository).Order("-id").Find(&pipelineLogList).Error
	if err != nil && !strings.Contains(err.Error(), "record not found") {
		log.Error("[pipeline's GetPipelineRunHistoryList]:error when get pipelineLog list from db:", err.Error())
		return nil, err
	}

	for _, pipelinelog := range pipelineLogList {
		if _, ok := pipelineLogIndexMap[pipelinelog.Pipeline]; !ok {
			pipelineInfoMap := make(map[string]interface{})
			pipelineVersionListInfoMap := make([]map[string]interface{}, 0)
			pipelineInfoMap["id"] = pipelinelog.FromPipeline
			pipelineInfoMap["name"] = pipelinelog.Pipeline
			pipelineInfoMap["versionList"] = pipelineVersionListInfoMap

			resultList = append(resultList, pipelineInfoMap)
			pipelineLogIndexMap[pipelinelog.Pipeline] = len(resultList) - 1

			versionIndexMap := make(map[string]int64)
			pipelineLogVersionIndexMap[pipelinelog.Pipeline] = versionIndexMap
		}

		pipelineInfoMap := resultList[pipelineLogIndexMap[pipelinelog.Pipeline]]
		if _, ok := pipelineLogVersionIndexMap[pipelinelog.Pipeline].(map[string]int64)[pipelinelog.Version]; !ok {
			pipelineVersionInfoMap := make(map[string]interface{})
			pipelineVersionSequenceListInfoMap := make([]map[string]interface{}, 0)
			pipelineVersionInfoMap["id"] = pipelinelog.ID
			pipelineVersionInfoMap["name"] = pipelinelog.Version
			pipelineVersionInfoMap["info"] = ""
			pipelineVersionInfoMap["total"] = int64(0)
			pipelineVersionInfoMap["success"] = int64(0)
			pipelineVersionInfoMap["sequenceList"] = pipelineVersionSequenceListInfoMap

			pipelineInfoMap["versionList"] = append(pipelineInfoMap["versionList"].([]map[string]interface{}), pipelineVersionInfoMap)
			pipelineLogVersionIndexMap[pipelinelog.Pipeline].(map[string]int64)[pipelinelog.Version] = int64(len(pipelineInfoMap["versionList"].([]map[string]interface{})) - 1)
		}

		pipelineVersionInfoMap := pipelineInfoMap["versionList"].([]map[string]interface{})[pipelineLogVersionIndexMap[pipelinelog.Pipeline].(map[string](int64))[pipelinelog.Version]]
		sequenceList := pipelineVersionInfoMap["sequenceList"].([]map[string]interface{})

		sequenceInfoMap := make(map[string]interface{})
		sequenceInfoMap["pipelineSequenceID"] = pipelinelog.ID
		sequenceInfoMap["sequence"] = pipelinelog.Sequence
		sequenceInfoMap["status"] = pipelinelog.RunState
		sequenceInfoMap["time"] = pipelinelog.CreatedAt.Format("2006-01-02 15:04:05")

		sequenceList = append(sequenceList, sequenceInfoMap)
		pipelineVersionInfoMap["sequenceList"] = sequenceList
		pipelineVersionInfoMap["total"] = pipelineVersionInfoMap["total"].(int64) + 1

		if pipelinelog.RunState == models.PipelineLogStateRunSuccess {
			pipelineVersionInfoMap["success"] = pipelineVersionInfoMap["success"].(int64) + 1
		}
	}

	for _, pipelineInfoMap := range resultList {
		for _, versionInfoMap := range pipelineInfoMap["versionList"].([]map[string]interface{}) {
			success := versionInfoMap["success"].(int64)
			total := versionInfoMap["total"].(int64)

			versionInfoMap["info"] = "Success: " + strconv.FormatInt(success, 10) + " Total: " + strconv.FormatInt(total, 10)
		}
	}

	return resultList, nil
}

func (pipelineInfo *Pipeline) CreateNewVersion(define map[string]interface{}, versionName string) error {
	var count int64
	err := new(models.Pipeline).GetPipeline().Where("namespace = ?", pipelineInfo.Namespace).Where("repository = ?", pipelineInfo.Repository).Where("pipeline = ?", pipelineInfo.Pipeline.Pipeline).Where("version = ?", versionName).Count(&count).Error
	if count > 0 {
		return errors.New("version code already exist!")
	}

	// get current least pipeline's version
	leastPipeline := new(models.Pipeline)
	err = leastPipeline.GetPipeline().Where("namespace = ? ", pipelineInfo.Namespace).Where("pipeline = ?", pipelineInfo.Pipeline.Pipeline).Order("-id").First(&leastPipeline).Error
	if err != nil {
		return errors.New("error when get least pipeline info :" + err.Error())
	}

	newPipelineInfo := new(models.Pipeline)
	newPipelineInfo.Namespace = pipelineInfo.Namespace
	newPipelineInfo.Repository = pipelineInfo.Repository
	newPipelineInfo.Pipeline = pipelineInfo.Pipeline.Pipeline
	newPipelineInfo.Event = pipelineInfo.Event
	newPipelineInfo.Version = versionName
	newPipelineInfo.VersionCode = leastPipeline.VersionCode + 1
	newPipelineInfo.State = models.PipelineStateDisable
	newPipelineInfo.Manifest = pipelineInfo.Manifest
	newPipelineInfo.Description = pipelineInfo.Description
	newPipelineInfo.SourceInfo = pipelineInfo.SourceInfo
	newPipelineInfo.Env = pipelineInfo.Env
	newPipelineInfo.Requires = pipelineInfo.Requires

	err = newPipelineInfo.GetPipeline().Save(newPipelineInfo).Error
	if err != nil {
		return err
	}

	return pipelineInfo.UpdatePipelineInfo(define)
}

func (pipelineInfo *Pipeline) GetPipelineToken() (map[string]interface{}, error) {
	result := make(map[string]interface{})

	if pipelineInfo.ID == 0 {
		log.Error("[pipeline's GetPipelineToken]:got an empty pipelin:", pipelineInfo)
		return nil, errors.New("pipeline's info is empty")
	}

	token := ""
	tokenMap := make(map[string]interface{})
	if pipelineInfo.SourceInfo == "" {
		// if sourceInfo is empty generate a token
		token = pipelineInfo.Pipeline.Pipeline
	} else {
		json.Unmarshal([]byte(pipelineInfo.SourceInfo), &tokenMap)

		if _, ok := tokenMap["token"].(string); !ok {
			token = pipelineInfo.Pipeline.Pipeline
		} else {
			token = tokenMap["token"].(string)
		}
	}

	tokenMap["token"] = token
	sourceInfo, _ := json.Marshal(tokenMap)
	pipelineInfo.SourceInfo = string(sourceInfo)
	err := pipelineInfo.GetPipeline().Save(pipelineInfo).Error

	if err != nil {
		log.Error("[pipeline's GetPipelineToken]:error when save pipeline's info to db:", err.Error())
		return nil, errors.New("error when get pipeline info from db:" + err.Error())
	}

	result["token"] = token

	url := ""

	projectAddr := ""
	if configure.GetString("projectaddr") == "" {
		projectAddr = "current-pipeline's-ip:port"
	} else {
		projectAddr = configure.GetString("projectaddr")
	}

	url += projectAddr
	url = strings.TrimSuffix(url, "/")
	url += "/v2" + "/" + pipelineInfo.Namespace + "/" + pipelineInfo.Repository + "/workflow/v1/exec/" + pipelineInfo.Pipeline.Pipeline

	result["url"] = url

	return result, nil
}

func (pipelineInfo *Pipeline) UpdatePipelineInfo(define map[string]interface{}) error {
	db := models.GetDB()
	err := db.Begin().Error
	if err != nil {
		log.Error("[pipeline's UpdatePipelineInfo]:when db.Begin():", err.Error())
		return err
	}

	pipelineOriginalManifestMap := make(map[string]interface{})
	if pipelineInfo.Manifest != "" {
		err := json.Unmarshal([]byte(pipelineInfo.Manifest), pipelineOriginalManifestMap)
		if err != nil {
			log.Error("[pipeline's UpdatePipelineInfo]:error unmarshal pipeline's manifest info:", err.Error(), " set it to empty")
			pipelineInfo.Manifest = ""
		}
	}

	pipelineOriginalManifestMap["define"] = define
	pipelineNewManifestBytes, err := json.Marshal(pipelineOriginalManifestMap)
	if err != nil {
		log.Error("[pipeline's UpdatePipelineInfo]:error when marshal pipeline's manifest info:", err.Error())
		rollbackErr := db.Rollback().Error
		if rollbackErr != nil {
			log.Error("[pipeline's UpdatePipelineInfo]:when rollback in save pipeline's info:", rollbackErr.Error())
			return errors.New("errors occur:\nerror1:" + err.Error() + "\nerror2:" + rollbackErr.Error())
		}
		return errors.New("error when save pipeline's define info:" + err.Error())
	}

	requestMap := make([]interface{}, 0)
	if request, ok := define["request"]; ok {
		if requestMap, ok = request.([]interface{}); !ok {
			log.Error("[pipeline's UpdatePipelineInfo]:error when get pipeline's request info:want a json array,got:", request)
			return errors.New("error when get pipeline's request info,want a json array")
		}
	} else {
		defaultRequestMap := make(map[string]interface{})
		defaultRequestMap["type"] = AuthTypePipelineDefault
		defaultRequestMap["token"] = AuthTokenDefault

		requestMap = append(requestMap, defaultRequestMap)
	}

	requestInfo, err := json.Marshal(requestMap)
	if err != nil {
		log.Error("[pipeline's UpdatePipelineInfo]:error when marshal pipeline's request info:", requestMap, " ===>error is:", err.Error())
		return errors.New("error when save pipeline's request info")
	}

	pipelineInfo.State = models.PipelineStateDisable
	pipelineInfo.Manifest = string(pipelineNewManifestBytes)
	pipelineInfo.Requires = string(requestInfo)
	err = db.Save(pipelineInfo).Error
	if err != nil {
		log.Error("[pipeline's UpdatePipelineInfo]:when save pipeline's info:", pipelineInfo, " ===>error is:", err.Error())
		rollbackErr := db.Rollback().Error
		if rollbackErr != nil {
			log.Error("[pipeline's UpdatePipelineInfo]:when rollback in save pipeline's info:", rollbackErr.Error())
			return errors.New("errors occur:\nerror1:" + err.Error() + "\nerror2:" + rollbackErr.Error())
		}
		return err
	}

	relationMap, stageDefineList, err := pipelineInfo.getPipelineDefineInfo(pipelineInfo.Pipeline)
	if err != nil {
		log.Error("[pipeline's UpdatePipelineInfo]:when get pipeline's define info:", err.Error())
		rollbackErr := db.Rollback().Error
		if rollbackErr != nil {
			log.Error("[pipeline's UpdatePipelineInfo]:when rollback after get pipeline define info:", rollbackErr.Error())
			return errors.New("errors occur:\nerror1:" + err.Error() + "\nerror2:" + rollbackErr.Error())
		}
		return err
	}

	// first delete old pipeline define
	err = db.Model(&models.Action{}).Where("pipeline = ?", pipelineInfo.ID).Delete(&models.Action{}).Error
	if err != nil {
		log.Error("[pipeline's UpdatePipelineInfo]:when delete action's that belong pipeline:", pipelineInfo, " ===>error is:", err.Error())
		rollbackErr := db.Rollback().Error
		if rollbackErr != nil {
			log.Error("[pipeline's UpdatePipelineInfo]:when rollback in delete action info:", rollbackErr.Error())
			return errors.New("errors occur:\nerror1:" + err.Error() + "\nerror2:" + rollbackErr.Error())
		}
		return errors.New("error when remove old action info:" + err.Error())
	}

	err = db.Model(&models.Stage{}).Where("pipeline = ?", pipelineInfo.ID).Delete(&models.Stage{}).Error
	if err != nil {
		log.Error("[pipeline's UpdatePipelineInfo]:when delete stage's that belong pipeline:", pipelineInfo, " ===>error is:", err.Error())
		rollbackErr := db.Rollback().Error
		if rollbackErr != nil {
			log.Error("[pipeline's UpdatePipelineInfo]:when rollback in delete stage info:", rollbackErr.Error())
			return errors.New("errors occur:\nerror1:" + err.Error() + "\nerror2:" + rollbackErr.Error())
		}
		return errors.New("error when update stage info:" + err.Error())
	}

	// then create new pipeline by define
	stageInfoMap := make(map[string]map[string]interface{})
	preStageId := int64(-1)
	allActionIdMap := make(map[string]int64)
	for _, stageDefine := range stageDefineList {
		stageId, stageTagId, actionMap, err := CreateNewStage(db, preStageId, pipelineInfo.Pipeline, stageDefine, relationMap)
		if err != nil {
			log.Error("[pipeline's UpdatePipelineInfo]:error when create new stage that pipeline define:", stageDefine, " preStage is :", preStageId, " pipeline is:", pipelineInfo, " relation is:", relationMap)
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
					log.Error("[pipeline's UpdatePipelineInfo]:error when get action's relation info in map:", allActionIdMap, " want :", fromActionOriginId)
					rollbackErr := db.Rollback().Error
					if rollbackErr != nil {
						log.Error("[pipeline's UpdatePipelineInfo]:when rollback in get action relation info:", rollbackErr.Error())
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
				log.Error("[pipeline's UpdatePipelineInfo]:error when get action info from db:", actionID, " ===>error is:", err.Error())
				rollbackErr := db.Rollback().Error
				if rollbackErr != nil {
					log.Error("[pipeline's UpdatePipelineInfo]:when rollback in get action info from db:", rollbackErr.Error())
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

			err = actionInfo.GetAction().Where("id = ?", actionID).UpdateColumn("manifest", actionInfo.Manifest).Error
			if err != nil {
				log.Error("[pipeline's UpdatePipelineInfo]:error when update action's column manifest:", actionInfo, " ===>error is:", err.Error())
				rollbackErr := db.Rollback().Error
				if rollbackErr != nil {
					log.Error("[pipeline's UpdatePipelineInfo]:when rollback in update action's column info from db:", rollbackErr.Error())
					return errors.New("errors occur:\nerror1:" + err.Error() + "\nerror2:" + rollbackErr.Error())
				}
				return err
			}
		}
	}

	return nil
}

func (pipeline *Pipeline) getPipelineDefineInfo(pipelineInfo *models.Pipeline) (map[string]interface{}, []map[string]interface{}, error) {
	lineList := make([]map[string]interface{}, 0)
	stageList := make([]map[string]interface{}, 0)

	manifestMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(pipelineInfo.Manifest), &manifestMap)
	if err != nil {
		log.Error("[pipeline's getPipelineDefineInfo]:error when unmarshal pipeline's manifes info:", pipeline.Manifest, " ===>error is:", err.Error())
		return nil, nil, errors.New("error when unmarshal pipeline manifes info:" + err.Error())
	}

	defineMap, ok := manifestMap["define"].(map[string]interface{})
	if !ok {
		log.Error("[pipeline's getPipelineDefineInfo]:pipeline's define is not a json obj:", manifestMap["define"])
		return nil, nil, errors.New("pipeline's define is not a json:" + err.Error())
	}

	realtionMap := make(map[string]interface{})
	if linesList, ok := defineMap["lineList"].([]interface{}); ok {
		if !ok {
			log.Error("[pipeline's getPipelineDefineInfo]:error in pipeline's lineList define,want a array,got:", defineMap["lineList"])
			return nil, nil, errors.New("pipeline's lineList define is not an array")
		}

		for _, lineInfo := range linesList {
			lineInfoMap, ok := lineInfo.(map[string]interface{})
			if !ok {
				log.Error("[pipeline's getPipelineDefineInfo]:error in pipeline's line define: want a json obj,got:", lineInfo)
				return nil, nil, errors.New("pipeline's line info is not a json")
			}

			lineList = append(lineList, lineInfoMap)
		}

		for _, lineInfo := range lineList {
			endData, ok := lineInfo["endData"].(map[string]interface{})
			if !ok {
				log.Error("[pipeline's getPipelineDefineInfo]:error in pipeline's line define:line doesn't define any end point info:", lineInfo)
				return nil, nil, errors.New("pipeline's line define is illegal,don't have a end point info")
			}

			endPointId, ok := endData["id"].(string)
			if !ok {
				log.Error("[pipeline's getPipelineDefineInfo]:error in pipeline's line define:end point's id is not a string:", endData)
				return nil, nil, errors.New("pipeline's line define is illegal,endPoint id is not a string")
			}

			if _, ok := realtionMap[endPointId]; !ok {
				realtionMap[endPointId] = make(map[string]interface{})
			}

			endPointMap := realtionMap[endPointId].(map[string]interface{})
			startData, ok := lineInfo["startData"].(map[string]interface{})
			if !ok {
				log.Error("[pipeline's getPipelineDefineInfo]:error in pipeline's line define:line doesn't define any start point info:", lineInfo)
				return nil, nil, errors.New("pipeline's line define is illegal,don;t have a start point info")
			}

			startDataId, ok := startData["id"].(string)
			if !ok {
				log.Error("[pipeline's getPipelineDefineInfo]:error in pipeline's line define:start point's id is not a string:", endData)
				return nil, nil, errors.New("pipeline's line define is illegal,startPoint id is not a string")
			}

			if startDataId == "start-stage" {
				if _, ok := endPointMap[startDataId]; !ok {
					endPointMap[startDataId] = make(map[string]interface{}, 0)
				}

				lineMap, ok := lineInfo["relation"].(map[string]interface{})
				if !ok {
					continue
				}

				lineOriginMap, ok := endPointMap[startDataId].(map[string]interface{})
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
		log.Error("[pipeline's getPipelineDefineInfo]:error in pipeline's define:pipeline doesn't define any stage info", defineMap)
		return nil, nil, errors.New("pipeline don't have a stage define")
	}

	stagesList, ok := stageListInfo.([]interface{})
	if !ok {
		log.Error("[pipeline's getPipelineDefineInfo]:error in stageList's define:want array,got:", stageListInfo)
		return nil, nil, errors.New("pipeline's stageList define is not an array")
	}

	for _, stageInfo := range stagesList {
		stageInfoMap, ok := stageInfo.(map[string]interface{})
		if !ok {
			log.Error("[pipeline's getPipelineDefineInfo]:error in stage's define,want a json obj,got:", stageInfo)
			return nil, nil, errors.New("pipeline's stage info is not a json")
		}

		stageList = append(stageList, stageInfoMap)
	}

	return realtionMap, stageList, nil
}

func (pipelineInfo *Pipeline) BeforeExecCheck(reqHeader http.Header, reqBody []byte) (bool, map[string]string, error) {
	if pipelineInfo.SourceInfo == "" {
		return false, nil, errors.New("pipeline's source info is empty")
	}

	sourceMap := make(map[string]interface{})
	sourceList := make([]interface{}, 0)
	err := json.Unmarshal([]byte(pipelineInfo.SourceInfo), &sourceMap)
	if err != nil {
		log.Error("[pipeline's BeforeExecCheck]:error when unmarshal pipeline source info, want json obj, got:", pipelineInfo.SourceInfo)
		return false, nil, errors.New("pipeline's source define error")
	}

	expectedToken, ok := sourceMap["token"].(string)
	if !ok {
		log.Error("[pipeline's BeforeExecCheck]:error when get source's expected token,want a string, got:", sourceMap["token"])
		return false, nil, errors.New("get token error")
	}

	sourceList, ok = sourceMap["sourceList"].([]interface{})
	if !ok {
		log.Error("[pipeline's BeforeExecCheck]:error when get sourceList:want json array, got:", sourceMap["sourceList"])
		return false, nil, errors.New("pipeline's sourceList define error")
	}

	eventInfoMap, err := getExecReqEventInfo(sourceList, reqHeader)
	if err != nil {
		log.Error("[pipeline's BeforeExecCheck]:error when get exec request's event type and event info:", err.Error())
		return false, nil, errors.New("get req's event info failed:" + err.Error())
	}

	passCheck := true

	checkerList, err := checker.GetWorkflowExecCheckerList()
	if err != nil {
		log.Error("[pipeline's BeforeExecCheck]:error when get checkerList:", err.Error())
		return false, nil, err
	}

	for _, checker := range checkerList {
		passCheck, err = checker.Check(eventInfoMap, expectedToken, reqHeader, reqBody)
		if !passCheck {
			log.Error("[pipeline's BeforeExecCheck]:check failed:", checker, "===>", err.Error(), "\neventInfoMap:", eventInfoMap, "\nreqHeader:", reqHeader, "\nreqBody:", reqBody)
			return false, nil, err
		}
	}

	return true, eventInfoMap, nil
}

func getExecReqEventInfo(sourceList []interface{}, reqHeader http.Header) (map[string]string, error) {
	result := make(map[string]string)
	for _, sourceConfigInfo := range sourceList {
		sourceConfig, ok := sourceConfigInfo.(map[string]interface{})
		if !ok {
			log.Error("[pipeline's getExecReqEventInfo]:error when parse sourceConfig,want a json obj, got :", sourceConfigInfo)
			return nil, errors.New("source config is not a json obj")
		}

		tokenKey, ok := sourceConfig["headerKey"].(string)
		if !ok {
			log.Error("[pipeline's getExecReqEventInfo]:error when get source's token key,want a string, got:", sourceConfig["headerKey"])
			return nil, errors.New("source's token key is not a string")
		}

		token := reqHeader.Get(tokenKey)
		if token != "" {
			supportEventList, ok := sourceConfig["eventList"].(string)
			if !ok {
				log.Error("[pipeline's getExecReqEventInfo]:error when get source's support event list,want a string, got:", sourceConfig["eventList"])
				continue
			}

			sourceType, ok := sourceConfig["sourceType"].(string)
			if !ok {
				log.Error("[pipeline's getExecReqEventInfo]:error when get source's sourceType,want a string, got:", sourceConfig["sourceType"])
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

func (pipelineInfo *Pipeline) GenerateNewLog(eventMap map[string]string) (*PipelineLog, error) {
	pipelinelogSequenceGenerateChan <- true
	result := new(PipelineLog)
	stageList := make([]models.Stage, 0)

	pipelineSequence := new(models.PipelineSequence)
	pipelineSequence.Pipeline = pipelineInfo.ID
	err := pipelineSequence.GetPipelineSequence().Save(pipelineSequence).Error
	if err != nil {
		<-pipelinelogSequenceGenerateChan
		log.Error("[pipeline's GenerateNewLog]:error when save pipeline sequence info to db", pipelineSequence, "===>error is :", err.Error())
		return nil, err
	}

	var count int64
	err = pipelineSequence.GetPipelineSequence().Where("id < ?", pipelineSequence.ID).Where("pipeline = ?", pipelineInfo.ID).Count(&count).Error
	if err != nil && !strings.Contains(err.Error(), "record not found") {
		<-pipelinelogSequenceGenerateChan
		log.Error("[pipeline's GenerateNewLog]:error when get pipeline sequence info to db:", err.Error())
		return nil, err
	}

	pipelineSequence.Sequence = count + 1
	err = pipelineSequence.GetPipelineSequence().Save(pipelineSequence).Error
	if err != nil {
		<-pipelinelogSequenceGenerateChan
		log.Error("[pipeline's GenerateNewLog]:error when save pipeline sequence info to db", pipelineSequence, "===>error is :", err.Error())
		return nil, err
	}

	err = new(models.Stage).GetStage().Where("pipeline = ?", pipelineInfo.ID).Find(&stageList).Error
	if err != nil {
		<-pipelinelogSequenceGenerateChan
		log.Error("[pipeline's GenerateNewLog]:error when get stage list by pipeline info", pipelineInfo, "===>error is :", err.Error())
		return nil, err
	}

	db := models.GetDB()
	db = db.Begin()

	eventInfoBytes, _ := json.Marshal(eventMap)

	// record pipeline's info
	pipelineLog := new(models.PipelineLog)
	pipelineLog.Namespace = pipelineInfo.Namespace
	pipelineLog.Repository = pipelineInfo.Repository
	pipelineLog.Pipeline = pipelineInfo.Pipeline.Pipeline
	pipelineLog.FromPipeline = pipelineInfo.ID
	pipelineLog.Version = pipelineInfo.Version
	pipelineLog.VersionCode = pipelineInfo.VersionCode
	pipelineLog.Sequence = pipelineSequence.Sequence
	pipelineLog.RunState = models.PipelineLogStateCanListen
	pipelineLog.Event = pipelineInfo.Event
	pipelineLog.Manifest = pipelineInfo.Manifest
	pipelineLog.Description = pipelineInfo.Description
	pipelineLog.SourceInfo = string(eventInfoBytes)
	pipelineLog.Env = pipelineInfo.Env
	pipelineLog.Requires = pipelineInfo.Requires
	pipelineLog.AuthList = ""

	err = db.Save(pipelineLog).Error
	if err != nil {
		<-pipelinelogSequenceGenerateChan
		log.Error("[pipeline's GenerateNewLog]:when save pipeline log to db:", pipelineLog, "===>error is :", err.Error())
		rollbackErr := db.Rollback().Error
		if rollbackErr != nil {
			log.Error("[pipeline's GenerateNewLog]:when rollback in save pipeline log:", rollbackErr.Error())
			return nil, errors.New("errors occur:\nerror1:" + err.Error() + "\nerror2:" + rollbackErr.Error())
		}
		return nil, err
	}

	preStageLogId := int64(-1)
	for _, stageInfo := range stageList {
		stage := new(Stage)
		stage.Stage = &stageInfo
		preStageLogId, err = stage.GenerateNewLog(db, pipelineLog, preStageLogId)
		if err != nil {
			<-pipelinelogSequenceGenerateChan
			log.Error("[pipeline's GenerateNewLog]:when generate stage log:", err.Error())
			return nil, err
		}
	}

	err = db.Commit().Error
	if err != nil {
		<-pipelinelogSequenceGenerateChan
		log.Error("[pipeline's GenerateNewLog]:when commit to db:", err.Error())
		return nil, errors.New("error when save pipeline info to db:" + err.Error())
	}
	result.PipelineLog = pipelineLog
	<-pipelinelogSequenceGenerateChan
	return result, nil
}

func (pipelineLog *PipelineLog) GetDefineInfo() (map[string]interface{}, error) {
	defineMap := make(map[string]interface{})
	stageListMap := make([]map[string]interface{}, 0)
	lineList := make([]map[string]interface{}, 0)

	stageList := make([]*models.StageLog, 0)
	err := new(models.StageLog).GetStageLog().Where("pipeline = ?", pipelineLog.ID).Find(&stageList).Error
	if err != nil {
		log.Error("[StageLog's GetStageLogDefineListByPipelineLogID]:error when get stage list from db:", err.Error())
		return nil, err
	}

	for _, stageInfo := range stageList {
		stage := new(StageLog)
		stage.StageLog = stageInfo
		stageDefineMap, err := stage.GetStageLogDefine()
		if err != nil {
			log.Error("[pipelineLog's GetDefineInfo]:error when get stagelog define:", stage, " ===>error is:", err.Error())
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
			log.Error("[pipelineLog's GetDefineInfo]:error when get actionlog list from db:", err.Error())
			continue
		}

		for _, actionInfo := range actionList {
			action := new(ActionLog)
			action.ActionLog = actionInfo
			actionLineInfo, err := action.GetActionLineInfo()
			if err != nil {
				log.Error("[pipelineLog's GetDefineInfo]:error when get actionlog line info:", err.Error())
				continue
			}

			lineList = append(lineList, actionLineInfo...)
		}
	}

	defineMap["pipeline"] = pipelineLog.Pipeline
	defineMap["version"] = pipelineLog.Version
	defineMap["sequence"] = pipelineLog.Sequence
	defineMap["status"] = pipelineLog.RunState
	defineMap["lineList"] = lineList
	defineMap["stageList"] = stageListMap

	return defineMap, nil
}

func (pipelineLog *PipelineLog) GetStartStageData() (map[string]interface{}, error) {
	dataMap := make(map[string]interface{})
	outCome := new(models.Outcome)
	err := outCome.GetOutcome().Where("pipeline = ?", pipelineLog.ID).Where("sequence = ?", pipelineLog.Sequence).Where("action = ?", models.OutcomeTypeStageStartActionID).First(outCome).Error
	if err != nil && !strings.Contains(err.Error(), "record not found") {
		log.Error("[pipelineLog's GetStartStageData]:error when get start stage data from db:", err.Error())
		return nil, err
	}

	err = json.Unmarshal([]byte(outCome.Output), &dataMap)
	if err != nil {
		log.Error("[pipelineLog's GetStartStageData]:error when unmarshal start stage's data:", outCome.Output, " ===>error is:", err.Error())
	}

	return dataMap, nil
}

func (pipelineLog *PipelineLog) Listen(startData string) error {
	pipelinelogListenChan <- true
	defer func() { <-pipelinelogListenChan }()

	err := pipelineLog.GetPipelineLog().Where("id = ?", pipelineLog.ID).First(pipelineLog).Error
	if err != nil {
		log.Error("[pipelineLog's Listen]:error when get pipelineLog info from db:", pipelineLog, " ===>error is:", err.Error())
		return errors.New("error when get pipelinelog's info from db:" + err.Error())
	}

	if pipelineLog.RunState != models.PipelineLogStateCanListen {
		log.Error("[pipelineLog's Listen]:error pipelinelog state:", pipelineLog)
		return errors.New("can't listen curren pipelinelog,current state is:" + strconv.FormatInt(pipelineLog.RunState, 10))
	}

	pipelineLog.RunState = models.PipelineLogStateWaitToStart
	err = pipelineLog.GetPipelineLog().Save(pipelineLog).Error
	if err != nil {
		log.Error("[pipelineLog's Listen]:error when change pipelinelog's run state to wait to start:", pipelineLog, " ===>error is:", err.Error())
		return errors.New("can't listen target pipeline,change pipeline's state failed")
	}

	canStartChan := make(chan bool, 1)
	go func() {
		for true {
			time.Sleep(1 * time.Second)

			err := pipelineLog.GetPipelineLog().Where("id = ?", pipelineLog.ID).First(pipelineLog).Error
			if err != nil {
				log.Error("[pipelineLog's Listen]:error when get pipelineLog's info:", pipelineLog, " ===>error is:", err.Error())
				canStartChan <- false
				break
			}
			if pipelineLog.Requires == "" || pipelineLog.Requires == "[]" {
				log.Info("[pipelineLog's Listen]:pipelineLog", pipelineLog, "is ready and will start")
				canStartChan <- true
				break
			}
		}
	}()

	go func() {
		canStart := <-canStartChan
		if !canStart {
			log.Error("[pipelineLog's Listen]:pipelineLog can't start", pipelineLog)
			pipelineLog.Stop(PipelineStopReasonRunFailed, models.PipelineLogStateRunFailed)
		} else {
			go pipelineLog.Start(startData)
		}

	}()

	return nil
}

func (pipelineLog *PipelineLog) Auth(authMap map[string]interface{}) error {
	pipelinelogAuthChan <- true
	defer func() { <-pipelinelogAuthChan }()

	authType, ok := authMap["type"].(string)
	if !ok {
		log.Error("[pipelineLog's Auth]:error when get authType from given authMap:", authMap, " ===>to pipelinelog:", pipelineLog)
		return errors.New("authType is illegal")
	}

	token, ok := authMap["token"].(string)
	if !ok {
		log.Error("[pipelineLog's Auth]:error when get token from given authMap:", authMap, " ===>to pipelinelog:", pipelineLog)
		return errors.New("token is illegal")
	}

	err := pipelineLog.GetPipelineLog().Where("id = ?", pipelineLog.ID).First(pipelineLog).Error
	if err != nil {
		log.Error("[pipelineLog's Auth]:error when get pipelineLog info from db:", pipelineLog, " ===>error is:", err.Error())
		return errors.New("error when get pipelinelog's info from db:" + err.Error())
	}

	if pipelineLog.Requires == "" || pipelineLog.Requires == "[]" {
		log.Error("[pipelineLog's Auth]:error when set auth info,pipelinelog's requires is empty", authMap, " ===>to pipelinelog:", pipelineLog)
		return errors.New("pipeline don't need any more auth")
	}

	requireList := make([]interface{}, 0)
	remainRequireList := make([]interface{}, 0)
	err = json.Unmarshal([]byte(pipelineLog.Requires), &requireList)
	if err != nil {
		log.Error("[pipelineLog's Auth]:error when unmarshal pipelinelog's require list:", pipelineLog, " ===>error is:", err.Error())
		return errors.New("error when get pipeline require auth info:" + err.Error())
	}

	hasAuthed := false
	for _, require := range requireList {
		requireMap, ok := require.(map[string]interface{})
		if !ok {
			log.Error("[pipelineLog's Auth]:error when get pipelinelog's require info:", pipelineLog, " ===> require is:", require)
			return errors.New("error when get pipeline require auth info,require is not a json object")
		}

		requireType, ok := requireMap["type"].(string)
		if !ok {
			log.Error("[pipelineLog's Auth]:error when get pipelinelog's require type:", pipelineLog, " ===> require map is:", requireMap)
			return errors.New("error when get pipeline require auth info,require don't have a type")
		}

		requireToken, ok := requireMap["token"].(string)
		if !ok {
			log.Error("[pipelineLog's Auth]:error when get pipelinelog's require token:", pipelineLog, " ===> require map is:", requireMap)
			return errors.New("error when get pipeline require auth info,require don't have a token")
		}

		if requireType == authType && requireToken == token {
			hasAuthed = true
			// record auth info to pipelinelog's auth info list
			pipelineLogAuthList := make([]interface{}, 0)
			if pipelineLog.AuthList != "" {
				err = json.Unmarshal([]byte(pipelineLog.AuthList), &pipelineLogAuthList)
				if err != nil {
					log.Error("[pipelineLog's Auth]:error when unmarshal pipelinelog's auth list:", pipelineLog, " ===>error is:", err.Error())
					return errors.New("error when set auth info to pipeline")
				}
			}

			pipelineLogAuthList = append(pipelineLogAuthList, authMap)

			authListInfo, err := json.Marshal(pipelineLogAuthList)
			if err != nil {
				log.Error("[pipelineLog's Auth]:error when marshal pipelinelog's auth list:", pipelineLogAuthList, " ===>error is:", err.Error())
				return errors.New("error when save pipeline auth info")
			}

			pipelineLog.AuthList = string(authListInfo)
			err = pipelineLog.GetPipelineLog().Save(pipelineLog).Error
			if err != nil {
				log.Error("[pipelineLog's Auth]:error when save pipelinelog's info to db:", pipelineLog, " ===>error is:", err.Error())
				return errors.New("error when save pipeline auth info")
			}
		} else {
			remainRequireList = append(remainRequireList, requireMap)
		}
	}

	if !hasAuthed {
		log.Error("[pipelineLog's Auth]:error when auth a pipelinelog to start, given auth:", authMap, " is not equal to any request one:", pipelineLog.Requires)
		return errors.New("illegal auth info, auth failed")
	}

	remainRequireAuthInfo, err := json.Marshal(remainRequireList)
	if err != nil {
		log.Error("[pipelineLog's Auth]:error when marshal pipelinelog's remainRequireAuth list:", remainRequireList, " ===>error is:", err.Error())
		return errors.New("error when sync remain require auth info")
	}

	pipelineLog.Requires = string(remainRequireAuthInfo)
	err = pipelineLog.GetPipelineLog().Save(pipelineLog).Error
	if err != nil {
		log.Error("[pipelineLog's Auth]:error when save pipelinelog's remain require auth info:", pipelineLog, " ===>error is:", err.Error())
		return errors.New("error when sync remain require auth info")
	}

	return nil
}

func (pipelineLog *PipelineLog) Start(startData string) {
	// get current pipelinelog's start stage
	startStageLog := new(models.StageLog)
	err := startStageLog.GetStageLog().Where("pipeline = ?", pipelineLog.ID).Where("pre_stage = ?", -1).Where("type = ?", models.StageTypeStart).First(startStageLog).Error
	if err != nil {
		log.Error("[pipelineLog's Start]:error when get pipelinelog's start stage info from db:", err.Error())
		pipelineLog.Stop(PipelineStopReasonRunFailed, models.PipelineLogStateRunFailed)
		return
	}

	stage := new(StageLog)
	stage.StageLog = startStageLog
	err = stage.Listen()
	if err != nil {
		log.Error("[pipelineLog's Start]:error when set pipeline", pipelineLog, " start stage:", startStageLog, "to listen:", err.Error())
		pipelineLog.Stop(PipelineStopReasonRunFailed, models.PipelineLogStateRunFailed)
		return
	}

	authMap := make(map[string]interface{})
	authMap["type"] = AuthTypePipelineStartDone
	authMap["token"] = AuthTokenDefault
	authMap["authorizer"] = "system - " + pipelineLog.Namespace + " - " + pipelineLog.Repository + " - " +
		pipelineLog.Pipeline + "(" + strconv.FormatInt(pipelineLog.FromPipeline, 10) + ")"
	authMap["time"] = time.Now().Format("2006-01-02 15:04:05")

	err = stage.Auth(authMap)
	if err != nil {
		log.Error("[pipelineLog's Start]:error when auth to start stage:", pipelineLog, " start stage is ", startStageLog, " ===>error is:", err.Error())
		pipelineLog.Stop(PipelineStopReasonRunFailed, models.PipelineLogStateRunFailed)
		return
	}

	err = pipelineLog.recordPipelineStartData(startData)
	if err != nil {
		log.Error("[pipelineLog's Start]:error when record pipeline's start data:", startData, " ===>error is:", err.Error())
		pipelineLog.Stop(PipelineStopReasonRunFailed, models.PipelineLogStateRunFailed)
		return
	}
}

func (pipelineLog *PipelineLog) Stop(reason string, runState int64) {
	err := pipelineLog.GetPipelineLog().Where("id = ?", pipelineLog.ID).First(pipelineLog).Error
	if err != nil {
		log.Error("[pipelineLog's Stop]:error when get pipelinelog info from db:", err.Error())
		return
	}

	notEndStageLogList := make([]models.StageLog, 0)
	new(models.StageLog).GetStageLog().Where("pipeline = ?", pipelineLog.ID).Where("run_state != ?", models.StageLogStateRunSuccess).Where("run_state != ?", models.StageLogStateRunFailed).Find(&notEndStageLogList)

	for _, stageLogInfo := range notEndStageLogList {
		stage := new(StageLog)
		stage.StageLog = &stageLogInfo
		stage.Stop(StageStopScopeAll, StageStopReasonRunFailed, models.StageLogStateRunFailed)
	}

	pipelineLog.RunState = runState
	err = pipelineLog.GetPipelineLog().Save(pipelineLog).Error
	if err != nil {
		log.Error("[pipelineLog's Stop]:error when change pipelinelog's run state:", pipelineLog, " ===>error is:", err.Error())
	}
}

func (pipelineLog *PipelineLog) recordPipelineStartData(startData string) error {
	startStage := new(models.StageLog)
	err := startStage.GetStageLog().Where("pipeline = ?", pipelineLog.ID).Where("type = ?", models.StageTypeStart).First(startStage).Error
	if err != nil {
		log.Error("[pipelineLog's recordPipelineStartData]:error when get pipeline startStage info:", startData, " ===>error is:", err.Error())
		return err
	}

	err = RecordOutcom(pipelineLog.ID, pipelineLog.FromPipeline, startStage.ID, startStage.FromStage, models.OutcomeTypeStageStartActionID, models.OutcomeTypeStageStartActionID, pipelineLog.Sequence, models.OutcomeTypeStageStartEventID, true, startData, startData)
	if err != nil {
		log.Error("[pipelineLog's recordPipelineStartData]:error when record pipeline startData info:", " ===>error is:", err.Error())
		return err
	}

	return nil
}
