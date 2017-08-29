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
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/Huawei/containerops/pilotage/model"
	. "github.com/logrusorgru/aurora"
	"gopkg.in/yaml.v2"
)

const (
	// Flow Run Model
	CliRun      = "CliRun"
	DaemonRun   = "DaemonRun"
	DaemonStart = "DaemonStart"
)

// Flow is DevOps orchestration flow struct.
type Flow struct {
	Model   string   `json:"-" yaml:"-"`
	URI     string   `json:"uri" yaml:"uri"`
	Number  int64    `json:",omitempty" yaml:",omitempty"`
	Title   string   `json:"title" yaml:"title"`
	Version int64    `json:"version" yaml:"version"`
	Tag     string   `json:"tag" yaml:"tag"`
	Timeout int64    `json:"timeout" yaml:"timeout"`
	Status  string   `json:"status,omitempty" yaml:"status,omitempty"`
	Logs    []string `json:"logs,omitempty" yaml:"logs,omitempty"`
	Stages  []Stage  `json:"stages,omitempty" yaml:"stages,omitempty"`
}

// JSON export flow data without
func (f *Flow) JSON() ([]byte, error) {
	return json.Marshal(&f)
}

//
func (f *Flow) YAML() ([]byte, error) {
	return yaml.Marshal(f)
}

//
func (f *Flow) URIs() (namespace, repository, name string, err error) {
	array := strings.Split(f.URI, "/")
	if len(array) != 3 {
		return "", "", "", fmt.Errorf("Invalid flow URI: %s", f.URI)
	}

	namespace, repository, name = array[0], array[1], array[2]
	return namespace, repository, name, nil
}

// TODO filter the log print with different color.
func (f *Flow) Log(log string, verbose, timestamp bool) {
	f.Logs = append(f.Logs, fmt.Sprintf("[%s] %s", time.Now().String(), log))
	l := new(model.LogV1)
	//TODO fill in phaseID
	l.Create(model.INFO, model.FLOW, 0, log)

	if verbose == true {
		if timestamp == true {
			fmt.Println(Cyan(fmt.Sprintf("[%s] %s", time.Now().String(), strings.TrimSpace(log))))
		} else {
			fmt.Println(Cyan(log))
		}
	}
}

// ParseFlowFromFile is init flow definition from a file.
// It's only used in CliRun or DaemonRun, and run with local kubectl.
func (f *Flow) ParseFlowFromFile(flowFile, runMode string, verbose, timestamp bool) error {
	// Init flow properties
	f.Model, f.Number, f.Status = runMode, 1, Pending

	if data, err := ioutil.ReadFile(flowFile); err != nil {
		f.Log(fmt.Sprintf("Read orchestration flow file %s error: %s", flowFile, err.Error()), verbose, timestamp)
		return err
	} else {
		if err := yaml.Unmarshal(data, &f); err != nil {
			f.Log(fmt.Sprintf("Unmarshal the flow file error: %s", err.Error()), verbose, timestamp)
			return err
		}
	}

	return nil
}

// LocalRun is run flow using Kubectl in the local.
func (f *Flow) LocalRun(verbose, timestamp bool) error {
	f.Status = Running
	f.Log(fmt.Sprintf("Flow [%s] status change to %s", f.URI, f.Status), verbose, timestamp)

	for i, _ := range f.Stages {
		stage := &f.Stages[i]

		f.Log(fmt.Sprintf("The Number [%d] stage is running: %s", i, stage.Title), verbose, timestamp)

		switch stage.T {
		case StartStage:
			f.Log("Start stage don't need any trigger in cli or daemon run mode.", verbose, timestamp)
		case NormalStage:
			switch stage.Sequencing {
			case Parallel:
				if status, err := stage.ParallelRun(verbose, timestamp, f); err != nil {
					f.Status = Failure
					f.Log(fmt.Sprintf("Stage [%s] run error: %s", stage.Name, err.Error()), verbose, timestamp)
				} else {
					f.Status = status
				}
			case Sequencing:
				if status, err := stage.SequencingRun(verbose, timestamp, f); err != nil {
					f.Status = Failure
					f.Log(fmt.Sprintf("Stage [%s] run error: %s", stage.Name, err.Error()), verbose, timestamp)
				} else {
					f.Status = status
				}
			default:
				f.Status = Failure
				f.Log(fmt.Sprintf("Stage [%s] has unknown sequencing type: %s", stage.Name, stage.T), verbose, timestamp)
			}
		case PauseStage:
			// TODO Pause running
		case EndStage:
			f.Log("End stage don't trigger any other flow.", verbose, timestamp)
		}

		// if status is failure or cancel, break the for loop.
		if f.Status == Failure || f.Status == Cancel {
			break
		}
	}

	return nil
}
