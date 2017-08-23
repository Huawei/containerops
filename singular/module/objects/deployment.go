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

package objects

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"

	. "github.com/logrusorgru/aurora"
	homeDir "github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v2"

	"github.com/Huawei/containerops/common"
	"github.com/Huawei/containerops/common/utils"
	"github.com/Huawei/containerops/singular/module/tools"
)

// Deployment is Singular base struct.
type Deployment struct {
	URI         string                 `json:"uri" yaml:"uri"`
	Title       string                 `json:"title" yaml:"title"`
	Version     int                    `json:"version" yaml:"version"`
	Tag         string                 `json:"tag" yaml:"tag"`
	Nodes       []Node                 `json:"nodes" yaml:"nodes"`
	Service     Service                `json:"service" yaml:"service"`
	Tools       Tools                  `json:"tools" yaml:"tools"`
	Infras      []Infra                `json:"infra" yaml:"infras"`
	Description string                 `json:"description" yaml:"description"`
	Short       string                 `json:"short" yaml:"short"`
	Logs        []string               `json:"logs,omitempty" yaml:"logs,omitempty"`
	Config      string                 `json:"-" yaml:"-"`
	Verbose     bool                   `json:"-" yaml:"-"`
	Timestamp   bool                   `json:"-" yaml:"-"`
	Outputs     map[string]interface{} `json:"-" yaml:"-"`
}

// ParseFromFile parse deploy template and set some configs value.
func (d *Deployment) ParseFromFile(t string, verbose, timestamp bool, output string) error {
	if data, err := ioutil.ReadFile(t); err != nil {
		return err
	} else {
		// Unmarshal template
		if err := yaml.Unmarshal(data, &d); err != nil {
			return err
		}

		// Set configs value
		d.Verbose, d.Timestamp = verbose, timestamp
		if err := d.InitConfigPath(output); err != nil {
			return err
		}
	}

	return nil
}

func (d *Deployment) DownloadBinaryFile(file, url string, nodes map[string]string) error {
	for _, ip := range nodes {
		chmodCmd := fmt.Sprintf("chmod +x %s", path.Join(tools.BinaryServerPath, file))

		if err := tools.DownloadComponent(url, path.Join(tools.BinaryServerPath, file), ip, d.Tools.SSH.Private, tools.DefaultSSHUser); err != nil {
			return err
		}

		if err := utils.SSHCommand("root", d.Tools.SSH.Private, ip, 22, chmodCmd, os.Stdout, os.Stderr); err != nil {
			return err
		}

	}

	return nil
}

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
		return "", "", "", fmt.Errorf("invalid deployment URI: %s", d.URI)
	}

	namespace, repository, name = array[0], array[1], array[2]

	return namespace, repository, name, nil
}

// TODO filter the log print with different color.
func (d *Deployment) Log(log string) {
	d.Logs = append(d.Logs, fmt.Sprintf("[%s] %s", time.Now().String(), log))

	if d.Verbose == true {
		if d.Timestamp == true {
			fmt.Println(Cyan(fmt.Sprintf("[%s] %s", time.Now().String(), strings.TrimSpace(log))))
		} else {
			fmt.Println(Cyan(log))
		}
	}
}

func (d *Deployment) Output(key, value string) {
	if d.Outputs == nil {
		d.Outputs = map[string]interface{}{}
	}

	d.Outputs[key] = value
}

func (d *Deployment) InitConfigPath(path string) error {
	if path == "" {
		namespace, repository, name, _ := d.URIs()
		home, _ := homeDir.Dir()
		d.Config = fmt.Sprintf("%s/.containerops/singular/%s/%s/%s/%d", home, namespace, repository, name, d.Version)
	}

	if utils.IsDirExist(d.Config) == false {
		if err := os.MkdirAll(d.Config, os.ModePerm); err != nil {
			return err
		}
	}

	return nil
}

// Check Sequence: CheckServiceAuth -> TODO Check Other?
func (d *Deployment) Check() error {
	if err := d.CheckServiceAuth(); err != nil {
		if len(d.Nodes) == 0 {
			return err
		}
	}

	return nil
}

// CheckServiceAuth
func (d *Deployment) CheckServiceAuth() error {
	if d.Service.Provider == "" || d.Service.Token == "" {
		if common.Singular.Provider == "" || common.Singular.Token == "" {
			return fmt.Errorf("Should provide infra service and auth token in %s", "deploy template, or configuration file")
		} else {
			d.Service.Provider, d.Service.Token = common.Singular.Provider, common.Singular.Token
		}
	}

	return nil
}
