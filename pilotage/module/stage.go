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
	"strconv"
	"time"

	"github.com/Huawei/containerops/pilotage/models"

	log "github.com/Sirupsen/logrus"
	"github.com/jinzhu/gorm"
)

const (
	StageStopReasonTimeout = "TIME_OUT"

	StageStopReasonRunSuccess = "RunSuccess"
	StageStopReasonRunFailed  = "RunFailed"

	StageStopScopeAll        = "all"
	StageStopScopeRecyclable = "recyclable"

	PipelineStageTypeStart = "pipeline-start"
	PipelineStageTypeRun   = "pipeline-stage"
	PipelineStageTypeAdd   = "pipeline-add-stage"
	PipelineStageTypeEnd   = "pipeline-end"
)

var (
	stagelogAuthChan   chan bool
	stagelogListenChan chan bool
)

type Stage struct {
	*models.Stage
}

type StageLog struct {
	*models.StageLog
}

func init() {
	stagelogAuthChan = make(chan bool, 1)
	stagelogListenChan = make(chan bool, 1)
}

func getStageEnvList(stageLogId int64) ([]map[string]interface{}, error) {
	resultList := make([]map[string]interface{}, 0)
	stageLog := new(models.StageLog)
	err := stageLog.GetStageLog().Where("id = ?", stageLogId).First(stageLog).Error
	if err != nil {
		log.Error("[stageLog's getStageEnvList]:error when get stageLog info from db:", err.Error())
		return nil, errors.New("error when get stage info from db:" + err.Error())
	}

	envMap := make(map[string]string)
	if stageLog.Env != "" {
		err = json.Unmarshal([]byte(stageLog.Env), &envMap)
		if err != nil {
			log.Error("[stageLog's getStageEnvList]:error when unmarshal stage's env setting:", stageLog.Env, " ===>error is:", err.Error())
			return nil, errors.New("error when unmarshal stage's env info" + err.Error())
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

func CreateNewStage(db *gorm.DB, preStageId int64, pipelineInfo *models.Pipeline, defineMap, relationMap map[string]interface{}) (int64, string, map[string]int64, error) {
	if db == nil {
		db = models.GetDB()
		err := db.Begin().Error
		if err != nil {
			log.Error("[stage's CreateNewStage]:when db.Begin():", err.Error())
			return 0, "", nil, err
		}
	}

	stageType := models.StageTypeRun
	actionIdMap := make(map[string]int64)
	stageName := ""
	timeout := int64(60 * 60 * 24 * 36)
	requestMapList := make([]interface{}, 0)
	authType := AuthTypePreStageDone

	idStr, ok := defineMap["id"].(string)
	if !ok {
		log.Error("[stage's CreateNewStage]:error in stage's define:want a string id, in define is:", defineMap)
		return 0, "", nil, errors.New("stage define does not have a string id")
	}

	stageDefineType, ok := defineMap["type"].(string)
	if !ok {
		log.Error("[stage's CreateNewStage]:error in stage's define:stage's type is not a string,in define is:", defineMap)
		return 0, "", nil, errors.New("stage type define is not a string")
	}

	if stageDefineType == PipelineStageTypeAdd {
		return 0, "", nil, nil
	} else if stageDefineType == PipelineStageTypeStart {
		authType = AuthTypePipelineStartDone
		stageType = models.StageTypeStart
		stageName = pipelineInfo.Pipeline + "-start-stage"
		timeout = 0

		if sourceMapList, ok := defineMap["outputJson"].([]interface{}); ok {
			sourceMap := make(map[string]interface{}, 0)
			json.Unmarshal([]byte(pipelineInfo.SourceInfo), &sourceMap)

			sourceList := make([]map[string]interface{}, 0)
			allSourceMap := make(map[string]interface{})
			for _, sourceInfo := range sourceMapList {
				sourceInfoMap, ok := sourceInfo.(map[string]interface{})
				if !ok {
					continue
				}

				sourceType, ok := sourceInfoMap["type"].(string)
				if !ok {
					continue
				}
				eventType, ok := sourceInfoMap["event"].(string)
				if !ok {
					continue
				}

				if _, ok := allSourceMap[sourceType].(map[string]bool); !ok {
					sourceEventMap := make(map[string]bool)
					allSourceMap[sourceType] = sourceEventMap
				}

				sourceEventMap := allSourceMap[sourceType].(map[string]bool)

				if exist, ok := sourceEventMap[eventType]; !ok || !exist {
					sourceEventMap[eventType] = true
				}
			}

			for sourceType, sourceEventMap := range allSourceMap {
				sourceEventMapList := sourceEventMap.(map[string]bool)
				sourceTypeKey := ""
				switch sourceType {
				case "github":
					sourceTypeKey = "X-Hub-Signature"
				case "customize":
					sourceTypeKey = "X-Workflow-Signature"
				case "gitlab":
					sourceTypeKey = "X-Gitlab-Token"
				}

				eventListStr := ","
				for eventName, _ := range sourceEventMapList {
					eventListStr += eventName + ","
				}

				tempSourceMap := make(map[string]interface{})
				tempSourceMap["sourceType"] = sourceType
				tempSourceMap["headerKey"] = sourceTypeKey
				tempSourceMap["eventList"] = eventListStr

				sourceList = append(sourceList, tempSourceMap)
			}

			sourceMap["sourceList"] = sourceList
			sourceMapBytes, _ := json.Marshal(sourceMap)
			pipelineInfo.SourceInfo = string(sourceMapBytes)
			err := db.Model(&models.Pipeline{}).Save(pipelineInfo).Error
			if err != nil {
				log.Error("[stage's CreateNewStage]:error when update pipeline's source info to db:", err.Error())
				rollbackErr := db.Rollback().Error
				if rollbackErr != nil {
					log.Error("[stage's CreateNewStage]:when rollback in update pipeline's source info:", rollbackErr.Error())
					return 0, "", nil, errors.New("errors occur:\nerror1:" + err.Error() + "\nerror2:" + rollbackErr.Error())
				}
				return 0, "", nil, err
			}

		}
	} else if stageDefineType == PipelineStageTypeEnd {
		stageType = models.StageTypeEnd
		stageName = pipelineInfo.Pipeline + "-end-stage"
		timeout = 0
	} else if stageDefineType == PipelineStageTypeRun {
		if setupDataMap, ok := defineMap["setupData"].(map[string]interface{}); ok {
			if defineName, ok := setupDataMap["name"]; ok {
				defineNameStr, ok := defineName.(string)
				if !ok {
					defineNameStr = ""
				}
				stageName = defineNameStr
			}

			defineTimeoutStr, ok := setupDataMap["timeout"].(string)
			if ok {
				var err error
				timeout, err = strconv.ParseInt(defineTimeoutStr, 10, 64)
				if err != nil {
					timeout = 0
				}
			}
		}
	} else {
		log.Error("[stage's CreateNewStage]:got an unknow stage type:", stageDefineType)
		return 0, "", nil, nil
	}

	stageRequest, ok := defineMap["request"].([]interface{})
	if !ok {
		defaultRequestMap := make(map[string]interface{})
		defaultRequestMap["type"] = authType
		defaultRequestMap["token"] = AuthTokenDefault

		requestMapList = append(requestMapList, defaultRequestMap)
	} else {
		requestMapList = stageRequest
	}
	requestInfos, _ := json.Marshal(requestMapList)

	stage := new(models.Stage)
	stage.Namespace = pipelineInfo.Namespace
	stage.Repository = pipelineInfo.Repository
	stage.Pipeline = pipelineInfo.ID
	stage.Type = int64(stageType)
	stage.PreStage = preStageId
	stage.Stage = stageName
	stage.Title = stageName
	stage.Description = stageName
	stage.Timeout = timeout
	stage.Requires = string(requestInfos)

	err := db.Model(&models.Stage{}).Save(stage).Error
	if err != nil {
		log.Error("[stage's CreateNewStage]:error when save stage info to db:", err.Error())
		rollbackErr := db.Rollback().Error
		if rollbackErr != nil {
			log.Error("[stage's CreateNewStage]:when rollback in save satge's info:", rollbackErr.Error())
			return 0, "", nil, errors.New("errors occur:\nerror1:" + err.Error() + "\nerror2:" + rollbackErr.Error())
		}
		return 0, "", nil, err
	}

	if stageDefineType == PipelineStageTypeStart {
		actionIdMap[idStr] = 0
	}

	if actionDefine, ok := defineMap["actions"]; ok {
		actionList, ok := actionDefine.([]interface{})
		if !ok {
			log.Error("[stage's CreateNewStage]:error when get action's define list,want array, got:", actionDefine)
			rollbackErr := db.Rollback().Error
			if rollbackErr != nil {
				log.Error("[stage's CreateNewStage]:when rollback when get action define list:", rollbackErr.Error())
				return 0, "", nil, errors.New("errors occur:\nerror1:" + err.Error() + "\nerror2:" + rollbackErr.Error())
			}
			return 0, "", nil, errors.New("action list is not an array")
		}

		actionDefineList := make([]map[string]interface{}, 0)
		for _, action := range actionList {
			actionDefineMap, ok := action.(map[string]interface{})
			if !ok {
				log.Error("[stage's CreateNewStage]:error when get action's define info,want a json obj, got:", action)
				rollbackErr := db.Rollback().Error
				if rollbackErr != nil {
					log.Error("[stage's CreateNewStage]:when rollback when get action define list:", rollbackErr.Error())
					return 0, "", nil, errors.New("errors occur:\nerror1:" + err.Error() + "\nerror2:" + rollbackErr.Error())
				}
				return 0, "", nil, errors.New("action's define is not a json")
			}
			actionDefineList = append(actionDefineList, actionDefineMap)
		}

		actionIdMap, err = CreateNewActions(db, pipelineInfo, stage, actionDefineList)
		// actionIdMap, err = createActionByDefine(actionDefineList, stage.ID)
		if err != nil {
			log.Error("[stage's CreateNewStage]:error when create actions by defineList:", actionDefineList, " ===>error is:", err.Error())
			return 0, "", nil, err
		}
	}

	return stage.ID, idStr, actionIdMap, err
}

func GetStageLogByName(namespace, repository, pipelineName string, sequence int64, stageName string) (*StageLog, error) {
	stage := new(StageLog)
	pipelineLog := new(models.PipelineLog)
	stageLog := new(models.StageLog)

	err := pipelineLog.GetPipelineLog().Where("namespace = ?", namespace).Where("repository = ?", repository).Where("pipeline = ?", pipelineName).Where("sequence = ?", sequence).First(pipelineLog).Error
	if err != nil {
		if err != nil {
			log.Error("[stageLog's GetStageLogByName]:error when get pipelineLog info from db:", err.Error())
			return nil, err
		}
	}

	err = stageLog.GetStageLog().Where("namespace = ?", namespace).Where("repository = ?", repository).Where("pipeline = ?", pipelineLog.ID).Where("sequence = ?", sequence).Where("stage = ?", stageName).First(stageLog).Error
	if err != nil {
		if err != nil {
			log.Error("[stageLog's GetStageLogByName]:error when get stageLog info from db:", err.Error())
			return nil, err
		}
	}

	stage.StageLog = stageLog

	return stage, nil
}

func (stageInfo *Stage) GenerateNewLog(db *gorm.DB, pipelineLog *models.PipelineLog, preStageLogID int64) (int64, error) {
	actionList := make([]models.Action, 0)

	err := new(models.Action).GetAction().Where("pipeline = ?", pipelineLog.FromPipeline).Where("stage = ?", stageInfo.ID).Find(&actionList).Error
	if err != nil {
		log.Error("[Stage's GenerateNewLog]:when get action list by stage info", stageInfo, "===>error is :", err.Error())
		return 0, err
	}

	if db == nil {
		db = models.GetDB()
		db = db.Begin()
	}

	// record stage's info
	stageLog := new(models.StageLog)
	stageLog.Namespace = stageInfo.Namespace
	stageLog.Repository = stageInfo.Repository
	stageLog.Pipeline = pipelineLog.ID
	stageLog.FromPipeline = pipelineLog.FromPipeline
	stageLog.Sequence = pipelineLog.Sequence
	stageLog.FromStage = stageInfo.ID
	stageLog.Type = stageInfo.Type
	stageLog.PreStage = preStageLogID
	stageLog.Stage = stageInfo.Stage.Stage
	stageLog.Title = stageInfo.Title
	stageLog.Description = stageInfo.Description
	stageLog.RunState = models.StageLogStateCanListen
	stageLog.Event = stageInfo.Event
	stageLog.Manifest = stageInfo.Manifest
	stageLog.Env = stageInfo.Env
	stageLog.Timeout = stageInfo.Timeout
	stageLog.Requires = stageInfo.Requires
	stageLog.AuthList = ""

	err = db.Save(stageLog).Error
	if err != nil {
		log.Error("[stage's GenerateNewLog]:when save stage log to db:", stageLog, " ===>error is:", err.Error())
		rollbackErr := db.Rollback().Error
		if rollbackErr != nil {
			log.Error("[stage's GenerateNewLog]:when rollback in save stage log:", rollbackErr.Error())
			return 0, errors.New("errors occur:\nerror1:" + err.Error() + "\nerror2:" + rollbackErr.Error())
		}
		return 0, err
	}

	for _, actionInfo := range actionList {
		action := new(Action)
		action.Action = &actionInfo
		err = action.GenerateNewLog(db, pipelineLog, stageLog)
		if err != nil {
			log.Error("[stage's GenerateNewLog]:when generate action log:", err.Error())
			return 0, err
		}
	}

	return stageLog.ID, nil
}

func (stageLog *StageLog) GetStageLogDefine() (map[string]interface{}, error) {
	actionList := make([]models.ActionLog, 0)
	err := new(models.ActionLog).GetActionLog().Where("stage = ?", stageLog.ID).Find(&actionList).Error
	if err != nil {
		log.Error("[StageLog's GetStageLogDefineListByPipelineLogID]:error when get action list from db:", err.Error())
	}

	stageInfoMap := make(map[string]interface{})
	stageInfoMap["id"] = "s-" + strconv.FormatInt(stageLog.ID, 10)
	stageInfoMap["name"] = stageLog.Stage
	stageInfoMap["setupData"] = map[string]string{"name": stageLog.Stage}
	stageInfoMap["status"] = stageLog.RunState
	stageInfoMap["type"] = models.StageTypeForWeb[stageLog.Type]
	stageInfoMap["runTime"] = stageLog.CreatedAt.Format("2006-01-02 15:04:05") + " - "

	endTimeStr := ""
	if stageLog.RunState == models.StageLogStateRunFailed || stageLog.RunState == models.StageLogStateRunSuccess {
		endTimeStr = stageLog.UpdatedAt.Format("2006-01-02 15:04:05")
	}
	stageInfoMap["runTime"] = stageInfoMap["runTime"].(string) + endTimeStr

	if len(actionList) > 0 {
		actionListMap := make([]map[string]interface{}, 0)
		for _, actionInfo := range actionList {
			tempActionInfoMap := make(map[string]interface{})
			tempActionInfoMap["id"] = "a-" + strconv.FormatInt(actionInfo.ID, 10)
			tempActionInfoMap["setupData"] = map[string]string{"name": actionInfo.Action}
			tempActionInfoMap["status"] = actionInfo.RunState
			tempActionInfoMap["type"] = "pipeline-action"

			actionListMap = append(actionListMap, tempActionInfoMap)
		}

		stageInfoMap["actions"] = actionListMap
	}

	return stageInfoMap, nil
}

func (stageLog *StageLog) Listen() error {
	stagelogListenChan <- true
	defer func() { <-stagelogListenChan }()

	err := stageLog.GetStageLog().Where("id = ?", stageLog.ID).First(stageLog).Error
	if err != nil {
		log.Error("[stageLog's Listen]:error when get stage info from db:", stageLog, " ===>error is:", err.Error())
		return errors.New("error when get stagelog's info from db:" + err.Error())
	}

	if stageLog.RunState != models.StageLogStateCanListen {
		log.Error("[stageLog's Listen]:error stagelog state:", stageLog)
		return errors.New("can't listen curren stagelog,current state is:" + strconv.FormatInt(stageLog.RunState, 10))
	}

	stageLog.RunState = models.StageLogStateWaitToStart
	err = stageLog.GetStageLog().Save(stageLog).Error
	if err != nil {
		log.Error("[stageLog's Listen]:error when change stageLog's run state to wait to start:", stageLog, " ===>error is:", err.Error())
		return errors.New("can't listen target stage,change stage's state failed")
	}

	canStartChan := make(chan bool, 1)
	go func() {
		for true {
			time.Sleep(1 * time.Second)

			err := stageLog.GetStageLog().Where("id = ?", stageLog.ID).First(stageLog).Error
			if err != nil {
				log.Error("[stageLog's Listen]:error when get stageLog's info:", stageLog, " ===>error is:", err.Error())
				canStartChan <- false
				break
			}
			if stageLog.Requires == "" || stageLog.Requires == "[]" {
				log.Info("[stageLog's Listen]:stageLog", stageLog, " is ready and will start")
				canStartChan <- true
				break
			}
		}
	}()

	go func() {
		canStart := <-canStartChan
		if !canStart {
			log.Error("[stageLog's Listen]:stageLog can't start", stageLog)
			stageLog.Stop(StageStopScopeAll, StageStopReasonRunFailed, models.StageLogStateRunFailed)
			return
		}

		log.Info("[stageLog's Listen]: will start stage", stageLog)
		go stageLog.Start()
	}()

	return nil
}

func (stageLog *StageLog) Auth(authMap map[string]interface{}) error {
	stagelogAuthChan <- true
	defer func() { <-stagelogAuthChan }()

	authType, ok := authMap["type"].(string)
	if !ok {
		log.Error("[stageLog's Auth]:error when get authType from given authMap:", authMap, " ===>to stagelog:", stageLog)
		return errors.New("authType is illegal")
	}

	token, ok := authMap["token"].(string)
	if !ok {
		log.Error("[stageLog's Auth]:error when get token from given authMap:", authMap, " ===>to stagelog:", stageLog)
		return errors.New("token is illegal")
	}

	err := stageLog.GetStageLog().Where("id = ?", stageLog.ID).First(stageLog).Error
	if err != nil {
		log.Error("[stageLog's Auth]:error when get stageLog info from db:", stageLog, " ===>error is:", err.Error())
		return errors.New("error when get stagelog's info from db:" + err.Error())
	}

	if stageLog.Requires == "" || stageLog.Requires == "[]" {
		log.Error("[stageLog's Auth]:error when set auth info,stagelog's requires is empty", authMap, " ===>to stagelog:", stageLog)
		return errors.New("stage don't need any more auth")
	}

	requireList := make([]interface{}, 0)
	remainRequireList := make([]interface{}, 0)
	err = json.Unmarshal([]byte(stageLog.Requires), &requireList)
	if err != nil {
		log.Error("[stageLog's Auth]:error when unmarshal stagelog's require list:", stageLog, " ===>error is:", err.Error())
		return errors.New("error when get stage require auth info:" + err.Error())
	}

	hasAuthed := false
	for _, require := range requireList {
		requireMap, ok := require.(map[string]interface{})
		if !ok {
			log.Error("[stageLog's Auth]:error when get stagelog's require info:", stageLog, " ===> require is:", require)
			return errors.New("error when get stage require auth info,require is not a json object")
		}

		requireType, ok := requireMap["type"].(string)
		if !ok {
			log.Error("[stageLog's Auth]:error when get stageLog's require type:", stageLog, " ===> require map is:", requireMap)
			return errors.New("error when get stage require auth info,require don't have a type")
		}

		requireToken, ok := requireMap["token"].(string)
		if !ok {
			log.Error("[stageLog's Auth]:error when get stageLog's require token:", stageLog, " ===> require map is:", requireMap)
			return errors.New("error when get stage require auth info,require don't have a token")
		}

		if requireType == authType && requireToken == token {
			hasAuthed = true
			// record auth info to stagelog's auth info list
			stageLogAuthList := make([]interface{}, 0)
			if stageLog.AuthList != "" {
				err = json.Unmarshal([]byte(stageLog.AuthList), &stageLogAuthList)
				if err != nil {
					log.Error("[stageLog's Auth]:error when unmarshal stagelog's auth list:", stageLog, " ===>error is:", err.Error())
					return errors.New("error when set auth info to stage")
				}
			}

			stageLogAuthList = append(stageLogAuthList, authMap)

			authListInfo, err := json.Marshal(stageLogAuthList)
			if err != nil {
				log.Error("[stageLog's Auth]:error when marshal stagelog's auth list:", stageLogAuthList, " ===>error is:", err.Error())
				return errors.New("error when save stage auth info")
			}

			stageLog.AuthList = string(authListInfo)
			err = stageLog.GetStageLog().Save(stageLog).Error
			if err != nil {
				log.Error("[stageLog's Auth]:error when save stageLog's info to db:", stageLog, " ===>error is:", err.Error())
				return errors.New("error when save stage auth info")
			}
		} else {
			remainRequireList = append(remainRequireList, requireMap)
		}
	}

	if !hasAuthed {
		log.Error("[stageLog's Auth]:error when auth a stagelog to start, given auth:", authMap, " is not equal to any request one:", stageLog.Requires)
		return errors.New("illegal auth info, auth failed")
	}

	remainRequireAuthInfo, err := json.Marshal(remainRequireList)
	if err != nil {
		log.Error("[stageLog's Auth]:error when marshal stageLog's remainRequireAuth list:", remainRequireList, " ===>error is:", err.Error())
		return errors.New("error when sync remain require auth info")
	}

	stageLog.Requires = string(remainRequireAuthInfo)
	err = stageLog.GetStageLog().Save(stageLog).Error
	if err != nil {
		log.Error("[stageLog's Auth]:error when save stageLog's remain require auth info:", stageLog, " ===>error is:", err.Error())
		return errors.New("error when sync remain require auth info")
	}

	return nil
}

func (stageLog *StageLog) Start() {
	nextStageCanStartChan := make(chan bool, 1)
	go stageLog.WaitAllActionDone(nextStageCanStartChan)

	if stageLog.Type == models.StageTypeEnd {
		log.Info("[stageLog's Start]:start end stage,all stage run success,so pipeline is success")
		pipelineLogInfo := new(models.PipelineLog)
		err := pipelineLogInfo.GetPipelineLog().Where("id = ?", stageLog.Pipeline).First(pipelineLogInfo).Error
		if err != nil {
			log.Error("[stageLog's Start]:error when get pipelinelog's info from db:", err.Error())
			stageLog.Stop(StageStopScopeAll, StageStopReasonRunFailed, models.StageLogStateRunFailed)
			return
		}

		stageLog.Stop(StageStopScopeAll, StageStopReasonRunSuccess, models.StageLogStateRunSuccess)
		pipelineLog := new(PipelineLog)
		pipelineLog.PipelineLog = pipelineLogInfo
		pipelineLog.Stop(PipelineStopReasonRunSuccess, models.PipelineLogStateRunSuccess)
		return
	}

	if stageLog.Type == models.StageTypeStart {
		log.Info("[stageLog's Start]:start pipeline's start stage ...")
		nextStageCanStartChan <- true
	}

	actionLogList := make([]models.ActionLog, 0)
	err := new(models.ActionLog).GetActionLog().Where("pipeline = ?", stageLog.Pipeline).Where("stage = ?", stageLog.ID).Find(&actionLogList).Error
	if err != nil {
		log.Error("[stageLog's Start]:error when get actionLog list fron db:", err.Error())
		stageLog.Stop(StageStopScopeAll, StageStopReasonRunFailed, models.StageLogStateRunFailed)
		return
	}

	go func() {
		canStart := <-nextStageCanStartChan
		if canStart {
			stageLog.Stop(StageStopScopeRecyclable, StageStopReasonRunSuccess, models.StageLogStateRunSuccess)
			nextStageLogInfo := new(models.StageLog)
			err := nextStageLogInfo.GetStageLog().Where("pipeline = ?", stageLog.Pipeline).Where("pre_stage = ?", stageLog.ID).First(nextStageLogInfo).Error
			if err != nil {
				log.Error("[stageLog's Start]:error when get next stageLog info from db:", err.Error())
				stageLog.Stop(StageStopScopeAll, StageStopReasonRunFailed, models.StageLogStateRunFailed)
				return
			}

			nextStage := new(StageLog)
			nextStage.StageLog = nextStageLogInfo

			err = nextStage.Listen()
			if err != nil {
				log.Error("[stageLog's Start]:error when set stage", nextStageLogInfo, "to listen:", err.Error())
				stageLog.Stop(StageStopScopeAll, StageStopReasonRunFailed, models.StageLogStateRunFailed)
				return
			}

			authMap := make(map[string]interface{})
			authMap["type"] = AuthTypePreStageDone
			authMap["token"] = AuthTokenDefault
			authMap["authorizer"] = "system - " + stageLog.Namespace + " - " + stageLog.Repository + " - " +
				strconv.FormatInt(stageLog.Pipeline, 10) + "(" + strconv.FormatInt(stageLog.FromPipeline, 10) + ") - " +
				strconv.FormatInt(stageLog.ID, 10) + "(" + strconv.FormatInt(stageLog.FromStage, 10) + ")"
			authMap["time"] = time.Now().Format("2006-01-02 15:04:05")

			err = nextStage.Auth(authMap)
			if err != nil {
				log.Error("[stageLog's Start]:error when auth to stage:", nextStageLogInfo, " ===>error is:", err.Error())
				stageLog.Stop(StageStopScopeAll, StageStopReasonRunFailed, models.StageLogStateRunFailed)
				return
			}
		}
	}()

	log.Info("got actionList:", actionLogList)
	for _, actionLog := range actionLogList {
		action := new(ActionLog)
		action.ActionLog = &actionLog
		err = action.Listen()
		if err != nil {
			log.Error("[stageLog's Start]:error when set action to listen:", err.Error())
			stageLog.Stop(StageStopScopeAll, StageStopReasonRunFailed, models.StageLogStateRunFailed)
			return
		}

		authMap := make(map[string]interface{})
		authMap["type"] = AuthTyptStageStartDone
		authMap["token"] = AuthTokenDefault
		authMap["authorizer"] = "system - " + stageLog.Namespace + " - " + stageLog.Repository + " - " +
			strconv.FormatInt(stageLog.Pipeline, 10) + "(" + strconv.FormatInt(stageLog.FromPipeline, 10) + ") - " +
			strconv.FormatInt(stageLog.ID, 10) + "(" + strconv.FormatInt(stageLog.FromStage, 10) + ")"
		authMap["time"] = time.Now().Format("2006-01-02 15:04:05")

		err = action.Auth(authMap)
		if err != nil {
			log.Error("[stageLog's Start]:error when auth to action:", actionLog, " ===>error is:", err.Error())
			stageLog.Stop(StageStopScopeAll, StageStopReasonRunFailed, models.StageLogStateRunFailed)
			return
		}
	}
}

func (stageLog *StageLog) Stop(scope, reason string, runState int64) {
	err := stageLog.GetStageLog().Where("id = ?", stageLog.ID).First(stageLog).Error
	if err != nil {
		log.Error("[stageLog's Stop]:error when get pipelinelog info from db:", err.Error())
		return
	}

	actionLogList := make([]models.ActionLog, 0)
	needStopActionList := make([]models.ActionLog, 0)
	err = new(models.ActionLog).GetActionLog().Where("pipeline = ?", stageLog.Pipeline).Where("stage = ?", stageLog.ID).Find(&actionLogList).Error
	if err != nil {
		log.Error("[stageLog's Stop]:error when get actionlog's list from db:", err.Error())
		return
	}

	for _, actionLog := range actionLogList {
		if scope == StageStopScopeRecyclable && actionLog.Timeout != 0 {
			needStopActionList = append(needStopActionList, actionLog)
		} else {
			needStopActionList = append(needStopActionList, actionLog)
		}
	}

	for _, actionLog := range needStopActionList {
		action := new(ActionLog)
		action.ActionLog = &actionLog
		action.Stop(reason, runState)
	}

	stageLog.RunState = runState
	err = stageLog.GetStageLog().Save(stageLog).Error
	if err != nil {
		log.Error("[stageLog's Stop]:error when change stage's run state:", stageLog, " ===>error is:", err.Error())
	}

	if reason == StageStopReasonRunFailed || reason == StageStopReasonTimeout {
		pipelineLogInfo := new(models.PipelineLog)
		err = pipelineLogInfo.GetPipelineLog().Where("id = ?", stageLog.Pipeline).First(pipelineLogInfo).Error
		if err != nil {
			log.Error("[stageLog's Stop]:error when get pipelineLog ingo from db:", err.Error())
			return
		}

		pipeline := new(PipelineLog)
		pipeline.PipelineLog = pipelineLogInfo

		pipeline.Stop(reason, runState)
	}
}

func (stageLog *StageLog) WaitAllActionDone(nextStageCanStartChan chan bool) {
	actionLogList := make([]models.ActionLog, 0)
	err := new(models.ActionLog).GetActionLog().Where("pipeline = ?", stageLog.Pipeline).Where("sequence = ?", stageLog.Sequence).Where("stage = ?", stageLog.ID).Find(&actionLogList).Error
	if err != nil {
		log.Error("[stageLog's WaitAllActionDone]:error when get action list from db:", err.Error())
		stageLog.Stop(StageStopScopeAll, StageStopReasonRunFailed, models.StageLogStateRunFailed)
		return
	}

	stopWait := false
	actionResultChan := make(chan bool, len(actionLogList))

	for _, actionLog := range actionLogList {
		go func(actionLog models.ActionLog, resultChan chan bool) {
			for !stopWait {
				actionLogInfo := new(models.ActionLog)
				err := actionLogInfo.GetActionLog().Where("id = ?", actionLog.ID).First(actionLogInfo).Error
				if err != nil {
					log.Error("[stageLog's WaitAllActionDone]:error when get actionlog's info from db:", err.Error())
					resultChan <- false
					return
				}

				if actionLogInfo.RunState == models.ActionLogStateRunFailed {
					resultChan <- false
					return
				} else if actionLogInfo.RunState == models.ActionLogStateRunSuccess {
					resultChan <- true
					return
				}

				time.Sleep(1 * time.Second)
			}

		}(actionLog, actionResultChan)
	}

	finalResultChan := make(chan bool, 1)

	go func() {
		count := 0
		for {
			runResult := <-actionResultChan
			if runResult {
				count++
			} else {
				finalResultChan <- false
				return
			}

			if count == len(actionLogList) {
				finalResultChan <- true
				return
			}
		}
	}()

	if stageLog.Timeout != 0 {
		duration, _ := time.ParseDuration(strconv.FormatInt(stageLog.Timeout, 10) + "s")
		select {
		case <-time.After(duration):
			log.Error("[stageLog's WaitAllActionDone]:got a timeout from stage", stageLog)
			nextStageCanStartChan <- false
			stageLog.Stop(StageStopScopeAll, StageStopReasonRunFailed, models.StageLogStateRunFailed)
			return
		case runResult := <-finalResultChan:
			nextStageCanStartChan <- runResult
		}
	} else {
		runResult := <-finalResultChan
		nextStageCanStartChan <- runResult
	}
}
