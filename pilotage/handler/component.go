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

package handler

import (
	"encoding/json"
	"github.com/Huawei/containerops/pilotage/models"
	"github.com/Huawei/containerops/pilotage/module"
	log "github.com/Sirupsen/logrus"
	"github.com/go-macaron/sockets"
	"github.com/golang/groupcache/lru"
	"github.com/gorilla/websocket"
	"gopkg.in/macaron.v1"
	"net/http"
	"strconv"
	"time"
)

var cache *lru.Cache

func init() {
	//TODO: may want to configure max entries number
	cache = lru.New(50)
	cache.OnEvicted = func(key lru.Key, value interface{}) {
		log.Warnf("Component message channel key %v evicted\n", key)
		channel, ok := value.(chan DebugEvent)
		if !ok {
			log.Warnf("Can't convert cache value %T to message channel\n", value)
			return
		}
		close(channel)
	}
}

func ListComponents(ctx *macaron.Context) (httpStatus int, result []byte) {
	var resp ListComponentsResp
	var fuzzy bool
	name := ctx.QueryTrim("name")
	version := ctx.QueryTrim("version")
	f := ctx.QueryTrim("fuzzy")
	if f == "" {
		fuzzy = false
	} else {
		var err error
		fuzzy, err = strconv.ParseBool(f)
		if err != nil {
			httpStatus = http.StatusBadRequest
			resp.OK = false
			resp.ErrorCode = ComponentError + ComponentReqBodyError
			resp.Message = "Parse query param fuzzy error: " + err.Error()

			result, err = json.Marshal(resp)
			if err != nil {
				log.Errorln("List components marshal data error: " + err.Error())
			}
			return
		}
	}
	pageNum := ctx.QueryInt("page_num")
	if pageNum <= 0 {
		pageNum = 5
	}
	versionNum := ctx.QueryInt("version_num")
	if versionNum <= 0 {
		versionNum = 5
	}
	offset := ctx.QueryInt("offset")
	if offset < 0 {
		offset = 0
	}

	components, err := module.GetComponents(name, version, fuzzy, pageNum, versionNum, offset)
	if err != nil {
		httpStatus = http.StatusBadRequest
		resp.OK = false
		resp.ErrorCode = ComponentError + ComponentListError
		resp.Message = "List components error: " + err.Error()

		result, err = json.Marshal(resp)
		if err != nil {
			log.Errorln("List components marshal data error: " + err.Error())
		}
		return
	}

	if len(components) == 0 {
		httpStatus = http.StatusNotFound
		resp.OK = false
		resp.ErrorCode = ComponentError + ComponentListError
		resp.Message = "components not found"

		result, err = json.Marshal(resp)
		if err != nil {
			log.Errorln("List components marshal data error: " + err.Error())
		}
		return
	}

	for _, component := range components {
		resp.Components = append(resp.Components, ComponentItem{
			ID: component.ID,
			Name: component.Name,
			Version: component.Version,
		})
	}

	httpStatus = http.StatusCreated
	resp.OK = true

	result, err = json.Marshal(resp)
	if err != nil {
		log.Errorln("List components marshal data error: " + err.Error())
	}
	return
	//result, _ := json.Marshal(map[string]string{"message": ""})
	//
	//namespace := ctx.Params(":namespace")
	//
	//if namespace == "" {
	//	result, _ = json.Marshal(map[string]string{"errMsg": "namespace can't be empty"})
	//	return http.StatusBadRequest, result
	//}
	//
	//componentList, err := module.GetComponentListByNamespace(namespace)
	//
	//if err != nil {
	//	result, _ = json.Marshal(map[string]string{"errMsg": "error when get component list:" + err.Error()})
	//	return http.StatusBadRequest, result
	//}
	//
	//result, _ = json.Marshal(map[string]interface{}{"list": componentList})
	//
	//return http.StatusOK, result
}

func CreateComponent(ctx *macaron.Context) (httpStatus int, result []byte) {
	var resp ComponentResp
	body, err := ctx.Req.Body().Bytes()
	if err != nil {
		httpStatus = http.StatusBadRequest
		resp.OK = false
		resp.ErrorCode = ComponentError + ComponentReqBodyError
		resp.Message = "Get requrest body error: " + err.Error()

		result, err = json.Marshal(resp)
		if err != nil {
			log.Errorln("Create component marshal data error: " + err.Error())
		}
		return
	}

	var req ComponentReq
	err = json.Unmarshal(body, &req)
	if err != nil {
		log.Errorln("CreateComponent unmarshal data error: ", err.Error())
		httpStatus = http.StatusMethodNotAllowed
		resp.OK = false
		resp.ErrorCode = ComponentError + ComponentUnmarshalError
		resp.Message = "unmarshal data error: " + err.Error()

		result, err = json.Marshal(resp)
		if err != nil {
			log.Errorln("Create component marshal data error: " + err.Error())
		}
		return
	}

	var component models.Component
	component.Name = req.Name
	component.Version = req.Version
	component.Type = 0
	for index, value := range models.ComponentTypes {
		if string(value) == req.Type {
			component.Type = index
			break
		}
	}
	component.ImageName = req.ImageName
	component.ImageTag = req.ImageTag
	component.Timeout = req.Timeout
	component.DataFrom = req.DataFrom
	component.UseAdvanced = req.UseAdvanced
	m := make(map[string]*json.RawMessage)
	m["pod"] = req.Pod
	m["service"] = req.Service
	data, err := json.Marshal(m)
	if err != nil {
		log.Errorln("Create component marshal KubeSetting data error: " + err.Error())
	}
	component.KubeSetting = string(data)
	data, err = json.Marshal(req.Input)
	if err != nil {
		log.Errorln("Create component marshal Input data error: " + err.Error())
	}
	component.Input = string(data)
	data, err = json.Marshal(req.Output)
	if err != nil {
		log.Errorln("Create component marshal Output data error: " + err.Error())
	}
	component.Output = string(data)
	data, err = json.Marshal(req.Env)
	if err != nil {
		log.Errorln("Create component marshal Env data error: " + err.Error())
	}
	component.Environment = string(data)

	if id, err := module.CreateComponent(&component); err != nil {
		httpStatus = http.StatusBadRequest
		resp.OK = false
		resp.ErrorCode = ComponentError + ComponentCreateError
		resp.Message = "Create component error: " + err.Error()

		result, err = json.Marshal(resp)
		if err != nil {
			log.Errorln("Create component marshal data error: " + err.Error())
		}
		return
	} else {
		httpStatus = http.StatusCreated
		resp.ID = id
		resp.OK = true
		resp.Message = "Component Created"
	}

	result, err = json.Marshal(resp)
	if err != nil {
		log.Errorln("Create component marshal data error: " + err.Error())
	}
	return
}

func GetComponent(ctx *macaron.Context) (httpStatus int, result []byte) {
	var resp ComponentResp
	componentID := ctx.Params(":component")
	id, err := strconv.ParseInt(componentID, 10, 64)
	if err != nil {
		httpStatus = http.StatusBadRequest
		resp.OK = false
		resp.ErrorCode = ComponentError + ComponentParseIDError
		resp.Message = "Parse component id error: " + err.Error()

		result, err = json.Marshal(resp)
		if err != nil {
			log.Errorln("Get component marshal data error: " + err.Error())
		}
		return
	}
	component, err := module.GetComponentByID(id)
	if err != nil {
		httpStatus = http.StatusBadRequest
		resp.OK = false
		resp.ErrorCode = ComponentError + ComponentGetError
		resp.Message = "get component by id error: " + err.Error()

		result, err = json.Marshal(resp)
		if err != nil {
			log.Errorln("Get component marshal data error: " + err.Error())
		}
		return
	}
	if component == nil {
		httpStatus = http.StatusNotFound
		resp.OK = false
		resp.ErrorCode = ComponentError + ComponentGetError
		resp.Message = "component not found"

		result, err = json.Marshal(resp)
		if err != nil {
			log.Errorln("Get component marshal data error: " + err.Error())
		}
		return
	}

	httpStatus = http.StatusOK
	resp.OK = true

	resp.ID = component.ID
	resp.Version = component.Version
	resp.ImageName = component.ImageName
	resp.ImageTag = component.ImageTag
	resp.Timeout = component.Timeout
	resp.Type = string(models.ComponentTypes[component.Type])
	resp.DataFrom = component.DataFrom
	resp.UseAdvanced = component.UseAdvanced
	if err := json.Unmarshal([]byte(component.Environment), &resp.Env); err != nil {
		log.Errorln("Get component unmarshal Environment data error: " + err.Error())
	}
	if err := json.Unmarshal([]byte(component.Input), &resp.Input); err != nil {
		log.Errorln("Get component unmarshal Input data error: " + err.Error())
	}
	if err := json.Unmarshal([]byte(component.Output), &resp.Output); err != nil {
		log.Errorln("Get component unmarshal Input data error: " + err.Error())
	}

	var setting map[string]*json.RawMessage
	if json.Unmarshal([]byte(component.KubeSetting), &setting); err != nil {
		log.Errorln("Get component unmarshal KubeSetting data error: " + err.Error())
	}
	resp.Pod = setting["pod"]
	resp.Service = setting["service"]

	result, err = json.Marshal(resp)
	if err != nil {
		log.Errorln("Get component marshal data error: " + err.Error())
	}
	return
}

func UpdateComponent(ctx *macaron.Context) (httpStatus int, result []byte) {
	var resp ComponentResp
	body, err := ctx.Req.Body().Bytes()
	if err != nil {
		httpStatus = http.StatusBadRequest
		resp.OK = false
		resp.ErrorCode = ComponentError + ComponentReqBodyError
		resp.Message = "Get requrest body error: " + err.Error()

		result, err = json.Marshal(resp)
		if err != nil {
			log.Errorln("Update component marshal data error: " + err.Error())
		}
		return
	}

	var req ComponentReq
	err = json.Unmarshal(body, &req)
	if err != nil {
		log.Errorln("UpdateComponent unmarshal data error: ", err.Error())
		httpStatus = http.StatusBadRequest
		resp.OK = false
		resp.ErrorCode = ComponentError + ComponentUnmarshalError
		resp.Message = "unmarshal data error: " + err.Error()

		result, err = json.Marshal(resp)
		if err != nil {
			log.Errorln("Update component marshal data error: " + err.Error())
		}
		return
	}

	componentID := ctx.Params(":component")
	id, err := strconv.ParseInt(componentID, 10, 64)
	if err != nil {
		httpStatus = http.StatusBadRequest
		resp.OK = false
		resp.ErrorCode = ComponentError + ComponentParseIDError
		resp.Message = "Parse component id error: " + err.Error()

		result, err = json.Marshal(resp)
		if err != nil {
			log.Errorln("Update component marshal data error: " + err.Error())
		}
		return
	}

	var component models.Component
	component.ID = id
	component.Name = req.Name
	component.Version = req.Version
	component.Type = 0
	for index, value := range models.ComponentTypes {
		if string(value) == req.Type {
			component.Type = index
			break
		}
	}
	component.ImageName = req.ImageName
	component.ImageTag = req.ImageTag
	component.Timeout = req.Timeout
	component.DataFrom = req.DataFrom
	component.UseAdvanced = req.UseAdvanced
	m := make(map[string]*json.RawMessage)
	m["pod"] = req.Pod
	m["service"] = req.Service
	data, err := json.Marshal(m)
	if err != nil {
		log.Errorln("Create component marshal KubeSetting data error: " + err.Error())
	}
	component.KubeSetting = string(data)
	data, err = json.Marshal(req.Input)
	if err != nil {
		log.Errorln("Create component marshal Input data error: " + err.Error())
	}
	component.Input = string(data)
	data, err = json.Marshal(req.Output)
	if err != nil {
		log.Errorln("Create component marshal Output data error: " + err.Error())
	}
	component.Output = string(data)
	data, err = json.Marshal(req.Env)
	if err != nil {
		log.Errorln("Create component marshal Env data error: " + err.Error())
	}
	component.Environment = string(data)

	if err := module.UpdateComponent(id, &component); err != nil {
		httpStatus = http.StatusBadRequest
		resp.OK = false
		resp.ErrorCode = ComponentError + ComponentUpdateError
		resp.Message = "update component error: " + err.Error()

		result, err = json.Marshal(resp)
		if err != nil {
			log.Errorln("Update component marshal data error: " + err.Error())
		}
		return
	}
	httpStatus = http.StatusOK
	resp.OK = true

	result, err = json.Marshal(resp)
	if err != nil {
		log.Errorln("Update component marshal data error: " + err.Error())
	}
	return
}

func DeleteComponent(ctx *macaron.Context) (httpStatus int, result []byte) {
	var resp ComponentResp

	componentID := ctx.Params(":component")
	id, err := strconv.ParseInt(componentID, 10, 64)
	if err != nil {
		httpStatus = http.StatusBadRequest
		resp.OK = false
		resp.ErrorCode = ComponentError + ComponentParseIDError
		resp.Message = "Parse component id error: " + err.Error()

		result, err = json.Marshal(resp)
		if err != nil {
			log.Errorln("Delete component marshal data error: " + err.Error())
		}
		return
	}

	if err := module.DeleteComponent(id); err != nil {
		httpStatus = http.StatusBadRequest
		resp.OK = false
		resp.ErrorCode = ComponentError + ComponentDeleteError
		resp.Message = "delete component error: " + err.Error()

		result, err = json.Marshal(resp)
		if err != nil {
			log.Errorln("Delete component marshal data error: " + err.Error())
		}
		return
	}

	httpStatus = http.StatusOK
	resp.OK = true

	result, err = json.Marshal(resp)
	if err != nil {
		log.Errorln("Delete component marshal data error: " + err.Error())
	}
	return
}

//func DebugComponent(ctx *macaron.Context) (httpStatus int, result []byte) {
//	var resp DebugComponentResp
//	body, err := ctx.Req.Body().Bytes()
//	if err != nil {
//		httpStatus = http.StatusBadRequest
//		resp.OK = false
//		resp.ErrorCode = componentErrCode + 5
//		resp.Message = "Get requrest body error: " + err.Error()
//
//		result, err = json.Marshal(resp)
//		if err != nil {
//			log.Errorln("Debug component marshal data error: " + err.Error())
//		}
//		return
//	}
//
//	componentID := ctx.Params(":component_id")
//	id, err := strconv.ParseInt(componentID, 10, 64)
//	if err != nil {
//		httpStatus = http.StatusBadRequest
//		resp.OK = false
//		resp.ErrorCode = componentErrCode + 10
//		resp.Message = "Parse component id error: " + err.Error()
//
//		result, err = json.Marshal(resp)
//		if err != nil {
//			log.Errorln("Debug component marshal data error: " + err.Error())
//		}
//		return
//	}
//
//	var req *DebugComponentReq
//	err = json.Unmarshal(body, req)
//	if err != nil {
//		log.Errorln("DebugComponent unmarshal data error: ", err.Error())
//		httpStatus = http.StatusBadRequest
//		resp.OK = false
//		resp.ErrorCode = componentErrCode + 5
//		resp.Message = "unmarshal data error: " + err.Error()
//
//		result, err = json.Marshal(resp)
//		if err != nil {
//			log.Errorln("Debug component marshal data error: " + err.Error())
//		}
//		return
//	}
//
//	if req.Kubernetes == "" {
//		httpStatus = http.StatusBadRequest
//		resp.OK = false
//		resp.ErrorCode = componentErrCode + 5
//		resp.Message = "should specify kubernetes api server"
//
//		result, err = json.Marshal(resp)
//		if err != nil {
//			log.Errorln("Debug component marshal data error: " + err.Error())
//		}
//		return
//	}
//
//	component, err := module.GetComponentByID(id)
//	if err != nil {
//		httpStatus = http.StatusBadRequest
//		resp.OK = false
//		resp.ErrorCode = componentErrCode + 4
//		resp.Message = "get component by id error: " + err.Error()
//
//		result, err = json.Marshal(resp)
//		if err != nil {
//			log.Errorln("Debug component marshal data error: " + err.Error())
//		}
//		return
//	}
//
//	logID, err := module.DebugComponent(component, req.Kubernetes, req.Input, req.Environment)
//	if err != nil {
//		httpStatus = http.StatusBadRequest
//		resp.OK = false
//		resp.ErrorCode = componentErrCode + 9
//		resp.Message = "debug component error: " + err.Error()
//
//		result, err = json.Marshal(resp)
//		if err != nil {
//			log.Errorln("Debug component marshal data error: " + err.Error())
//		}
//		return
//	}
//
//	httpStatus = http.StatusOK
//	resp.OK = true
//	resp.LogID = logID
//	result, err = json.Marshal(resp)
//	if err != nil {
//		log.Errorln("Debug component marshal data error: " + err.Error())
//	}
//	return
//}

func DebugComponentJson() macaron.Handler {
	//TODO: add socket options
	return sockets.JSON(DebugComponentMessage{})
}

func DebugComponentLog(ctx *macaron.Context,
	receiver <-chan *DebugComponentMessage,
	sender chan<- *DebugComponentMessage,
	done <-chan bool,
	disconnect chan<- int,
	errChan <-chan error) {
	id, err := strconv.ParseInt(ctx.Params(":component"), 10, 64)
	if err != nil {
		sender <- &DebugComponentMessage{
			CommonResp: CommonResp{
				OK:        false,
				ErrorCode: ComponentError + ComponentParseIDError,
				Message:   "Parse component id error: " + err.Error(),
			},
		}
		disconnect <- websocket.CloseUnsupportedData
		return
	}
	component, err := module.GetComponentByID(id)
	if err != nil {
		sender <- &DebugComponentMessage{
			CommonResp: CommonResp{
				OK:        false,
				ErrorCode: ComponentError + ComponentGetError,
				Message:   "get component by id error: " + err.Error(),
			},
		}
		disconnect <- websocket.CloseUnsupportedData
		return
	}

	var actionLog *module.ActionLog
	eventChan := make(chan DebugEvent)
	ticker := time.Tick(time.Duration(component.Timeout + 30) * time.Second)
	for {
		select {
		case event, ok := <-eventChan:
			if !ok {
				break
			}
			sender <- &DebugComponentMessage{
				DebugID: actionLog.ID,
				Event:   event,
				CommonResp: CommonResp{
					OK: true,
				},
			}
			if event.Type == module.COMPONENT_STOP {
				if output, err := actionLog.GetOutcome(); err != nil {
					log.Errorf("DebugComponent get output data error: %s\n", err)
				} else {
					var outputMap map[string]interface{}
					if err := json.Unmarshal([]byte(output), &outputMap); err == nil {
						sender <- &DebugComponentMessage{
							DebugID: actionLog.ID,
							Output:  outputMap,
							CommonResp: CommonResp{
								OK: true,
							},
						}
					}
				}
				disconnect <- websocket.CloseNormalClosure
				return
			}
		case msg := <-receiver:
			if msg.DebugID > 0 {
				cache.Remove(msg.DebugID)
				eventChan = make(chan DebugEvent)
			}
			if msg.Kubernetes == "" {
				sender <- &DebugComponentMessage{
					CommonResp: CommonResp{
						OK:        false,
						ErrorCode: ComponentError + ComponentDebugError,
						Message:   "should specify kubernetes api server",
					},
				}
				disconnect <- websocket.CloseUnsupportedData
				return
			}
			envMap := make(map[string]string)
			for _, item := range msg.Env {
				envMap[item.Key] = item.Value
			}
			env, err := json.Marshal(envMap)
			if err != nil {
				sender <- &DebugComponentMessage{
					CommonResp: CommonResp{
						OK:        false,
						ErrorCode: ComponentError + ComponentmarshalError,
						Message:   "debug component error: " + err.Error(),
					},
				}
				disconnect <- websocket.CloseInternalServerErr
				return
			}
			actionLog, err = module.DebugComponent(component, msg.Kubernetes, msg.Input, string(env))
			if err != nil {
				sender <- &DebugComponentMessage{
					CommonResp: CommonResp{
						OK:        false,
						ErrorCode: ComponentError + ComponentDebugError,
						Message:   "debug component error: " + err.Error(),
					},
				}
				disconnect <- websocket.CloseInternalServerErr
				return
			}
			cache.Add(actionLog.ID, eventChan)
			sender <- &DebugComponentMessage{
				DebugID: actionLog.ID,
				Input:   msg.Input,
				CommonResp: CommonResp{
					OK: true,
				},
			}

		case <-done:
			log.Debug("DebugComponent socket closed by client")
			return
		case <-ticker:
			log.Debug("DebugComponent socket closed by server")
			disconnect <- websocket.CloseNormalClosure
		case err := <-errChan:
			log.Errorf("Debug Component socket error: %s\n", err)
		}
	}
}
