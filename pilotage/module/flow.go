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
	"regexp"
	"time"

	"bufio"
	"github.com/mitchellh/colorstring"
	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	apiv1 "k8s.io/client-go/pkg/api/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/tools/clientcmd"
)

// ExecuteFlowFromFile
func (f *Flow) ExecuteFlowFromFile(flowFile string, verbose, timestamp bool) error {
	f.Model = RunModelCli

	if data, err := ioutil.ReadFile(flowFile); err != nil {
		colorstring.Println(fmt.Sprintf("[red]Read orchestration flow file %s error: %s", flowFile, err.Error()))
		return err
	} else {
		if err := yaml.Unmarshal(data, &f); err != nil {
			colorstring.Println(fmt.Sprintf("[red]Unmarshal the flow file error: %s", err.Error()))
			return err
		} else {
			f.LocalRun(verbose, timestamp)
		}
	}

	return nil
}

// LocalRun
func (f *Flow) LocalRun(verbose, timestamp bool) error {
	// Print Flow Title
	colorstring.Println(fmt.Sprintf("[magenta]The [light_green]\"%s\" [magenta]is running:", f.Title))
	colorstring.Println("")

	colorstring.Println(fmt.Sprintf("[magenta]Version Number: [cyan]%d", f.Version))
	colorstring.Println(fmt.Sprintf("[magenta]Tag: [cyan]%s", f.Tag))
	if f.Timeout == 0 {
		colorstring.Println("[magenta]Timeout: [cyan]Unlimited")
	} else {
		colorstring.Println(fmt.Sprintf("[magenta]Timeout: [cyan]%d", f.Timeout))
	}

	colorstring.Println("")

	for key, stage := range f.Stages {
		colorstring.Println(fmt.Sprintf("[magenta]Stage Number: [cyan]%d", key))

		switch stage.T {
		case StageTypeStart:

			colorstring.Println("[magenta]Start Stage: [cyan]cli mode don't need trigger.")

		case StageTypeNormal:
			colorstring.Println(fmt.Sprintf("[magenta]Stage: [cyan]%s", stage.Title))
			colorstring.Println(fmt.Sprintf("[magenta]Stage Sequencing: [cyan]%s", stage.Sequencing))
			colorstring.Println("")

			switch stage.Sequencing {
			case StageTypeParallel:
				if err := stage.ParallelRun(verbose, timestamp); err != nil {
					return err
				}
			case StageTypeSequencing:
				if err := stage.SequencingRun(verbose, timestamp); err != nil {
					return err
				}
			default:
				return fmt.Errorf("Unknown sequencing type.")
			}
		case StageTypePause:
			//fmt.Println(fmt.Sprintf("Pause Stage: %s", stage.Title))
		case StageTypeEnd:
			colorstring.Println("[magenta]End Stage: [cyan]Flow execute end.")
		}

		colorstring.Println("")
		colorstring.Println("")
	}

	return nil
}

// SequencingRun
func (s *Stage) SequencingRun(verbose, timestamp bool) error {

	for key, action := range s.Actions {
		colorstring.Println(fmt.Sprintf("[magenta]\tAction Number: [cyan]%d", key))
		colorstring.Println(fmt.Sprintf("[magenta]\tAction: [cyan]%s", action.Title))

		if result, err := action.Run(verbose, timestamp); err != nil {
			return err
		} else if result == false {
			colorstring.Println(fmt.Sprintf("[magenta]Job End: [cyan]%s", time.Now().String()))
			colorstring.Println(fmt.Sprintf("[magenta]Action Result: [cyan]%s", "failure"))
		} else if result == true {
			colorstring.Println(fmt.Sprintf("[magenta]Job End: [cyan]%s", time.Now().String()))
			colorstring.Println(fmt.Sprintf("[magenta]Action Result: [cyan]%s", "successfully"))
		}

		colorstring.Println("")
	}

	return nil
}

// ParallelRun
func (s *Stage) ParallelRun(verbose, timestamp bool) error {
	return nil
}

func (a *Action) Run(verbose, timestamp bool) (bool, error) {

	for key, job := range a.Jobs {
		colorstring.Println(fmt.Sprintf("[magenta]\t\tJob: [cyan]%d", key))
		colorstring.Println(fmt.Sprintf("[magenta]\t\tJob Type: [cyan]%s", job.T))
		colorstring.Println(fmt.Sprintf("[magenta]\t\tJob Verbose: [cyan]%t", verbose))
		colorstring.Println(fmt.Sprintf("[magenta]\t\tJob Timestamp: [cyan]%t", timestamp))
		colorstring.Println(fmt.Sprintf("[magenta]\t\tJob Start: [cyan]%s", time.Now().String()))
		colorstring.Println("[magenta]\t\tJob Running:")

		if result, err := job.Run(a.Name, verbose, timestamp); err != nil {
			return false, err
		} else {
			return result, nil
		}

	}

	return true, nil
}

func (j *Job) Run(name string, verbose, timestamp bool) (bool, error) {
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
				Timestamps: timestamp,
			})

			if read, err := req.Stream(); err != nil {
				return false, err
			} else {

				reader := bufio.NewReader(read)
				for {
					line, err := reader.ReadString('\n')

					if err != nil {
						if err == io.EOF {
							colorPrint(line, verbose)
							break
						}
						return false, nil
					}

					colorPrint(line, verbose)
				}
			}

		}
	}

	return true, nil
}

func colorPrint(line string, verbose bool) {
	if has, _ := regexp.Match("CO_RESULT = false", []byte(line)); has == true {
		colorstring.Print(fmt.Sprintf("[red]\t\t\t %s", line))
	} else {
		colorstring.Print(fmt.Sprintf("[green]\t\t\t %s", line))
	}
}
