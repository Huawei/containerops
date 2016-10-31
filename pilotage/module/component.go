/*
Copyright 2014 Huawei Technologies Co., Ltd. All rights reserved.

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
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/Huawei/containerops/pilotage/models"
)

const (
	// RUNENV_KUBE means component will run in k8s
	RUNENV_KUBE = "KUBERNETES"
	// RUNENV_SWARM means component will run in swarm
	RUNENV_SWARM = "SWARM"
)

var (
	createComponentChan chan bool
)

type component interface {
	// start a component
	// id contains pipelineId stageId actionId pipelineSequence
	Start(id string, eventList []models.EventDefinition, envMap map[string]string) error
	Stop(string) (string, error)
	GetIp(...interface{}) (string, error)
}

func init() {
	createComponentChan = make(chan bool, 1)
}

func InitComponet(actionInfo models.ActionLog, runenv, host, namespace string) (component, error) {
	if actionInfo.Component == 0 {
		return nil, errors.New("component's id is 0")
	}

	if runenv == RUNENV_KUBE {
		kubeCom := new(kubeComponent)

		ComponentConfigMap := make(map[string]interface{}, 0)
		err := json.Unmarshal([]byte(actionInfo.Kubernetes), &ComponentConfigMap)
		if err != nil {
			return kubeCom, errors.New("component kube config error:" + err.Error())
		}

		// if _, ok := ComponentConfigMap["port"]; !ok {
		// 	return kubeCom, errors.New("component kube config error ,port is not set")
		// }
		// portF, ok := ComponentConfigMap["port"].(float64)
		// if !ok {
		// 	return kubeCom, errors.New("component kube config error ,port is not a int")
		// }
		// port := int64(portF)

		if _, ok := ComponentConfigMap["reachableIPs"]; !ok {
			return kubeCom, errors.New("component kube config error ,reachableIPs is not set")
		}
		reachableIPs, ok := ComponentConfigMap["reachableIPs"].([]interface{})
		if !ok {
			return kubeCom, errors.New("component kube config error ,reachableIPs is not a array")
		}
		reachableIPStrs := make([]string, 0)
		for _, ip := range reachableIPs {
			ipStr, ok := ip.(string)
			if !ok {
				return kubeCom, errors.New("reachable ip is not a string")
			}

			reachableIPStrs = append(reachableIPStrs, ipStr)
		}

		podConfig := make(map[string]interface{})
		if _, ok := ComponentConfigMap["podConfig"]; ok {
			podConfig, ok = ComponentConfigMap["podConfig"].(map[string]interface{})
			if !ok {
				return kubeCom, errors.New("component kube config error ,podConfig is not a json data")
			}
		}

		serviceConfig := make(map[string]interface{})
		if _, ok := ComponentConfigMap["serviceConfig"]; ok {
			serviceConfig, ok = ComponentConfigMap["serviceConfig"].(map[string]interface{})
			if !ok {
				return kubeCom, errors.New("component kube config error ,serviceConfig is not a json data")
			}
		}

		kubeCom.host = host
		// kubeCom.port = port
		kubeCom.namespace = namespace
		kubeCom.reachableIPs = reachableIPStrs
		kubeCom.podConfig = podConfig
		kubeCom.serviceConfig = serviceConfig
		kubeCom.componentInfo = actionInfo

		return kubeCom, nil
	}

	return nil, errors.New("can't create a component in" + runenv)
}

// GetComponentInfo is
func GetComponentInfo(namespace, componentName string, componentId int64) (map[string]interface{}, error) {
	resultMap := make(map[string]interface{})
	componentInfo := new(models.Component)
	err := componentInfo.GetComponent().Where("id = ?", componentId).First(&componentInfo).Error
	if err != nil {
		return nil, errors.New("error when get component info from db:" + err.Error())
	}

	if componentInfo.Component != componentName {
		return nil, errors.New("component's name is not equal to target component")
	}

	// get component define json first, if has a define json,return it
	if componentInfo.Manifest != "" {
		defineMap := make(map[string]interface{})
		json.Unmarshal([]byte(componentInfo.Manifest), &defineMap)
		if defineInfo, ok := defineMap["define"]; ok {
			if defineInfoMap, ok := defineInfo.(map[string]interface{}); ok {
				return defineInfoMap, nil
			}
		}
	}

	resultMap["setupData"] = make(map[string]interface{})
	resultMap["inputJson"] = make(map[string]interface{})
	resultMap["outputJson"] = make(map[string]interface{})

	return resultMap, nil
}

// CreateNewComponent is
func CreateNewComponent(namespace, componentName, componentVersion string) (string, error) {
	createComponentChan <- true
	defer func() {
		<-createComponentChan
	}()

	var count int64
	err := new(models.Component).GetComponent().Where("namespace = ?", namespace).Where("component = ?", componentName).Order("-id").Count(&count).Error
	if err != nil {
		return "", errors.New("error when query component data in database:" + err.Error())
	}

	if count > 0 {
		return "", errors.New("component name is exist!")
	}

	componentInfo := new(models.Component)
	componentInfo.Namespace = namespace
	componentInfo.Component = componentName
	componentInfo.Version = componentVersion
	componentInfo.VersionCode = 1

	err = componentInfo.GetComponent().Save(componentInfo).Error
	if err != nil {
		return "", errors.New("error when save component info:" + err.Error())
	}

	return "create new component success", nil
}

// CreateNewComponentVersion is
func CreateNewComponentVersion(componentInfo models.Component, versionName string) error {
	var count int64
	new(models.Component).GetComponent().Where("namespace = ?", componentInfo.Namespace).Where("component = ?", componentInfo.Component).Where("version = ?", versionName).Count(&count)
	if count > 0 {
		return errors.New("version code already exist!")
	}

	// get current least component's version
	leastComponent := new(models.Component)
	err := leastComponent.GetComponent().Where("namespace = ? ", componentInfo.Namespace).Where("component = ?", componentInfo.Component).Order("-id").First(&leastComponent).Error
	if err != nil {
		return errors.New("error when get least component info :" + err.Error())
	}

	newComponentInfo := new(models.Component)
	newComponentInfo.Namespace = componentInfo.Namespace
	newComponentInfo.Version = versionName
	newComponentInfo.VersionCode = leastComponent.VersionCode + 1
	newComponentInfo.Component = componentInfo.Component
	newComponentInfo.Type = componentInfo.Type
	newComponentInfo.Title = componentInfo.Title
	newComponentInfo.Gravatar = componentInfo.Gravatar
	newComponentInfo.Description = componentInfo.Description
	newComponentInfo.Endpoint = componentInfo.Endpoint
	newComponentInfo.Source = componentInfo.Source
	newComponentInfo.Environment = componentInfo.Environment
	newComponentInfo.Tag = componentInfo.Tag
	newComponentInfo.VolumeLocation = componentInfo.VolumeLocation
	newComponentInfo.VolumeData = componentInfo.VolumeData
	newComponentInfo.Makefile = componentInfo.Makefile
	newComponentInfo.Kubernetes = componentInfo.Kubernetes
	newComponentInfo.Swarm = componentInfo.Swarm
	newComponentInfo.Input = componentInfo.Input
	newComponentInfo.Output = componentInfo.Output
	newComponentInfo.Manifest = componentInfo.Manifest

	return newComponentInfo.GetComponent().Save(newComponentInfo).Error
}

func UpdateComponentInfo(componentInfo models.Component) error {
	manifestMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(componentInfo.Manifest), &manifestMap)
	if err != nil {
		return errors.New("error when unmarshal component's define info:" + err.Error())
	}

	defineMap, ok := manifestMap["define"].(map[string]interface{})
	if !ok {
		return errors.New("component define is not a json:" + err.Error())
	}

	inputMap, ok := defineMap["inputJson"].(map[string]interface{})
	if ok {
		inputDescribe, err := describeJSON(inputMap, "")
		if err != nil {
			return errors.New("error in component output json define:" + err.Error())
		}

		inputDescBytes, _ := json.Marshal(inputDescribe)
		componentInfo.Input = string(inputDescBytes)
	}

	outputMap, ok := defineMap["outputJson"].(map[string]interface{})
	if ok {
		outputDescribe, err := describeJSON(outputMap, "")
		if err != nil {
			return errors.New("error in component output json define:" + err.Error())
		}

		outputDescBytes, _ := json.Marshal(outputDescribe)
		componentInfo.Output = string(outputDescBytes)
	}

	setupDataMap, ok := defineMap["setupData"].(map[string]interface{})
	if !ok {
		return errors.New("error in component setup data: setup data is not a json")
	}

	envMap := make(map[string]interface{})
	envDefind, ok := defineMap["env"].([]interface{})
	if ok {
		for _, env := range envDefind {
			if tempEnvMap, ok := env.(map[string]interface{}); ok {
				envMap[tempEnvMap["key"].(string)] = tempEnvMap["value"].(string)
			}
		}

		envByte, _ := json.Marshal(envMap)
		componentInfo.Environment = string(envByte)
	}

	componentSetupDetail, ok := setupDataMap["action"].(map[string]interface{})
	if ok {
		// if name, ok := componentSetupDetail["name"].(string); ok {
		// componentInfo.Component = name
		// }

		if imageInfo, ok := componentSetupDetail["image"].(map[string]interface{}); ok {
			imageName := ""
			if name, ok := imageInfo["name"].(string); ok {
				imageName = name + ":"
				if tag, ok := imageInfo["tag"].(string); ok {
					imageName += tag
				} else {
					imageName += "latest"
				}
			}

			componentInfo.Endpoint = imageName
		}

		if env, ok := componentSetupDetail["env"].(string); ok {
			componentInfo.Environment = env
		}

		if timeout, ok := componentSetupDetail["timeout"].(string); ok {
			timeoutInt, err := strconv.ParseInt(timeout, 10, 64)
			if err != nil {
				return errors.New("component's timeout is not a string")
			}
			componentInfo.Timeout = timeoutInt
		}

		// unmarshal k8s info
		if useAdvanced, ok := componentSetupDetail["useAdvanced"].(bool); ok {
			configMap := make(map[string]interface{})
			podConfigKey := "pod"
			serviceConfigKey := "service"
			if useAdvanced {
				podConfigKey = "pod_advanced"
				serviceConfigKey = "service_advanced"
			}

			podConfig, ok := setupDataMap[podConfigKey].(map[string]interface{})
			if !ok {
				configMap["podConfig"] = make(map[string]interface{})
			} else {
				configMap["podConfig"] = podConfig
			}

			serviceConfig, ok := setupDataMap[serviceConfigKey].(map[string]interface{})
			if !ok {
				configMap["serviceConfig"] = make(map[string]interface{})
			} else {
				configMap["serviceConfig"] = serviceConfig
			}

			kuberSetting, _ := json.Marshal(configMap)
			componentInfo.Kubernetes = string(kuberSetting)
		}
	}

	return componentInfo.GetComponent().Save(&componentInfo).Error
}

type kubeComponent struct {
	host string
	// port          int64
	namespace     string
	reachableIPs  []string
	podConfig     map[string]interface{}
	serviceConfig map[string]interface{}
	componentInfo models.ActionLog
}

// start a component in kube env
func (kube *kubeComponent) Start(id string, eventList []models.EventDefinition, envMap map[string]string) error {
	// get componentName for service selector
	reqMap := kube.podConfig
	if len(reqMap) == 0 {
		reqMap = make(map[string]interface{})
	}

	componentName := strings.Replace(id, ",", "-", -1)
	if len(componentName) > 250 {
		componentName = componentName[len(componentName)-250:]
	}
	componentName = "pod-" + componentName

	serviceAddr, err := kube.startService(id, componentName)
	if err != nil {
		return err
	}

	err = kube.startPod(id, serviceAddr, eventList, envMap)

	return err
}

// to start a pod, except podConfig,also need add some extra info
// for example , env call back url or event list need to callback or current pod's run id use to id itself when make a call back
func (kube *kubeComponent) startPod(id, serviceAddr string, eventList []models.EventDefinition, envMap map[string]string) error {

	podInfoMap, err := kube.getPodInfo(id, serviceAddr, eventList, envMap)
	if err != nil {
		return err
	}

	rcName := strings.Replace(id, ",", "-", -1)
	if len(rcName) > 250 {
		rcName = rcName[len(rcName)-250:]
	}
	rcName = "rc-" + rcName

	rcMap := make(map[string]interface{})
	rcMetaData := make(map[string]interface{})
	rcSpecMap := make(map[string]interface{})

	rcMetaData["name"] = rcName
	rcMap["metadata"] = rcMetaData

	rcSpecMap["replicas"] = 1
	rcSpecMap["template"] = podInfoMap
	rcMap["spec"] = rcSpecMap

	reqBody, _ := json.Marshal(rcMap)

	fmt.Println("=====================replicationcontrollers=============================")
	fmt.Println(string(reqBody))
	fmt.Println(kube.host + "/api/v1/namespaces/" + kube.namespace + "/replicationcontrollers")
	fmt.Println("=====================replicationcontrollers=============================")
	resp, err := http.Post(kube.host+"/api/v1/namespaces/"+kube.namespace+"/replicationcontrollers", "application/json", bytes.NewReader(reqBody))
	if err != nil {
		return err
	}

	if resp.StatusCode != 201 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errors.New("error when get pod err body:" + err.Error() + "==== error code is :" + strconv.FormatInt(int64(resp.StatusCode), 10))
		} else {
			return errors.New("error when create pod: msg is :" + string(body))
		}
	}
	return nil
}

func (kube *kubeComponent) getPodInfo(id, serviceAddr string, eventList []models.EventDefinition, envMap map[string]string) (map[string]interface{}, error) {
	reqMap := kube.podConfig
	if len(reqMap) == 0 {
		reqMap = make(map[string]interface{})
	}

	// // first set kube metadata
	metaInfoMap := make(map[string]interface{})
	metaInfo, ok := reqMap["metadata"]
	if ok {
		metaInfoMap, ok = metaInfo.(map[string]interface{})
		if !ok {
			return nil, errors.New("component's kube config error, metadata is not a json!")
		}
	}

	podName := strings.Replace(id, ",", "-", -1)
	if len(podName) > 250 {
		podName = podName[len(podName)-250:]
	}
	podName = "pod-" + podName

	metaInfoMap["name"] = podName

	labelsMap := make(map[string]interface{})
	labelInfo, ok := metaInfoMap["labels"]
	if ok {
		labelsMap, ok = labelInfo.(map[string]interface{})
		if !ok {
			return nil, errors.New("component's kube config error, labels is not a json!")
		}
	}
	labelsMap["PIPELINE_DEFAULT_POD_LABLE"] = podName
	metaInfoMap["labels"] = labelsMap

	reqMap["metadata"] = metaInfoMap

	// set container spec info
	specMap := make(map[string]interface{})
	specInfo, ok := reqMap["spec"]
	if ok {
		specMap, ok = specInfo.(map[string]interface{})
		if !ok {
			return nil, errors.New("component's kube config error, specInfo is not a json!")
		}
	}

	// set containers info
	containers := make([]map[string]interface{}, 0)
	containersInfo, ok := specMap["containers"]
	if ok {
		tempContainersInfo, ok := containersInfo.([]interface{})
		if !ok {
			return nil, errors.New("component's kube config error, container info in containers is not a array")
		}

		for _, tempContainerInfo := range tempContainersInfo {
			tempContainer, ok := tempContainerInfo.(map[string]interface{})
			if !ok {
				return nil, errors.New("component's kube config error, container info in containers is not a json")
			}
			containers = append(containers, tempContainer)
		}
	}

	if len(containers) < 1 {
		containerInfo := make(map[string]interface{})
		containerInfo["name"] = strings.Replace(id, ",", "-", -1) + "-pod"
		containerInfo["image"] = kube.componentInfo.Endpoint

		containers = append(containers, containerInfo)
	}

	// add event envMap
	envList := make([]interface{}, 0)
	for key, value := range envMap {
		tempEnv := make(map[string]string)
		tempEnv["name"] = key
		tempEnv["value"] = value

		envList = append(envList, tempEnv)
	}

	eventListStr := ""
	for _, event := range eventList {
		tempEnv := make(map[string]string)
		tempEnv["name"] = event.Title
		tempEnv["value"] = event.Definition

		envList = append(envList, tempEnv)

		eventListStr += ";" + event.Title + "," + strconv.FormatInt(event.ID, 10)
	}

	if serviceAddr != "" {
		envList = append(envList, map[string]string{"name": "SERVICE_ADDR", "value": serviceAddr})
	}

	envList = append(envList, map[string]string{"name": "POD_NAME", "value": podName})
	envList = append(envList, map[string]string{"name": "RUN_ID", "value": id})
	envList = append(envList, map[string]string{"name": "EVENT_LIST", "value": strings.TrimPrefix(eventListStr, ";")})

	// set env to each container
	for _, container := range containers {
		if _, ok := container["name"]; !ok {
			container["name"] = strings.Replace(id, ",", "-", -1) + "-pod"
		}
		if _, ok := container["image"]; !ok {
			container["image"] = kube.componentInfo.Endpoint
		}

		env, ok := container["env"]
		if !ok {
			container["env"] = envList
		} else {
			cEnvMap, ok := env.([]interface{})
			if !ok {
				return nil, errors.New("component's kube config error, container's env is not a array")
			}

			container["env"] = append(cEnvMap, envList...)
		}

		ports := make([]map[string]interface{}, 0)
		for _, serviceAddrInfo := range strings.Split(serviceAddr, ",") {
			if len(strings.Split(serviceAddrInfo, ":")) > 2 {
				tempPort := make(map[string]interface{})
				protInt, err := strconv.ParseInt(strings.Split(serviceAddrInfo, ":")[2], 10, 64)
				if err != nil {
					return nil, errors.New("error when parse ports info to number :" + err.Error())
				}
				tempPort["containerPort"] = protInt
				ports = append(ports, tempPort)
			}
		}

		if len(ports) > 0 {
			container["ports"] = ports
		}

	}

	specMap["containers"] = containers
	reqMap["spec"] = specMap

	return reqMap, nil
}

func (kube *kubeComponent) startService(id, podName string) (string, error) {
	reqMap := kube.serviceConfig
	if len(reqMap) == 0 {
		// if don't have a service config ,will not start a service,start a port directly
		return "", nil
	}

	// first set kube metadata
	metaInfoMap := make(map[string]interface{})
	metaInfo, ok := reqMap["metadata"]
	if ok {
		metaInfoMap, ok = metaInfo.(map[string]interface{})
		if !ok {
			return "", errors.New("component's kube config error, metadata is not a json!")
		}
	}

	serviceName := strings.Replace(id, ",", "-", -1)
	if len(serviceName) > 20 {
		serviceName = serviceName[len(serviceName)-20:]
	}
	serviceName = "ser-" + serviceName
	metaInfoMap["name"] = serviceName

	reqMap["metadata"] = metaInfoMap

	// if service config has config a service ip&host ,then ,use config ,otherwise use system allocate
	// set service spec info
	specMap := make(map[string]interface{})
	specInfo, ok := reqMap["spec"]
	if ok {
		specMap, ok = specInfo.(map[string]interface{})
		if !ok {
			return "", errors.New("component's kube config error, specInfo is not a json!")
		}
	}

	// set ports info
	ports := make([]map[string]interface{}, 0)
	portsInfo, ok := specMap["ports"]
	if ok {
		tPorts, ok := portsInfo.([]interface{})
		if !ok {
			return "", errors.New("component's kube config error, container info in ports is not a array")
		}

		for _, port := range tPorts {
			tempPort, ok := port.(map[string]interface{})
			if !ok {
				return "", errors.New("component's kube config error, container info in ports is not a json")
			}

			ports = append(ports, tempPort)
		}
	}

	if len(ports) < 1 {
		return "", errors.New("component's kube config error, service must set at least one port！")
	}

	// set selector
	selectorMap := make(map[string]string)
	selectorMap["PIPELINE_DEFAULT_POD_LABLE"] = podName

	specMap["selector"] = selectorMap
	specMap["type"] = "NodePort"

	// if len(kube.reachableIPs) > 0 {
	// 	specMap["externalIPs"] = kube.reachableIPs
	// }

	// create service
	reqBody, _ := json.Marshal(reqMap)
	fmt.Println("=====================service=============================")
	fmt.Println(string(reqBody))
	fmt.Println(kube.host + "/api/v1/namespaces/" + kube.namespace + "/services")
	fmt.Println("=====================service=============================")
	resp, err := http.Post(kube.host+"/api/v1/namespaces/"+kube.namespace+"/services", "application/json", bytes.NewReader(reqBody))
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 201 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", errors.New("error when get service err body:" + err.Error() + "==== error code is :" + strconv.FormatInt(int64(resp.StatusCode), 10))
		} else {
			return "", errors.New("error when create service: msg is :" + string(body))
		}
	}

	// get resp of create service ,to get service's ip and ports
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("error when read create service resp:" + err.Error())
	}

	respMap := make(map[string]interface{})
	err = json.Unmarshal(respBody, &respMap)
	if err != nil {
		return "", errors.New("error when unmarshal create service resp:" + err.Error())
	}

	respSpecMap := make(map[string]interface{})
	respSpecInfo, ok := respMap["spec"]
	if ok {
		respSpecMap, ok = respSpecInfo.(map[string]interface{})
		if !ok {
			return "", errors.New("error when read create service resp:, specInfo is not a json!")
		}
	}

	clusterIp := ""
	cluIp, ok := respSpecMap["clusterIP"]
	if !ok {
		return "", errors.New("error when read create service resp: clusterIP is nil!")
	} else {
		cIp, ok := cluIp.(string)
		if !ok {
			return "", errors.New("error when read create service resp:, compoent name is not a string!")
		} else {
			clusterIp = cIp
		}
	}

	// set ports info
	respPorts := make([]map[string]interface{}, 0)
	respPortsInfo, ok := respSpecMap["ports"]
	if ok {
		tempRespPorts, ok := respPortsInfo.([]interface{})
		if !ok {
			return "", errors.New("error when read create service resp:, container info in ports is not a array")
		}

		for _, port := range tempRespPorts {
			tempPortMap, ok := port.(map[string]interface{})
			if !ok {
				return "", errors.New("error when read create service resp:, container info in ports is not a json")
			}
			respPorts = append(respPorts, tempPortMap)
		}
	}

	if len(respPorts) < 1 {
		return "", errors.New("error when read create service resp:, service has not set a port！")
	}

	portsStr := ""
	listenPortStr := ""

	for _, port := range respPorts {
		portF, ok := port["nodePort"].(float64)
		if !ok {
			return "", errors.New("error when parse create service resp: nodePort is not a number!")
		}
		portStr := strconv.FormatFloat(portF, 'f', 0, 64)

		listenPort, ok := port["targetPort"].(float64)
		if !ok {
			return "", errors.New("error when parse create service resp: targetPort is not a number!")
		}

		listenPortStr = strconv.FormatFloat(listenPort, 'f', 0, 64)

		portsStr += "," + portStr
	}
	portsStr = strings.TrimPrefix(portsStr, ",")

	serviceAddr := clusterIp + ":" + portsStr + ":" + listenPortStr

	return strings.TrimPrefix(serviceAddr, ","), nil
}

func (kube *kubeComponent) Stop(id string) (string, error) {
	client := &http.Client{}

	// set rc's replicas = 0 then delete rc
	podInfoMap, err := kube.getPodInfo(id, "", nil, nil)
	if err != nil {
		return "", err
	}

	rcName := strings.Replace(id, ",", "-", -1)
	if len(rcName) > 250 {
		rcName = rcName[len(rcName)-250:]
	}
	rcName = "rc-" + rcName

	rcMap := make(map[string]interface{})
	rcMetaData := make(map[string]interface{})
	rcSpecMap := make(map[string]interface{})

	rcMetaData["name"] = rcName
	rcMap["metadata"] = rcMetaData

	rcSpecMap["replicas"] = 0
	rcSpecMap["template"] = podInfoMap
	rcMap["spec"] = rcSpecMap

	reqBody, _ := json.Marshal(rcMap)
	fmt.Println("=====================replicationcontrollers=============================")
	fmt.Println(string(reqBody))
	fmt.Println(kube.host + "/api/v1/namespaces/" + kube.namespace + "/replicationcontrollers/" + rcName)
	fmt.Println("=====================replicationcontrollers=============================")
	req, err := http.NewRequest("PUT", kube.host+"/api/v1/namespaces/"+kube.namespace+"/replicationcontrollers/"+rcName, bytes.NewReader(reqBody))
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.New(strconv.FormatInt(int64(resp.StatusCode), 10))
	}

	// delete rc
	fmt.Println("=====================replicationcontrollers=============================")
	fmt.Println(kube.host + "/api/v1/namespaces/" + kube.namespace + "/replicationcontrollers/" + rcName)
	fmt.Println("=====================replicationcontrollers=============================")
	req, err = http.NewRequest("DELETE", kube.host+"/api/v1/namespaces/"+kube.namespace+"/replicationcontrollers/"+rcName, nil)
	if err != nil {
		return "", err
	}

	resp, err = client.Do(req)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.New(strconv.FormatInt(int64(resp.StatusCode), 10))
	}

	// delete service
	serviceName := strings.Replace(id, ",", "-", -1)
	if len(serviceName) > 20 {
		serviceName = serviceName[len(serviceName)-20:]
	}
	serviceName = "ser-" + serviceName

	fmt.Println("=====================service=============================")
	fmt.Println(kube.host + "/api/v1/namespaces/" + kube.namespace + "/services/" + serviceName)
	fmt.Println("=====================service=============================")
	req, err = http.NewRequest("DELETE", kube.host+"/api/v1/namespaces/"+kube.namespace+"/services/"+serviceName, nil)
	if err != nil {
		return "", err
	}

	resp, err = client.Do(req)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.New(strconv.FormatInt(int64(resp.StatusCode), 10))
	}

	return "success", nil
}

func (kube *kubeComponent) GetIp(args ...interface{}) (string, error) {
	// return kube.host, nil
	if len(kube.reachableIPs) < 0 {
		return "", errors.New("no reachable ips for this kube cluster")
	}

	return kube.reachableIPs[0], nil
}
