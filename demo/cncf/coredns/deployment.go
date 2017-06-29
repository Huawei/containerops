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

package main

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"time"

	"github.com/mitchellh/colorstring"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	apiv1 "k8s.io/client-go/pkg/api/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {

	if config, err := clientcmd.BuildConfigFromFlags("", "/Users/meaglith/.kube/config"); err != nil {
		panic(err.Error())
	} else {
		if clientSet, err := kubernetes.NewForConfig(config); err != nil {
			panic(err.Error())
		} else {
			p := clientSet.CoreV1Client.Pods(apiv1.NamespaceDefault)

			if _, err := p.Create(
				&apiv1.Pod{
					TypeMeta: metav1.TypeMeta{
						Kind:       "Pod",
						APIVersion: "v1",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name: "cncf-demo-coredns-release",
					},
					Spec: apiv1.PodSpec{
						Containers: []apiv1.Container{
							{
								Name:  "cncf-demo-coredns-release",
								Image: "docker.io/containerops/cncf-demo-coredns:latest",
								Env: []apiv1.EnvVar{
									{
										Name:  "CO_DATA",
										Value: "coredns=https://github.com/coredns/coredns.git action=release release=test.opshub.sh/containerops/cncf-demo/demo",
									},
								},
								Resources: apiv1.ResourceRequirements{
									Requests: apiv1.ResourceList{
										apiv1.ResourceCPU:    resource.MustParse("2"),
										apiv1.ResourceMemory: resource.MustParse("4G"),
									},
								},
							},
						},
						RestartPolicy: apiv1.RestartPolicyNever,
					},
				},
			); err != nil {
				panic(err.Error())
			}

			time.Sleep(time.Second * 5)

			req := p.GetLogs("cncf-demo-coredns-release", &apiv1.PodLogOptions{
				Follow:     true,
				Timestamps: false,
			})

			if read, err := req.Stream(); err != nil {
				panic(err.Error())
			} else {

				reader := bufio.NewReader(read)
				for {
					line, err := reader.ReadString('\n')
					if err != nil {
						if err == io.EOF {
							if has, _ := regexp.Match("[COUT]", []byte(line)); has == true {
								if has, _ := regexp.Match("CO_RESULT = false", []byte(line)); has == true {
									colorstring.Print(fmt.Sprintf("[red] %s", line))
								} else {
									colorstring.Print(fmt.Sprintf("[green] %s", line))
								}
							}

							break
						}
						panic(err)
					}

					if has, _ := regexp.Match("[COUT]", []byte(line)); has == true {
						if has, _ := regexp.Match("CO_RESULT = false", []byte(line)); has == true {
							colorstring.Print(fmt.Sprintf("[red] %s", line))
						} else {
							colorstring.Print(fmt.Sprintf("[green] %s", line))
						}
					}

				}
			}

		}
	}
}
