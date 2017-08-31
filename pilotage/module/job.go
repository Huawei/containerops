package module

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	. "github.com/logrusorgru/aurora"
	homeDir "github.com/mitchellh/go-homedir"

	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/Huawei/containerops/common/utils"
	"github.com/Huawei/containerops/pilotage/model"
)

// Job is
type Job struct {
	ID            int64               `json:"-" yaml:"-"`
	T             string              `json:"type" yaml:"type"`
	Name          string              `json:"name" yaml:"name,omitempty"`
	Kubectl       string              `json:"kubectl" yaml:"kubectl"`
	Endpoint      string              `json:"endpoint" yaml:"endpoint"`
	Timeout       int64               `json:"timeout" yaml:"timeout"`
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
	l := new(model.LogV1)
	//TODO fill in phaseID
	l.Create(model.INFO, model.JOB, 0, log)

	if verbose == true {
		if timestamp == true {
			fmt.Println(Cyan(fmt.Sprintf("[%s] %s", time.Now().String(), strings.TrimSpace(log))))
		} else {
			fmt.Println(Cyan(log))
		}
	}
}

func (j *Job) Run(name string, verbose, timestamp bool, f *Flow, stageIndex, actionIndex int) (string, error) {

	// Save Job into database
	job := new(model.JobV1)
	resources, _ := j.Resources.JSON()
	environments, _ := json.Marshal(j.Environments)
	outputs, _ := json.Marshal(j.Outputs)
	subscriptions, _ := json.Marshal(j.Subscriptions)
	//tmpjson,_:=f.JSON()
	//fmt.Println(string(tmpjson))
	jobID, err := job.Put(f.Stages[stageIndex].Actions[actionIndex].ID, j.Timeout, j.Name, j.T, j.Endpoint, string(resources), string(environments), string(outputs), string(subscriptions))
	if err != nil {
		j.Log(fmt.Sprintf("Save Job [%s] error: %s", j.Name, err.Error()), false, timestamp)
	}
	j.ID = jobID

	// Record Job data
	jobData := new(model.JobDataV1)
	startTime := time.Now()

	defer func() {
		currentNumber, err := jobData.GetNumbers(j.ID)
		if err != nil {
			j.Log(fmt.Sprintf("Get Job Data [%s] Numbers error: %s", j.Name, err.Error()), verbose, timestamp)
		}
		if err := jobData.Put(j.ID, currentNumber+1, j.Status, startTime, time.Now()); err != nil {
			j.Log(fmt.Sprintf("Save Job Data [%s] error: %s", j.Name, err.Error()), false, timestamp)
		}
	}()

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

func (r *Resource) JSON() ([]byte, error) {
	return json.Marshal(&r)
}
