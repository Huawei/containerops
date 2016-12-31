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

const (
	//ComponentTypeDocker means the component is Docker container.
	ComponentTypeDocker = iota
	//ComponentTypeRkt means the component is rkt container.
	ComponentTypeRkt
	//ComponentTypeOCI reserved for OCI format container.
	ComponentTypeOCI
)

//Component is customized container(docker or rkt) for executing DevOps tasks.
type Component struct {
	BaseIDField
	Name        string `json:"name" sql:"not null;type:varchar(100);unique_index:uix_component_1"` //Component name for query.
	Version     string `json:"version" sql:"not null;type:varchar(30);unique_index:uix_component_1"`   // component version for display
	Type        int64  `json:"type" sql:"not null;default:0"`                                      //Container type: docker or rkt.
	ImageName   string `json:"image_name" sql:"not null;varchar(100);index:idx_component_1"`
	ImageTag    string `json:"image_tag" sql:"varchar(30)";index:idx_component_1`
	Timeout     int64  `json:"timeout"`                           //
	KubeSetting string `json:"kubernetes" sql:"null;type:text"`   //Kubernetes execute script.
	Input       string `json:"input" sql:"null;type:text"`        //component input
	Output      string `json:"output" sql:"null;type:text"`       //component output
	Environment string `json:"environment" sql:"null;type:text"`  //Environment parameters.
	Manifest    string `json:"manifest" sql:"null;type:longtext"` //
	BaseModel
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
	err = db.Where(condition).First(component).Error
	return
}

func (component *Component) Save() error {
	return db.Save(component).Error
}

func (component *Component) Delete() error {
	return db.Unscoped().Delete(component).Error
}