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

import (
	"time"

	"github.com/jinzhu/gorm"
)

const (
	// PipelineStateDisable is the state that pipeline is disabled ,can't start
	PipelineStateDisable = iota
	// PipelineStateAble is the state that pipeline can start
	PipelineStateAble
)

const (
	// PipelineLogStateCanListen is the state that current pipelineLog can be listen (one pipelinelog only can listen one time)
	PipelineLogStateCanListen = iota
	// PipelineLogStateWaitToStart is the state that pipeline is wait to start(may because lack some condition to start)
	PipelineLogStateWaitToStart
	// PipelineLogStateDoing is the state that pipeline is working
	PipelineLogStateDoing
	// PipelineLogStateRunSuccess is the state that at this time, pipeline run result is success
	PipelineLogStateRunSuccess
	// PipelineLogStateRunFailed is the state that at this time, pipeline run result is failed
	PipelineLogStateRunFailed
)

const (
	//StageTypeStart is the Stage type being the start the pipeline.
	StageTypeStart = iota
	//StageTypeEnd is the Stage type being the end of the pipeline.
	StageTypeEnd
	//StageTypeRun is the Stage type being the running stage of pipeline.
	StageTypeRun
)

const (
	// StageLogStateCanListen is the state that current stageLog can be listen (one Stagelog only can listen one time)
	StageLogStateCanListen = iota
	// StageLogStateWaitToStart is the state that pipeline is wait to start(may because lack some condition to start)
	StageLogStateWaitToStart
	// StageLogStateDoing is the state that pipeline is working
	StageLogStateDoing
	// StageLogStateRunSuccess is the state that at this time, pipeline run result is success
	StageLogStateRunSuccess
	// StageLogStateRunFailed is the state that at this time, pipeline run result is failed
	StageLogStateRunFailed
)

const (
	// ActionLogStateCanListen is the state that current actionLog can be listen (one Actionlog only can listen one time)
	ActionLogStateCanListen = iota
	// ActionLogStateWaitToStart is the state that pipeline is wait to start(may because lack some condition to start)
	ActionLogStateWaitToStart
	// ActionLogStateDoing is the state that pipeline is working
	ActionLogStateDoing
	// ActionLogStateRunSuccess is the state that at this time, pipeline run result is success
	ActionLogStateRunSuccess
	// ActionLogStateRunFailed is the state that at this time, pipeline run result is failed
	ActionLogStateRunFailed
)

const (
	//When StageID point to the StageTypeStart , the Action ID is 0.
	//When StageID point to the StageTypeEnd , the Action ID is -1.
	OutcomeTypeStageStartActionID = 0
	OutcomeTypeStageEndActionID   = -1

	OutcomeTypeStageStartEventID = 0
	OutcomeTypeStageEndEventID   = -1
)

var (
	//StageTypeForWeb is the stage type that use for web display
	StageTypeForWeb = []string{"pipeline-start", "pipeline-end", "pipeline-stage"}
)

//Pipeline is DevOps workflow definition unit.
type Pipeline struct {
	ID          int64      `json:"id" gorm:"primary_key"`                       //
	Namespace   string     `json:"namespace" sql:"not null;type:varchar(255)"`  //Username or organization
	Repository  string     `json:"repository" sql:"not null;type:varchar(255)"` //
	Pipeline    string     `json:"pipeline" sql:"not null;type:varchar(255)"`   //pipeline name
	Event       int64      `json:"event" sql:"null;default:0"`                  //
	Version     string     `json:"version" sql:"null;type:varchar(255)"`        //User define Pipeline version
	VersionCode int64      `json:"versionCode" sql:"null;type:varchar(255)"`    //System define Pipeline version,unique,for query
	State       int64      `json:"state" sql:"null;type:bigint"`                //pipeline state
	Manifest    string     `json:"manifest" sql:"null;type:longtext"`           //
	Description string     `json:"description" sql:"null;type:text"`            //
	SourceInfo  string     `json:"source"`                                      // define of source like : {"token":"","sourceList":[{"sourceType":"Github","headerKey":"X-Hub-Signature","eventList":",pull request,"]}
	Env         string     `json:"env" sql:"null;type:longtext"`                // env that all action in this pipeline will get
	Requires    string     `json:"requires" sql:"type:longtext"`                // pipeline run requires auth
	CreatedAt   time.Time  `json:"created" sql:""`                              //
	UpdatedAt   time.Time  `json:"updated" sql:""`                              //
	DeletedAt   *time.Time `json:"deleted" sql:"index"`                         //
}

//TableName is return the table name of Pipeline in MySQL database.
func (p *Pipeline) TableName() string {
	return "pipeline"
}

func (p *Pipeline) GetPipeline() *gorm.DB {
	return db.Model(&Pipeline{})
}

//PipelineLog is pipeline run history log.
type PipelineLog struct {
	ID           int64      `json:"id" gorm:"primary_key"`                       //
	Namespace    string     `json:"namespace" sql:"not null;type:varchar(255)"`  //Username or organization
	Repository   string     `json:"repository" sql:"not null;type:varchar(255)"` //
	Pipeline     string     `json:"pipeline" sql:"not null;type:varchar(255)"`   //pipeline name
	FromPipeline int64      `json:"fromPipeline" sql:"not null;default:0"`       //
	Version      string     `json:"version" sql:"null;type:varchar(255)"`        //User define Pipeline version
	VersionCode  int64      `json:"versionCode" sql:"null;type:varchar(255)"`    //System define Pipeline version,unique,for query
	Sequence     int64      `json:"sequence" sql:"not null;default:0"`           //pipeline run sequence
	RunState     int64      `json:"runState" sql:"null;type:bigint"`             //pipeline run state
	Event        int64      `json:"event" sql:"null;default:0"`                  //
	Manifest     string     `json:"manifest"sql:"null;type:longtext"`            //
	Description  string     `json:"description" sql:"null;type:text"`            //
	SourceInfo   string     `json:"source"`                                      // define of source like : [{"sourceType":"Github","headerKey":"X-Hub-Signature","eventList":",pull request,","secretKey":"asdfFDSA!@d12"}]
	Env          string     `json:"env" sql:"null;type:longtext"`                // env that all action in this pipeline will get
	Requires     string     `json:"requires" sql:"type:longtext"`                // pipeline run requires auth
	AuthList     string     `json:"authList" sql:"type:longtext"`                //
	CreatedAt    time.Time  `json:"created" sql:""`                              //
	UpdatedAt    time.Time  `json:"updated" sql:""`                              //
	DeletedAt    *time.Time `json:"deleted" sql:"index"`                         //
}

//TableName is return the table name of Pipeline in MySQL database.
func (p *PipelineLog) TableName() string {
	return "pipeline_log"
}

func (p *PipelineLog) GetPipelineLog() *gorm.DB {
	return db.Model(&PipelineLog{})
}

//Stage is Pipeline unit.
type Stage struct {
	ID          int64      `json:"id" gorm:"primary_key"`                       //
	Namespace   string     `json:"namespace" sql:"not null;type:varchar(255)"`  //Username or organization
	Repository  string     `json:"repository" sql:"not null;type:varchar(255)"` //
	Pipeline    int64      `json:"pipeline" sql:"not null;default:0"`           //Pipeline's ID.
	Type        int64      `json:"type" sql:"not null;default:0"`               //StageTypeStart, StageTypeEnd or StageTypeRun
	PreStage    int64      `json:"preStage" sql:"not null;default:0"`           //Pre stage ID ,first stage is -1
	Stage       string     `json:"stage" sql:"not null;type:varchar(255)"`      //Stage name for query.
	Title       string     `json:"title" sql:"not null;type:varchar(255)"`      //Stage title for display
	Description string     `json:"description" sql:"null;type:text"`            //
	Event       int64      `json:"event" sql:"null;default:0"`                  //
	Manifest    string     `json:"manifest" sql:"null;type:longtext"`           //
	Env         string     `json:"env" sql:"null;type:longtext"`                //
	Timeout     int64      `json:"timeout"`                                     //
	Requires    string     `json:"requires" sql:"type:longtext"`                // pipeline run requires auth
	CreatedAt   time.Time  `json:"created" sql:""`                              //
	UpdatedAt   time.Time  `json:"updated" sql:""`                              //
	DeletedAt   *time.Time `json:"deleted" sql:"index"`                         //
}

//TableName is return the table name of Stage in MySQL database.
func (s *Stage) TableName() string {
	return "stage"
}

func (s *Stage) GetStage() *gorm.DB {
	return db.Model(&Stage{})
}

//StageLog is stage run log.
type StageLog struct {
	ID           int64      `json:"id" gorm:"primary_key"`                       //
	Namespace    string     `json:"namespace" sql:"not null;type:varchar(255)"`  //Username or organization
	Repository   string     `json:"repository" sql:"not null;type:varchar(255)"` //
	Pipeline     int64      `json:"pipeline" sql:"not null;default:0"`           //PipelineLog's ID.
	FromPipeline int64      `json:"fromPipeline" sql:"not null;default:0"`       //pipeline's ID.
	Sequence     int64      `json:"sequence" sql:"not null;default:0"`           //pipeline run sequence
	FromStage    int64      `json:"fromStage" sql:"not null;default:0"`          //
	Type         int64      `json:"type" sql:"not null;default:0"`               //StageTypeStart, StageTypeEnd or StageTypeRun
	PreStage     int64      `json:"preStage" sql:"not null;default:0"`           //Pre stage ID ,first stage is -1
	Stage        string     `json:"stage" sql:"not null;type:varchar(255)"`      //Stage name for query.
	Title        string     `json:"title" sql:"not null;type:varchar(255)"`      //Stage title for display
	Description  string     `json:"description" sql:"null;type:text"`            //
	RunState     int64      `json:"runState" sql:"null;type:bigint"`             //stage run state
	Event        int64      `json:"event" sql:"null;default:0"`                  //
	Manifest     string     `json:"manifest" sql:"null;type:longtext"`           //
	Env          string     `json:"env" sql:"null;type:longtext"`                //
	Timeout      int64      `json:"timeout"`                                     //
	Requires     string     `json:"requires" sql:"type:longtext"`                // pipeline run requires auth
	AuthList     string     `json:"authList" sql:"type:longtext"`                //
	CreatedAt    time.Time  `json:"created" sql:""`                              //
	UpdatedAt    time.Time  `json:"updated" sql:""`                              //
	DeletedAt    *time.Time `json:"deleted" sql:"index"`                         //
}

//TableName is return the table name of Stage in MySQL database.
func (s *StageLog) TableName() string {
	return "stage_log"
}

func (s *StageLog) GetStageLog() *gorm.DB {
	return db.Model(&StageLog{})
}

//Action is Stage unit.
type Action struct {
	ID          int64      `json:"id" gorm:"primary_key"`                       //
	Namespace   string     `json:"namespace" sql:"not null;type:varchar(255)"`  //Username or organization
	Repository  string     `json:"repository" sql:"not null;type:varchar(255)"` //
	Pipeline    int64      `json:"pipeline" sql:"not null;default:0"`           //PipelineLog's ID.
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
	Endpoint    string     `json:"endpoint"`                                    //
	Timeout     int64      `json:"timeout"`                                     //
	Requires    string     `json:"requires" sql:"type:longtext"`                // pipeline run requires auth
	CreatedAt   time.Time  `json:"created" sql:""`                              //
	UpdatedAt   time.Time  `json:"updated" sql:""`                              //
	DeletedAt   *time.Time `json:"deleted" sql:"index"`                         //
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
	ID           int64      `json:"id" gorm:"primary_key"`                       //
	Namespace    string     `json:"namespace" sql:"not null;type:varchar(255)"`  //Username or organization
	Repository   string     `json:"repository" sql:"not null;type:varchar(255)"` //
	Pipeline     int64      `json:"pipeline" sql:"not null;default:0"`           //PipelineLog's ID.
	FromPipeline int64      `json:"fromPipeline" sql:"not null;default:0"`       //
	Sequence     int64      `json:"sequence" sql:"not null;default:0"`           //pipeline run sequence
	Stage        int64      `json:"stage" sql:"not null;default:0"`              //
	FromStage    int64      `json:"fromStage" sql:"not null;default:0"`          //
	FromAction   int64      `json:"fromAction" sql:"not null;default:0"`         //
	RunState     int64      `json:"runState" sql:"null;type:bigint"`             //action run state
	Component    int64      `json:"component" sql:"not null;default:0"`          //
	Service      int64      `json:"service" sql:"not null;default:0"`            //
	Action       string     `json:"action" sql:"not null;varchar(255)"`          //
	Title        string     `json:"title" sql:"not null;type:varchar(255)"`      //
	Description  string     `json:"description" sql:"null;type:text"`            //
	Event        int64      `json:"event" sql:"null;default:0"`                  //
	Manifest     string     `json:"manifest" sql:"null;type:longtext"`           //
	Environment  string     `json:"environment" sql:"null;type:text"`            // Environment parameters.
	Kubernetes   string     `json:"kubernetes" sql:"null;type:text"`             //
	Swarm        string     `json:"swarm" sql:"null;type:text"`                  //
	Input        string     `json:"input" sql:"null;type:text"`                  //
	Output       string     `json:"input" sql:"null;type:text"`                  //
	Endpoint     string     `json:"endpoint"`                                    //
	Timeout      int64      `json:"timeout"`                                     //
	Requires     string     `json:"requires" sql:"type:longtext"`                // pipeline run requires auth
	AuthList     string     `json:"authList" sql:"type:longtext"`                //
	CreatedAt    time.Time  `json:"created" sql:""`                              //
	UpdatedAt    time.Time  `json:"updated" sql:""`                              //
	DeletedAt    *time.Time `json:"deleted" sql:"index"`                         //
}

//TableName is return the name of Action in MySQL database.
func (a *ActionLog) TableName() string {
	return "action_log"
}

func (a *ActionLog) GetActionLog() *gorm.DB {
	return db.Model(&ActionLog{})
}

//Outcome is Stage running results.
//When StageID point to the StageTypeStart , the Action ID is 0.
//When StageID point to the StageTypeEnd , the Action ID is -1.
type Outcome struct {
	ID           int64      `json:"id" gorm:"primary_key"`                 //
	Pipeline     int64      `json:"pipeline" sql:"not null;default:0"`     //PipelineLog id
	RealPipeline int64      `json:"realPipeline" sql:"not null;default:0"` //Pipeline id
	Stage        int64      `json:"stage" sql:"not null;default:0"`        //stageLog id
	RealStage    int64      `json:"realStage" sql:"not null;default:0"`    //stage id
	Action       int64      `json:"action" sql:"not null;default:0"`       //actionLog id
	RealAction   int64      `json:"realAction" sql:"not null;default:0"`   //
	Event        int64      `json:"event" sql:"null;default:0"`            //event id
	Sequence     int64      `json:"sequence" sql:"not null;default:0"`     //pipeline run sequence
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

func (o *Outcome) GetOutcome() *gorm.DB {
	return db.Model(&Outcome{})
}

// pipelineSequence is a table describe pipeline's run sequence
type PipelineSequence struct {
	ID        int64      `json:"id" gorm:"primary_key"`             //
	Pipeline  int64      `json:"pipeline" sql:"not null;default:0"` // pipeline name
	Sequence  int64      `json:"sequence" sql:"not null;default:0"` //pipeline run sequence
	CreatedAt time.Time  `json:"created" sql:""`                    //
	UpdatedAt time.Time  `json:"updated" sql:""`                    //
	DeletedAt *time.Time `json:"deleted" sql:"index"`               //
}

//TableName is return the table name of Pipeline in MySQL database.
func (p *PipelineSequence) TableName() string {
	return "pipeline_sequence"
}

func (p *PipelineSequence) GetPipelineSequence() *gorm.DB {
	return db.Model(&PipelineSequence{})
}
