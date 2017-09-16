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

	"gopkg.in/yaml.v2"

	"github.com/Huawei/containerops/singular/model"
)

//Infra is Singular deploy unit.
type Infra struct {
	Name       string                 `json:"name" yaml:"name"`
	Version    string                 `json:"version" yaml:"version" `
	Master     int                    `json:"master" yaml:"master"`
	Minion     int                    `json:"minion" yaml:"minion"`
	Result     bool                   `json:"result,omitempty" yaml:"result,omitempty"`
	Components []*Component           `json:"components" yaml:"components"`
	Logs       []string               `json:"logs,omitempty" yaml:"logs,omitempty"`
	Outputs    map[string]interface{} `json:"-" yaml:"-"`
}

//WriteLog implement Logger interface.
func (i *Infra) WriteLog(log string, writer io.Writer, output bool) error {
	i.Logs = append(i.Logs, log)

	if output == true {
		if _, err := io.WriteString(writer, fmt.Sprintf("%s\n", log)); err != nil {
			return err
		}
	}

	return nil
}

//JSON export Infra data of JSON format
func (i *Infra) JSON() ([]byte, error) {
	return json.Marshal(&i)
}

//YAML export Infra data of YAML format
func (i *Infra) YAML() ([]byte, error) {
	return yaml.Marshal(&i)
}

//Output gather export data of infra
func (i *Infra) Output(key, value string) {
	if i.Outputs == nil {
		i.Outputs = map[string]interface{}{}
	}

	i.Outputs[key] = value
}

//Save infra deployment data
func (i *Infra) Save(deployment int64) error {
	infra := new(model.InfraV1)

	if err := infra.Put(deployment, i.Name, i.Version); err != nil {
		return err
	}

	log, _ := i.YAML()

	if err := infra.Update(infra.ID, i.Master, i.Minion, string(log)); err != nil {
		return err
	}

	for _, c := range i.Components {
		if err := c.Save(infra.ID); err != nil {
			return err
		}
	}

	return nil
}

//Component is part of infra
type Component struct {
	Binary  string   `json:"binary" yaml:"binary"`
	URL     string   `json:"url" yaml:"url"`
	Package bool     `json:"package" yaml:"package"`
	Systemd string   `json:"systemd" yaml:"systemd"`
	CA      string   `json:"ca" yaml:"ca"`
	Before  []string `json:"before" yaml:"before"`
	After   []string `json:"after" yaml:"after"`
	Logs    []string `json:"logs,omitempty" yaml:"logs,omitempty"`
}

//WriteLog implement Logger interface.
func (c *Component) WriteLog(log string, writer io.Writer, output bool) error {
	c.Logs = append(c.Logs, log)

	if output == true {
		if _, err := io.WriteString(writer, fmt.Sprintf("%s\n", log)); err != nil {
			return err
		}
	}

	return nil
}

//Save component deploy data.
func (c *Component) Save(infra int64) error {
	component := new(model.ComponentV1)

	if err := component.Put(infra, c.Binary, c.URL); err != nil {
		return err
	}

	before, _ := yaml.Marshal(c.Before)
	after, _ := yaml.Marshal(c.After)

	if err := component.Update(component.ID, c.URL, string(before), string(after), c.Package); err != nil {
		return err
	}

	return nil
}
