package module

import (
	"time"
	"bufio"
	"io"
	"fmt"
	"strings"

	. "github.com/logrusorgru/aurora"
	homeDir "github.com/mitchellh/go-homedir"

	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/Huawei/containerops/common/utils")

// Job is
type Job struct {
	T             string              `json:"type" yaml:"type"`
	Kubectl       string              `json:"kubectl" yaml:"kubectl"`
	Endpoint      string              `json:"endpoint" yaml:"endpoint"`
	Timeout       string              `json:"timeout" yaml:"timeout"`
	Status        string              `json:"status,omitempty" yaml:"status,omitempty"`
	Resources     Resource            `json:"resources" yaml:"resources"`
	Logs          []string            `json:"logs,omitempty" yaml:"logs,omitempty"`
	Environments  []map[string]string `json:"environments" yaml:"environments"`
	Outputs       []map[string]string `json:"outputs,omitempty" yaml:"outputs,omitempty"`
	Subscriptions []map[string]string `json:"subscriptions,omitempty" yaml:"subscriptions,omitempty"`
}

// Resources is
type Resource struct {
	CPU    string `json:"cpu" yaml:"cpu"`
	Memory string `json:"memory" yaml:"memory"`
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

func (j *Job) Run(name string, verbose, timestamp bool, f *Flow) (string, error) {
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

					j.Log(line, false, timestamp)
					f.Log(line, verbose, timestamp)
				}
			}
		}
	}

	j.Status = Success
	return Success, nil
}

