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
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"time"

	. "github.com/logrusorgru/aurora"
	homeDir "github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v2"

	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/Huawei/containerops/common/utils"
)

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

	if verbose == true {
		if timestamp == true {
			fmt.Println(Cyan(fmt.Sprintf("[%s] %s", time.Now().String(), log)))
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
				// TODO Parallel running
			case Sequencing:
				if status, err := stage.SequencingRun(verbose, timestamp); err != nil {
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

// TODO filter the log print with different color.
func (s *Stage) Log(log string, verbose, timestamp bool) {
	s.Logs = append(s.Logs, fmt.Sprintf("[%s] %s", time.Now().String(), log))

	if verbose == true {
		if timestamp == true {
			fmt.Println(Cyan(fmt.Sprintf("[%s] %s", time.Now().String(), log)))
		} else {
			fmt.Println(Cyan(log))
		}
	}
}

func (s *Stage) SequencingRun(verbose, timestamp bool) (string, error) {
	s.Status = Running
	s.Log(fmt.Sprintf("Stage [%s] status change to %s", s.Name, s.Status), verbose, timestamp)

	for i, _ := range s.Actions {
		action := &s.Actions[i]

		s.Log(fmt.Sprintf("The Number [%d] action is running: %s", i, s.Title), verbose, timestamp)

		if status, err := action.Run(verbose, timestamp); err != nil {
			s.Status = Failure
			s.Log(fmt.Sprintf("Action [%s] run error: %s", action.Name, err.Error()), verbose, timestamp)
		} else {
			s.Status = status
		}

		if s.Status == Failure || s.Status == Cancel {
			break
		}
	}

	return s.Status, nil
}

// TODO filter the log print with different color.
func (a *Action) Log(log string, verbose, timestamp bool) {
	a.Logs = append(a.Logs, fmt.Sprintf("[%s] %s", time.Now().String(), log))

	if verbose == true {
		if timestamp == true {
			fmt.Println(Cyan(fmt.Sprintf("[%s] %s", time.Now().String(), log)))
		} else {
			fmt.Println(Cyan(log))
		}
	}
}

func (a *Action) Run(verbose, timestamp bool) (string, error) {
	a.Status = Running
	a.Log(fmt.Sprintf("Action [%s] status change to %s", a.Name, a.Status), verbose, timestamp)

	for i, _ := range a.Jobs {
		job := &a.Jobs[i]

		a.Log(fmt.Sprintf("The Number [%d] job is running: %s", i, a.Title), verbose, timestamp)

		if status, err := job.Run(a.Name, verbose, timestamp); err != nil {
			a.Status = Failure
			a.Log(fmt.Sprintf("Job [%d] run error: %s", i, err.Error()), verbose, timestamp)
		} else {
			a.Status = status
		}

		if a.Status == Failure || a.Status == Cancel {
			break
		}

	}

	return a.Status, nil
}

// TODO filter the log print with different color.
func (j *Job) Log(log string, verbose, timestamp bool) {
	j.Logs = append(j.Logs, fmt.Sprintf("[%s] %s", time.Now().String(), log))

	if verbose == true {
		if timestamp == true {
			fmt.Println(Cyan(fmt.Sprintf("[%s] %s", time.Now().String(), strings.TrimSpace(log))))
		} else {
			fmt.Println(Cyan(log))
		}
	}
}

func (j *Job) Run(name string, verbose, timestamp bool) (string, error) {
	home, _ := homeDir.Dir()

	randomContainerName := fmt.Sprintf("%s-%s", name, utils.RandomString(10))

	if config, err := clientcmd.BuildConfigFromFlags("", fmt.Sprintf("%s/.kube/config", home)); err != nil {
		return Failure, err
	} else {
		if clientSet, err := kubernetes.NewForConfig(config); err != nil {
			return Failure, err
		} else {
			p := clientSet.CoreV1Client.Pods(apiv1.NamespaceDefault)

			if _, err := p.Create(
				&apiv1.Pod{
					TypeMeta: metav1.TypeMeta{
						Kind:       "Pod",
						APIVersion: "v1",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name: randomContainerName,
					},
					Spec: apiv1.PodSpec{
						Containers: []apiv1.Container{
							{
								Name:  randomContainerName,
								Image: j.Endpoint,
								Env: []apiv1.EnvVar{
									{
										Name:  "CO_DATA",
										Value: j.Environments[0]["CO_DATA"],
									},
								},
								Resources: apiv1.ResourceRequirements{
									Requests: apiv1.ResourceList{
										apiv1.ResourceCPU:    resource.MustParse(j.Resources.CPU),
										apiv1.ResourceMemory: resource.MustParse(j.Resources.Memory),
									},
								},
							},
						},
						RestartPolicy: apiv1.RestartPolicyNever,
					},
				},
			); err != nil {
				j.Status = Failure
				return Failure, err
			}

			j.Status = Pending
			time.Sleep(time.Second * 5)

			req := p.GetLogs(randomContainerName, &apiv1.PodLogOptions{
				Follow:     true,
				Timestamps: false,
			})

			if read, err := req.Stream(); err != nil {
				// TODO Parse ContainerCreating error
			} else {
				reader := bufio.NewReader(read)
				for {
					line, err := reader.ReadString('\n')

					if err != nil {
						if err == io.EOF {
							break
						}

						j.Status = Failure
						return Failure, nil
					}
					j.Status = Running
					j.Log(line, verbose, timestamp)
				}
			}
		}
	}

	j.Status = Success
	return Success, nil
}
