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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	. "github.com/logrusorgru/aurora"
	"gopkg.in/yaml.v2"

	"github.com/Huawei/containerops/common"
)

// JSON export deployment data
func (d *Deployment) JSON() ([]byte, error) {
	return json.Marshal(&d)
}

//
func (d *Deployment) YAML() ([]byte, error) {
	return yaml.Marshal(&d)
}

//
func (d *Deployment) URIs() (namespace, repository, name string, err error) {
	array := strings.Split(d.URI, "/")

	if len(array) != 3 {
		return "", "", "", fmt.Errorf("Invalid deployment URI: %s", d.URI)
	}

	namespace, repository, name = array[0], array[1], array[2]

	return namespace, repository, name, nil
}

// TODO filter the log print with different color.
func (d *Deployment) Log(log string, verbose, timestamp bool) {
	d.Logs = append(d.Logs, fmt.Sprintf("[%s] %s", time.Now().String(), log))

	if verbose == true {
		if timestamp == true {
			fmt.Println(Cyan(fmt.Sprintf("[%s] %s", time.Now().String(), strings.TrimSpace(log))))
		} else {
			fmt.Println(Cyan(log))
		}
	}
}

// ParseFromFile
func (d *Deployment) ParseFromFile(t string, verbose, timestamp bool) error {
	if data, err := ioutil.ReadFile(t); err != nil {
		d.Log(fmt.Sprintf("Read deployment template file %s error: %s", t, err.Error()), verbose, timestamp)
		return err
	} else {
		if err := yaml.Unmarshal(data, &d); err != nil {
			d.Log(fmt.Sprintf("Unmarshal the template file error: %s", err.Error()), verbose, timestamp)
			return err
		}
	}

	return nil
}

// Deploy
func (d *Deployment) Deploy() error {
	if err := d.CheckServiceAuth(); err != nil {
		return fmt.Errorf("check template or configuration error: %s ", err.Error())
	}

	return nil
}

// CheckServiceAuth
func (d *Deployment) CheckServiceAuth() error {
	if d.Service.Provider == "" || d.Service.Token == "" {
		if common.Singular.Provider == "" || common.Singular.Provider == "" {
			return fmt.Errorf("Should provide infra service and auth token in %s", "deploy template, or configuration file")
		} else {
			d.Service.Provider, d.Service.Token = common.Singular.Provider, common.Singular.Token
		}
	}
	return nil
}
