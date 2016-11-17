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

var eventList = map[string]string{"COMPONENT_START": "event", "COMPONENT_STOP": "event", "TASK_START": "event", "TASK_RESULT": "event", "TASK_STATUS": "event", "REGISTER_URL": "register"}
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
		db = models.GetDB()
		err := db.Begin().Error
		if err != nil {
			log.Error("[setSystemEvent]:when db.Begin():", err.Error())
			return err
		}
	}

	pipelineLog := new(models.PipelineLog)
	err := db.Model(&models.PipelineLog{}).Where("id = ?", actionLog.Pipeline).First(pipelineLog).Error
	if err != nil {
		log.Error("[setSystemEvent]:error when get pipelinelog info from db:", err.Error())
		rollbackErr := db.Rollback().Error
		if rollbackErr != nil {
			log.Error("[setSystemEvent]:when rollback in get pipelinelog's info:", rollbackErr.Error())
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
		tempEvent.Pipeline = actionLog.Pipeline
		tempEvent.Stage = actionLog.Stage
		tempEvent.Action = actionLog.ID
		tempEvent.Character = models.CharacterComponentEvent
		tempEvent.Type = models.TypeSystemEvent
		tempEvent.Source = models.SourceInnerEvent
		tempEvent.Definition = projectAddr + "/pipeline/v1/" + actionLog.Namespace + "/" + actionLog.Repository + "/" + pipelineLog.Pipeline + "/" + value

		err := db.Save(tempEvent).Error
		if err != nil {
			log.Error("[setSystemEvent]:error when save event definition to db:", err.Error())
			rollbackErr := db.Rollback().Error
			if rollbackErr != nil {
				log.Error("[setSystemEvent]:when rollback in get pipelinelog's info:", rollbackErr.Error())
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
			eventDefine.Pipeline = pipelinInt
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
	event.Pipeline = eventDefine.Pipeline
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
