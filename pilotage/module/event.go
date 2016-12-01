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
	"errors"
	"strconv"
	"strings"

	"github.com/Huawei/containerops/pilotage/models"

	log "github.com/Sirupsen/logrus"
	"github.com/containerops/configure"
	"github.com/jinzhu/gorm"
)

var eventList = map[string]string{
	"CO_COMPONENT_START": "CO_COMPONENT_START",
	"CO_COMPONENT_STOP":  "CO_COMPONENT_STOP",
	"CO_TASK_START":      "CO_TASK_START",
	"CO_TASK_RESULT":     "CO_TASK_RESULT",
	"CO_TASK_STATUS":     "CO_TASK_STATUS",
	"CO_REGISTER_URL":    "CO_register"}
var projectAddr = ""

func init() {
	if configure.GetString("projectaddr") == "" {
		projectAddr = "http://localhost"
	} else {
		projectAddr = configure.GetString("projectaddr")
	}
	projectAddr = strings.TrimSuffix(projectAddr, "/")
}

func setSystemEvent(db *gorm.DB, actionLog *models.ActionLog) error {
	if db == nil {
		db = models.GetDB().Begin()
		err := db.Error
		if err != nil {
			log.Error("[setSystemEvent]:when db.Begin():", err.Error())
			return err
		}
	}

	workflowLog := new(models.WorkflowLog)
	err := db.Model(&models.WorkflowLog{}).Where("id = ?", actionLog.Workflow).First(workflowLog).Error
	if err != nil {
		log.Error("[setSystemEvent]:error when get workflowlog info from db:", err.Error())
		rollbackErr := db.Rollback().Error
		if rollbackErr != nil {
			log.Error("[setSystemEvent]:when rollback in get workflowlog's info:", rollbackErr.Error())
			return errors.New("errors occur:\nerror1:" + err.Error() + "\nerror2:" + rollbackErr.Error())
		}
		return err
	}

	for key, value := range eventList {
		tempEvent := new(models.EventDefinition)
		tempEvent.Event = key
		tempEvent.Title = key
		tempEvent.Namespace = actionLog.Namespace
		tempEvent.Repository = actionLog.Repository
		tempEvent.Workflow = actionLog.Workflow
		tempEvent.Stage = actionLog.Stage
		tempEvent.Action = actionLog.ID
		tempEvent.Character = models.CharacterComponentEvent
		tempEvent.Type = models.TypeSystemEvent
		tempEvent.Source = models.SourceInnerEvent
		tempEvent.Definition = projectAddr + "/v2/" + actionLog.Namespace + "/" + actionLog.Repository + "/workflow/v1/event/" + workflowLog.Workflow + "/" + value

		err := db.Save(tempEvent).Error
		if err != nil {
			log.Error("[setSystemEvent]:error when save event definition to db:", err.Error())
			rollbackErr := db.Rollback().Error
			if rollbackErr != nil {
				log.Error("[setSystemEvent]:when rollback in get workflowlog's info:", rollbackErr.Error())
				return errors.New("errors occur:\nerror1:" + err.Error() + "\nerror2:" + rollbackErr.Error())
			}
			return err
		}
	}

	return nil
}

func getSystemEventList(actionID int64) ([]map[string]interface{}, error) {
	result := make([]map[string]interface{}, 0)

	eventDefineList := make([]models.EventDefinition, 0)

	err := new(models.EventDefinition).GetEventDefinition().Where("action = ?", actionID).Find(&eventDefineList).Error

	if err != nil {
		log.Error("[getSystemEventList]:error when get systemEventList from db:", err.Error())
		return nil, err
	}

	for _, eventDefine := range eventDefineList {
		tempMap := make(map[string]interface{})
		tempMap["name"] = eventDefine.Title
		tempMap["value"] = eventDefine.Definition
		tempMap["Title"] = eventDefine.Title
		tempMap["ID"] = eventDefine.ID

		result = append(result, tempMap)
	}

	return result, nil
}

func RecordEventInfo(eventDefineId, sequence int64, headerInfo, payload, authInfo string, eventDefineInfo ...string) error {
	eventDefine := new(models.EventDefinition)
	if eventDefineId < 0 {
		eventDefine.Type = models.TypeSystemEvent
		eventDefine.Source = models.SourceInnerEvent

		if len(eventDefineInfo) > 0 {
			eventDefine.Title = eventDefineInfo[0]
		}

		if len(eventDefineInfo) > 1 {
			characterInt, _ := strconv.ParseInt(eventDefineInfo[1], 10, 64)
			eventDefine.Character = characterInt
		}

		if len(eventDefineInfo) > 2 {
			eventDefine.Namespace = eventDefineInfo[2]
		}

		if len(eventDefineInfo) > 3 {
			eventDefine.Repository = eventDefineInfo[3]
		}

		if len(eventDefineInfo) > 4 {
			pipelinInt, _ := strconv.ParseInt(eventDefineInfo[4], 10, 64)
			eventDefine.Workflow = pipelinInt
		}

		if len(eventDefineInfo) > 5 {
			stageInt, _ := strconv.ParseInt(eventDefineInfo[5], 10, 64)
			eventDefine.Stage = stageInt
		}

		if len(eventDefineInfo) > 6 {
			actionInt, _ := strconv.ParseInt(eventDefineInfo[6], 10, 64)
			eventDefine.Action = actionInt
		}
	} else {
		err := eventDefine.GetEventDefinition().Where("id = ?", eventDefineId).First(eventDefine).Error
		if err != nil {
			log.Error("[event's RecordEventInfo]:error when get event define from db:", err.Error())
			return err
		}
	}

	event := new(models.Event)
	event.Definition = eventDefineId
	event.Title = eventDefine.Title
	event.Header = headerInfo
	event.Payload = payload
	event.Authorization = authInfo
	event.Type = eventDefine.Type
	event.Source = eventDefine.Source
	event.Character = eventDefine.Character
	event.Namespace = eventDefine.Namespace
	event.Repository = eventDefine.Repository
	event.Workflow = eventDefine.Workflow
	event.Stage = eventDefine.Stage
	event.Action = eventDefine.Action
	event.Sequence = sequence

	err := event.GetEvent().Save(event).Error
	if err != nil {
		log.Error("[event's RecordEventInfo]:error when save event info to db:", err.Error())
		return err
	}

	return nil
}

func SetWorkflowVarInfo(id int64, varMap map[string]interface{}) error {
	db := models.GetDB().Begin()
	err := db.Error
	if err != nil {
		log.Error("[workflowVar's SetWorkflowVarInfo]:when db.Begin():", err.Error())
		return errors.New("error when db.Begin")
	}

	err = db.Model(&models.WorkflowVar{}).Where("workflow = ?", id).Unscoped().Delete(&models.WorkflowVar{}).Error
	if err != nil {
		log.Error("[workflowVar's SetWorkflowVarInfo]:when delet var info from db:", err.Error())
		rollbackErr := db.Rollback().Error
		if rollbackErr != nil {
			log.Error("[workflowVar's SetWorkflowVarInfo]:when rollback in delet var info got err:", rollbackErr.Error())
			return errors.New("errors occur:\nerror1:" + err.Error() + "\nerror2:" + rollbackErr.Error())
		}
		return errors.New("error when delete var info from db:" + err.Error())
	}

	for key, defaultValue := range varMap {
		varSet := new(models.WorkflowVar)

		defaultValueStr, ok := defaultValue.(string)
		if !ok {
			log.Error("[workflowVar's SetWorkflowVarInfo]:error when pase default value, want a string,got:", defaultValue)
			return errors.New("var's vaule is not string")
		}

		varSet.Workflow = id
		varSet.Key = key
		varSet.Default = defaultValueStr

		err = db.Model(&models.WorkflowVar{}).Save(varSet).Error
		if err != nil {
			log.Error("[workflowVar's SetWorkflowVarInfo]:when save var info from db:", err.Error())
			rollbackErr := db.Rollback().Error
			if rollbackErr != nil {
				log.Error("[workflowVar's SetWorkflowVarInfo]:when rollback in save var info got err:", rollbackErr.Error())
				return errors.New("errors occur:\nerror1:" + err.Error() + "\nerror2:" + rollbackErr.Error())
			}
			return errors.New("error when save var info")
		}
	}

	db.Commit()
	return nil
}

func GetWorkflowVarInfo(id int64) (map[string]string, error) {
	resultMap := make(map[string]string)
	varList := make([]models.WorkflowVar, 0)

	err := new(models.WorkflowVar).GetWorkflowVar().Where("workflow = ?", id).Find(&varList).Error
	if err != nil {
		log.Error("[workflowVar's GetWorkflowVarInfo]:error when get var list from db:", err.Error())
		return nil, errors.New("error when get var info from db")
	}

	for _, varInfo := range varList {
		resultMap[varInfo.Key] = varInfo.Default
	}

	return resultMap, nil
}
