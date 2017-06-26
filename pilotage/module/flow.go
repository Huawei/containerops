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
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"

	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	apiv1 "k8s.io/client-go/pkg/api/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/tools/clientcmd"
)

// ExecuteFlowFromFile
func (f *Flow) ExecuteFlowFromFile(flowFile string, verbose bool) error {
	f.Model = RunModelCli

	if data, err := ioutil.ReadFile(flowFile); err != nil {
		fmt.Println("Read ", flowFile, " error: ", err.Error())
		return err
	} else {
		if err := yaml.Unmarshal(data, &f); err != nil {
			fmt.Println("Unmarshal the flow file error:", err.Error())
			return err
		} else {
			f.LocalRun(verbose)
		}
	}

	return nil
}

// LocalRun
func (f *Flow) LocalRun(verbose bool) error {
	for _, stage := range f.Stages {
		switch stage.T {
		case StageTypeStart:
			fmt.Println(fmt.Sprintf("Start Stage: %s", stage.Title))
			fmt.Println(fmt.Sprintf("Running Orchestration flow now."))
		case StageTypeNormal:
			fmt.Println(fmt.Sprintf("Normal Stage: %s", stage.Title))

			switch stage.Sequencing {
			case StageTypeParallel:
				if err := stage.ParallelRun(); err != nil {
					return err
				}
			case StageTypeSequencing:
				if err := stage.SequencingRun(); err != nil {
					return err
				}
			default:
				return fmt.Errorf("Unknown sequencing type.")
			}
		case StageTypePause:
			fmt.Println(fmt.Sprintf("Pause Stage: %s", stage.Title))
		case StageTypeEnd:
			fmt.Println(fmt.Sprintf("End Stage: %s", stage.Title))
		}
	}

	return nil
}

// SequencingRun
func (s *Stage) SequencingRun() error {
	for _, action := range s.Actions {
		fmt.Println(action.Name)
		if result, err := action.Run(); err != nil {
			return err
		} else if result == false {
			return fmt.Errorf("Action execute error.")
		} else if result == true {
			fmt.Println("Action execute successfully.")
		}
	}

	return nil
}

// ParallelRun
func (s *Stage) ParallelRun() error {
	return nil
}

func (a *Action) Run() (bool, error) {
	for _, job := range a.Jobs {
		job.Run(a.Name)
	}

	return true, nil
}

func (j *Job) Run(name string) (bool, error) {
	if config, err := clientcmd.BuildConfigFromFlags("", "/Users/meaglith/.kube/config"); err != nil {
		return false, err
	} else {
		if clientSet, err := kubernetes.NewForConfig(config); err != nil {
			return false, err
		} else {
			p := clientSet.CoreV1Client.Pods(apiv1.NamespaceDefault)

			if _, err := p.Create(
				&apiv1.Pod{
					TypeMeta: metav1.TypeMeta{
						Kind:       "Pod",
						APIVersion: "v1",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name: name,
					},
					Spec: apiv1.PodSpec{
						Containers: []apiv1.Container{
							{
								Name:  name,
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
				return false, err
			}

			time.Sleep(time.Second * 5)

			req := p.GetLogs(name, &apiv1.PodLogOptions{
				Follow:     true,
				Timestamps: true,
			})

			if read, err := req.Stream(); err != nil {
				return false, err
			} else {
				io.Copy(os.Stdout, read)
			}

		}
	}

	return true, nil
}
