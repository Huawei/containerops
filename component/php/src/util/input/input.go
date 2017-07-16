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

package input

import (
	"strings"
	"fmt"
)

func HandleInput(data string, keys []string, result map[string]string) error {
	fields:= strings.Fields(data)
	if len(fields) == 0 {
		return fmt.Errorf("CO_DATA value is null")
	}

	kvs := map[string]string{}

	for _, v := range fields {
		s := strings.Split(v, "=")
		key, value := s[0], s[1]

		kvs[key] = value
	}

	for _, v := range keys {
		result[v] = kvs[v]
	}

	return nil;
}