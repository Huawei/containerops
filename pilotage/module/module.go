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
