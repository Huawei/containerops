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
	"regexp"
	"strings"
	"time"

	. "github.com/logrusorgru/aurora"
	"gopkg.in/yaml.v2"

	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	apiv1 "k8s.io/client-go/pkg/api/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/tools/clientcmd"
)

//
func (f *Flow) JSON() ([]byte, error) {
	return json.Marshal(f)
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

// ExecuteFlowFromFile is init flow definition from a file.
// It's only used in CliRun or DaemonRun, and run with local kubectl.
func (f *Flow) ExecuteFlowFromFile(flowFile, runMode string, verbose, timestamp bool) error {
	// Init flow properties
	f.Model, f.Number, f.Status = runMode, 1, Pending

	if data, err := ioutil.ReadFile(flowFile); err != nil {
		fmt.Println(Red(fmt.Sprintf("[red]Read orchestration flow file %s error: %s", flowFile, err.Error())))
		return err
	} else {
		if err := yaml.Unmarshal(data, &f); err != nil {
			fmt.Println(Red(fmt.Sprintf("[red]Unmarshal the flow file error: %s", err.Error())))
			return err
		} else {
			f.LocalRun(verbose, timestamp)
		}
	}

	return nil
}

// LocalRun is run flow using Kubectl in the local.
func (f *Flow) LocalRun(verbose, timestamp bool) error {
	// Print Flow Title
	fmt.Println(fmt.Sprintf("[magenta]The [light_green]\"%s\" [magenta]is running:", f.Title))
	fmt.Println("")

	fmt.Println(fmt.Sprintf("[magenta]Version Number: [cyan]%d", f.Version))
	fmt.Println(fmt.Sprintf("[magenta]Tag: [cyan]%s", f.Tag))
	if f.Timeout == 0 {
		fmt.Println("[magenta]Timeout: [cyan]Unlimited")
	} else {
		fmt.Println(fmt.Sprintf("[magenta]Timeout: [cyan]%d", f.Timeout))
	}

	fmt.Println("")

	for key, stage := range f.Stages {
		fmt.Println(fmt.Sprintf("[magenta]Stage Number: [cyan]%d", key))

		switch stage.T {
		case StartStage:

			fmt.Println("[magenta]Start Stage: [cyan]cli mode don't need trigger.")

		case NormalStage:
			fmt.Println(fmt.Sprintf("[magenta]Stage: [cyan]%s", stage.Title))
			fmt.Println(fmt.Sprintf("[magenta]Stage Sequencing: [cyan]%s", stage.Sequencing))
			fmt.Println("")

			switch stage.Sequencing {
			case Parallel:
				if err := stage.ParallelRun(verbose, timestamp); err != nil {
					return err
				}
			case Sequencing:
				if err := stage.SequencingRun(verbose, timestamp); err != nil {
					return err
				}
			default:
				return fmt.Errorf("Unknown sequencing type.")
			}
		case PauseStage:
			//fmt.Println(fmt.Sprintf("Pause Stage: %s", stage.Title))
		case EndStage:
			fmt.Println("[magenta]End Stage: [cyan]Flow execute end.")
		}

		fmt.Println("")
		fmt.Println("")
	}

	return nil
}

// SequencingRun
func (s *Stage) SequencingRun(verbose, timestamp bool) error {

	for key, action := range s.Actions {
		fmt.Println(fmt.Sprintf("[magenta]\tAction Number: [cyan]%d", key))
		fmt.Println(fmt.Sprintf("[magenta]\tAction: [cyan]%s", action.Title))

		if result, err := action.Run(verbose, timestamp); err != nil {
			return err
		} else if result == false {
			fmt.Println(fmt.Sprintf("[magenta]Job End: [cyan]%s", time.Now().String()))
			fmt.Println(fmt.Sprintf("[magenta]\tAction Result: [cyan]%s", "failure"))
		} else if result == true {
			fmt.Println(fmt.Sprintf("[magenta]Job End: [cyan]%s", time.Now().String()))
			fmt.Println(fmt.Sprintf("[magenta]\tAction Result: [cyan]%s", "successfully"))
		}

		fmt.Println("")
	}

	return nil
}

// ParallelRun
func (s *Stage) ParallelRun(verbose, timestamp bool) error {
	return nil
}

func (a *Action) Run(verbose, timestamp bool) (bool, error) {

	for key, job := range a.Jobs {
		fmt.Println(fmt.Sprintf("[magenta]\t\tJob: [cyan]%d", key))
		fmt.Println(fmt.Sprintf("[magenta]\t\tJob Type: [cyan]%s", job.T))
		fmt.Println(fmt.Sprintf("[magenta]\t\tJob Verbose: [cyan]%t", verbose))
		fmt.Println(fmt.Sprintf("[magenta]\t\tJob Timestamp: [cyan]%t", timestamp))
		fmt.Println(fmt.Sprintf("[magenta]\t\tJob Start: [cyan]%s", time.Now().String()))
		fmt.Println("[magenta]\t\tJob Running:")

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
		fmt.Print(fmt.Sprintf("[red]\t\t\t %s", line))
	} else {
		fmt.Print(fmt.Sprintf("[green]\t\t\t %s", line))
	}
}
