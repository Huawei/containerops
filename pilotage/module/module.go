/*
Copyright 2016 - 2017 Huawei Technologies Co., Ltd. All rights reserved.

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
	"errors"
	"strings"

	log "github.com/Sirupsen/logrus"
)

func describeJSON(jsonObj map[string]interface{}, path string) ([]map[string]string, error) {
	resultList := make([]map[string]string, 0)

	for key, value := range jsonObj {
		temp := make(map[string]string)

		typeStr := ""
		switch value.(type) {
		case string:
			typeStr = "string"
		case float64:
			typeStr = "float64"
		case bool:
			typeStr = "boolean"
		case []interface{}:
			typeStr = "array"
		case map[string]interface{}:
			typeStr = "object"
			childResult, err := describeJSON(value.(map[string]interface{}), temp["path"])
			if err != nil {
				return nil, err
			}

			resultList = append(resultList, childResult...)
		}

		temp["key"] = key
		temp["path"] = path + "." + key
		temp["type"] = typeStr

		resultList = append(resultList, temp)
	}

	return resultList, nil
}

// getJsonDataByPath is get a value from a map by give path
func getJsonDataByPath(path string, data map[string]interface{}) (interface{}, error) {
	depth := len(strings.Split(path, "."))
	if depth == 1 {
		if info, ok := data[path]; !ok {
			log.Error("[module's getJsonDataByPath]:error when get data from action,action's key not exist :" + path)
			return "", nil
		} else {
			return info, nil
		}
	}

	childDataInterface, ok := data[strings.Split(path, ".")[0]]
	if !ok {
		log.Error("error when get data from action,action's key not exist :" + path)
		return "", nil
	}
	childData, ok := childDataInterface.(map[string]interface{})
	if !ok {
		log.Error("[module's getJsonDataByPath]:error when get data from output, want a json,got:", childDataInterface)
		return nil, errors.New("child data is not a json!")
	}

	childPath := strings.Join(strings.Split(path, ".")[1:], ".")
	return getJsonDataByPath(childPath, childData)
}

// setDataToMapByPath is set a data to a map by give path ,if parent path not exist,it will auto creat
func setDataToMapByPath(data interface{}, result map[string]interface{}, path string) {
	depth := len(strings.Split(path, "."))
	if depth == 1 {
		result[path] = data
		return
	}

	currentPath := strings.Split(path, ".")[0]
	currentMap := make(map[string]interface{})
	if _, ok := result[currentPath]; !ok {
		result[currentPath] = currentMap
	}

	var ok bool
	currentMap, ok = result[currentPath].(map[string]interface{})
	if !ok {
		return
	}

	childPath := strings.Join(strings.Split(path, ".")[1:], ".")
	setDataToMapByPath(data, currentMap, childPath)
	return
}
