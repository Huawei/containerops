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
	"strconv"
	"time"

	"github.com/Huawei/containerops/pilotage/models"
	cron "github.com/Huawei/containerops/pilotage/utils/timer"

	log "github.com/Sirupsen/logrus"
)

var workflowsTimer map[string]*cron.Cron
var taskNextStartList map[string][]models.Timer
var updateTaskChan chan bool

func init() {
	workflowsTimer = make(map[string]*cron.Cron)
	taskNextStartList = make(map[string][]models.Timer)
	updateTaskChan = make(chan bool, 1)
}

// InitTimerTask is
func InitTimerTask() {
	taskList := make([]models.Timer, 0)

	new(models.Timer).GetTimer().Where("available = ?", true).Find(&taskList)

	for _, task := range taskList {
		index := getTimerMapIndex(task.Namespace, task.Repository, task.Workflow)

		var timer *cron.Cron
		if t, ok := workflowsTimer[index]; !ok {
			t = cron.New()
			t.Start()

			workflowsTimer[index] = t
			timer = t
		} else {
			timer = t
		}

		schedule, err := cron.Parse(task.Cron)
		if err != nil {
			log.Error("[timer's InitTimerTask]:error when parse task:", task, " error is:", err.Error())
			continue
		}

		next := schedule.Next(time.Now()).Format("2006-01-02 15:04:05")
		if list, ok := taskNextStartList[next]; ok {
			list = append(list, task)
			taskNextStartList[next] = list
		} else {
			list = make([]models.Timer, 0)
			list = append(list, task)
			taskNextStartList[next] = list
		}

		timer.AddFunc(task.Cron, func() { StartTask() })
		log.Info("add task ", task.Cron, " to ", task.Namespace, "-", task.Repository, "-", task.Workflow)
	}
}

// UpdateWorkflowTimer is
func UpdateWorkflowTimer(namespace, repository string, workflowID int64) {
	index := getTimerMapIndex(namespace, repository, workflowID)

	taskList := make([]models.Timer, 0)

	new(models.Timer).GetTimer().Where("namespace = ?", namespace).Where("repository = ?", repository).Where("workflow = ?", workflowID).Where("available = ?", true).Find(&taskList)

	var timer *cron.Cron
	if t, ok := workflowsTimer[index]; ok {
		t.Stop()
	}

	timer = cron.New()
	workflowsTimer[index] = timer
	timer.Start()

	for _, task := range taskList {
		schedule, err := cron.Parse(task.Cron)
		if err != nil {
			log.Error("[timer's InitTimerTask]:error when parse task:", task, " error is:", err.Error())
			continue
		}

		next := schedule.Next(time.Now()).Format("2006-01-02 15:04:05")
		if list, ok := taskNextStartList[next]; ok {
			list = append(list, task)
			taskNextStartList[next] = list
		} else {
			list = make([]models.Timer, 0)
			list = append(list, task)
			taskNextStartList[next] = list
		}

		timer.AddFunc(task.Cron, func() { StartTask() })
	}
}

// StartTask is
func StartTask() {
	now := time.Now().Format("2006-01-02 15:04:05")

	shouldStartList := updateTaskNextStartList(now)

	for _, task := range shouldStartList {
		authMap := make(map[string]interface{})
		authMap["type"] = AuthTypeWorkflowDefault
		authMap["token"] = AuthTokenDefault
		authMap["eventName"] = task.EventName
		authMap["eventType"] = task.EventType
		authMap["time"] = time.Now().Format("2006-01-02 15:04:05")

		pass, err := checkInstanceNum(task.Workflow)
		if !pass {
			log.Error("[timer's StartTask]:error when check instance num:", err.Error())
			continue
		}

		_, err = Run(task.Workflow, authMap, task.StartJson)
		if err != nil {
			log.Error("[timer's StartTask]:error when run workflow:", err.Error())
			continue
		}
	}
}

func getTimerMapIndex(namespace, repository string, workflowID int64) string {
	return namespace + "-" + repository + "-" + strconv.FormatInt(workflowID, 10)
}

func updateTaskNextStartList(timeStr string) []models.Timer {
	updateTaskChan <- true
	defer func() {
		<-updateTaskChan
	}()

	list, ok := taskNextStartList[timeStr]
	delete(taskNextStartList, timeStr)

	if ok {
		for _, task := range list {
			schedule, err := cron.Parse(task.Cron)
			if err != nil {
				log.Error("[timer's updateTaskNextStartList]:error when parse task:", task, " error is:", err.Error())
				continue
			}

			next := schedule.Next(time.Now()).Format("2006-01-02 15:04:05")
			if list, ok := taskNextStartList[next]; ok {
				list = append(list, task)
				taskNextStartList[next] = list
			} else {
				list = make([]models.Timer, 0)
				list = append(list, task)
				taskNextStartList[next] = list
			}
		}

		return list
	}

	return make([]models.Timer, 0)
}
