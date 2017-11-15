package module

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
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

var (
	RWlock        sync.RWMutex
	GlobalOutputs map[string]string
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
	Outputs       []string            `json:"outputs,omitempty" yaml:"outputs,omitempty"`
	Subscriptions []map[string]string `json:"subscriptions,omitempty" yaml:"subscriptions,omitempty"`
}

// Resources is
type Resource struct {
	CPU    string `json:"cpu" yaml:"cpu"`
	Memory string `json:"memory" yaml:"memory"`
}

func init() {
	GlobalOutputs = make(map[string]string)
}

// TODO filter the log print with different color.
func (j *Job) Log(log string, verbose, timestamp bool) {
	j.Logs = append(j.Logs, fmt.Sprintf("[%s] %s", time.Now().String(), log))
	l := new(model.LogV1)
	l.Create(model.INFO, model.JOB, j.ID, log)

	if verbose == true {
		if timestamp == true {
			fmt.Println(Cyan(fmt.Sprintf("[%s] %s", time.Now().String(), strings.TrimSpace(log))))
		} else {
			fmt.Println(Cyan(log))
		}
	}
}

func (j *Job) Run(name string, verbose, timestamp bool, f *Flow, stageIndex, actionIndex int) (string, error) {

	j.SaveDatabase(verbose, timestamp, f, stageIndex, actionIndex)

	randomContainerName := fmt.Sprintf("%s-%s", name, utils.RandomString(10))
	podTemplate := j.PodTemplates(randomContainerName, f)

	if err := j.InvokePod(podTemplate, randomContainerName, verbose, timestamp, f, stageIndex, actionIndex); err != nil {
		return Failure, err
	}

	j.Status = Success

	return Success, nil
}

func (j *Job) RunKubectl(name string, verbose, timestamp bool, f *Flow, stageIndex, actionIndex int) (string, error) {

	j.SaveDatabase(verbose, timestamp, f, stageIndex, actionIndex)

	originYaml := []byte{}
	if u, err := url.Parse(j.Kubectl); err != nil {
		return Failure, err
	} else {
		if u.Scheme == "" {
			if utils.IsFileExist(j.Kubectl) == true {
				// Read YAML file from local
				data, err := ioutil.ReadFile(j.Kubectl)
				if err != nil {
					return Failure, err
				}
				originYaml = data
			} else {
				return Failure, errors.New("Kubectl PATH is invalid")
			}
		} else {
			// Download YAML from URL
			resp, err := http.Get(j.Kubectl)
			if err != nil {
				return Failure, err
			}
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return Failure, err
			}
			originYaml = body
		}
	}
	base64Yaml := base64.StdEncoding.EncodeToString(originYaml)

	//TODO port and ip address can set from setting
	home, _ := homeDir.Dir()
	configFile, err := clientcmd.BuildConfigFromFlags("", fmt.Sprintf("%s/.kube/config", home))
	if err != nil {
		return Failure, err
	}
	apiServerInsecure := fmt.Sprintf("http:%s:8080", strings.Split(configFile.Host, ":")[1])
	namespace := "default"
	if f.Namespace != "" {
		namespace = f.Namespace
	}
	randomContainerName := fmt.Sprintf("kubectl-create-%s", utils.RandomString(10))
	podTemplate := j.KubectlPodTemplates(randomContainerName, apiServerInsecure, namespace, base64Yaml, f)

	if err := j.InvokePod(podTemplate, randomContainerName, verbose, timestamp, f, stageIndex, actionIndex); err != nil {
		return Failure, err
	}

	j.Status = Success
	return Success, nil
}

func (j *Job) InvokePod(podTemplate *apiv1.Pod, randomContainerName string, verbose, timestamp bool, f *Flow, stageIndex, actionIndex int) error {
	home, _ := homeDir.Dir()
	if config, err := clientcmd.BuildConfigFromFlags("", fmt.Sprintf("%s/.kube/config", home)); err != nil {
		return err
	} else {
		if clientSet, err := kubernetes.NewForConfig(config); err != nil {
			return err
		} else {
			p := clientSet.CoreV1().Pods(apiv1.NamespaceDefault)
			if _, err := p.Create(podTemplate); err != nil {
				j.Status = Failure
				return err
			}

			j.Status = Pending
			time.Sleep(time.Second * 2)

			start := time.Now()
		ForLoop:
			for {
				pod, err := p.Get(randomContainerName, metav1.GetOptions{})
				if err != nil {
					j.Log(err.Error(), false, timestamp)
					return err
				}
				switch pod.Status.Phase {
				case apiv1.PodPending:
					j.Log(fmt.Sprintf("Job %s is %s", j.Name, pod.Status.Phase), verbose, timestamp)
				case apiv1.PodRunning, apiv1.PodSucceeded:
					break ForLoop
				case apiv1.PodUnknown:
					j.Log(fmt.Sprintf("Job %s is %s, Detail:[%s] \n", j.Name, pod.Status.Phase, pod.Status.ContainerStatuses[0].State.String()), verbose, timestamp)
				case apiv1.PodFailed:
					j.Log(fmt.Sprintf("Job %s is %s, Detail:[%s] \n", j.Name, pod.Status.Phase, pod.Status.ContainerStatuses[0].State.String()), verbose, timestamp)
					break ForLoop
				}
				duration := time.Now().Sub(start)
				if duration.Minutes() > 3 {
					return errors.New(fmt.Sprintf("Job %s Pending more than 3 minutes", j.Name))
				}
				time.Sleep(time.Second * 2)
			}

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
						return err
					}
					if strings.Contains(line, "[COUT]") && len(j.Outputs) != 0 {
						j.FetchOutputs(f.Stages[stageIndex].Name, f.Stages[stageIndex].Actions[actionIndex].Name, line)
					}

					j.Status = Running

					j.Log(line, false, timestamp)
					f.Log(line, verbose, timestamp)
				}
			}
		}
	}
	return nil
}

func (j *Job) SaveDatabase(verbose, timestamp bool, f *Flow, stageIndex, actionIndex int) {
	// Save Job into database
	job := new(model.JobV1)
	resources, _ := j.Resources.JSON()
	environments, _ := json.Marshal(j.Environments)
	outputs, _ := json.Marshal(j.Outputs)
	subscriptions, _ := json.Marshal(j.Subscriptions)
	jobID, err := job.Put(f.Stages[stageIndex].Actions[actionIndex].ID, j.Timeout, j.Name, j.T, j.Endpoint, string(resources), string(environments), string(outputs), string(subscriptions))
	if err != nil {
		j.Log(fmt.Sprintf("Save Job [%s] errorK: %s", j.Name, err.Error()), false, timestamp)
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
}

func (j *Job) FetchOutputs(stageName, actionName, log string) error {
	output := strings.TrimPrefix(log, "[COUT]")
	splits := strings.Split(output, "=")
	for _, o := range j.Outputs {
		if strings.TrimSpace(o) == strings.TrimSpace(splits[0]) {
			key := fmt.Sprintf("%s.%s.%s[%s]", stageName, actionName, j.Name, o)
			RWlock.Lock()
			GlobalOutputs[key] = strings.TrimSpace(splits[1])
			RWlock.Unlock()
		}
	}
	return nil
}

func (j *Job) KubectlPodTemplates(randomContainerName, apiServer, namespace, yamlContent string, f *Flow) *apiv1.Pod {
	result := &apiv1.Pod{
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
					Name: randomContainerName,
					//TODO can config from settings
					Image: "hub.opshub.sh/containerops/kubectl-create:1.7.4",
				},
			},
			RestartPolicy: apiv1.RestartPolicyNever,
		},
	}
	//Add api-server address, namespace & yaml content
	coDataValue := fmt.Sprintf(" api-server-url=%s namespace=%s", apiServer, namespace)
	result.Spec.Containers[0].Env = append(result.Spec.Containers[0].Env, apiv1.EnvVar{Name: "CO_DATA", Value: coDataValue})
	result.Spec.Containers[0].Env = append(result.Spec.Containers[0].Env, apiv1.EnvVar{Name: "YAML", Value: yamlContent})

	//Add user defined enviroments
	if len(j.Environments) > 0 {
		for _, environment := range j.Environments {
			for k, v := range environment {
				env := apiv1.EnvVar{
					Name:  k,
					Value: v,
				}
				result.Spec.Containers[0].Env = append(result.Spec.Containers[0].Env, env)
			}
		}
	}

	//Add flow enviroments
	if len(f.Environments) > 0 {
		for _, environment := range f.Environments {
			for k, v := range environment {
				env := apiv1.EnvVar{
					Name:  k,
					Value: v,
				}
				result.Spec.Containers[0].Env = append(result.Spec.Containers[0].Env, env)
			}
		}
	}

	return result
}

func (j *Job) PodTemplates(randomContainerName string, f *Flow) *apiv1.Pod {
	result := &apiv1.Pod{
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
	}
	//Add user defined enviroments
	if len(j.Environments) > 0 {
		for _, environment := range j.Environments {
			for k, v := range environment {
				env := apiv1.EnvVar{
					Name:  k,
					Value: v,
				}
				result.Spec.Containers[0].Env = append(result.Spec.Containers[0].Env, env)
			}
		}
	}
	//Add flow enviroments
	if len(f.Environments) > 0 {
		for _, environment := range f.Environments {
			for k, v := range environment {
				env := apiv1.EnvVar{
					Name:  k,
					Value: v,
				}
				result.Spec.Containers[0].Env = append(result.Spec.Containers[0].Env, env)
			}
		}
	}

	//Add user defined subscrptions
	if len(j.Subscriptions) > 0 {
		for _, subscription := range j.Subscriptions {
			for k, env_key := range subscription {
				if env_value, ok := GlobalOutputs[k]; ok {
					env := apiv1.EnvVar{
						Name:  env_key,
						Value: env_value,
					}
					result.Spec.Containers[0].Env = append(result.Spec.Containers[0].Env, env)
				}
			}
		}
	}
	return result
}

func (r *Resource) JSON() ([]byte, error) {
	return json.Marshal(&r)
}
