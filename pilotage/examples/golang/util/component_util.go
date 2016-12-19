package util

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	CO_POD_NAME     = "CO_POD_NAME"
	CO_RUN_ID       = "CO_RUN_ID"
	CO_EVENT_LIST   = "CO_EVENT_LIST"
	CO_SERVICE_ADDR = "CO_SERVICE_ADDR"

	CO_COMPONENT_START = "CO_COMPONENT_START"
	CO_COMPONENT_STOP  = "CO_COMPONENT_STOP"

	CO_ACTION_TIMEOUT = "CO_ACTION_TIMEOUT"

	CO_TASK_START  = "CO_TASK_START"
	CO_TASK_RESULT = "CO_TASK_RESULT"
	CO_TASK_STATUS = "CO_TASK_STATUS"

	CO_REGISTER_URL = "CO_register"

	CO_DATA = "CO_DATA"

	CO_SET_GLOBAL_VAR_URL = "CO_SET_GLOBAL_VAR_URL"

	CO_LINKSTART_TOKEN = "CO_LINKSTART_TOKEN"
	CO_LINKSTART_URL   = "CO_LINKSTART_URL"
)

var (
	eventIDMap   map[string]int64
	eventURLMap  map[string]string
	eventINFOMap map[string]string

	isWaitData      bool
	receiveDataChan chan map[string]interface{}
)

func init() {
	isWaitData = false
	receiveDataChan = make(chan map[string]interface{}, 1)
	eventIDMap = make(map[string]int64)
	eventURLMap = make(map[string]string)
	eventINFOMap = make(map[string]string)

	// init component env info
	log.Println("[component util]", "===>start init...")

	eventINFOMap[CO_POD_NAME] = os.Getenv(CO_POD_NAME)
	eventINFOMap[CO_RUN_ID] = os.Getenv(CO_RUN_ID)
	eventINFOMap[CO_DATA] = os.Getenv(CO_DATA)
	eventINFOMap[CO_SERVICE_ADDR] = os.Getenv(CO_SERVICE_ADDR)

	eventINFOMap[CO_EVENT_LIST] = os.Getenv(CO_EVENT_LIST)
	eventINFOMap[CO_COMPONENT_START] = os.Getenv(CO_COMPONENT_START)
	eventINFOMap[CO_COMPONENT_STOP] = os.Getenv(CO_COMPONENT_STOP)
	eventINFOMap[CO_ACTION_TIMEOUT] = os.Getenv(CO_ACTION_TIMEOUT)
	eventINFOMap[CO_TASK_START] = os.Getenv(CO_TASK_START)
	eventINFOMap[CO_TASK_RESULT] = os.Getenv(CO_TASK_RESULT)
	eventINFOMap[CO_TASK_STATUS] = os.Getenv(CO_TASK_STATUS)
	eventINFOMap[CO_REGISTER_URL] = os.Getenv(CO_REGISTER_URL)

	for _, eventInfo := range strings.Split(eventINFOMap[CO_EVENT_LIST], ";") {
		if len(strings.Split(eventInfo, ",")) > 1 {
			eventKey := strings.Split(eventInfo, ",")[0]
			eventId := strings.Split(eventInfo, ",")[1]

			eventIdInt, _ := strconv.ParseInt(eventId, 10, 64)
			eventIDMap[eventKey] = eventIdInt
			eventURLMap[eventKey] = os.Getenv(eventKey)
		}
	}

	eventURLMap[CO_SET_GLOBAL_VAR_URL] = os.Getenv(CO_SET_GLOBAL_VAR_URL)
	eventINFOMap[CO_LINKSTART_TOKEN] = os.Getenv(CO_LINKSTART_TOKEN)
	eventURLMap[CO_LINKSTART_URL] = os.Getenv(CO_LINKSTART_URL)

	log.Println("[component util]", "<===init done")
	log.Println("[component util]", "<===got event map:", eventINFOMap)
}

func NotifyEvent(eventName string, status bool, result, output string) error {
	if eventURLMap[eventName] == "" || eventIDMap[eventName] == int64(0) {
		log.Println("[component util]", "===>error when notify event:", eventName, " because event info is illegal, got evnet id:", eventIDMap[eventName], " and event url:", eventURLMap[eventName])
		return errors.New("event is illegal")
	}

	reqBody := make(map[string]interface{})
	reqBody["EVENT"] = eventName
	reqBody["EVENT_ID"] = eventIDMap[eventName]
	reqBody["RUN_ID"] = eventINFOMap[CO_RUN_ID]
	reqBody["INFO"] = map[string]interface{}{"status": status, "result": result, "output": output}

	reqBodyBytes, _ := json.Marshal(reqBody)

	log.Println("[component util]", "===>component start notify event:", eventName, " to ", eventURLMap[eventName], "reqBody:", reqBody)
	resp, err := http.Post(eventURLMap[eventName], "application/json", bytes.NewReader(reqBodyBytes))

	if err != nil {
		log.Println("[component util]", "===>component send event:", eventName, " to:", eventURLMap[eventName], " \t error, error is:", err.Error())
		return errors.New("error when send req to workflow")
	}

	respBody, _ := ioutil.ReadAll(resp.Body)

	log.Println("[component util]", "===>component send event:", eventName, " got resp:\n", string(respBody), "\n")
	return nil
}

func ComponentStart(info string) error {
	return NotifyEvent(CO_COMPONENT_START, true, info, "")
}

func ComponentStop(info string) error {
	return NotifyEvent(CO_COMPONENT_STOP, true, info, "")
}

func TaskStart(info string) error {
	return NotifyEvent(CO_TASK_START, true, info, "")
}

func TaskResult(info string) error {
	return NotifyEvent(CO_TASK_RESULT, true, info, "")
}

func TaskStatus(status bool, info, output string) error {
	return NotifyEvent(CO_TASK_STATUS, status, info, output)
}

func GetData(port int64, forceRefresh bool, dataChan chan map[string]interface{}) error {
	dataMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(eventINFOMap[CO_DATA]), &dataMap)
	if err == nil {
		dataChan <- dataMap
		return nil
	}

	log.Println("[component util]", "===>error when get CO_DATA:", err.Error())

	dataChan <- dataMap
	return errors.New("error when get data")
}

func HoldProj() {
	timeout, err := strconv.Atoi(eventINFOMap[CO_ACTION_TIMEOUT])
	if err != nil {
		timeout = 0
	}
	if timeout <= 0 {
		c := make(chan bool, 1)
		<-c
	} else {
		time.Sleep(time.Duration(timeout) * time.Second)
	}
}

func ChangeGlobalVar(varName, value string) error {
	reqBody := make(map[string]interface{})

	reqBody["RUN_ID"] = eventINFOMap[CO_RUN_ID]
	reqBody["varMap"] = map[string]interface{}{"KEY": varName, "VALUE": value}

	reqBodyBytes, _ := json.Marshal(reqBody)

	log.Println("[component util]", "===>component start change global var info, \nvar name:", varName, "\nvalue:", value, "\nreqBody:", reqBody, "\nto:", eventURLMap[CO_SET_GLOBAL_VAR_URL])
	resp, err := http.Post(eventURLMap[CO_SET_GLOBAL_VAR_URL], "application/json", bytes.NewReader(reqBodyBytes))

	if err != nil {
		log.Println("[component util]", "===>component send event:", CO_SET_GLOBAL_VAR_URL, " to:", eventURLMap[CO_SET_GLOBAL_VAR_URL], " \t error, error is:", err.Error())
		return errors.New("error when send req to workflow")
	}

	respBody, _ := ioutil.ReadAll(resp.Body)

	log.Println("[component util]", "===>component send event:", CO_SET_GLOBAL_VAR_URL, " got resp:\n", string(respBody), "\n")
	return nil
}

func LinkStart(workflowName, workflowVersion, eventName, eventType string, startJson map[string]interface{}) error {
	reqBody := make(map[string]interface{})

	startJsonBytes, _ := json.Marshal(startJson)

	reqBody["RUN_ID"] = eventINFOMap[CO_RUN_ID]
	reqBody["linkInfoMap"] = map[string]interface{}{
		"token":           eventINFOMap[CO_LINKSTART_TOKEN],
		"workflowName":    workflowName,
		"workflowVersion": workflowVersion,
		"eventName":       eventName,
		"eventType":       eventType,
		"startJson":       string(startJsonBytes)}

	reqBodyBytes, _ := json.Marshal(reqBody)

	log.Println("[component util]", "===>component start link start, \ntoken:", eventINFOMap[CO_LINKSTART_TOKEN], "\nworkflow:", workflowName, ":", workflowVersion, "\nstartjson:", startJson, "\nbody:", string(reqBodyBytes), "\nto:", eventURLMap[CO_LINKSTART_URL]+workflowName)
	resp, err := http.Post(eventURLMap[CO_LINKSTART_URL]+workflowName, "application/json", bytes.NewReader(reqBodyBytes))
	if err != nil {
		log.Println("[component util]", "===>component send event:", CO_LINKSTART_URL, " to:", eventURLMap[CO_LINKSTART_URL]+workflowName, " \t error, error is:", err.Error())
		return errors.New("error when send req to workflow")
	}

	respBody, _ := ioutil.ReadAll(resp.Body)

	log.Println("[component util]", "===>component send event:", CO_LINKSTART_URL, " got resp:\n", string(respBody), "\n")
	return nil
}

func waitData(port int64) {
	if isWaitData {
		return
	}
	isWaitData = true

	http.HandleFunc("/receivedata", receiveDataHandler)
	http.ListenAndServe(":"+strconv.FormatInt(port, 10), nil)
}

func receiveDataHandler(w http.ResponseWriter, r *http.Request) {
	result, _ := json.Marshal(map[string]string{"message": "ok"})

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("error when get data body:" + err.Error())
	}

	codePathMap := make(map[string]interface{})
	json.Unmarshal([]byte(body), &codePathMap)
	receiveDataChan <- codePathMap

	w.Write(result)
}
