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

package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

const (
	// WorkflowStateDisable is the state that workflow is disabled ,can't start
	WorkflowStateDisable = iota
	// WorkflowStateAble is the state that workflow can start
	WorkflowStateAble
)

const (
	// WorkflowLogStateCanListen is the state that current workflowLog can be listen (one workflowlog only can listen one time)
	WorkflowLogStateCanListen = iota
	// WorkflowLogStateWaitToStart is the state that workflow is wait to start(may because lack some condition to start)
	WorkflowLogStateWaitToStart
	// WorkflowLogStateDoing is the state that workflow is working
	WorkflowLogStateDoing
	// WorkflowLogStateRunSuccess is the state that at this time, workflow run result is success
	WorkflowLogStateRunSuccess
	// WorkflowLogStateRunFailed is the state that at this time, workflow run result is failed
	WorkflowLogStateRunFailed
)

const (
	//StageTypeStart is the Stage type being the start the workflow.
	StageTypeStart = iota
	//StageTypeEnd is the Stage type being the end of the workflow.
	StageTypeEnd
	//StageTypeRun is the Stage type being the running stage of workflow.
	StageTypeRun
)

const (
	// StageLogStateCanListen is the state that current stageLog can be listen (one Stagelog only can listen one time)
	StageLogStateCanListen = iota
	// StageLogStateWaitToStart is the state that workflow is wait to start(may because lack some condition to start)
	StageLogStateWaitToStart
	// StageLogStateDoing is the state that workflow is working
	StageLogStateDoing
	// StageLogStateRunSuccess is the state that at this time, workflow run result is success
	StageLogStateRunSuccess
	// StageLogStateRunFailed is the state that at this time, workflow run result is failed
	StageLogStateRunFailed
)

const (
	// ActionLogStateCanListen is the state that current actionLog can be listen (one Actionlog only can listen one time)
	ActionLogStateCanListen = iota
	// ActionLogStateWaitToStart is the state that workflow is wait to start(may because lack some condition to start)
	ActionLogStateWaitToStart
	// ActionLogStateDoing is the state that workflow is working
	ActionLogStateDoing
	// ActionLogStateRunSuccess is the state that at this time, workflow run result is success
	ActionLogStateRunSuccess
	// ActionLogStateRunFailed is the state that at this time, workflow run result is failed
	ActionLogStateRunFailed
)

const (
	//When StageID point to the StageTypeStart , the Action ID is 0.
	//When StageID point to the StageTypeEnd , the Action ID is -1.

	// OutcomeTypeStageStartActionID is
	OutcomeTypeStageStartActionID = 0
	// OutcomeTypeStageEndActionID is
	OutcomeTypeStageEndActionID = -1

	// OutcomeTypeStageStartEventID is
	OutcomeTypeStageStartEventID = 0
	// OutcomeTypeStageEndEventID is
	OutcomeTypeStageEndEventID = -1
)

var (
	//StageTypeForWeb is the stage type that use for web display
	StageTypeForWeb = []string{"workflow-start", "workflow-end", "workflow-stage"}
)

//Workflow is DevOps workflow definition unit.
type Workflow struct {
	ID              int64      `json:"id" gorm:"primary_key"`                       //
	Namespace       string     `json:"namespace" sql:"not null;type:varchar(255)"`  //Username or organization
	Repository      string     `json:"repository" sql:"not null;type:varchar(255)"` //
	Workflow        string     `json:"workflow" sql:"not null;type:varchar(255)"`   //workflow name
	Event           int64      `json:"event" sql:"null;default:0"`                  //
	Version         string     `json:"version" sql:"null;type:varchar(255)"`        //User define Workflow version
	VersionCode     int64      `json:"versionCode" sql:"null;type:varchar(255)"`    //System define Workflow version,unique,for query
	State           int64      `json:"state" sql:"null;type:bigint"`                //workflow state
	Manifest        string     `json:"manifest" sql:"null;type:longtext"`           //
	Description     string     `json:"description" sql:"null;type:text"`            //
	SourceInfo      string     `json:"source" sql:"null;type:longtext"`             // define of source like : {"token":"","sourceList":[{"sourceType":"Github","headerKey":"X-Hub-Signature","eventList":",pull request,"]}
	Env             string     `json:"env" sql:"null;type:longtext"`                // env that all action in this workflow will get
	Requires        string     `json:"requires" sql:"type:longtext"`                // workflow run requires auth
	IsLimitInstance bool       `json:"isLimitInstance"`                             //
	LimitInstance   int64      `json:"limitInstance"`                               //
	CurrentInstance int64      `json:"currentInstance"`                             //
	CreatedAt       time.Time  `json:"created" sql:""`                              //
	UpdatedAt       time.Time  `json:"updated" sql:""`                              //
	DeletedAt       *time.Time `json:"deleted" sql:"index"`                         //
}

//TableName is return the table name of Workflow in MySQL database.
func (p *Workflow) TableName() string {
	return "workflow"
}

// GetWorkflow is
func (p *Workflow) GetWorkflow() *gorm.DB {
	return conn.Model(&Workflow{})
}

//WorkflowLog is workflow run history log.
type WorkflowLog struct {
	ID              int64      `json:"id" gorm:"primary_key"`                       //
	Namespace       string     `json:"namespace" sql:"not null;type:varchar(255)"`  //Username or organization
	Repository      string     `json:"repository" sql:"not null;type:varchar(255)"` //
	Workflow        string     `json:"workflow" sql:"not null;type:varchar(255)"`   //workflow name
	FromWorkflow    int64      `json:"fromWorkflow" sql:"not null;default:0"`       //
	PreWorkflow     int64      `json:"preWorkflow" sql:"not null;default:0"`
	PreStage        int64      `json:"preStage" sql:"not null;default:0"`
	PreAction       int64      `json:"preAction" sql:"not null;default:0"`
	PreWorkflowInfo string     `json:"preWorkflowInfo"`
	Version         string     `json:"version" sql:"null;type:varchar(255)"`     //User define Workflow version
	VersionCode     int64      `json:"versionCode" sql:"null;type:varchar(255)"` //System define Workflow version,unique,for query
	Sequence        int64      `json:"sequence" sql:"not null;default:0"`        //workflow run sequence
	RunState        int64      `json:"runState" sql:"null;type:bigint"`          //workflow run state
	FailReason      string     `json:"failReason"`                               //
	Event           int64      `json:"event" sql:"null;default:0"`               //
	Manifest        string     `json:"manifest"sql:"null;type:longtext"`         //
	Description     string     `json:"description" sql:"null;type:text"`         //
	SourceInfo      string     `json:"source" sql:"null;type:longtext"`          // define of source like : [{"sourceType":"Github","headerKey":"X-Hub-Signature","eventList":",pull request,","secretKey":"asdfFDSA!@d12"}]
	Env             string     `json:"env" sql:"null;type:longtext"`             // env that all action in this workflow will get
	Requires        string     `json:"requires" sql:"type:longtext"`             // workflow run requires auth
	AuthList        string     `json:"authList" sql:"type:longtext"`             //
	CreatedAt       time.Time  `json:"created" sql:""`                           //
	UpdatedAt       time.Time  `json:"updated" sql:""`                           //
	DeletedAt       *time.Time `json:"deleted" sql:"index"`                      //
}

//TableName is return the table name of Workflow in MySQL database.
func (p *WorkflowLog) TableName() string {
	return "workflow_log"
}

// GetWorkflowLog is
func (p *WorkflowLog) GetWorkflowLog() *gorm.DB {
	return conn.Model(&WorkflowLog{})
}

//Stage is Workflow unit.
type Stage struct {
	ID          int64      `json:"id" gorm:"primary_key"`                       //
	Namespace   string     `json:"namespace" sql:"not null;type:varchar(255)"`  //Username or organization
	Repository  string     `json:"repository" sql:"not null;type:varchar(255)"` //
	Workflow    int64      `json:"workflow" sql:"not null;default:0"`           //Workflow's ID.
	Type        int64      `json:"type" sql:"not null;default:0"`               //StageTypeStart, StageTypeEnd or StageTypeRun
	PreStage    int64      `json:"preStage" sql:"not null;default:0"`           //Pre stage ID ,first stage is -1
	Stage       string     `json:"stage" sql:"not null;type:varchar(255)"`      //Stage name for query.
	Title       string     `json:"title" sql:"not null;type:varchar(255)"`      //Stage title for display
	Description string     `json:"description" sql:"null;type:text"`            //
	Event       int64      `json:"event" sql:"null;default:0"`                  //
	Manifest    string     `json:"manifest" sql:"null;type:longtext"`           //
	Env         string     `json:"env" sql:"null;type:longtext"`                //
	Timeout     string     `json:"timeout"`                                     //
	Requires    string     `json:"requires" sql:"type:longtext"`                // workflow run requires auth
	CreatedAt   time.Time  `json:"created" sql:""`                              //
	UpdatedAt   time.Time  `json:"updated" sql:""`                              //
	DeletedAt   *time.Time `json:"deleted" sql:"index"`                         //
}

//TableName is return the table name of Stage in MySQL database.
func (s *Stage) TableName() string {
	return "stage"
}

// GetStage is
func (s *Stage) GetStage() *gorm.DB {
	return conn.Model(&Stage{})
}

//StageLog is stage run log.
type StageLog struct {
	ID           int64      `json:"id" gorm:"primary_key"`                       //
	Namespace    string     `json:"namespace" sql:"not null;type:varchar(255)"`  //Username or organization
	Repository   string     `json:"repository" sql:"not null;type:varchar(255)"` //
	Workflow     int64      `json:"workflow" sql:"not null;default:0"`           //WorkflowLog's ID.
	FromWorkflow int64      `json:"fromWorkflow" sql:"not null;default:0"`       //workflow's ID.
	Sequence     int64      `json:"sequence" sql:"not null;default:0"`           //workflow run sequence
	FromStage    int64      `json:"fromStage" sql:"not null;default:0"`          //
	Type         int64      `json:"type" sql:"not null;default:0"`               //StageTypeStart, StageTypeEnd or StageTypeRun
	PreStage     int64      `json:"preStage" sql:"not null;default:0"`           //Pre stage ID ,first stage is -1
	Stage        string     `json:"stage" sql:"not null;type:varchar(255)"`      //Stage name for query.
	Title        string     `json:"title" sql:"not null;type:varchar(255)"`      //Stage title for display
	Description  string     `json:"description" sql:"null;type:text"`            //
	RunState     int64      `json:"runState" sql:"null;type:bigint"`             //stage run state
	FailReason   string     `json:"failReason"`                                  //
	Event        int64      `json:"event" sql:"null;default:0"`                  //
	Manifest     string     `json:"manifest" sql:"null;type:longtext"`           //
	Env          string     `json:"env" sql:"null;type:longtext"`                //
	Timeout      string     `json:"timeout"`                                     //
	Requires     string     `json:"requires" sql:"type:longtext"`                // workflow run requires auth
	AuthList     string     `json:"authList" sql:"type:longtext"`                //
	CreatedAt    time.Time  `json:"created" sql:""`                              //
	UpdatedAt    time.Time  `json:"updated" sql:""`                              //
	DeletedAt    *time.Time `json:"deleted" sql:"index"`                         //
}

//TableName is return the table name of Stage in MySQL database.
func (s *StageLog) TableName() string {
	return "stage_log"
}

// GetStageLog is
func (s *StageLog) GetStageLog() *gorm.DB {
	return conn.Model(&StageLog{})
}

//Action is Stage unit.
type Action struct {
	ID          int64      `json:"id" gorm:"primary_key"`                       //
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
	CreatedAt   time.Time  `json:"created" sql:""`                              //
	UpdatedAt   time.Time  `json:"updated" sql:""`                              //
	DeletedAt   *time.Time `json:"deleted" sql:"index"`                         //
}

//TableName is return the name of Action in MySQL database.
func (a *Action) TableName() string {
	return "action"
}

// GetAction is
func (a *Action) GetAction() *gorm.DB {
	return conn.Model(&Action{})
}

//ActionLog is action run history.
type ActionLog struct {
	ID           int64      `json:"id" gorm:"primary_key"`                       //
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
	CreatedAt    time.Time  `json:"created" sql:""`                              //
	UpdatedAt    time.Time  `json:"updated" sql:""`                              //
	DeletedAt    *time.Time `json:"deleted" sql:"index"`                         //
}

//TableName is return the name of Action in MySQL database.
func (a *ActionLog) TableName() string {
	return "action_log"
}

// GetActionLog is
func (a *ActionLog) GetActionLog() *gorm.DB {
	return conn.Model(&ActionLog{})
}

//Outcome is Stage running results.
//When StageID point to the StageTypeStart , the Action ID is 0.
//When StageID point to the StageTypeEnd , the Action ID is -1.
type Outcome struct {
	ID           int64      `json:"id" gorm:"primary_key"`                 //
	Workflow     int64      `json:"workflow" sql:"not null;default:0"`     //WorkflowLog id
	RealWorkflow int64      `json:"realWorkflow" sql:"not null;default:0"` //Workflow id
	Stage        int64      `json:"stage" sql:"not null;default:0"`        //stageLog id
	RealStage    int64      `json:"realStage" sql:"not null;default:0"`    //stage id
	Action       int64      `json:"action" sql:"not null;default:0"`       //actionLog id
	RealAction   int64      `json:"realAction" sql:"not null;default:0"`   //
	Event        int64      `json:"event" sql:"null;default:0"`            //event id
	Sequence     int64      `json:"sequence" sql:"not null;default:0"`     //workflow run sequence
	Status       bool       `json:"status" sql:"null;varchar(255)"`        //
	Result       string     `json:"result" sql:"null;type:longtext"`       //
	Output       string     `json:"output" sql:"null;type:longtext"`       //
	CreatedAt    time.Time  `json:"created" sql:""`                        //
	UpdatedAt    time.Time  `json:"updated" sql:""`                        //
	DeletedAt    *time.Time `json:"deleted" sql:"index"`                   //
}

//TableName is return the name of Outcome in MySQL database.
func (o *Outcome) TableName() string {
	return "outcome"
}

// GetOutcome is
func (o *Outcome) GetOutcome() *gorm.DB {
	return conn.Model(&Outcome{})
}

// WorkflowSequence is a table describe workflow's run sequence
type WorkflowSequence struct {
	ID        int64      `json:"id" gorm:"primary_key"`             //
	Workflow  int64      `json:"workflow" sql:"not null;default:0"` // workflow name
	Sequence  int64      `json:"sequence" sql:"not null;default:0"` //workflow run sequence
	CreatedAt time.Time  `json:"created" sql:""`                    //
	UpdatedAt time.Time  `json:"updated" sql:""`                    //
	DeletedAt *time.Time `json:"deleted" sql:"index"`               //
}

//TableName is return the table name of Workflow in MySQL database.
func (p *WorkflowSequence) TableName() string {
	return "workflow_sequence"
}

// GetWorkflowSequence is
func (p *WorkflowSequence) GetWorkflowSequence() *gorm.DB {
	return conn.Model(&WorkflowSequence{})
}
