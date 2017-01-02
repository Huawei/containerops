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
package models

import "github.com/jinzhu/gorm"

//Action is Stage unit.
type Action struct {
	BaseIDField
	Namespace   string     `json:"namespace" sql:"not null;type:varchar(255)"`  //Username or organization
	Repository  string     `json:"repository" sql:"not null;type:varchar(255)"` //
	Workflow    int64      `json:"workflow" sql:"not null;default:0"`           //WorkflowLog's ID.
	Stage       int64      `json:"stage" sql:"not null;default:0"`              //
	Component   int64      `json:"component" sql:"not null;default:0"`          //
	Service     int64      `json:"service" sql:"not null;default:0"`            //
	Action      string     `json:"action" sql:"not null;varchar(255)"`          //
	Title       string     `json:"title" sql:"not null;type:varchar(255)"`      //
	Description string     `json:"description" sql:"null;type:text"`            //
	Event       int64      `json:"event" sql:"null;default:0"`                  //
	Manifest    string     `json:"manifest" sql:"null;type:longtext"`           // has run platform's type and platform setting
	Environment string     `json:"environment" sql:"null;type:text"`            // Environment parameters.
	Kubernetes  string     `json:"kubernetes" sql:"null;type:text"`             //
	Swarm       string     `json:"swarm" sql:"null;type:text"`                  //
	Input       string     `json:"input" sql:"null;type:text"`                  //
	Output      string     `json:"input" sql:"null;type:text"`                  //
	ImageName   string     `json:"imageName"`                                   //
	ImageTag    string     `json:"imageTag"`                                    //
	Timeout     string     `json:"timeout"`                                     //
	Requires    string     `json:"requires" sql:"type:longtext"`                // workflow run requires auth
	BaseModel1
}

//TableName is return the name of Action in MySQL database.
func (a *Action) TableName() string {
	return "action"
}

func (a *Action) GetAction() *gorm.DB {
	return db.Model(&Action{})
}

//ActionLog is action run history.
type ActionLog struct {
	BaseIDField
	Namespace    string     `json:"namespace" sql:"not null;type:varchar(255)"`  //Username or organization
	Repository   string     `json:"repository" sql:"not null;type:varchar(255)"` //
	Workflow     int64      `json:"workflow" sql:"not null;default:0"`           //WorkflowLog's ID.
	FromWorkflow int64      `json:"fromWorkflow" sql:"not null;default:0"`       //
	Sequence     int64      `json:"sequence" sql:"not null;default:0"`           //workflow run sequence
	Stage        int64      `json:"stage" sql:"not null;default:0"`              //
	FromStage    int64      `json:"fromStage" sql:"not null;default:0"`          //
	FromAction   int64      `json:"fromAction" sql:"not null;default:0"`         //
	RunState     int64      `json:"runState" sql:"null;type:bigint"`             //action run state
	FailReason   string     `json:"failReason"`                                  //
	Component    int64      `json:"component" sql:"not null;default:0"`          //
	Service      int64      `json:"service" sql:"not null;default:0"`            //
	Action       string     `json:"action" sql:"not null;varchar(255)"`          //
	ContainerId  string     `json:"containerID"`                                 //
	Title        string     `json:"title" sql:"not null;type:varchar(255)"`      //
	Description  string     `json:"description" sql:"null;type:text"`            //
	Event        int64      `json:"event" sql:"null;default:0"`                  //
	Manifest     string     `json:"manifest" sql:"null;type:longtext"`           //
	Environment  string     `json:"environment" sql:"null;type:text"`            // Environment parameters.
	Kubernetes   string     `json:"kubernetes" sql:"null;type:text"`             //
	Swarm        string     `json:"swarm" sql:"null;type:text"`                  //
	Input        string     `json:"input" sql:"null;type:text"`                  //
	Output       string     `json:"input" sql:"null;type:text"`                  //
	ImageName    string     `json:"imageName"`                                   //
	ImageTag     string     `json:"imageTag"`                                    //
	Timeout      string     `json:"timeout"`                                     //
	Requires     string     `json:"requires" sql:"type:longtext"`                // workflow run requires auth
	AuthList     string     `json:"authList" sql:"type:longtext"`                //
	BaseModel2
}

func (log *ActionLog) TableName() string {
	return "action_log"
}

func (log *ActionLog) GetActionLog() *gorm.DB {
	return db.Model(&ActionLog{})
}

func SelectActionLogFromID(id uint64) (actionLog *ActionLog, err error) {
	var condition ActionLog
	condition.ID = id
	err = db.Where(&condition).First(actionLog).Error
	return
}
