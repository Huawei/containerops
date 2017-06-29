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

const (
	// Flow Run Model
	RunModelCli    = "CLI"
	RunModelDaemon = "DAEMON"

	// Stage Type
	StageTypeStart  = "start"
	StageTypeEnd    = "end"
	StageTypeNormal = "normal"
	StageTypePause  = "pause"

	// Action Type
	StageTypeSequencing = "sequence"
	StageTypeParallel   = "parallel"
)

// Flow is DevOps orchestration flow struct.
type Flow struct {
	Model   string  `json:"model" yaml:"model"`
	Name    string  `json:"name" yaml:"name"`
	Title   string  `json:"title" yaml:"title"`
	Version int64   `json:"version" yaml:"version"`
	Tag     string  `json:"tag" yaml:"tag"`
	Timeout int64   `json:"timeout" yaml:"timeout"`
	Stages  []Stage `json:"stages" yaml:"stages"`
}

// Stage is
type Stage struct {
	T          string   `json:"type" yaml:"type"`
	Name       string   `json:"name" yaml:"name"`
	Title      string   `json:"title" yaml:"title"`
	Sequencing string   `json:"sequencing" yaml:"sequencing"`
	Actions    []Action `json:"actions" yaml:"actions"`
}

// Action is
type Action struct {
	Name  string `json:"name" yaml:"name"`
	Title string `json:"title" yaml:"title"`
	Jobs  []Job  `json:"jobs" yaml:"jobs"`
}

// Job is
type Job struct {
	T            string              `json:"type" yaml:"type"`
	Kubectl      string              `json:"kubectl" yaml:"kubectl"`
	Endpoint     string              `json:"endpoint" yaml:"endpoint"`
	Timeout      string              `json:"timeout" yaml:"timeout"`
	Resources    Resource            `json:"resources" yaml:"resources"`
	Environments []map[string]string `json:"environments" yaml:"environments"`
	Output       []string            `json:"output" yaml:"output"`
}

// Resources is
type Resource struct {
	CPU    string `json:"cpu" yaml:"cpu"`
	Memory string `json:"memory" yaml:"memory"`
}
