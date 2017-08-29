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
	"fmt"
	"io"
	"time"
)

type Logger interface {
	WriteLog(log string, writer io.Writer) error
}

func WriteLog(obj Logger, log string, writer io.Writer, timestamp bool) error {
	if timestamp == true {
		log = fmt.Sprintf("[%d] %s", time.Now().Unix(), log)
	}

	if err := obj.WriteLog(log, writer); err != nil {
		return err
	}

	return nil
}
