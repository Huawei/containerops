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
	"github.com/Huawei/containerops/singular/model"
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
	Nodes       []*Node                `json:"nodes" yaml:"nodes"`
	Service     *Service               `json:"service" yaml:"service"`
	Tools       Tools                  `json:"tools" yaml:"tools"`
	Infras      []*Infra               `json:"infra" yaml:"infras"`
	Description string                 `json:"description" yaml:"description"`
	Short       string                 `json:"short" yaml:"short"`
	Logs        []string               `json:"logs,omitempty" yaml:"logs,omitempty"`
	Template    string                 `json:"template" yaml:"template"`
	Config      string                 `json:"-" yaml:"-"`
	Outputs     map[string]interface{} `json:"-" yaml:"-"`
}

//WriteLog implement Logger interface.
func (d *Deployment) WriteLog(log string, writer io.Writer, output bool) error {
	d.Logs = append(d.Logs, log)

	if output == true {
		if _, err := io.WriteString(writer, fmt.Sprintf("%s\n", log)); err != nil {
			return err
		}
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

		if d.Service == nil {
			d.Service = &Service{}
		}

		//Set configs value
		d.Namespace, d.Repository, d.Name, _ = d.URIs()
		d.Template = string(data)
		if d.Config, err = initConfigPath(d.Namespace, d.Repository, d.Name, output, d.Version); err != nil {
			return err
		}
	}

	return nil
}

//Download binary file and change mode +x
func (d *Deployment) DownloadBinaryFile(file, url string, nodes []*Node, stdout io.Writer, timestamp bool) error {
	for _, node := range nodes {
		files := []map[string]string{
			{
				"src":  url,
				"dest": path.Join(tools.BinaryServerPath, file),
			},
		}

		if err := tools.DownloadComponent(files, node.IP, d.Tools.SSH.Private, node.User, stdout); err != nil {
			return err
		}

		chmodCmd := fmt.Sprintf("chmod +x %s", path.Join(tools.BinaryServerPath, file))
		if err := utils.SSHCommand(node.User, d.Tools.SSH.Private, node.IP, tools.DefaultSSHPort, []string{chmodCmd}, stdout, os.Stderr); err != nil {
			return err
		}
		WriteLog(fmt.Sprintf("%s exec in %s node", chmodCmd, node.IP), stdout, timestamp, d, node)
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

//Output gather data of deployment.
func (d *Deployment) Output(key, value string) {
	if d.Outputs == nil {
		d.Outputs = map[string]interface{}{}
	}

	d.Outputs[key] = value
}

//Check sequence: CheckServiceAuth
func (d *Deployment) Check() error {
	if err := d.CheckServiceAuth(); err != nil {
		if len(d.Nodes) == 0 {
			return err
		}
	}

	return nil
}

//CheckServiceAuth check has service token in deploy template file or environment variables.
func (d *Deployment) CheckServiceAuth() error {
	if d.Service.Provider == "" || d.Service.Token == "" {
		if common.Singular.Provider == "" || common.Singular.Token == "" {
			return fmt.Errorf("should provide infra service and auth token in deploy template, or configuration file")
		} else {
			d.Service.Provider, d.Service.Token = common.Singular.Provider, common.Singular.Token
		}
	}

	return nil
}

//Save func save deployment data into database
func (d *Deployment) Save() error {
	singular := new(model.SingularV1)
	deploy := new(model.DeploymentV1)

	if err := singular.Put(d.Namespace, d.Repository, d.Name); err != nil {
		return err
	}

	if err := deploy.Put(singular.ID, d.Tag); err != nil {
		return err
	}

	s, _ := d.Service.YAML()
	log, _ := yaml.Marshal(d.Logs)

	files := map[string]string{}
	for _, key := range []string{tools.CARootConfigFile, tools.CARootCSRConfigFile, tools.CARootPemFile, tools.CARootCSRFile, tools.CARootKeyFile} {
		if data, err := ioutil.ReadFile(d.Outputs[key].(string)); err != nil {
			return err
		} else {
			files[key] = string(data)
		}
	}
	caData, _ := yaml.Marshal(files)

	if err := deploy.Update(deploy.ID, string(s), string(log), d.Short, d.Description, len(d.Nodes), d.Template, string(caData)); err != nil {
		return err
	}

	for _, infra := range d.Infras {
		if err := infra.Save(deploy.ID); err != nil {
			return err
		}
	}

	if err := deploy.UpdateResult(deploy.ID, true); err != nil {
		return err
	}

	return nil
}

//initConfigPath init config files and log files folder.
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

//Node used for deploy with server already exist.
type Node struct {
	ID      int      `json:"id" yaml:"id"`
	IP      string   `json:"ip" yaml:"ip"`
	Private string   `json:"private" yaml:"private"`
	User    string   `json:"user" yaml:"user"`
	Distro  string   `json:"distro" yaml:"distro"`
	Logs    []string `json:"logs,omitempty" yaml:"logs,omitempty"`
}

//WriteLog implement Logger interface.
func (n *Node) WriteLog(log string, writer io.Writer, output bool) error {
	n.Logs = append(n.Logs, log)

	if output == true {
		if _, err := io.WriteString(writer, fmt.Sprintf("%s\n", log)); err != nil {
			return err
		}
	}

	return nil
}

//Service is cloud provider.
type Service struct {
	Provider string   `json:"provider" yaml:"provider"`
	Token    string   `json:"token" yaml:"token"`
	Region   string   `json:"region" yaml:"region"`
	Size     string   `json:"size" yaml:"size"`
	Image    string   `json:"image" yaml:"image"`
	Nodes    int      `json:"nodes" yaml:"nodes"`
	Logs     []string `json:"logs,omitempty" yaml:"logs,omitempty"`
}

//WriteLog implement Logger interface.
func (s *Service) WriteLog(log string, writer io.Writer, output bool) error {
	s.Logs = append(s.Logs, log)

	if output == true {
		if _, err := io.WriteString(writer, fmt.Sprintf("%s\n", log)); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) YAML() ([]byte, error) {
	return yaml.Marshal(s)
}

func (s *Service) JSON() ([]byte, error) {
	return json.Marshal(&s)
}

//Tools is part of deployment, include SSH and others.
type Tools struct {
	SSH SSH `json:"ssh" yaml:"ssh"`
}

//SSH is public, private files and fingerprint data.
type SSH struct {
	Private     string `json:"private" yaml:"private"`
	Public      string `json:"public" yaml:"public"`
	Fingerprint string `json:"fingerprint" yaml:"fingerprint"`
}
