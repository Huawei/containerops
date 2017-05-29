/*
Copyright 2014 - 2017 Huawei Technologies Co., Ltd. All rights reserved.

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

package models

import (
	"time"

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
	ID             int64      `json:"id" gorm:"primary_key"`
	Namespace      string     `json:"namespace" sql:"not null;type:varchar(255)"` //User or organization.
	Version        string     `json:"version" sql:"null;type:text"`               // component version for display
	VersionCode    int64      `json:"versionCode" sql:"null;type:bigint"`         // component version code system set
	Component      string     `json:"component" sql:"not null;type:varchar(255)"` //Component name for query.
	Type           int64      `json:"type" sql:"not null;default:0"`              //Container type: docker or rkt.
	Title          string     `json:"title" sql:"null;type:varchar(255)"`         //Component name for display.
	Gravatar       string     `json:"gravatar" sql:"null;type:text"`              //Logo.
	Description    string     `json:"description" sql:"null;type:text"`           //Description with markdown style.
	Endpoint       string     `json:"endpoint" sql:"null;type:text"`              //Contaienr location like: `dockyard.sh/genedna/cloudnativeday:1.0`.
	Source         string     `json:"source" sql:"not null;type:text"`            //Component source location like: `git@github.com/containerops/components`.
	Environment    string     `json:"environment" sql:"null;type:text"`           //Environment parameters.
	Tag            string     `json:"tag" sql:"null;type:varchar(255)"`           //Tag for version.
	VolumeLocation string     `json:"volume_location" sql:"null;type:text"`       //Volume path in the container.
	VolumeData     string     `json:"volume_data" sql:"null;type:text"`           //Volume data source.
	Makefile       string     `json:"makefile" sql:"null;type:text"`              //Like Dockerfile or acbuild script.
	Kubernetes     string     `json:"kubernetes" sql:"null;type:text"`            //Kubernetes execute script.
	Swarm          string     `json:"swarm" sql:"null;type:text"`                 //Docker Swarm execute script.
	Input          string     `json:"input" sql:"null;type:text"`                 //component input
	Output         string     `json:"output" sql:"null;type:text"`                //component output
	Timeout        int64      `json:"timeout"`                                    //
	Manifest       string     `json:"manifest" sql:"null;type:longtext"`          //
	CreatedAt      time.Time  `json:"created" sql:""`                             //
	UpdatedAt      time.Time  `json:"updated" sql:""`                             //
	DeletedAt      *time.Time `json:"deleted" sql:"index"`                        //
}

//TableName is return the table name of Component in MySQL database.
func (c *Component) TableName() string {
	return "component"
}

// GetComponent is
func (c *Component) GetComponent() *gorm.DB {
	return conn.Model(&Component{})
}
