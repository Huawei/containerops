// /*
// Copyright 2014 Huawei Technologies Co., Ltd. All rights reserved.

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// */

package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

const (
	// ComponentContainerTypeDocker is the Container type of component
	ComponentContainerTypeDocker = iota
	// ComponentContainerTypeRKT is the Container type of component
	ComponentContainerTypeRKT
)

const (
	// ComponentEngineTypeKube is the Engine type of component
	ComponentEngineTypeKube = iota
	// ComponentEngineTypeSwarm is the Engine type of component
	ComponentEngineTypeSwarm
)

var (
	// ComponentEngineMap define all component engine type
	ComponentEngineMap = map[string]int{
		"Kubernetes": ComponentEngineTypeKube,
		"Swarm":      ComponentEngineTypeSwarm,
	}
)

//Component is customized container(docker or rkt) for executing DevOps tasks.
type Component struct {
	ID            int64      `json:"id" gorm:"primary_key"`                    //
	Namespace     string     `json:"_" sql:"not null;type:varchar(255)"`       //
	Name          string     `json:"name" sql:"not null;type:varchar(100)"`    // Component name for query.
	Version       string     `json:"version" sql:"not null;type:varchar(30)"`  // Component version for display
	Containertype int        `json:"-" sql:"not null;default:0"`               // Container type: docker or rkt.
	Runengine     int        `json:"-" sql:"not null;defalut:0"`               // Component engine tye : kube or swarm
	ImageName     string     `json:"imageName"`                                //
	ImageTag      string     `json:"imageTag"`                                 //
	Timeout       int        `json:"timeout" sql:"default 0"`                  //
	UseAdvanced   bool       `json:"useAdvanced" sql:"not null;default:false"` // is component use advance setting, if true, when send create msg to k8s/swarm use setting directly
	KubeSetting   string     `json:"-" sql:"null;type:text"`                   // Kubernetes execute script.
	Input         string     `json:"-" sql:"null;type:text"`                   // Component input
	Output        string     `json:"-" sql:"null;type:text"`                   // Component output
	Environment   string     `json:"-" sql:"null;type:text"`                   // Environment parameters.
	BaseImageName string     `json:"-"`                                        // base Img that component image use
	BaseImageTag  string     `json:"-"`                                        // base Img's tag that component image use
	EventShell    string     `json:"-"`                                        // shell info that component image incloud
	CreatedAt     time.Time  `json:"created" sql:""`                           //
	UpdatedAt     time.Time  `json:"updated" sql:""`                           //
	DeletedAt     *time.Time `json:"deleted" sql:"index"`                      //
}

//TableName is return the table name of Component in MySQL database.
func (c *Component) TableName() string {
	return "component"
}

// GetComponent is return the db conn with model is component
func (c *Component) GetComponent() *gorm.DB {
	if db == nil {
		OpenDatabase()
	}
	return db.Model(&Component{})
}
