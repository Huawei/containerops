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

package model

import "time"

type SingularV1 struct {
	ID          int64      `json:"id" yaml:"id" gorm:"column:id;primary_key"`
	Namespace   string     `json:"namespace" yaml:"namespace" sql:"not null;type:varchar(255)" gorm:"column:namespace;unique_index:singular_repository"`
	Repository  string     `json:"repository" yaml:"repository" sql:"not null;type:varchar(255)" gorm:"column:repository;unique_index:singular_repository"`
	Name        string     `json:"name" yaml:"name" sql:"not null;type:varchar(255)" gorm:"column:name;unique_index:singular_repository"`
	Short       string     `json:"short" yaml:"short" sql:"null;type:text" gorm:"column:short"`
	Description string     `json:"description" yaml:"description" sql:"null;type:text" gorm:"column:description"`
	CreatedAt   time.Time  `json:"create_at" sql:"" gorm:"column:create_at"`
	UpdatedAt   time.Time  `json:"update_at" sql:"" gorm:"column:update_at"`
	DeletedAt   *time.Time `json:"delete_at" sql:"index" gorm:"column:delete_at"`
}

func (b *SingularV1) TableName() string {
	return "singular_v1"
}

type DeploymentV1 struct {
	ID         int64      `json:"id" yaml:"id" gorm:"column:id;primary_key"`
	SingularV1 int64      `json:"singular_v1" yaml:"singular_v1" sql:"not null;default:0" gorm:"column:singular_v1;unique_index:singular_deployment"`
	Tag        string     `json:"tag" yaml:"tag" sql:"not null;type:varchar(255)" gorm:"column:tag;unique_index:singular_deployment"`
	Version    int64      `json:"version" yaml:"version" sql:"not null;default:0" gorm:"column:version;unique_index:singular_deployment"`
	Service    string     `json:"service" yaml:"service" sql:"not null;type:varchar(255)" gorm:"column:service;unique_index:singular_deployment"`
	Result     bool       `json:"result" yaml:"result" sql:"not null" gorm:"column:result"`
	Log        string     `json:"log" yaml:"log" sql:"null;type:text" gorm:"column:log"`
	CreatedAt  time.Time  `json:"create_at" sql:"" gorm:"column:create_at"`
	UpdatedAt  time.Time  `json:"update_at" sql:"" gorm:"column:update_at"`
	DeletedAt  *time.Time `json:"delete_at" sql:"index" gorm:"column:delete_at"`
}

func (d *DeploymentV1) TableName() string {
	return "deployment_v1"
}

type InfraV1 struct {
	ID           int64      `json:"id" yaml:"id"yaml:"id"  gorm:"column:id;primary_key"`
	DeploymentV1 int64      `json:"deployment_v1" yaml:"deployment_v1" sql:"not null;default:0" gorm:"column:deployment_v1"`
	Name         string     `json:"name" yaml:"name" sql:"not null;type:varchar(255)" gorm:"column:name"`
	Version      string     `json:"version" yaml:"version" sql:"not null;type:varchar(255)" gorm:"column:version"`
	Master       int64      `json:"master" yaml:"master" sql:"not null" gorm:"column:master"`
	Node         int64      `json:"node" yaml:"node" sql:"not null" gorm:"column:node"`
	Log          string     `json:"log" yaml:"log" sql:"null;type:text" gorm:"column:log"`
	CreatedAt    time.Time  `json:"create_at" sql:"" gorm:"column:create_at"`
	UpdatedAt    time.Time  `json:"update_at" sql:"" gorm:"column:update_at"`
	DeletedAt    *time.Time `json:"delete_at" sql:"index" gorm:"column:delete_at"`
}

func (i *InfraV1) TableName() string {
	return "infra_v1"
}

type ComponentV1 struct {
	ID        int64      `json:"id" yaml:"id"yaml:"id"  gorm:"column:id;primary_key"`
	InfraV1   int64      `json:"infra_v1" yaml:"infra_v1" sql:"not null;default:0" gorm:"column:infra_v1"`
	Name      string     `json:"name" yaml:"name" sql:"not null;type:varchar(255)" gorm:"column:name"`
	URL       string     `json:"url" yaml:"url" sql:"not null;type:text" gorm:"column:url"`
	Package   bool       `json:"package" yaml:"package" sql:"not null;default:false" gorm:"column:package"`
	Install   string     `json:"install" yaml:"install" sql:"null;type:text" gorm:"column:install"`
	Systemd   string     `json:"systemd" yaml:"systemd" sql:"null;type:text" gorm:"column:systemd"`
	Setting   string     `json:"setting" yaml:"setting" sql:"null;type:text" gorm:"column:setting"`
	SSL       string     `json:"ssl" yaml:"ssl" sql:"null;type:text" gorm:"column:ssl"`
	Before    string     `json:"before" yaml:"before" sql:"null;type:text" gorm:"column:before"`
	After     string     `json:"after" yaml:"after" sql:"null;type:text" gorm:"column:after"`
	Log       string     `json:"log" yaml:"log" sql:"null;type:text" gorm:"column:log"`
	CreatedAt time.Time  `json:"create_at" sql:"" gorm:"column:create_at"`
	UpdatedAt time.Time  `json:"update_at" sql:"" gorm:"column:update_at"`
	DeletedAt *time.Time `json:"delete_at" sql:"index" gorm:"column:delete_at"`
}

func (c *ComponentV1) TableName() string {
	return "component_singular_v1"
}
