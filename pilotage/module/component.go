package module

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/containerops/pilotage/models"
)

const (
	// RUNENV_KUBE means component will run in k8s
	RUNENV_KUBE = "KUBE"
	// RUNENV_SWARM means component will run in swarm
	RUNENV_SWARM = "SWARM"
)

type component interface {
	// start a component
	// id contains pipelineId stageId actionId pipelineSequence
	Start(id string, eventList []models.EventDefinition) error
	Stop(string) (string, error)
	GetIp(...interface{}) (string, error)
}

func InitComponet(componentInfo models.ComponentLog, runenv string) (component, error) {
	if componentInfo.ID == 0 {
		return nil, errors.New("component's id is 0")
	}

	if runenv == RUNENV_KUBE {
		kubeCom := new(kubeComponent)

		ComponentConfigMap := make(map[string]interface{}, 0)
		err := json.Unmarshal([]byte(componentInfo.Kubernetes), &ComponentConfigMap)
		if err != nil {
			return kubeCom, errors.New("component kube config error:" + err.Error())
		}

		if _, ok := ComponentConfigMap["host"]; !ok {
			return kubeCom, errors.New("component kube config error ,host is not set")
		}
		host, ok := ComponentConfigMap["host"].(string)
		if !ok {
			return kubeCom, errors.New("component kube config error ,host is not a string")
		}

		if _, ok := ComponentConfigMap["port"]; !ok {
			return kubeCom, errors.New("component kube config error ,port is not set")
		}
		portF, ok := ComponentConfigMap["port"].(float64)
		if !ok {
			return kubeCom, errors.New("component kube config error ,port is not a int")
		}
		port := int64(portF)

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
		kubeCom.port = port
		kubeCom.reachableIPs = reachableIPStrs
		kubeCom.podConfig = podConfig
		kubeCom.serviceConfig = serviceConfig
		kubeCom.componentInfo = componentInfo

		return kubeCom, nil
	}

	return nil, errors.New("can't get a component in" + runenv)
}

type kubeComponent struct {
	host          string
	port          int64
	reachableIPs  []string
	podConfig     map[string]interface{}
	serviceConfig map[string]interface{}
	componentInfo models.ComponentLog
}

// start a component in kube env
func (kube *kubeComponent) Start(id string, eventList []models.EventDefinition) error {
	// get componentName for service selector
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
			return errors.New("component's kube config error, metadata is not a json!")
		}
	}

	componentName := ""
	comName, ok := metaInfoMap["name"]
	if ok {
		cName, ok := comName.(string)
		if !ok {
			return errors.New("component's kube config error, compoent name is not a string!")
		} else {
			componentName = cName + strings.Replace(id, ",", "-", -1) + "-pod"
		}
	} else {
		componentName = kube.componentInfo.Component + strings.Replace(id, ",", "-", -1) + "-pod"
	}
	if len(componentName) > 24 {
		componentName = componentName[len(componentName)-24:]
	}

	serviceAddr, err := kube.startService(id, componentName)
	if err != nil {
		return err
	}

	err = kube.startPod(id, serviceAddr, eventList)

	return err
}

// to start a pod, except podConfig,also need add some extra info
// for example , env call back url or event list need to callback or current pod's run id use to id itself when make a call back
func (kube *kubeComponent) startPod(id, serviceAddr string, eventList []models.EventDefinition) error {

	podInfoMap, err := kube.getPodInfo(id, serviceAddr, eventList)
	if err != nil {
		return err
	}

	rcName := kube.componentInfo.Component + strings.Replace(id, ",", "-", -1) + "-rc"
	if len(rcName) > 254 {
		rcName = rcName[len(rcName)-254:]
	}

	rcMap := make(map[string]interface{})
	rcMetaData := make(map[string]interface{})
	rcSpecMap := make(map[string]interface{})

	rcMetaData["name"] = rcName
	rcMap["metadata"] = rcMetaData

	rcSpecMap["replicas"] = 1
	rcSpecMap["template"] = podInfoMap
	rcMap["spec"] = rcSpecMap

	reqBody, _ := json.Marshal(rcMap)
	resp, err := http.Post(kube.componentInfo.Source+"/api/v1/namespaces/"+kube.componentInfo.Namespace+"/replicationcontrollers", "application/json", bytes.NewReader(reqBody))
	if err != nil {
		return err
	}

	if resp.StatusCode != 201 {
		return errors.New(strconv.FormatInt(int64(resp.StatusCode), 10))
	}

	return nil
}

func (kube *kubeComponent) getPodInfo(id, serviceAddr string, eventList []models.EventDefinition) (map[string]interface{}, error) {
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

	podName := kube.componentInfo.Component + strings.Replace(id, ",", "-", -1) + "-pod"
	if len(podName) > 254 {
		podName = podName[len(podName)-254:]
	}

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
		containerInfo["name"] = kube.componentInfo.Component + strings.Replace(id, ",", "-", -1) + "-pod"
		containerInfo["image"] = kube.componentInfo.Endpoint

		containers = append(containers, containerInfo)
	}

	// prepare envMap
	envMap := make([]interface{}, 0)

	eventListStr := ""
	for _, event := range eventList {
		tempEnv := make(map[string]string)
		tempEnv["name"] = event.Title
		tempEnv["value"] = event.Definition

		envMap = append(envMap, tempEnv)

		eventListStr += ";" + event.Title + "," + strconv.FormatInt(event.ID, 10)
	}

	if serviceAddr != "" {
		envMap = append(envMap, map[string]string{"name": "SERVICE_ADDR", "value": serviceAddr})
	}

	envMap = append(envMap, map[string]string{"name": "POD_NAME", "value": podName})
	envMap = append(envMap, map[string]string{"name": "RUN_ID", "value": id})
	envMap = append(envMap, map[string]string{"name": "EVENT_LIST", "value": strings.TrimPrefix(eventListStr, ";")})

	// set env to each container
	for _, container := range containers {
		if _, ok := container["name"]; !ok {
			container["name"] = kube.componentInfo.Component + strings.Replace(id, ",", "-", -1) + "-pod"
		}
		if _, ok := container["image"]; !ok {
			container["image"] = kube.componentInfo.Endpoint
		}

		env, ok := container["env"]
		if !ok {
			container["env"] = envMap
		} else {
			cEnvMap, ok := env.([]interface{})
			if !ok {
				return nil, errors.New("component's kube config error, container's env is not a array")
			}

			container["env"] = append(cEnvMap, envMap...)
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
	serviceAddr := ""

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

	serviceName := kube.componentInfo.Component + strings.Replace(id, ",", "-", -1) + "-service"
	if len(serviceName) > 24 {
		serviceName = serviceName[len(serviceName)-24:]
	}
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

	if len(kube.reachableIPs) > 0 {
		specMap["externalIPs"] = kube.reachableIPs
	}

	// create service
	reqBody, _ := json.Marshal(reqMap)
	resp, err := http.Post(kube.componentInfo.Source+"/api/v1/namespaces/"+kube.componentInfo.Namespace+"/services", "application/json", bytes.NewReader(reqBody))
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 201 {
		return "", errors.New(strconv.FormatInt(int64(resp.StatusCode), 10))
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

	for _, port := range respPorts {
		portF, ok := port["port"].(float64)
		if !ok {
			return "", errors.New("error when parse create service resp: nodePort is not a number!")
		}
		portStr := strconv.FormatFloat(portF, 'f', 0, 64)

		listenPort, ok := port["targetPort"].(float64)
		if !ok {
			return "", errors.New("error when parse create service resp: targetPort is not a number!")
		}

		listenPortStr := strconv.FormatFloat(listenPort, 'f', 0, 64)

		serviceAddr += "," + clusterIp + ":" + portStr + ":" + listenPortStr
	}

	return strings.TrimPrefix(serviceAddr, ","), nil
}

func (kube *kubeComponent) Stop(id string) (string, error) {
	client := &http.Client{}

	// set rc's replicas = 0 then delete rc
	podInfoMap, err := kube.getPodInfo(id, "", nil)
	if err != nil {
		return "", err
	}

	rcName := kube.componentInfo.Component + strings.Replace(id, ",", "-", -1) + "-rc"
	if len(rcName) > 254 {
		rcName = rcName[len(rcName)-254:]
	}

	rcMap := make(map[string]interface{})
	rcMetaData := make(map[string]interface{})
	rcSpecMap := make(map[string]interface{})

	rcMetaData["name"] = rcName
	rcMap["metadata"] = rcMetaData

	rcSpecMap["replicas"] = 0
	rcSpecMap["template"] = podInfoMap
	rcMap["spec"] = rcSpecMap

	reqBody, _ := json.Marshal(rcMap)
	req, err := http.NewRequest("PUT", kube.componentInfo.Source+"/api/v1/namespaces/"+kube.componentInfo.Namespace+"/replicationcontrollers/"+rcName, bytes.NewReader(reqBody))
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
	req, err = http.NewRequest("DELETE", kube.componentInfo.Source+"/api/v1/namespaces/"+kube.componentInfo.Namespace+"/replicationcontrollers/"+rcName, nil)
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
	serviceName := kube.componentInfo.Component + strings.Replace(id, ",", "-", -1) + "-service"
	if len(serviceName) > 24 {
		serviceName = serviceName[len(serviceName)-24:]
	}

	req, err = http.NewRequest("DELETE", kube.componentInfo.Source+"/api/v1/namespaces/"+kube.componentInfo.Namespace+"/services/"+serviceName, nil)
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
