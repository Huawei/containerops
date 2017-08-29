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
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"

	homeDir "github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v2"

	"github.com/Huawei/containerops/common"
	"github.com/Huawei/containerops/common/utils"
	"github.com/Huawei/containerops/singular/module/tools"
)

//Deployment is Singular base struct.
type Deployment struct {
	Namespace   string                 `json:"namespace" yaml:"namespace"`
	Repository  string                 `json:"repository" yaml:"repository"`
	Name        string                 `json:"name" yaml:"name"`
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
	Outputs     map[string]interface{} `json:"-" yaml:"-"`
}

//WriteLog implement Logger interface.
func (d *Deployment) WriteLog(log string, writer io.Writer) error {
	d.Logs = append(d.Logs, log)

	if _, err := io.WriteString(writer, fmt.Sprintf("%s\n", log)); err != nil {
		return err
	}

	return nil
}

//ParseFromFile parse deploy template and set some configs value.
func (d *Deployment) ParseFromFile(t string, output string) error {
	if data, err := ioutil.ReadFile(t); err != nil {
		return err
	} else {
		//Unmarshal template
		if err := yaml.Unmarshal(data, &d); err != nil {
			return err
		}

		//Set configs value
		d.Namespace, d.Repository, d.Name, _ = d.URIs()
		if d.Config, err = initConfigPath(d.Namespace, d.Repository, d.Name, output, d.Version); err != nil {
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

//JSON export deployment data
func (d *Deployment) JSON() ([]byte, error) {
	return json.Marshal(&d)
}

//YAML export deployment data
func (d *Deployment) YAML() ([]byte, error) {
	return yaml.Marshal(&d)
}

//URIs return namespace, repository and deploy name of template.
func (d *Deployment) URIs() (namespace, repository, name string, err error) {
	array := strings.Split(d.URI, "/")

	if len(array) != 3 {
		return "", "", "", fmt.Errorf("invalid deployment URI: %s", d.URI)
	}

	namespace, repository, name = array[0], array[1], array[2]

	return namespace, repository, name, nil
}

func (d *Deployment) Output(key, value string) {
	if d.Outputs == nil {
		d.Outputs = map[string]interface{}{}
	}

	d.Outputs[key] = value
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

// initConfigPath init config files and log files folder.
func initConfigPath(namespace, repository, name, path string, version int) (string, error) {
	var config string

	if path == "" {
		home, _ := homeDir.Dir()
		config = fmt.Sprintf("%s/.containerops/singular/%s/%s/%s/%d", home, namespace, repository, name, version)
	} else {
		config = path
	}

	if utils.IsDirExist(config) == true {
		os.RemoveAll(config)
	}

	if err := os.MkdirAll(config, os.ModePerm); err != nil {
		return "", err
	}

	return config, nil
}
