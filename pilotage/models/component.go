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

package models

import (
	"github.com/jinzhu/gorm"
)

type ComponentType string

const (
	ComponentTypeKubernetes ComponentType = "Kubernetes"
	ComponentTypeMesos      ComponentType = "Mesos"
	ComponentTypeSwarm      ComponentType = "Swarm"
)

var ComponentTypes = []ComponentType{ComponentTypeKubernetes, ComponentTypeMesos, ComponentTypeSwarm}

//Component is customized container(docker or rkt) for executing DevOps tasks.
type Component struct {
	BaseIDField
	Name        string `sql:"not null;type:varchar(100);unique_index:uix_component_1"` //Component name for query.
	Version     string `sql:"not null;type:varchar(30);unique_index:uix_component_1"`  // component version for display
	Type        int    `sql:"not null;default:0"`                                      //Container type: docker or rkt.
	ImageName   string `sql:"not null;varchar(100);index:idx_component_1"`
	ImageTag    string `sql:"varchar(30)";index:idx_component_1`
	Timeout     int    `` //
	DataFrom    string
	UseAdvanced bool   `sql:"not null;default:false"`
	KubeSetting string `sql:"null;type:text"` //Kubernetes execute script.
	Input       string `sql:"null;type:text"` //component input
	Output      string `sql:"null;type:text"` //component output
	Environment string `sql:"null;type:text"` //Environment parameters.
	//Manifest    string `json:"manifest" sql:"null;type:longtext"` //
	BaseModel2
}

//TableName is return the table name of Component in MySQL database.
func (c *Component) TableName() string {
	return "component"
}

func (c *Component) GetComponent() *gorm.DB {
	return db.Model(&Component{})
}

func (component *Component) Create() error {
	return db.Create(component).Error
}

func (condition *Component) SelectComponent() (component *Component, err error) {
	var result Component
	err = db.Where(condition).First(&result).Error
	component = &result
	return
}

func (condition *Component) SelectComponents(offset int) (components []Component, err error) {
	err = db.Offset(offset).Where(condition).
		Order("name").Order("version").
		Find(&components).Error
	return
}

func (component *Component) Save() error {
	return db.Save(component).Error
}

func (component *Component) Delete() error {
	return db.Delete(component).Error
}
