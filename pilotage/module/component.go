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
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"fmt"
	"github.com/Huawei/containerops/pilotage/models"
	"github.com/Huawei/containerops/pilotage/utils"
	log "github.com/Sirupsen/logrus"
)

const (
	// K8SCOMPONENT means component will run in k8s
	K8SCOMPONENT = "KUBERNETES"
	// SWARMCOMPONENT means component will run in swarm
	SWARMCOMPONENT = "SWARM"
)

type component interface {
	Start() error
	Update()
	Stop() error
	SendData(receiveDataUri string, data []byte) ([]*http.Response, error)
}

type kubeComponent struct {
	runID         string
	apiServerUri  string
	namespace     string
	nodeIP        string
	podConfig     map[string]interface{}
	serviceConfig map[string]interface{}
	componentInfo models.ActionLog
}

// GetComponentListByNamespace is get component list by given namespace
func GetComponentListByNamespace(namespace string) ([]map[string]interface{}, error) {
	resultMap := make([]map[string]interface{}, 0)
	componentList := make([]models.Component, 0)
	componentsMap := make(map[string]interface{})

	err := new(models.Component).GetComponent().Where("namespace = ?", namespace).Order("-id").Find(&componentList).Error
	if err != nil {
		return nil, errors.New("error when get component infos by namespace:" + namespace + ",error:" + err.Error())
	}

	for _, componentInfo := range componentList {
		if _, ok := componentsMap[componentInfo.Name]; !ok {
			tempMap := make(map[string]interface{})
			tempMap["version"] = make(map[int64]interface{})
			componentsMap[componentInfo.Name] = tempMap
		}

		componentMap := componentsMap[componentInfo.Name].(map[string]interface{})
		versionMap := componentMap["version"].(map[int64]interface{})

		versionMap[componentInfo.VersionCode] = componentInfo
		componentMap["id"] = componentInfo.ID
		componentMap["name"] = componentInfo.Name
		componentMap["version"] = versionMap
	}

	for _, component := range componentList {
		componentInfo := componentsMap[component.Name].(map[string]interface{})

		if isSign, ok := componentInfo["isSign"].(bool); ok && isSign {
			continue
		}

		componentInfo["isSign"] = true
		componentsMap[component.Name] = componentInfo

		versionList := make([]map[string]interface{}, 0)
		for _, componentVersion := range componentList {
			if componentVersion.Name == componentInfo["name"].(string) {
				versionMap := make(map[string]interface{})
				versionMap["id"] = componentVersion.ID
				versionMap["version"] = componentVersion.Version
				versionMap["versionCode"] = componentVersion.VersionCode

				versionList = append(versionList, versionMap)
			}
		}

		tempResult := make(map[string]interface{})
		tempResult["id"] = componentInfo["id"]
		tempResult["name"] = componentInfo["name"]
		tempResult["version"] = versionList

		resultMap = append(resultMap, tempResult)
	}

	return resultMap, nil
}

func CreateComponent(data string) (uint64, error) {
	var component *models.Component
	err := json.Unmarshal([]byte(data), component)
	if err != nil {
		log.Errorln("CreateComponent unmarshal data error: ", err.Error())
		return 0, errors.New("unmarshal data error: " + err.Error())
	}
	if component.ID != 0 {
		return 0, fmt.Errorf("specify a component id: %d", component.ID)
	}
	condition := &models.Component{
		Name: component.Name,
		Version: component.Version,
	}
	if result, err := condition.SelectComponent(); err != nil {
		log.Errorln("CreateComponent query component error: ", err.Error())
		return 0, errors.New("query component error: " + err.Error())
	} else {
		if result != nil {
			return 0, errors.New("component exists, id is: " + result.ID)
		}
	}

	if err = component.Create(); err != nil {
		log.Errorln("CreateComponent query component error: ", err.Error())
		return 0, errors.New("query component error: " + err.Error())
	}
	return component.ID, nil
}

// GetComponentInfo is get component info by given namespace and componentname and componentId
func GetComponentInfo(namespace, componentName string, componentId int64) (map[string]interface{}, error) {
	resultMap := make(map[string]interface{})
	componentInfo := new(models.Component)

	err := componentInfo.GetComponent().Where("id = ?", componentId).First(&componentInfo).Error
	if err != nil {
		return nil, errors.New("error when get component info from db:" + err.Error())
	}

	if componentInfo.Name != componentName {
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

// UpdateComponentInfo is update a component define by give info
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

	if inputMap, ok := defineMap["inputJson"].(map[string]interface{}); ok {
		inputDescribe, err := describeJSON(inputMap, "")
		if err != nil {
			return errors.New("error in component output json define:" + err.Error())
		}

		inputDescBytes, _ := json.Marshal(inputDescribe)
		componentInfo.Input = string(inputDescBytes)
	}

	if outputMap, ok := defineMap["outputJson"].(map[string]interface{}); ok {
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
	if envDefind, ok := defineMap["env"].([]interface{}); ok {
		for _, env := range envDefind {
			if tempEnvMap, ok := env.(map[string]interface{}); ok {
				envMap[tempEnvMap["key"].(string)] = tempEnvMap["value"].(string)
			}
		}

		envByte, _ := json.Marshal(envMap)
		componentInfo.Environment = string(envByte)
	}

	if componentSetupDetail, ok := setupDataMap["action"].(map[string]interface{}); ok {
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
			timeoutInt := int64(0)
			if timeout != "" {
				timeoutInt, err = strconv.ParseInt(timeout, 10, 64)
				if err != nil {
					return errors.New("component's timeout is not a string")
				}
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
			componentInfo.KubeSetting = string(kuberSetting)
		}
	}

	return componentInfo.GetComponent().Save(&componentInfo).Error
}

// CreateNewComponentVersion is copy current component info to a new component with diff version name
func CreateNewComponentVersion(componentInfo models.Component, versionName string) error {
	var count int64
	err := new(models.Component).GetComponent().Where("namespace = ?", componentInfo.Namespace).Where("component = ?", componentInfo.Name).Where("version = ?", versionName).Count(&count).Error
	if err != nil {
		return errors.New("error when get component version info:" + err.Error())
	}

	if count > 0 {
		return errors.New("version already exist!")
	}

	// get current least component's version
	leastComponent := new(models.Component)
	err = leastComponent.GetComponent().Where("namespace = ? ", componentInfo.Namespace).Where("component = ?", componentInfo.Name).Order("-id").First(&leastComponent).Error
	if err != nil {
		return errors.New("error when get least component info :" + err.Error())
	}

	newComponentInfo := new(models.Component)
	newComponentInfo.Namespace = componentInfo.Namespace
	newComponentInfo.Version = strings.TrimSpace(versionName)
	newComponentInfo.VersionCode = leastComponent.VersionCode + 1
	newComponentInfo.Name = componentInfo.Name
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
	newComponentInfo.KubeSetting = componentInfo.KubeSetting
	newComponentInfo.Swarm = componentInfo.Swarm
	newComponentInfo.Input = componentInfo.Input
	newComponentInfo.Output = componentInfo.Output
	newComponentInfo.Manifest = componentInfo.Manifest

	return newComponentInfo.GetComponent().Save(newComponentInfo).Error
}

func NewComponent(actionLog *ActionLog) (component, error) {
	platformSetting, err := actionLog.GetActionPlatformInfo()
	if err != nil {
		log.Errorln("[component's InitComponent]:error when get given actionLog's platformSetting:", actionLog, " ===>error is:", err.Error())
		return nil, err
	}

	switch platformSetting["platformType"] {
	case K8SCOMPONENT:
		k8sComp := new(kubeComponent)

		ComponentConfigMap := make(map[string]interface{}, 0)
		err := json.Unmarshal([]byte(actionLog.Kubernetes), &ComponentConfigMap)
		if err != nil {
			log.Errorln("[component's InitComponent]:error when get action's kubernetes setting:", actionLog, " ===>error is:", err.Error())
			return k8sComp, errors.New("get action's kube config error:" + err.Error())
		}

		nodeIP, ok := ComponentConfigMap["nodeIP"].(string)
		if !ok {
			log.Errorln("[component's InitComponent]:error when get component's nodeIP:",
				ComponentConfigMap)
			return nil, errors.New("get action's kube config error,kube's nodeIP is not set")
		}

		podConfig, ok := ComponentConfigMap["podConfig"].(map[string]interface{})
		if !ok {
			log.Errorln("[component's InitComponent]:error when get component's podConfig:", ComponentConfigMap, " podConfig is not a json obj")
			return k8sComp, errors.New("component kube config error ,podConfig is not a json obj")
		}

		serviceConfig, ok := ComponentConfigMap["serviceConfig"].(map[string]interface{})
		if !ok {
			log.Errorln("[component's InitComponent]:error when get component's serviceConfig:", ComponentConfigMap, " serviceConfig is not a json obj")
			return k8sComp, errors.New("component kube config error ,serviceConfig is not a json data")
		}

		k8sComp.runID = fmt.Sprintf("%d-%d-%d", actionLog.Workflow, actionLog.Stage, actionLog.ID)
		k8sComp.apiServerUri = platformSetting["platformHost"]
		k8sComp.namespace = actionLog.Namespace
		k8sComp.nodeIP = nodeIP
		k8sComp.podConfig = podConfig
		k8sComp.serviceConfig = serviceConfig
		k8sComp.componentInfo = *actionLog.ActionLog

		return k8sComp, nil
	case SWARMCOMPONENT:
		return nil, fmt.Errorf("Component type %s isn't supported", platformSetting["platformType"])
	default:
		return nil, fmt.Errorf("Component type %s isn't supported", platformSetting["platformType"])
	}
}

func (c *kubeComponent) Start() error {
	exist, err := c.IsNamespaceExist()
	if err != nil {
		log.Errorln("[kubeComponent Start]: query namespace info error, ", err)
		return err
	}

	if !exist {
		err = c.CreateNamespace()
		if err != nil {
			log.Errorln("[kubeComponent Start]: create namespace error, ", err)
			return err
		}
	}

	serviceAddr, err := c.StartService()
	if err != nil {
		log.Errorln("[kubeComponent Start]:start service error, ", err)
		return err
	}

	err = c.StartRC(serviceAddr)
	if err != nil {
		log.Errorln("[kubeComponent Start]:start RC error, ", err)
		return err
	}

	return nil
}

func (kube *kubeComponent) Stop() error {
	client := &http.Client{}

	// delete service
	serviceName := "ser-" + kube.runID
	if len(serviceName) > 253 {
		serviceName = serviceName[len(serviceName)-253:]
	}

	kubeReqDeleteServiceUrl := kube.apiServerUri + "/api/v1/namespaces/" + kube.namespace + "/services/" + serviceName
	log.Info("[kubeComponent's Stop]:send delete service req to:", kubeReqDeleteServiceUrl)

	req, err := http.NewRequest("DELETE", kubeReqDeleteServiceUrl, nil)
	if err != nil {
		log.Error("[kubeComponent's Stop]:error when generate new req to:", kubeReqDeleteServiceUrl, " ===>error is:", err.Error())
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Error("[kubeComponent's Stop]:error when send req to:", kubeReqDeleteServiceUrl, " ===>error is:", err.Error())
		return err
	}

	if resp.StatusCode != http.StatusOK {
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Error("[kubeComponent's Stop]:error when get resp body from delete service:", err.Error(), " ===>error code:", resp.StatusCode)
			return err
		}
		return errors.New("error when delete service:" + string(respBody))
	}

	// set rc's replicas = 0 then delete rc
	podInfoMap, err := kube.GetPodDefine("")
	if err != nil {
		log.Error("[kubeComponent's Stop]:error when get component's pod define:", err.Error())
		return err
	}

	rcName := "rc-" + kube.runID
	if len(rcName) > 253 {
		rcName = rcName[len(rcName)-253:]
	}

	rcMap := make(map[string]interface{})
	rcSpecMap := make(map[string]interface{})

	rcMap["metadata"] = map[string]interface{}{"name": rcName}

	rcSpecMap["replicas"] = 0
	rcSpecMap["template"] = podInfoMap
	rcMap["spec"] = rcSpecMap

	kubeReqModifyRCUrl := kube.apiServerUri + "/api/v1/namespaces/" + kube.namespace + "/replicationcontrollers/" + rcName
	reqBody, _ := json.Marshal(rcMap)
	log.Info("[kubeComponent's Stop]:send to:", kubeReqModifyRCUrl, " reqBody is:", string(reqBody))

	req, err = http.NewRequest("PUT", kubeReqModifyRCUrl, bytes.NewReader(reqBody))
	if err != nil {
		log.Error("[kubeComponent's Stop]:error when generate new req to:", kubeReqModifyRCUrl, " ===>error is:", err.Error())
		return err
	}

	resp, err = client.Do(req)
	if err != nil {
		log.Error("[kubeComponent's Stop]:error when send req to:", kubeReqModifyRCUrl, " reqBody is:", string(reqBody), " ===>error is:", err.Error())
		return err
	}

	if resp.StatusCode != http.StatusOK {
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Error("[kubeComponent's Stop]:error when get resp body from modify rc:", err.Error(), " ===>error code:", resp.StatusCode)
			return err
		}
		return errors.New("error when modify rc:" + string(respBody))
	}

	// delete rc
	kubeReqDeleteRCUrl := kube.apiServerUri + "/api/v1/namespaces/" + kube.namespace + "/replicationcontrollers/" + rcName

	log.Info("[kubeComponent's Stop]:send delete rc req to:", kubeReqDeleteRCUrl)
	req, err = http.NewRequest("DELETE", kubeReqDeleteRCUrl, nil)
	if err != nil {
		return err
	}

	resp, err = client.Do(req)
	if err != nil {
		log.Error("[kubeComponent's Stop]:error when generate new req to:", kubeReqModifyRCUrl, " ===>error is:", err.Error())
		return err
	}

	if resp.StatusCode != http.StatusOK {
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Error("[kubeComponent's Stop]:error when get resp body from delete rc:", err.Error(), " ===>error code:", resp.StatusCode)
			return err
		}
		return errors.New("error when delete rc:" + string(respBody))
	}

	return nil
}

func (kube *kubeComponent) SendData(receiveDataUri string, reqBody []byte) (respList []*http.Response, err error) {
	var resp *http.Response
	kubeReqUrl := ""

	if strings.HasPrefix(kube.nodeIP, "http://") || strings.HasPrefix(kube.nodeIP, "https://") {
		kubeReqUrl = kube.nodeIP + receiveDataUri
	} else {
		kubeReqUrl = "http://" + kube.nodeIP + receiveDataUri
	}

	log.Info("[kubeComponent's SendData]:send data:", string(reqBody), " to:", kubeReqUrl)

	sendSuccess := false

	for count := 0; count < 10 && !sendSuccess; count++ {
		resp, err = http.Post(kubeReqUrl, "application/json", bytes.NewReader(reqBody))
		if err == nil && resp != nil && resp.StatusCode == http.StatusOK {
			sendSuccess = true
		}
		log.Info("[kubeComponent's SendData]:send data:", string(reqBody), " to:", kubeReqUrl, " count:", count, "\n resp:", resp, " err:", err)

		time.Sleep(1 * time.Second)
	}

	respList = append(respList, resp)
	return respList, err
}

func (kube *kubeComponent) StartService() (string, error) {
	reqMap := kube.serviceConfig

	serviceName := "co-svc-" + kube.runID
	if len(serviceName) > 253 {
		serviceName = serviceName[:253]
	}

	if metadataMap, ok := reqMap["metadata"].(map[string]interface{}); ok {
		metadataMap["name"] = serviceName
	} else {
		reqMap["metadata"] = map[string]interface{}{"name": serviceName}
	}

	// if service config has config a service ip&host ,then ,use config ,otherwise use system allocate
	// set service spec info
	specMap := make(map[string]interface{})
	specInfo, ok := reqMap["spec"]
	if ok {
		specMap, ok = specInfo.(map[string]interface{})
		if !ok {
			log.Error("[kubeComponent's StartService]:error when get service's spec config: config is not a json obj:", specInfo)
			return "", errors.New("component's kube config error, specInfo is not a json!")
		}
	}

	// set ports info
	ports := make([]map[string]interface{}, 0)
	portsInfo, ok := specMap["ports"]
	if ok {
		tPorts, ok := portsInfo.([]interface{})
		if !ok {
			log.Error("[kubeComponent's StartService]:error when get kube config, container info in ports is not a array:", portsInfo)
			return "", errors.New("component's kube config error, container info in ports is not a array")
		}

		for i, port := range tPorts {
			tempPort, ok := port.(map[string]interface{})
			if !ok {
				log.Error("[kubeComponent's StartService]:error when get port info, port is not a json obj:", port)
				return "", errors.New("component's kube config error, container info in ports is not a json")
			}
			if _, ok := tempPort["name"]; !ok {
				tempPort["name"] = "port-" + strconv.FormatInt(kube.componentInfo.ID, 10) + "-" + strconv.FormatInt(int64(i), 10)
			}

			ports = append(ports, tempPort)
		}
	}

	if len(ports) == 0 {
		log.Error("[kubeComponent's StartService]:service config doesn't has any port:", specMap)
	}

	// set selector
	selectorMap, ok := specMap["selector"].(map[string]interface{})
	if !ok {
		selectorMap = make(map[string]interface{})
	}

	selectorMap["WORKFLOW_DEFAULT_POD_LABLE"] = "pod-" + kube.runID

	specMap["selector"] = selectorMap

	// create service
	kubeReqUrl := kube.apiServerUri + "/api/v1/namespaces/" + kube.namespace + "/services"
	reqBody, err := json.Marshal(reqMap)
	if err != nil {

	}
	log.Info("[kubeComponent's StartService]:send request to:", kubeReqUrl, " ,reqBody is:", string(reqBody))

	resp, err := http.Post(kubeReqUrl, "application/json", bytes.NewReader(reqBody))
	if err != nil {
		log.Error("[kubeComponent's StartService]:error when send req to kube:", err.Error())
		return "", err
	}

	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("[kubeComponent's StartService]:error when read resp body info:", err.Error(), " ===>error code:", resp.StatusCode)
		return "", errors.New("error when get service err body:" + err.Error() + "==== error code is :" + strconv.FormatInt(int64(resp.StatusCode), 10))
	}

	if resp.StatusCode != http.StatusCreated {
		log.Error("[kubeComponent's StartService]:error when create service(", resp.StatusCode, "):", string(respBody))
		return "", errors.New("error when create service: msg is :" + string(respBody))
	}

	respMap := make(map[string]interface{})
	err = json.Unmarshal(respBody, &respMap)
	if err != nil {
		log.Error("[kubeComponent's StartService]:error when unmarshal resp:", string(respBody), " ===>error is:", err.Error())
		return "", errors.New("error when unmarshal create service resp:" + err.Error())
	}

	respSpecMap := make(map[string]interface{})
	respSpecInfo, ok := respMap["spec"]
	if ok {
		respSpecMap, ok = respSpecInfo.(map[string]interface{})
		if !ok {
			log.Error("[kubeComponent's StartService]:error when get resp's specInfo: specInfo is not a json obj:", respSpecInfo)
			return "", errors.New("error when read create service resp:, specInfo is not a json!")
		}
	}

	clusterIp := ""
	cluIp, ok := respSpecMap["clusterIP"].(string)
	if !ok {
		log.Error("[kubeComponent's StartService]:error when get service's clusterIP from resp:", respSpecMap)
		return "", errors.New("error when read create service resp: clusterIP is illegal!")
	} else {
		clusterIp = cluIp
	}

	// set ports info
	respPorts := make([]map[string]interface{}, 0)
	respPortsInfo, ok := respSpecMap["ports"]
	if ok {
		tempRespPorts, ok := respPortsInfo.([]interface{})
		if !ok {
			log.Error("[kubeComponent's StartService]:error when get ports info:", respPortsInfo)
			return "", errors.New("error when read create service resp:, container info in ports is not a array")
		}

		for _, port := range tempRespPorts {
			tempPortMap, ok := port.(map[string]interface{})
			if !ok {
				log.Error("[kubeComponent's StartService]:error when get port info:", port)
				return "", errors.New("error when read create service resp:, container info in ports is not a json")
			}
			respPorts = append(respPorts, tempPortMap)
		}
	}

	if len(respPorts) < 1 {
		log.Error("[kubeComponent's StartService]:error when get resp ports,resp ports is null,resp is:", respSpecMap)
	}

	portsStr := ""
	listenPortsStr := ""

	for _, port := range respPorts {
		portF, ok := port["nodePort"].(float64)
		if !ok {
			portF = float64(0)
		}
		portStr := strconv.FormatFloat(portF, 'f', 0, 64)

		listenPort, ok := port["targetPort"].(float64)
		if !ok {
			log.Error("[kubeComponent's StartService]:error when get service's target is not a number:", port)
			return "", errors.New("error when parse create service resp: targetPort is not a number!")
		}

		listenPortStr := strconv.FormatFloat(listenPort, 'f', 0, 64)

		portsStr += "," + portStr
		listenPortsStr += "," + listenPortStr
	}
	portsStr = strings.TrimPrefix(portsStr, ",")
	listenPortsStr = strings.TrimPrefix(listenPortsStr, ",")

	serviceAddr := clusterIp + ":" + portsStr + ":" + listenPortsStr

	return serviceAddr, nil
}

func (kube *kubeComponent) StartRC(serviceAddr string) error {
	podInfoMap, err := kube.GetPodDefine(serviceAddr)
	if err != nil {
		return err
	}

	rcName := "co-rc-" + kube.runID
	if len(rcName) > 253 {
		rcName = rcName[len(rcName)-253:]
	}

	rcMap := make(map[string]interface{})
	rcSpecMap := make(map[string]interface{})

	rcMap["metadata"] = map[string]interface{}{"name": rcName}

	rcSpecMap["replicas"] = 1
	rcSpecMap["template"] = podInfoMap
	rcMap["spec"] = rcSpecMap

	kubeReqUrl := kube.apiServerUri + "/api/v1/namespaces/" + kube.namespace + "/replicationcontrollers"
	reqBody, _ := json.Marshal(rcMap)

	log.Info("[kubeComponent's StartRC]:send req to:", kubeReqUrl, " body is:", string(reqBody))

	resp, err := http.Post(kubeReqUrl, "application/json", bytes.NewReader(reqBody))
	if err != nil {
		log.Error("[kubeComponent's StartRC]:error when send req to:", kubeReqUrl, " body is:", string(reqBody), " ===>error is:", err.Error())
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Error("[kubeComponent's StartRC]:error when get resp body:", err.Error(), " error code:", resp.StatusCode)
			return errors.New("error when get create rc resp body:" + err.Error() + "==== error code is :" + strconv.FormatInt(int64(resp.StatusCode), 10))
		} else {
			log.Error("[kubeComponent's StartRC]:error when create rc:", string(body), " ===>error code:", resp.StatusCode)
			return errors.New("error when create rc(" + strconv.FormatInt(int64(resp.StatusCode), 10) + "):" + string(body))
		}
	}

	go kube.Update()

	return nil
}

// IsNamespaceExist is test is kubecomponent's namespace exist
func (kube *kubeComponent) IsNamespaceExist() (bool, error) {
	kubeReqUrl := kube.apiServerUri + "/api/v1/namespaces/" + kube.namespace
	log.Debugf("[kubeComponent IsNamespaceExist]: send request to %s\n", kubeReqUrl)
	resp, err := http.Get(kubeReqUrl)
	if err != nil {
		log.Errorf("[kubeComponent IsNamespaceExist]: send request to %s error: %s\n", kubeReqUrl, err)
		return false, fmt.Errorf("get k8s namespace error: %s", err)
	}

	defer resp.Body.Close()
	switch resp.StatusCode {
	case http.StatusOK:
		return true, nil
	case http.StatusNotFound:
		log.Debugf("[kubeComponent IsNamespaceExist]: namespace %s not found\n", kube.namespace)
		return false, nil
	default:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return false, errors.New("error when get resp's err body:" + err.Error() + "==== error code is :" + strconv.FormatInt(int64(resp.StatusCode), 10))
		}
		log.Debugf("[kubeComponent IsNamespaceExist]: get k8s namespace error, status code is %d, response body is %s\n",
			resp.StatusCode, string(body))
		return false, fmt.Errorf("get k8s namespace error, status code is: %d", resp.StatusCode)
	}

}

// CreateNamespace is create kubecomponent's namespace
func (kube *kubeComponent) CreateNamespace() error {
	kubeReqUrl := kube.apiServerUri + "/api/v1/namespaces"
	reqMap := make(map[string]interface{})
	reqMap["metadata"] = map[string]string{"name": kube.namespace}
	reqBody, err := json.Marshal(reqMap)
	if err != nil {
		return errors.New("marshal data for creating namespace error: " + err.Error())
	}
	log.Debugf("[kubeComponent CreateNamespace]: send request to %s with body%s\n", kubeReqUrl, string(reqBody))

	resp, err := http.Post(kubeReqUrl, "application/json", bytes.NewReader(reqBody))
	if err != nil {
		log.Errorf("[kubeComponent CreateNamespace]: send request to %s with body %s error: %s\n", kubeReqUrl, reqBody, err)
		return fmt.Errorf("create k8s namespace error: %s", err)
	}

	defer resp.Body.Close()

	//TODO: handle error when namespace exists
	switch resp.StatusCode {
	case http.StatusCreated:
		return nil
	default:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errors.New("read response body from create namespace error: " + err.Error())
		} else {
			return fmt.Errorf("create namespace error, status code is %d, body is %s",
				resp.StatusCode, string(body))
		}
	}
}

func (kube *kubeComponent) GetPodDefine(serviceAddr string) (map[string]interface{}, error) {
	reqMap := kube.podConfig
	if len(reqMap) == 0 {
		reqMap = make(map[string]interface{})
	}

	// first set kube metadata
	metaInfoMap := make(map[string]interface{})
	metaInfo, ok := reqMap["metadata"]
	if ok {
		metaInfoMap, ok = metaInfo.(map[string]interface{})
		if !ok {
			log.Error("[kubeComponent's GetPodDefine]:error when get component's podconfig, metadata is not a json obj:", metaInfo)
			return nil, errors.New("component's kube config error, metadata is not a json!")
		}
	}

	labelsMap := make(map[string]interface{})
	labelInfo, ok := metaInfoMap["labels"]
	if ok {
		labelsMap, ok = labelInfo.(map[string]interface{})
		if !ok {
			log.Error("[kubeComponent's GetPodDefine]:error when get component's labels define,define is not a json obj:", labelInfo)
		}
	}

	podName := "pod-" + kube.runID

	labelsMap["WORKFLOW_DEFAULT_POD_LABLE"] = podName
	metaInfoMap["labels"] = labelsMap

	reqMap["metadata"] = metaInfoMap

	// set container spec info
	specMap := make(map[string]interface{})
	specInfo, ok := reqMap["spec"]
	if ok {
		specMap, ok = specInfo.(map[string]interface{})
		if !ok {
			log.Error("[kubeComponent's GetPodDefine]:error when get component's specInfo,specInfo is not a json obj:", specInfo)
			return nil, errors.New("component's kube config error, specInfo is not a json!")
		}
	}

	// set containers info
	containers := make([]map[string]interface{}, 0)
	containersInfo, ok := specMap["containers"]
	if ok {
		tempContainersInfo, ok := containersInfo.([]interface{})
		if !ok {
			log.Error("[kubeComponent's GetPodDefine]:error when get component's containers config,want an array, got:", containersInfo)
			return nil, errors.New("component's kube config error, container info in containers is not a array")
		}

		for _, tempContainerInfo := range tempContainersInfo {
			tempContainer, ok := tempContainerInfo.(map[string]interface{})
			if !ok {
				log.Error("[kubeComponent's GetPodDefine]:error when get component's container info,container info is not a json obj:", tempContainerInfo)
				return nil, errors.New("component's kube config error, container info in containers is not a json")
			}
			containers = append(containers, tempContainer)
		}
	}

	if len(containers) < 1 {
		containerInfo := make(map[string]interface{})
		containerInfo["name"] = kube.runID + "-pod"

		imageName := kube.componentInfo.ImageName
		if kube.componentInfo.ImageTag != "" {
			imageName += ":" + kube.componentInfo.ImageTag
		} else {
			imageName += ":latest"
		}

		containerInfo["image"] = imageName

		containers = append(containers, containerInfo)
	}

	// get action's data
	actionLog := new(ActionLog)
	actionLog.ActionLog = &kube.componentInfo
	manifestMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(actionLog.Manifest), &manifestMap)
	if err != nil {
		log.Error("[kubeComponent's GetPodDefine]:error when get action manifest info:" + err.Error())
		return nil, errors.New("error when unmarshal action's manifestMap")
	}

	dataMap := make(map[string]interface{})
	relations, ok := manifestMap["relation"]
	if ok {
		relationInfo, ok := relations.([]interface{})
		if !ok {
			log.Error("[kubeComponent's GetPodDefine]:error when parse relations,want an array,got:", relations)
			return nil, errors.New("error when parse relations")
		}

		dataMap, err = actionLog.merageFromActionsOutputData(relationInfo)
		if err != nil {
			log.Error("[kubeComponent's GetPodDefine]:error when get data map from action: " + err.Error())
		}
	}
	dataMapBytes, _ := json.Marshal(dataMap)

	// add event envMap
	allEventMap := make(map[string]interface{})
	envList := make([]map[string]interface{}, 0)
	workflowEnvList, err := getWorkflowEnvList(kube.componentInfo.Workflow)
	if err != nil {
		log.Error("[kubeComponent's GetPodDefine]:error when get workflow's env list:", err.Error())
		return nil, err
	}

	for _, env := range workflowEnvList {
		allEventMap[env["name"].(string)] = env["value"]
	}

	stageEnvList, err := getStageEnvList(kube.componentInfo.Stage)
	if err != nil {
		log.Error("[kubeComponent's GetPodDefine]:error when get stage's env list:", err.Error())
		return nil, err
	}

	for _, env := range stageEnvList {
		allEventMap[env["name"].(string)] = env["value"]
	}

	actionEnvList, err := getActionEnvList(kube.componentInfo.ID)
	if err != nil {
		log.Error("[kubeComponent's GetPodDefine]:error when get action's env list:", err.Error())
		return nil, err
	}

	for _, env := range actionEnvList {
		allEventMap[env["name"].(string)] = env["value"]
	}

	systemEventList, err := getSystemEventList(kube.componentInfo.ID)
	if err != nil {
		log.Error("[kubeComponent's GetPodDefine]:error when get system event define from db:", err.Error())
		return nil, err
	}

	for _, env := range systemEventList {
		allEventMap[env["name"].(string)] = env["value"]
	}

	eventListStr := ""
	for _, event := range systemEventList {
		eventListStr += ";" + event["Title"].(string) + "," + strconv.FormatInt(event["ID"].(int64), 10)
	}

	if serviceAddr != "" {
		allEventMap["CO_SERVICE_ADDR"] = serviceAddr
	}

	allEventMap["CO_POD_NAME"] = podName
	allEventMap["CO_RUN_ID"] = kube.runID
	allEventMap["CO_EVENT_LIST"] = strings.TrimPrefix(eventListStr, ";")
	allEventMap["CO_DATA"] = string(dataMapBytes)
	allEventMap["CO_SET_GLOBAL_VAR_URL"] = projectAddr + "/v2/" + actionLog.Namespace + "/" + actionLog.Repository + "/workflow/v1/runtime/var/" + strconv.FormatInt(actionLog.Workflow, 10)
	allEventMap["CO_LINKSTART_TOKEN"] = utils.MD5(actionLog.Action + kube.runID)
	allEventMap["CO_LINKSTART_URL"] = projectAddr + "/v2/" + actionLog.Namespace + "/" + actionLog.Repository + "/workflow/v1/runtime/linkstart/" + strconv.FormatInt(actionLog.Workflow, 10) + "/"
	allEventMap["CO_ACTION_TIMEOUT"] = actionLog.Timeout

	for key, value := range allEventMap {
		tempEnv := make(map[string]interface{})
		tempEnv["name"] = key
		tempEnv["value"] = value

		envList = append(envList, tempEnv)
	}

	// set env to each container
	for _, container := range containers {
		if _, ok := container["name"]; !ok {
			container["name"] = kube.runID + "-pod"
		}
		if _, ok := container["image"]; !ok {
			imageName := kube.componentInfo.ImageName
			if kube.componentInfo.ImageTag != "" {
				imageName += ":" + kube.componentInfo.ImageTag
			} else {
				imageName += ":leatest"
			}
			container["image"] = imageName
		}

		if env, ok := container["env"]; ok {
			cEnvMap, ok := env.([]map[string]interface{})
			if !ok {
				log.Error("[kubeComponent's GetPodDefine]:error when get container's env info,want an array,got:", env.(int64))
				return nil, errors.New("component's kube config error, container's env is not a array")
			}

			for _, tempEnvMap := range cEnvMap {
				container["env"] = append(container["env"].([]map[string]interface{}), tempEnvMap)
			}
		}
		container["env"] = envList
		container["imagePullPolicy"] = "Always"

		ports := make([]map[string]interface{}, 0)
		serviceAddrInfo := strings.Split(serviceAddr, ":")
		if len(serviceAddrInfo) > 2 {
			for _, port := range strings.Split(serviceAddrInfo[2], ",") {
				tempPort := make(map[string]interface{})
				portInt, err := strconv.ParseInt(port, 10, 64)
				if err != nil {
					log.Error("[kubeComponent's GetPodDefine]:error when get port info from serviceAddr,port is not a number:", serviceAddr)
					return nil, errors.New("error when parse ports info to number:" + err.Error())
				}
				// tempPort["name"] = "port" + "-" + strconv.FormatInt(kube.componentInfo.ID, 10) + "-" + strconv.FormatInt(int64(i), 10)
				tempPort["containerPort"] = portInt
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

func (kube *kubeComponent) GetServiceInfo() (map[string]interface{}, error) {
	serviceName := "co-svc-" + kube.runID
	if len(serviceName) > 253 {
		serviceName = serviceName[len(serviceName)-253:]
	}

	kubeReqUrl := kube.apiServerUri + "/api/v1/namespaces/" + kube.componentInfo.Namespace + "/services/" + serviceName

	resp, err := http.Get(kubeReqUrl)
	if err != nil {
		log.Error("[kubeComponent's GetServiceInfo]:error when send req to kube:", err.Error())
		return nil, err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Error("[kubeComponent's GetServiceInfo]:error when read resp body info:", err.Error(), " ===>error code:", resp.StatusCode)
		return nil, errors.New("error when get service err body:" + err.Error() + "==== error code is :" + strconv.FormatInt(int64(resp.StatusCode), 10))
	}

	result := make(map[string]interface{})
	err = json.Unmarshal(respBody, &result)
	if err != nil {
		log.Error("[kubeComponent's GetServiceInfo]:error when unmarshal respBody, want a json obj, got :", string(respBody), "\n ===>error is:", err.Error())
		return nil, err
	}

	return result, nil
}

func (kube *kubeComponent) Update() {
	for i := 0; i < 10; i++ {
		info, err := kube.GetPodInfo()
		if err != nil {
			log.Error("[kubeComponent's UpdatePodInfo]:error when get pod info:", err.Error())
			return
		}

		if info != nil {
			containerLogName := ""
			if metadataMap, ok := info["metadata"].(map[string]interface{}); ok {
				containerName := metadataMap["name"].(string)
				containerLogName = containerName
			}

			containerLogName = containerLogName + "_" + kube.namespace + "_" + kube.runID + "-pod-"

			if statusMap, ok := info["status"].(map[string]interface{}); ok {
				if containerStatuses, ok := statusMap["containerStatuses"].([]interface{}); ok {
					if len(containerStatuses) > 0 {
						containerInfo := containerStatuses[0].(map[string]interface{})
						if containerID, ok := containerInfo["containerID"].(string); ok {
							containerLogName = containerLogName + strings.TrimPrefix(containerID, "docker://") + ".log"
						} else {
							continue
						}
					}
				}
			}

			err := kube.componentInfo.GetActionLog().Where("id = ?", kube.componentInfo.ID).UpdateColumn("container_id", containerLogName).Error
			if err != nil {
				log.Error("[kubeComponent's UpdatePodInfo]:error when update action's container_id:", err.Error())
			}
			break
		}
		time.Sleep(1 * time.Second)
	}
}

func (kube *kubeComponent) GetPodInfo() (map[string]interface{}, error) {
	podName := "pod-" + kube.runID

	podLable := "WORKFLOW_DEFAULT_POD_LABLE%3D" + podName

	kubeReqUrl := kube.apiServerUri + "/api/v1/namespaces/" + kube.componentInfo.Namespace + "/pods?labelSelector=" + podLable

	resp, err := http.Get(kubeReqUrl)
	if err != nil {
		log.Error("[kubeComponent's GetPodInfo]:error when send req to kube:", err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("[kubeComponent's GetPodInfo]:error when read resp body info:", err.Error(), " ===>error code:", resp.StatusCode)
		return nil, errors.New("error when get service err body:" + err.Error() + "==== error code is :" + strconv.FormatInt(int64(resp.StatusCode), 10))
	}

	result := make(map[string]interface{})
	err = json.Unmarshal(respBody, &result)
	if err != nil {
		log.Error("[kubeComponent's GetPodInfo]:error when unmarshal respBody, want a json obj, got :", string(respBody), "\n ===>error is:", err.Error())
		return nil, err
	}

	pods := result["items"].([]interface{})
	if len(pods) == 0 {
		return nil, nil
	}

	podInfo := pods[0].(map[string]interface{})
	return podInfo, nil
}
