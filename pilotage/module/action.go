package module

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Huawei/containerops/pilotage/models"

	log "github.com/Sirupsen/logrus"
	"github.com/jinzhu/gorm"
)

const (
	ActionStopReasonTimeout = "TIME_OUT"

	ActionStopReasonSendDataFailed = "SendDataFailed"

	ActionStopReasonRunSuccess = "RunSuccess"
	ActionStopReasonRunFailed  = "RunFailed"
)

var (
	actionlogAuthChan   chan bool
	actionlogListenChan chan bool
)

type Action struct {
	*models.Action
}

type ActionLog struct {
	*models.ActionLog
}

type Relation struct {
	From string
	To   string
}

func init() {
	actionlogAuthChan = make(chan bool, 1)
	actionlogListenChan = make(chan bool, 1)
}

func getActionEnvList(actionLogId int64) ([]map[string]interface{}, error) {
	resultList := make([]map[string]interface{}, 0)
	actionLog := new(models.ActionLog)
	err := actionLog.GetActionLog().Where("id = ?", actionLogId).First(actionLog).Error
	if err != nil {
		log.Error("[actionLog's getActionEnvList]:error when get actionLog info from db:", err.Error())
		return nil, errors.New("error when get action info from db:" + err.Error())
	}

	envMap := make(map[string]string)
	if actionLog.Environment != "" {
		err = json.Unmarshal([]byte(actionLog.Environment), &envMap)
		if err != nil {
			log.Error("[actionLog's getActionEnvList]:error when unmarshal action's env setting:", actionLog.Environment, " ===>error is:", err.Error())
			return nil, errors.New("error when unmarshal action's env info" + err.Error())
		}
	}

	for key, value := range envMap {
		tempEnvMap := make(map[string]interface{})
		tempEnvMap["name"] = key
		tempEnvMap["value"] = value

		resultList = append(resultList, tempEnvMap)
	}

	return resultList, nil
}

func CreateNewActions(db *gorm.DB, pipelineInfo *models.Pipeline, stageInfo *models.Stage, defineList []map[string]interface{}) (map[string]int64, error) {
	if db == nil {
		db = models.GetDB()
		db = db.Begin()
	}

	actionIdMap := make(map[string]int64)
	for _, actionDefine := range defineList {
		actionName := ""
		actionImage := ""
		kubernetesSetting := ""
		inputStr := ""
		outputStr := ""
		actionTimeout := int64(60 * 60 * 24 * 36)
		componentId := int64(0)
		serviceId := int64(0)
		platformMap := make(map[string]string)
		requestMapList := make([]interface{}, 0)

		// get component info
		component, ok := actionDefine["component"]
		if ok {
			componentMap, ok := component.(map[string]interface{})
			if !ok {
				log.Error("[action's CreateNewActions]:error when get action's component info, want a json obj, got:", component)
				return nil, errors.New("action's component is not a json")
			}

			componentVersion, ok := componentMap["versionid"].(float64)
			if !ok {
				log.Error("[action's CreateNewActions]:error when get action's component info,compoent doesn't has a versionid,component define is:", componentMap)
				return nil, errors.New("action's component info error !")
			}

			componentId = int64(componentVersion)
		}

		// get action setup data info map
		if setupDataMap, ok := actionDefine["setupData"].(map[string]interface{}); ok {
			if actionSetupDataMap, ok := setupDataMap["action"].(map[string]interface{}); ok {
				if name, ok := actionSetupDataMap["name"].(string); ok {
					actionName = name
				}

				if image, ok := actionSetupDataMap["image"].(map[string]interface{}); ok {
					actionImage = ""
					if name, ok := image["name"]; ok {
						actionImage = name.(string) + ":"
						if tag, ok := image["tag"]; ok {
							actionImage += tag.(string)
						} else {
							actionImage += "latest"
						}
					}
				}

				if timeoutStr, ok := actionSetupDataMap["timeout"].(string); ok {
					var err error
					if actionTimeout, err = strconv.ParseInt(timeoutStr, 10, 64); err != nil {
						log.Error("[action's CreateNewActions]:error when get action's timeout value,want a string, got:", timeoutStr)
						return nil, errors.New("action's timeout is not string")
					}
				}

				configMap := make(map[string]interface{})
				// record platform info
				if platFormType, ok := actionSetupDataMap["type"].(string); ok {
					platformMap["platformType"] = strings.ToUpper(platFormType)
				}

				if platformHost, ok := actionSetupDataMap["apiserver"].(string); ok {
					platformHost = strings.TrimSuffix(platformHost, "/")
					platformMap["platformHost"] = platformHost
				}

				if ip, ok := actionSetupDataMap["ip"].(string); ok {
					configMap["nodeIP"] = ip
				}

				// unmarshal k8s info
				if useAdvanced, ok := actionSetupDataMap["useAdvanced"].(bool); ok {
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

					kuberSettingBytes, _ := json.Marshal(configMap)
					kubernetesSetting = string(kuberSettingBytes)
				}
			}
		}

		inputMap, ok := actionDefine["inputJson"].(map[string]interface{})
		if ok {
			inputDescribe, err := describeJSON(inputMap, "")
			if err != nil {
				log.Error("[action's CreateNewActions]:error when describe action's input json define:", inputMap, " ===>error is:", err.Error())
				return nil, errors.New("error in component output json define:" + err.Error())
			}

			inputDescBytes, _ := json.Marshal(inputDescribe)
			inputStr = string(inputDescBytes)
		}

		outputMap, ok := actionDefine["outputJson"].(map[string]interface{})
		if ok {
			outputDescribe, err := describeJSON(outputMap, "")
			if err != nil {
				log.Error("[action's CreateNewActions]:error when describe action's output json define:", inputMap, " ===>error is:", err.Error())
				return nil, errors.New("error in component output json define:" + err.Error())
			}

			outputDescBytes, _ := json.Marshal(outputDescribe)
			outputStr = string(outputDescBytes)
		}

		allEnvMap := make(map[string]string)
		if envMap, ok := actionDefine["env"].([]interface{}); ok {
			for _, envInfo := range envMap {
				envInfoMap, ok := envInfo.(map[string]interface{})
				if !ok {
					log.Error("[action's CreateNewActions]:error when get action's env setting, want a json obj,got:", envInfo)
					return nil, errors.New("action's env set is not a json")
				}

				key, ok := envInfoMap["key"].(string)
				if !ok {
					log.Error("[action's CreateNewActions]:error when get action's env setting, want string key,got:", envInfoMap)
					return nil, errors.New("action's key is not a string")
				}

				value, ok := envInfoMap["value"].(string)
				if !ok {
					log.Error("[action's CreateNewActions]:error when get action's env setting, want string value,got:", envInfoMap)
					return nil, errors.New("action's value is not a string")
				}
				allEnvMap[key] = value
			}
		}
		envBytes, _ := json.Marshal(allEnvMap)

		// get aciont line info
		actionId, ok := actionDefine["id"].(string)
		if !ok {
			log.Error("[action's CreateNewActions]:error when action's id from action define, want string, got:", actionDefine)
			return nil, errors.New("action's id is not a string")
		}

		manifestMap := make(map[string]interface{})
		manifestMap["platform"] = platformMap
		manifestBytes, _ := json.Marshal(manifestMap)

		stageRequest, ok := actionDefine["request"].([]interface{})
		if !ok {
			defaultRequestMap := make(map[string]interface{})
			defaultRequestMap["type"] = AuthTyptStageStartDone
			defaultRequestMap["token"] = AuthTokenDefault

			requestMapList = append(requestMapList, defaultRequestMap)
		} else {
			requestMapList = stageRequest
		}
		requestInfos, _ := json.Marshal(requestMapList)

		action := new(models.Action)
		action.Namespace = pipelineInfo.Namespace
		action.Repository = pipelineInfo.Repository
		action.Pipeline = stageInfo.Pipeline
		action.Stage = stageInfo.ID
		action.Component = componentId
		action.Service = serviceId
		action.Action = actionName
		action.Title = actionName
		action.Description = actionName
		action.Manifest = string(manifestBytes)
		action.Environment = string(envBytes)
		action.Kubernetes = kubernetesSetting
		action.Input = inputStr
		action.Output = outputStr
		action.Endpoint = actionImage
		action.Timeout = actionTimeout
		action.Requires = string(requestInfos)

		err := db.Model(&models.Action{}).Save(action).Error
		if err != nil {
			log.Error("[action's CreateNewActions]:error when save action info to db:", err.Error())
			rollbackErr := db.Rollback().Error
			if rollbackErr != nil {
				log.Error("[action's CreateNewActions]:when rollback in save action's info:", rollbackErr.Error())
				return nil, errors.New("errors occur:\nerror1:" + err.Error() + "\nerror2:" + rollbackErr.Error())
			}
			return nil, errors.New("error when save action info to db:" + err.Error())
		}
		actionIdMap[actionId] = action.ID
	}

	return actionIdMap, nil
}

func GetActionLog(actionLogId int64) (*ActionLog, error) {
	action := new(ActionLog)
	actionLog := new(models.ActionLog)
	err := actionLog.GetActionLog().Where("id = ?", actionLogId).First(actionLog).Error
	if err != nil {
		log.Error("[actionLog's GetActionLog]:error when get action log info from db:", err.Error())
		return nil, err
	}

	action.ActionLog = actionLog
	return action, nil
}

func GetActionLogByName(namespace, repository, pipelineName string, sequence int64, stageName, actionName string) (*ActionLog, error) {
	action := new(ActionLog)
	pipelineLog := new(models.PipelineLog)
	stageLog := new(models.StageLog)
	actionLog := new(models.ActionLog)

	err := pipelineLog.GetPipelineLog().Where("namespace = ?", namespace).Where("repository = ?", repository).Where("pipeline = ?", pipelineName).Where("sequence = ?", sequence).First(pipelineLog).Error
	if err != nil {
		if err != nil {
			log.Error("[actionLog's GetActionLog]:error when get pipelineLog info from db:", err.Error())
			return nil, err
		}
	}

	err = stageLog.GetStageLog().Where("namespace = ?", namespace).Where("repository = ?", repository).Where("pipeline = ?", pipelineLog.ID).Where("sequence = ?", sequence).Where("stage = ?", stageName).First(stageLog).Error
	if err != nil {
		if err != nil {
			log.Error("[actionLog's GetActionLog]:error when get stageLog info from db:", err.Error())
			return nil, err
		}
	}

	err = actionLog.GetActionLog().Where("namespace = ?", namespace).Where("repository = ?", repository).Where("pipeline = ?", pipelineLog.ID).Where("sequence = ?", sequence).Where("stage = ?", stageLog.ID).Where("action = ?", actionName).First(actionLog).Error
	if err != nil {
		log.Error("[actionLog's GetActionLog]:error when get action log info from db:", err.Error())
		return nil, err
	}

	action.ActionLog = actionLog
	return action, nil
}

func (actionInfo *Action) GenerateNewLog(db *gorm.DB, pipelineLog *models.PipelineLog, stageLog *models.StageLog) error {
	if db == nil {
		db = models.GetDB()
		err := db.Begin().Error
		if err != nil {
			log.Error("[action's GenerateNewLog]:when db.Begin():", err.Error())
			return err
		}
	}

	// record action's info
	actionLog := new(models.ActionLog)
	actionLog.Namespace = actionInfo.Namespace
	actionLog.Repository = actionInfo.Repository
	actionLog.Pipeline = pipelineLog.ID
	actionLog.FromPipeline = pipelineLog.FromPipeline
	actionLog.Sequence = pipelineLog.Sequence
	actionLog.Stage = stageLog.ID
	actionLog.FromStage = stageLog.FromStage
	actionLog.FromAction = actionInfo.ID
	actionLog.RunState = models.ActionLogStateCanListen
	actionLog.Component = actionInfo.Component
	actionLog.Service = actionInfo.Service
	actionLog.Action = actionInfo.Action.Action
	actionLog.Title = actionInfo.Title
	actionLog.Description = actionInfo.Description
	actionLog.Event = actionInfo.Event
	actionLog.Manifest = actionInfo.Manifest
	actionLog.Environment = actionInfo.Environment
	actionLog.Kubernetes = actionInfo.Kubernetes
	actionLog.Swarm = actionInfo.Swarm
	actionLog.Input = actionInfo.Input
	actionLog.Output = actionInfo.Output
	actionLog.Endpoint = actionInfo.Endpoint
	actionLog.Timeout = actionInfo.Timeout
	actionLog.Requires = actionInfo.Requires
	actionLog.AuthList = ""

	err := db.Save(actionLog).Error
	if err != nil {
		log.Error("[action's GenerateNewLog]:when save action log to db:", actionLog, " ===>error is:", err.Error())
		rollbackErr := db.Rollback().Error
		if rollbackErr != nil {
			log.Error("[action's GenerateNewLog]:when rollback in save action log:", rollbackErr.Error())
			return errors.New("errors occur:\nerror1:" + err.Error() + "\nerror2:" + rollbackErr.Error())
		}
		return err
	}

	err = setSystemEvent(db, actionLog)
	if err != nil {
		log.Error("[action's GenerateNewLog]:when save action log to db:", err.Error())
		return err
	}

	return nil
}

func (actionLog *ActionLog) GetActionLineInfo() ([]map[string]interface{}, error) {
	lineList := make([]map[string]interface{}, 0)

	manifestMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(actionLog.Manifest), &manifestMap)
	if err != nil {
		log.Error("[actionLog's GetActionLineInfo]:error when unmarshal action's manifest, want a json obj,got:", actionLog.Manifest)
		return nil, err
	}

	relationList, ok := manifestMap["relation"].([]interface{})
	if !ok {
		log.Error("[actionLog's GetActionLineInfo]:error when get action's relation from action's manifestMap:", manifestMap, manifestMap["relation"].(int64))
		return lineList, nil
	}

	for _, relation := range relationList {
		relationInfo, ok := relation.(map[string]interface{})
		if err != nil {
			log.Error("[actionLog's GetActionLineInfo]:error when get action's relation info, want a json obj,got:", relation)
			continue
		}

		fromRealActionIDF, ok := relationInfo["fromAction"].(float64)
		if !ok {
			log.Error("[actionLog's GetActionLineInfo]:error when get fromRealActionID from action's relation,want a number,got:", relationInfo["fromAction"])
			continue
		}

		fromRealActionID := int64(fromRealActionIDF)
		fromActionInfoMap := make(map[string]string)
		if fromRealActionID == int64(0) {
			// if action's id == 0 ,this is a realtion from pipeline's start stage
			startStage := new(models.StageLog)
			err := startStage.GetStageLog().Where("namespace = ?", actionLog.Namespace).Where("repository = ?", actionLog.Repository).Where("pipeline = ?", actionLog.Pipeline).Where("type = ?", models.StageTypeStart).First(startStage).Error
			if err != nil {
				log.Error("[actionLog's GetActionLineInfo]:error when get pipline's start stage from db:", err.Error())
				continue
			}

			fromActionInfoMap["id"] = "s-" + strconv.FormatInt(startStage.ID, 10)
			fromActionInfoMap["type"] = models.StageTypeForWeb[startStage.Type]
		} else {
			fromActionInfo := new(models.ActionLog)
			err = fromActionInfo.GetActionLog().Where("namespace = ?", actionLog.Namespace).Where("repository = ?", actionLog.Repository).Where("pipeline = ?", actionLog.Pipeline).Where("sequence = ?", actionLog.Sequence).Where("from_action = ?", fromRealActionID).First(fromActionInfo).Error
			if err != nil {
				log.Error("[actionLog's GetActionLineInfo]:error when get preActionlog info from db:", err.Error())
				continue
			}

			fromActionInfoMap["id"] = "a-" + strconv.FormatInt(fromActionInfo.ID, 10)
			fromActionInfoMap["type"] = "pipeline-action"
		}

		toActionInfoMap := make(map[string]string)
		toActionInfoMap["id"] = "a-" + strconv.FormatInt(actionLog.ID, 10)
		toActionInfoMap["type"] = "pipeline-action"
		toActionInfoMap["name"] = actionLog.Action

		lineMap := make(map[string]interface{})
		lineMap["id"] = fromActionInfoMap["id"] + "-" + toActionInfoMap["id"]
		lineMap["pipelineLineViewId"] = "pipeline-line-view"

		lineMap["startData"] = map[string]string{
			"id":   fromActionInfoMap["id"],
			"type": fromActionInfoMap["type"],
		}

		lineMap["endData"] = map[string]interface{}{
			"id": toActionInfoMap["id"],
			"setupData": map[string]interface{}{
				"action": map[string]string{
					"name": toActionInfoMap["name"],
				},
			},
		}

		lineList = append(lineList, lineMap)
	}

	return lineList, nil
}

func (actionLog *ActionLog) GetActionHistoryInfo() (map[string]interface{}, error) {
	result := make(map[string]interface{})

	inputMap, err := actionLog.GetInputData()
	if err != nil {
		log.Error("[actionLog's GetActionHistoryInfo]:error when get action's input data:", err.Error())
		return nil, err
	}

	outputMap, err := actionLog.GetOutputData()
	if err != nil {
		log.Error("[actionLog's GetActionHistoryInfo]:error when get action's output data:", err.Error())
		return nil, err
	}

	dataMap := make(map[string]interface{})
	dataMap["input"] = inputMap
	dataMap["output"] = outputMap

	logList := make([]models.Event, 0)
	err = new(models.Event).GetEvent().Where("namespace = ?", actionLog.Namespace).Where("repository = ?", actionLog.Repository).Where("pipeline = ?", actionLog.Pipeline).Where("stage = ?", actionLog.Stage).Where("action = ?", actionLog.ID).Order("id").Find(&logList).Error
	if err != nil && !strings.Contains(err.Error(), "record not found") {
		log.Error("[actionLog's GetActionHistoryInfo]:error when get actionlog's log from db:", err.Error())
		return nil, err
	}

	logListStr := make([]string, 0)
	for _, log := range logList {
		logStr := log.CreatedAt.Format("2006-01-02 15:04:05") + " -> " + log.Payload

		logListStr = append(logListStr, logStr)
	}

	result["data"] = dataMap
	result["logList"] = logListStr
	return result, nil
}

func (actionLog *ActionLog) GetInputData() (map[string]interface{}, error) {
	inputMap := make(map[string]interface{})

	inputInfo := new(models.Event)
	err := inputInfo.GetEvent().Where("namespace = ?", actionLog.Namespace).Where("repository = ?", actionLog.Repository).Where("pipeline = ?", actionLog.Pipeline).Where("sequence = ?", actionLog.Sequence).Where("action = ?", actionLog.ID).Where("title = ?", "SEND_DATA").First(inputInfo).Error
	if err != nil && !strings.Contains(err.Error(), "record not found") {
		log.Error("[actionLog's GetInputData]:error when get actionlog's input info from db:", err.Error())
		return nil, err
	}

	err = json.Unmarshal([]byte(inputInfo.Payload), &inputMap)
	if err != nil {
		log.Error("[actionLog's GetInputData]:error when unmarshal input info:", inputInfo.Payload, " ===>error is:"+err.Error())
	}

	inputStr, ok := inputMap["data"].(string)
	if !ok {
		log.Error("[actionLog's GetInputData]:error when get inputMap str, want a string, got:", inputMap["data"])
		inputStr = ""
	}

	realinputMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(inputStr), &realinputMap)
	if err != nil {
		log.Error("[actionLog's GetInputData]:error when unmarshal real input info:", inputStr, "===>error is:", err.Error())
	}

	return inputMap, nil
}

func (actionLog *ActionLog) GetOutputData() (map[string]interface{}, error) {
	outputMap := make(map[string]interface{})

	outputInfo := new(models.Outcome)
	err := outputInfo.GetOutcome().Where("pipeline = ?", actionLog.Pipeline).Where("sequence = ?", actionLog.Sequence).Where("stage = ?", actionLog.Stage).Where("action = ?", actionLog.ID).First(outputInfo).Error
	if err != nil && !strings.Contains(err.Error(), "record not found") {
		log.Error("[actionLog's GetOutputData]:error when get actionlog's output info from db:", err.Error())
		return nil, err
	}

	err = json.Unmarshal([]byte(outputInfo.Output), &outputMap)
	if err != nil {
		log.Error("[actionLog's GetOutputData]:error when unmarshal output info:", outputInfo.Output, "===>error is:", err.Error())
	}

	return outputMap, nil
}

func (actionLog *ActionLog) Listen() error {
	actionlogListenChan <- true
	defer func() { <-actionlogListenChan }()

	err := actionLog.GetActionLog().Where("id = ?", actionLog.ID).First(actionLog).Error
	if err != nil {
		log.Error("[actionLog's Listen]:error when get action info from db:", actionLog, " ===>error is:", err.Error())
		return errors.New("error when get actionlog's info from db:" + err.Error())
	}

	if actionLog.RunState != models.ActionLogStateCanListen {
		log.Error("[actionLog's Listen]:error actionlog state:", actionLog)
		return errors.New("can't listen curren actionlog,current state is:" + strconv.FormatInt(actionLog.RunState, 10))
	}

	actionLog.RunState = models.ActionLogStateWaitToStart
	err = actionLog.GetActionLog().Save(actionLog).Error
	if err != nil {
		log.Error("[actionLog's Listen]:error when change actionLog's run state to wait to start:", actionLog, " ===>error is:", err.Error())
		return errors.New("can't listen target action,change action's state failed")
	}

	canStartChan := make(chan bool, 1)
	go func() {
		for true {
			time.Sleep(1 * time.Second)

			err := actionLog.GetActionLog().Where("id = ?", actionLog.ID).First(actionLog).Error
			if err != nil {
				log.Error("[actionLog's Listen]:error when get actionLog's info:", actionLog, " ===>error is:", err.Error())
				canStartChan <- false
				break
			}
			if actionLog.Requires == "" || actionLog.Requires == "[]" {
				log.Info("[actionLog's Listen]:actionLog", actionLog, " is ready and will start")
				canStartChan <- true
				break
			}
		}
	}()

	go func() {
		canStart := <-canStartChan
		if !canStart {
			log.Error("[actionLog's Listen]:actionLog can't start", actionLog)
			actionLog.Stop(StageStopReasonRunFailed, models.ActionLogStateRunFailed)
			return
		}

		go actionLog.Start()
	}()

	return nil
}

func (actionLog *ActionLog) Auth(authMap map[string]interface{}) error {
	actionlogAuthChan <- true
	defer func() { <-actionlogAuthChan }()

	authType, ok := authMap["type"].(string)
	if !ok {
		log.Error("[actionLog's Auth]:error when get authType from given authMap:", authMap, " ===>to actionLog:", actionLog)
		return errors.New("authType is illegal")
	}

	token, ok := authMap["token"].(string)
	if !ok {
		log.Error("[actionLog's Auth]:error when get token from given authMap:", authMap, " ===>to actionLog:", actionLog)
		return errors.New("token is illegal")
	}

	err := actionLog.GetActionLog().Where("id = ?", actionLog.ID).First(actionLog).Error
	if err != nil {
		log.Error("[actionLog's Auth]:error when get actionLog info from db:", actionLog, " ===>error is:", err.Error())
		return errors.New("error when get stagelog's info from db:" + err.Error())
	}

	if actionLog.Requires == "" || actionLog.Requires == "[]" {
		log.Error("[actionLog's Auth]:error when set auth info,actionLog's requires is empty", authMap, " ===>to actionLog:", actionLog)
		return errors.New("action don't need any more auth")
	}

	requireList := make([]interface{}, 0)
	remainRequireList := make([]interface{}, 0)
	err = json.Unmarshal([]byte(actionLog.Requires), &requireList)
	if err != nil {
		log.Error("[actionLog's Auth]:error when unmarshal actionLog's require list:", actionLog, " ===>error is:", err.Error())
		return errors.New("error when get action require auth info:" + err.Error())
	}

	hasAuthed := false
	for _, require := range requireList {
		requireMap, ok := require.(map[string]interface{})
		if !ok {
			log.Error("[actionLog's Auth]:error when get actionLog's require info:", actionLog, " ===> require is:", require)
			return errors.New("error when get actionLog require auth info,require is not a json object")
		}

		requireType, ok := requireMap["type"].(string)
		if !ok {
			log.Error("[actionLog's Auth]:error when get actionLog's require type:", actionLog, " ===> require map is:", requireMap)
			return errors.New("error when get action require auth info,require don't have a type")
		}

		requireToken, ok := requireMap["token"].(string)
		if !ok {
			log.Error("[actionLog's Auth]:error when get actionLog's require token:", actionLog, " ===> require map is:", requireMap)
			return errors.New("error when get action require auth info,require don't have a token")
		}

		if requireType == authType && requireToken == token {
			hasAuthed = true
			// record auth info to actionLog's Auth info list
			actionLogAuthList := make([]interface{}, 0)
			if actionLog.AuthList != "" {
				err = json.Unmarshal([]byte(actionLog.AuthList), &actionLogAuthList)
				if err != nil {
					log.Error("[actionLog's Auth]:error when unmarshal actionLog's Auth list:", actionLog, " ===>error is:", err.Error())
					return errors.New("error when set auth info to action")
				}
			}

			actionLogAuthList = append(actionLogAuthList, authMap)

			authListInfo, err := json.Marshal(actionLogAuthList)
			if err != nil {
				log.Error("[actionLog's Auth]:error when marshal actionLog's Auth list:", actionLogAuthList, " ===>error is:", err.Error())
				return errors.New("error when save action auth info")
			}

			actionLog.AuthList = string(authListInfo)
			err = actionLog.GetActionLog().Save(actionLog).Error
			if err != nil {
				log.Error("[actionLog's Auth]:error when save actionLog's info to db:", actionLog, " ===>error is:", err.Error())
				return errors.New("error when save action auth info")
			}
		} else {
			remainRequireList = append(remainRequireList, requireMap)
		}
	}

	if !hasAuthed {
		log.Error("[actionLog's Auth]:error when auth a actionLog to start, given auth:", authMap, " is not equal to any request one:", actionLog.Requires)
		return errors.New("illegal auth info, auth failed")
	}

	remainRequireAuthInfo, err := json.Marshal(remainRequireList)
	if err != nil {
		log.Error("[actionLog's Auth]:error when marshal actionLog's remainRequireAuth list:", remainRequireList, " ===>error is:", err.Error())
		return errors.New("error when sync remain require auth info")
	}

	actionLog.Requires = string(remainRequireAuthInfo)
	err = actionLog.GetActionLog().Save(actionLog).Error
	if err != nil {
		log.Error("[actionLog's Auth]:error when save actionLog's remain require auth info:", actionLog, " ===>error is:", err.Error())
		return errors.New("error when sync remain require auth info")
	}

	return nil
}

func (actionLog *ActionLog) Start() {
	if actionLog.Timeout != 0 {
		go actionLog.WaitActionDone()
	}

	if actionLog.Component != 0 {
		c, err := InitComponetNew(actionLog)
		if err != nil {
			log.Error("[actionLog's Start]:error when init component:", err.Error())
			actionLog.Stop(ActionStopReasonRunFailed, models.ActionLogStateRunFailed)
			return
		}

		err = c.Start()
		if err != nil {
			log.Error("[actionLog's Start]:error when start component:", err.Error())
			RecordOutcom(actionLog.Pipeline, actionLog.FromPipeline, actionLog.Stage, actionLog.FromStage, actionLog.ID, actionLog.FromAction, actionLog.Sequence, 0, false, "start action error", err.Error())
			actionLog.Stop(ActionStopReasonRunFailed, models.ActionLogStateRunFailed)
		}
	} else if actionLog.Service != 0 {
		log.Info("[actionLog's Start]:start an action that use service:", actionLog)
		RecordOutcom(actionLog.Pipeline, actionLog.FromPipeline, actionLog.Stage, actionLog.FromStage, actionLog.ID, actionLog.FromAction, actionLog.Sequence, 0, false, "start action error", "use service but not support")
		actionLog.Stop(ActionStopReasonRunSuccess, models.ActionLogStateRunSuccess)
	} else {
		log.Error("[actionLog's Start]:error when start action,action doesn't spec a type", actionLog)
		RecordOutcom(actionLog.Pipeline, actionLog.FromPipeline, actionLog.Stage, actionLog.FromStage, actionLog.ID, actionLog.FromAction, actionLog.Sequence, 0, false, "start action error", "action doesn't spec a component or a service")
		actionLog.Stop(ActionStopReasonRunFailed, models.ActionLogStateRunFailed)
	}
}

func (actionLog *ActionLog) Stop(reason string, runState int64) {
	err := actionLog.GetActionLog().Where("id = ?", actionLog.ID).First(actionLog).Error
	if err != nil {
		log.Error("[actionLog's Stop]:error when get actionLog's info from db:", err.Error())
		return
	}

	if actionLog.RunState == models.ActionLogStateRunFailed || actionLog.RunState == models.ActionLogStateRunSuccess {
		return
	}

	actionLog.RunState = runState
	err = actionLog.GetActionLog().Save(actionLog).Error
	if err != nil {
		log.Error("[actionLog's Stop]:error when change action state:", actionLog, " ===>error is:", err.Error())
		return
	}

	if actionLog.Component != 0 {
		c, err := InitComponetNew(actionLog)
		if err != nil {
			log.Error("[actionLog's Stop]:error when init component:", err.Error())
			return
		}

		err = c.Stop()
		if err != nil {
			log.Error("[actionLog's Stop]:error when stop component:", err.Error())
		}
	} else if actionLog.Service != 0 {
		log.Info("[actionLog's Stop]:stop an action that use service:", actionLog)
	} else {
		log.Error("[actionLog's Stop]:error when stop action,action doesn't spec a type", actionLog)
	}
}

func (actionLog *ActionLog) RecordEvent(eventId int64, eventKey string, reqBody map[string]interface{}, headerInfo http.Header) error {
	c, err := InitComponetNew(actionLog)
	if err != nil {
		recordErr := RecordOutcom(actionLog.Pipeline, actionLog.FromPipeline, actionLog.Stage, actionLog.FromStage, actionLog.ID, actionLog.FromAction, actionLog.Sequence, eventId, false, "component init error:"+err.Error(), "")
		if recordErr != nil {
			log.Error("[actionLog's RecordEvent]:error when record outcome info:", recordErr.Error())
			return recordErr
		}

		log.Error("[actionLog's RecordEvent]:error when get action's platformInfo:", actionLog, " ===>error is:", err.Error())
		return err
	}

	if eventKey == models.EVENT_TASK_RESULT {
		resultReqBody, ok := reqBody["INFO"].(map[string]interface{})
		if !ok {
			log.Error("[actionLog's RecordEvent]:error when get request's info body, want a json obj, got:", reqBody["INFO"])
			return errors.New("request body's info is not a json obj")
		}

		status, ok := resultReqBody["status"].(bool)
		if !ok {
			status = false
		}

		result, ok := resultReqBody["result"].(string)
		if !ok {
			result = ""
		}

		output, ok := resultReqBody["output"].(map[string]interface{})
		outputStr := ""
		if !ok {
			outputStr = ""
		} else {
			outputBytes, _ := json.Marshal(output)
			outputStr = string(outputBytes)
		}

		recordErr := RecordOutcom(actionLog.Pipeline, actionLog.FromPipeline, actionLog.Stage, actionLog.FromStage, actionLog.ID, actionLog.FromAction, actionLog.Sequence, eventId, status, result, outputStr)
		if recordErr != nil {
			log.Error("[actionLog's RecordEvent]:error when record outcome info:", recordErr.Error())
			return recordErr
		}

		stopStatus := models.ActionLogStateRunFailed
		stopReason := ActionStopReasonRunFailed
		if status {
			stopStatus = models.ActionLogStateRunSuccess
			stopReason = ActionStopReasonRunSuccess
		}
		actionLog.Stop(stopReason, int64(stopStatus))
	}

	if eventKey == models.EVENT_COMPONENT_STOP {
		c.Stop()
	}

	headerMap := make(map[string]interface{})
	for key, value := range headerInfo {
		headerMap[key] = value
	}
	headerBytes, _ := json.Marshal(headerMap)

	eventDefine := new(models.EventDefinition)
	err = eventDefine.GetEventDefinition().Where("id = ?", eventId).First(&eventDefine).Error
	if err != nil {
		log.Error("[actionLog's RecordEvent]:error when get eventDefine from db:", err.Error())
		return err
	}

	authStr := ""
	auths, ok := headerMap["Authorization"].([]string)
	if !ok {
		authStr = ""
	} else {
		authStr = strings.Join(auths, ";")
	}

	bodyBytes, _ := json.Marshal(reqBody)

	// log evnet
	err = RecordEventInfo(eventId, actionLog.Sequence, string(headerBytes), string(bodyBytes), authStr)
	if err != nil {
		log.Error("[actionLog's RecordEvent]:error when save event to db:", err.Error())
		return err
	}

	return nil
}

func (actionLog *ActionLog) SendDataToAction(targetUrl string) {
	manifestMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(actionLog.Manifest), &manifestMap)
	if err != nil {
		log.Error("[actionLog's SendDataToAction]:error when get action manifest info:" + err.Error())
		return
	}

	dataMap := make(map[string]interface{})
	relations, ok := manifestMap["relation"]
	if ok {
		relationInfo, ok := relations.([]interface{})
		if !ok {
			log.Error("[actionLog's SendDataToAction]:error when parse relations,want an array,got:", relations)
			return
		}

		dataMap, err = actionLog.merageFromActionsOutputData(relationInfo)
		if err != nil {
			log.Error("[actionLog's SendDataToAction]:error when get data map from action: " + err.Error())
		}
	}

	log.Info("[actionLog's SendDataToAction]:action", actionLog, " got data:", dataMap)

	var dataByte []byte

	if len(dataMap) == 0 {
		dataByte = make([]byte, 0)
	} else {
		dataByte, err = json.Marshal(dataMap)
		if err != nil {
			log.Error("[actionLog's SendDataToAction]:error when marshal dataMap:", dataMap, " ===>error is:", err.Error())
			return
		}
	}

	character := int(0)
	// send data to component or service
	resps := make([]*http.Response, 0)
	if actionLog.Component != 0 {
		character = models.CharacterComponentEvent
		resps, err = actionLog.sendDataToComponent(targetUrl, dataByte)
	} else {
		character = models.CharacterServiceEvent
		resps, err = actionLog.sendDataToService(dataByte)
	}

	resultStr := ""
	status := false
	payload := make(map[string]interface{})
	if err != nil {
		resultStr = err.Error()
		status = false
		go actionLog.Stop(ActionStopReasonSendDataFailed, models.ActionLogStateRunFailed)
	} else {
		respMap := make(map[int64]string, len(resps))
		for count, resp := range resps {
			if resp != nil {
				respBody, _ := ioutil.ReadAll(resp.Body)
				respStr := string(respBody)

				respMap[int64(count)] = respStr
			}
		}

		result, _ := json.Marshal(respMap)
		resultStr = string(result)
		status = true
	}

	payload["EVENT"] = "SEND_DATA"
	payload["EVENTID"] = "SEND_DATA"
	payload["INFO"] = map[string]interface{}{"output": string(dataByte), "result": resultStr, "status": status}
	payload["RUN_ID"] = strconv.FormatInt(actionLog.Pipeline, 10) + "-" + strconv.FormatInt(actionLog.Stage, 10) + "-" + strconv.FormatInt(actionLog.ID, 10)

	payloadInfo, marshalErr := json.Marshal(payload)
	if marshalErr != nil {
		go actionLog.Stop(ActionStopReasonSendDataFailed, models.ActionLogStateRunFailed)
		log.Error("[actionLog's SendDataToAction]:error when marshal payload info:" + marshalErr.Error())
	}

	if err != nil {
		go actionLog.Stop(ActionStopReasonSendDataFailed, models.ActionLogStateRunFailed)
		log.Error("[actionLog's SendDataToAction]:error when send data to action:" + err.Error())
	}

	err = RecordEventInfo(models.EventDefineIDSendDataToAction, actionLog.Sequence, "", string(payloadInfo), "", "SEND_DATA", strconv.FormatInt(int64(character), 10), actionLog.Namespace, actionLog.Repository, strconv.FormatInt(actionLog.Pipeline, 10), strconv.FormatInt(actionLog.Stage, 10), strconv.FormatInt(actionLog.ID, 10))
	if err != nil {
		go actionLog.Stop(ActionStopReasonSendDataFailed, models.ActionLogStateRunFailed)
		log.Error("[actionLog's SendDataToAction]:error when save send data info :" + err.Error())
	}
}

func (actionLog *ActionLog) merageFromActionsOutputData(relationInfo []interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	for _, relation := range relationInfo {
		relationMap, ok := relation.(map[string]interface{})
		if !ok {
			log.Println("[actionLog's merageFromActionsOutputData]:error when get relation info:want json obj,got:", relation)
			return nil, errors.New("error when parse relation info,relation is not a json!")
		}

		fromOutcome := new(models.Outcome)
		err := fromOutcome.GetOutcome().Where("real_pipeline = ?", actionLog.FromPipeline).Where("pipeline = ? ", actionLog.Pipeline).Where("real_action = ?", relationMap["fromAction"]).First(&fromOutcome).Order("-id").Error
		if err != nil {
			log.Error("[actionLog's merageFromActionsOutputData]:error when get request action's output from db:want get action(", actionLog.ID, ")'s output, ===>error is:", err.Error())
			return nil, errors.New("error when get from outcome, error:" + err.Error())
		}

		tempData := make(map[string]interface{})
		err = json.Unmarshal([]byte(fromOutcome.Output), &tempData)
		if err != nil {
			log.Error("[actionLog's merageFromActionsOutputData]:error when unmarshal action(", actionLog.ID, ")'s output info:", fromOutcome.Output, " ===>error is:", err.Error())
			return nil, errors.New("error when parse from action data1:" + err.Error() + "\n" + fromOutcome.Output)
		}

		relationArray, ok := relationMap["relation"].([]interface{})
		if !ok {
			log.Error("[actionLog's merageFromActionsOutputData]:error when get relation from relationMap:", relationMap)
			return nil, errors.New("relation doesn't have a relation info")
		}

		relationList := make([]Relation, 0)
		if len(relationArray) > 0 {
			for _, realationDefines := range relationArray {
				relationByte, err := json.Marshal(realationDefines)
				if err != nil {
					log.Error("[actionLog's merageFromActionsOutputData]:error went marshal relation array:", realationDefines, " ===>error is:", err.Error())
					return nil, errors.New("error when marshal relation array:" + err.Error())
				}

				var r Relation
				err = json.Unmarshal(relationByte, &r)
				if err != nil {
					log.Error("[actionLog's merageFromActionsOutputData]:error when parse relation info:", string(relationByte), " ===>error is:", err.Error())
					return nil, errors.New("error when parse relation info:" + err.Error())
				}

				relationList = append(relationList, r)
			}
		}

		actionResult := make(map[string]interface{})
		err = getResultFromRelation(fromOutcome.Output, relationList, actionResult)
		if err != nil {
			log.Error("[actionLog's merageFromActionsOutputData]:error when get result from action's relation:", err.Error())
			return nil, errors.New("error when get from data:" + err.Error())
		}

		for key, value := range actionResult {
			result[key] = value
		}
	}

	return result, nil
}

func (actionLog *ActionLog) sendDataToComponent(targetUrl string, data []byte) ([]*http.Response, error) {
	c, err := InitComponetNew(actionLog)
	if err != nil {
		log.Error("[actionLog's sendDataToComponent]:error when init component info:", err.Error())
		return nil, errors.New("error when init component info:" + err.Error())
	}

	log.Info("start send data to component...")
	return c.SendData(targetUrl, data)
}

func (actionLog *ActionLog) sendDataToService(data []byte) ([]*http.Response, error) {
	return nil, nil
}

func (actionLog *ActionLog) WaitActionDone() {
	canStop := false
	actionRunResultChan := make(chan bool, 1)
	go func() {
		for !canStop {
			actionLogInfo := new(models.ActionLog)
			err := actionLogInfo.GetActionLog().Where("id = ?", actionLog.ID).First(actionLogInfo).Error
			if err != nil {
				log.Error("[actionLog's WaitActionDone]:error when get actionLog's info from db:", err.Error())
				actionRunResultChan <- false
				return
			}

			if actionLogInfo.RunState == models.ActionLogStateRunFailed {
				actionRunResultChan <- false
				return
			} else if actionLogInfo.RunState == models.ActionLogStateRunSuccess {
				actionRunResultChan <- true
				return
			}

			time.Sleep(1 * time.Second)
		}
	}()

	duration, _ := time.ParseDuration(strconv.FormatInt(actionLog.Timeout, 10) + "s")
	select {
	case <-time.After(duration):
		canStop = true
		actionLog.Stop(ActionStopReasonTimeout, models.ActionLogStateRunFailed)
	case runSuccess := <-actionRunResultChan:
		if runSuccess {
			actionLog.Stop(ActionStopReasonRunSuccess, models.ActionLogStateRunSuccess)
		} else {
			actionLog.Stop(ActionStopReasonRunFailed, models.ActionLogStateRunFailed)
		}
	}
}

func (actionLog *ActionLog) GetActionPlatformInfo() (map[string]string, error) {
	manifestMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(actionLog.Manifest), &manifestMap)
	if err != nil {
		log.Error("[actionLog's GetActionPlatformInfo]:error when unmarshal action's manifest:", actionLog.Manifest, " ===>error is:", err.Error())
		return nil, errors.New("action " + actionLog.ActionLog.Action + "'s manifest is illegal")
	}

	platformSetting, ok := manifestMap["platform"].(map[string]interface{})
	if !ok {
		log.Error("[actionLog's GetActionPlatformInfo]:error when unmarshal action's platform info:", manifestMap, " platform setting is not a map[string]interface{}")
		return nil, errors.New("action " + actionLog.ActionLog.Action + "'s platform setting is illegal")
	}

	platformType, ok := platformSetting["platformType"].(string)
	if !ok {
		log.Error("[actionLog's GetActionPlatformInfo]:error when get action's platformType:", platformSetting, " platformType is not a string")
		return nil, errors.New("action " + actionLog.ActionLog.Action + "'s platform type is illegal")
	}

	platformHost, ok := platformSetting["platformHost"].(string)
	if !ok {
		log.Error("[actionLog's GetActionPlatformInfo]:error when get action's platformHost:", platformSetting, " platformHost is not a string")
		return nil, errors.New("action " + actionLog.ActionLog.Action + "'s platform host is illegal")
	}

	result := make(map[string]string)
	result["platformType"] = platformType
	result["platformHost"] = platformHost

	return result, nil
}

func getResultFromRelation(outputJson string, relationList []Relation, result map[string]interface{}) error {
	fromActionData := make(map[string]interface{})

	err := json.Unmarshal([]byte(outputJson), &fromActionData)
	if err != nil {
		log.Error("[actionLog's getResultFromRelation]:error when unmarshal action's output json:", outputJson, " ===>error is:", err.Error())
		return errors.New("error when parse from action data2:" + err.Error() + "\n" + outputJson)
	}

	for _, relation := range relationList {
		fromData, err := getJsonDataByPath(strings.TrimPrefix(relation.From, "."), fromActionData)
		if err != nil {
			return errors.New("error when get fromData :" + err.Error())
		}

		setDataToMapByPath(fromData, result, strings.TrimPrefix(relation.To, "."))
	}

	return nil
}
