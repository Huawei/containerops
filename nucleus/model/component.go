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

package model

import (
	"time"
)

//The Component is a container image, and there are three image types in the community: Docker, Appc, OCI.
//Just now, the ContainerOps platform only supports the Docker image spec, will support Appc and OCI spec in the feature.
//
//Reference:
//1. Docker Image Spec - https://github.com/docker/distribution/blob/master/docs/spec/manifest-v2-2.md
//2. Appc              - https://github.com/appc/spec/blob/master/spec/aci.md
//3. OCI Image Spec    - https://github.com/opencontainers/image-spec
//4. Rkt               - https://github.com/coreos/rkt
const (
	//ComponentTypeDocker is the docker image type
	ComponentTypeDocker = "DOCKER"
	//ComponentTypeAppc is the Appc image type used by Rkt.
	ComponentTypeAppc = "APPC"
	//ComponentTypeOCI is the OCI image type
	ComponentTypeOCI = "OCI"
)

//The Component is a container image encapsulated DevOps program written in any programming language like Bush, Python or Ruby.
type Component struct {
	ID          int64      `json:"id" gorm:"primary_key"`                                                                    //
	Namespace   string     `json:"namespace" gorm:"not null;type:varchar(128);unique_index:idx_component_namespace_version"` //Namespace is User or Organization name, the regex grammar is "/[a-z0-9]+(?:[-][a-z0-9]+){6,127}/". The User or Organization have 1:n relationships with components. If the component only creates, update or delete by the system administrator, we call it is a library. And the field value is "LIBRARY."
	Name        string     `json:"name" gorm:"not null;type:varchar(128);unique_index:idx_component_namespace_version"`      //The Component name, the regex grammar is "/[a-z0-9]+(?:[-][a-z0-9]+){6,127}/".
	Version     string     `json:"version" gorm:"null;type:varchar(128);unique_index:idx_component_namespace_version"`       //The Component version, the regex grammar is "/[\w][\w.-]{0,127}/". The version doesn't link tag of Docker or rkt.
	Type        string     `json:"type" gorm:"not null;type:varchar(128)"`                                                   //The Component type link to the [ComponentTypeDocker, ComponentTypeAppc, ComponentTypeOCI]
	Endpoint    string     `json:"endpoint" gorm:"null;type:text"`                                                           //The Component endpoint is the container URI like dockyard.sh/genedna/cloudnativeday:1.0.
	Gravatar    string     `json:"gravatar" gorm:"null;type:text"`                                                           //The component also has a gravatar like a user or an organization. The system will crawl the registry/hub of the endpoint. If the crawling failure, the system will generate a gravataing picture. And user also could upload a picture take place the crawling or generate picture.
	ImageName   string     `sql:"not null;type:varchar(256);index:idx_component_iamge_name"`                                 //Image name it must match the regular expression [a-z0-9]+(?:[._-][a-z0-9]+)* , and must be less than 256 characters. Specification at [Docker Registry V2 Sepcification](https://github.com/docker/distribution/blob/master/docs/spec/api.md#overview)
	ImageTag    string     `sql:"null;type:varchar(255);index:idx_component_image_tag"`                                      //
	Timeout     int        `sql:"not null;default:0"`                                                                        //
	UseAdvanced bool       `sql:"not null;default:false"`                                                                    //
	KubeSetting string     `sql:"null;type:text"`                                                                            //Kubernetes execute script.
	Input       string     `sql:"null;type:text"`                                                                            //component input
	Output      string     `sql:"null;type:text"`                                                                            //component output
	Environment string     `sql:"null;type:text"`                                                                            //Environment parameters.
	Manifest    string     `sql:"null;type:longtext"`                                                                        //
	CreatedAt   time.Time  ``                                                                                                //
	UpdatedAt   time.Time  ``                                                                                                //
	DeletedAt   *time.Time ``                                                                                                //
}

//TableName is return the table name of Component in MySQL database.
func (c *Component) TableName() string {
	return "component"
}
