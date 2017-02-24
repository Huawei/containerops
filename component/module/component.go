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
	"encoding/json"
	"fmt"

	"github.com/Huawei/containerops/component/models"

	"github.com/jinzhu/gorm"
)

// GetComponents is return component list with given conditions
// returns :
// map[string]interface{}{
// 	"id":XX,
// 	"name":"name",
// 	"version":"version"
// }
func GetComponents(namespace, name string, fuzzy bool, pageNum, versionNum, offset int) ([]ComponentBaseData, error) {
	result := make([]ComponentBaseData, 0)

	componentNames, err := getComponentNames(namespace, name, fuzzy, pageNum, offset)
	if err != nil {
		return nil, fmt.Errorf("get component error:%s", err.Error())
	}

	for _, cName := range componentNames {
		versionInfos, err := getComponentVersion(namespace, cName, fuzzy, versionNum, offset)
		if err != nil {
			continue
		}

		for _, vInfo := range versionInfos {
			var tempCInfo ComponentBaseData
			tempCInfo.ID = vInfo.ID
			tempCInfo.Name = vInfo.Name
			tempCInfo.Version = vInfo.Version

			result = append(result, tempCInfo)
		}
	}

	return result, nil
}

// getComponentNames is only return componet's name slice with given conditions
func getComponentNames(namespace, name string, fuzzy bool, pageNum, offset int) ([]string, error) {
	result := make([]string, 0)
	names := make([]struct {
		Name string
	}, 0)

	// if not fuzzy, means require a specific component
	if !fuzzy {
		db := models.GetDB().Table("component").Select("DISTINCT(name)").Where("namespace = ?", namespace)
		count := int64(0)
		err := db.Where("name = ?", name).Count(&count).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			// error when query from db
			return result, err
		} else if err != nil && err == gorm.ErrRecordNotFound {
			// don't have result, return empty result
			return result, nil
		} else {
			// only return this name
			result = append(result, name)
			return result, nil
		}
	}

	queryStr := " AND namespace = ? "
	queryArray := make([]interface{}, 0)
	queryArray = append(queryArray, namespace)

	if name != "" {
		queryStr += " AND name like ? "
		queryArray = append(queryArray, "%"+name+"%")
	}

	queryLimitStr := "LIMIT ? OFFSET ?"
	queryArray = append(queryArray, pageNum)
	queryArray = append(queryArray, offset)

	err := models.GetDB().Raw("SELECT DISTINCT(A.name) as name FROM ( SELECT B.* FROM component B WHERE 1=1 "+queryStr+" ORDER BY updated_at desc "+queryLimitStr+") A", queryArray...).Scan(&names).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return result, err
	}

	for _, n := range names {
		result = append(result, n.Name)
	}

	return result, nil
}

// getComponentVersion is only return componet's name slice with given conditions
func getComponentVersion(namespace, name string, fuzzy bool, versionNum, offset int) ([]models.Component, error) {
	result := make([]models.Component, 0)

	db := new(models.Component).GetComponent().Where("namespace = ?", namespace).Where("name = ?", name)

	// if not fuzzy, means require a specific component's vresion,so offset is should be consider
	if !fuzzy {
		db = db.Offset(offset)
	}

	err := db.Order("updated_at desc").Limit(versionNum).Find(&result).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return result, nil
}

// CreateComponent is create a component with given msg, this will only check is component's name and version is empty,will not check other info
func CreateComponent(namespace string, componentInfo []byte) (int64, error) {
	kubeSetting := ""
	kubeSettingMap := make(map[string]interface{})
	inputStr := ""
	outputStr := ""
	envStr := ""
	shellStr := ""

	cInfo := new(ComponentData)
	err := json.Unmarshal(componentInfo, &cInfo)
	if err != nil {
		return 0, fmt.Errorf("CreateComponent error: component data error:%s", err.Error())
	}

	if cInfo.Name == "" {
		return 0, fmt.Errorf("CreateComponent error: should specify component name")
	}

	if cInfo.Version == "" {
		return 0, fmt.Errorf("CreateComponent error: should specify component version")
	}

	kubeSettingMap["podConfig"] = cInfo.Pod
	kubeSettingMap["serviceConfig"] = cInfo.Service
	setting, _ := json.Marshal(kubeSettingMap)
	kubeSetting = string(setting)

	if cInfo.Input != nil {
		input, err := json.Marshal(cInfo.Input)
		if err != nil {
			return 0, fmt.Errorf("CreateComponent error: format input error:%s", err.Error())
		}
		inputStr = string(input)
	}

	if cInfo.Output != nil {
		output, err := json.Marshal(cInfo.Output)
		if err != nil {
			return 0, fmt.Errorf("CreateComponent error: format output error:%s", err.Error())
		}
		outputStr = string(output)
	}

	if len(cInfo.Env) > 0 {
		env, err := json.Marshal(cInfo.Env)
		if err != nil {
			return 0, fmt.Errorf("CreateComponent error: format env error:%s", err.Error())
		}
		envStr = string(env)
	}

	if cInfo.ImageSetting != nil {
		if shell, err := json.Marshal(cInfo.ImageSetting.EventShell); err != nil {
			return 0, fmt.Errorf("CreateComponent error: format eventShell error: %s", err.Error())
		} else {
			shellStr = string(shell)
		}
	}

	component := new(models.Component)
	component.Namespace = namespace
	component.Name = cInfo.Name
	component.Version = cInfo.Version
	component.Containertype = models.ComponentContainerTypeDocker
	component.Runengine = models.ComponentEngineMap[cInfo.Type]
	component.ImageName = cInfo.ImageName
	component.ImageTag = cInfo.ImageTag
	component.Timeout = cInfo.Timeout
	component.UseAdvanced = cInfo.UseAdvanced

	component.KubeSetting = kubeSetting
	component.Input = inputStr
	component.Output = outputStr
	component.Environment = envStr

	if cInfo.ImageSetting != nil {
		component.BaseImageName = cInfo.ImageSetting.From.ImageName
		component.BaseImageTag = cInfo.ImageSetting.From.ImageTag
		component.EventShell = shellStr
	}

	count := int64(0)
	err = component.GetComponent().Where("namespace = ?", namespace).Where("name = ?", cInfo.Name).Where("version = ?", cInfo.Version).Count(&count).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return 0, fmt.Errorf("error when query component info:", err.Error())
	} else if count > 0 {
		return 0, fmt.Errorf("CreateComponent error: component exist")
	}

	err = component.GetComponent().Save(component).Error
	if err != nil {
		log.Errorln("CreateComponent error:when save component info:", err.Error())
		return 0, fmt.Errorf("CreateCompoent error: when save component:%s", err.Error())
	}

	return component.ID, nil
}

// GetComponentByID is get component detail info by given conditions
func GetComponentByID(namespace string, id int64) (*ComponentData, error) {
	result := new(ComponentData)
	result.ImageSetting = new(BuildInfo)
	result.ImageSetting.From = new(ImageInfo)
	result.ImageSetting.Push = new(ImageInfo)

	if id <= 0 {
		return nil, fmt.Errorf("invalid component id")
	}

	component := new(models.Component)
	err := component.GetComponent().Where("namespace = ?", namespace).Where("id = ?", id).First(&component).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("error when get component info: %s", err.Error())
	} else if err != nil && err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	cInfo, err := json.Marshal(component)
	if err != nil {
		return nil, fmt.Errorf("error when marshal component info: %s", err.Error())
	}

	if err = json.Unmarshal(cInfo, &result); err != nil {
		return nil, fmt.Errorf("error when unmarshal component info: %s", err.Error())
	}

	switch component.Runengine {
	case models.ComponentEngineTypeKube:
		result.Type = "Kubernetes"
	case models.ComponentEngineTypeSwarm:
		result.Type = "Swarm"
	}

	if component.Input != "" {
		inputMap := make(map[string]interface{})
		err := json.Unmarshal([]byte(component.Input), &inputMap)
		if err != nil {
			return nil, fmt.Errorf("error when unmarshal component input: %s", err.Error())
		}

		result.Input = inputMap
	}

	if component.Output != "" {
		outputMap := make(map[string]interface{})
		err := json.Unmarshal([]byte(component.Output), &outputMap)
		if err != nil {
			return nil, fmt.Errorf("error when unmarshal component output: %s", err.Error())
		}

		result.Output = outputMap
	}

	if component.Environment != "" {
		err = json.Unmarshal([]byte(component.Environment), &result.Env)
		if err != nil {
			return nil, fmt.Errorf("error when unmarshal component env2 info: %s", err.Error())
		}
	}

	if component.KubeSetting != "" {
		kubeSettingMap := make(map[string]interface{})
		if err = json.Unmarshal([]byte(component.KubeSetting), &kubeSettingMap); err != nil {
			return nil, fmt.Errorf("error when unmarshal component kubesetting: %s", err.Error())
		}

		if podSetting, ok := kubeSettingMap["podConfig"]; ok {
			setting, err := json.Marshal(podSetting)
			if err != nil {
				return nil, fmt.Errorf("error when unmarshal component kubesetting[podConfig]: %s", err.Error())
			}

			json.Unmarshal(setting, result.Pod)
		}

		if serviceSetting, ok := kubeSettingMap["serviceConfig"]; ok {
			setting, err := json.Marshal(serviceSetting)
			if err != nil {
				return nil, fmt.Errorf("error when unmarshal component kubesetting[serviceConfig]: %s", err.Error())
			}

			json.Unmarshal(setting, result.Service)
		}
	}

	if component.EventShell != "" {
		var eventMap map[string]interface{}
		err := json.Unmarshal([]byte(component.EventShell), &eventMap)
		if err != nil {
			return nil, fmt.Errorf("error when unmarshal compoent event shell: %s", err.Error())
		}

		result.ImageSetting.EventShell.ComponentStart = eventMap["componentStart"]
		result.ImageSetting.EventShell.ComponentResult = eventMap["componentResult"]
		result.ImageSetting.EventShell.ComponentStop = eventMap["componentStop"]
	}

	if component.BaseImageName != "" {
		result.ImageSetting.From.ImageName = component.BaseImageName
	}

	if component.BaseImageTag != "" {
		result.ImageSetting.From.ImageTag = component.BaseImageTag
	}

	result.ImageSetting.Push.ImageName = component.ImageName
	result.ImageSetting.Push.ImageTag = component.ImageTag

	return result, nil
}

// UpdateComponent is update component by give conditions
func UpdateComponent(namespace string, id int64, componentInfo []byte) error {
	if id <= 0 {
		return fmt.Errorf("invalid component id")
	}

	kubeSetting := ""
	kubeSettingMap := make(map[string]interface{})
	inputStr := ""
	outputStr := ""
	envStr := ""

	cInfo := new(ComponentData)
	err := json.Unmarshal(componentInfo, &cInfo)
	if err != nil {
		return fmt.Errorf("UpdateComponent error: component data error:%s", err.Error())
	}

	if cInfo.Name == "" {
		return fmt.Errorf("UpdateComponent error: should specify component name")
	}

	if cInfo.Version == "" {
		return fmt.Errorf("UpdateComponent error: should specify component version")
	}

	kubeSettingMap["podConfig"] = cInfo.Pod
	kubeSettingMap["serviceConfig"] = cInfo.Service
	setting, _ := json.Marshal(kubeSettingMap)
	kubeSetting = string(setting)

	if cInfo.Input != nil {
		input, err := json.Marshal(cInfo.Input)
		if err != nil {
			return fmt.Errorf("UpdateComponent error: format input error:%s", err.Error())
		}
		inputStr = string(input)
	}

	if cInfo.Output != nil {
		output, err := json.Marshal(cInfo.Output)
		if err != nil {
			return fmt.Errorf("UpdateComponent error: format output error:%s", err.Error())
		}
		outputStr = string(output)
	}

	if len(cInfo.Env) > 0 {
		env, err := json.Marshal(cInfo.Env)
		if err != nil {
			return fmt.Errorf("UpdateComponent error: format env error:%s", err.Error())
		}
		envStr = string(env)
	}

	component := new(models.Component)
	err = component.GetComponent().Where("namespace = ?", namespace).Where("id = ?", id).First(&component).Error
	if err != nil {
		return fmt.Errorf("error when get component info: %s", err.Error())
	}

	component.Runengine = models.ComponentEngineMap[cInfo.Type]
	component.ImageName = cInfo.ImageName
	component.ImageTag = cInfo.ImageTag
	component.Timeout = cInfo.Timeout
	component.UseAdvanced = cInfo.UseAdvanced
	component.KubeSetting = kubeSetting
	component.Input = inputStr
	component.Output = outputStr
	component.Environment = envStr

	err = component.GetComponent().Save(component).Error
	if err != nil {
		log.Errorln("UpdateComponent error:when save component info:", err.Error())
		return fmt.Errorf("CreateCompoent error: when save component:%s", err.Error())
	}

	return nil
}

// DeleteComponentByID is delete component by give conditions
func DeleteComponentByID(namespace string, id int64) error {
	if id <= 0 {
		return fmt.Errorf("invaild component id")
	}

	err := new(models.Component).GetComponent().Where("namespace = ?", namespace).Where("id = ?", id).Delete(&models.Component{}).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return fmt.Errorf("error when delete component: %s", err.Error())
	}

	return nil
}
