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

import (
	"sync"
	"time"
)

type SingularV1 struct {
	ID         int64      `json:"id" yaml:"id" gorm:"column:id;primary_key"`
	Namespace  string     `json:"namespace" yaml:"namespace" sql:"not null;type:varchar(255)" gorm:"column:namespace;unique_index:singular_repository"`
	Repository string     `json:"repository" yaml:"repository" sql:"not null;type:varchar(255)" gorm:"column:repository;unique_index:singular_repository"`
	Name       string     `json:"name" yaml:"name" sql:"not null;type:varchar(255)" gorm:"column:name;unique_index:singular_repository"`
	CreatedAt  time.Time  `json:"create_at" sql:"" gorm:"column:create_at"`
	UpdatedAt  time.Time  `json:"update_at" sql:"" gorm:"column:update_at"`
	DeletedAt  *time.Time `json:"delete_at" sql:"index" gorm:"column:delete_at"`
}

func (s *SingularV1) TableName() string {
	return "singular_v1"
}

var singularV1Mutex sync.Mutex

//Put get or create singular data
func (s *SingularV1) Put(namespace, repository, name string) error {
	s.Namespace, s.Repository, s.Name = namespace, repository, name

	singularV1Mutex.Lock()
	defer singularV1Mutex.Unlock()

	tx := DB.Begin()
	if err := tx.Where("namespace = ? AND repository = ? AND name = ?", namespace, repository, name).FirstOrCreate(&s).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

type DeploymentV1 struct {
	ID          int64      `json:"id" yaml:"id" gorm:"column:id;primary_key"`
	SingularV1  int64      `json:"singular_v1" yaml:"singular_v1" sql:"not null;default:0" gorm:"column:singular_v1;unique_index:singular_deployment"`
	Tag         string     `json:"tag" yaml:"tag" sql:"not null;type:varchar(255)" gorm:"column:tag;unique_index:singular_deployment"`
	Version     int64      `json:"version" yaml:"version" sql:"null;default:0" gorm:"column:version;unique_index:singular_deployment"`
	Service     string     `json:"service" yaml:"service" sql:"null;type:text" gorm:"column:service"`
	Node        int        `json:"node" yaml:"node" sql:"not null;default:0" gorm:"column:node"`
	Log         string     `json:"log" yaml:"log" sql:"null;type:text" gorm:"column:log"`
	Description string     `json:"description" yaml:"description" sql:"null;type:text" gorm:"column:description"`
	Data        string     `json:"data" yaml:"data" sql:"null;type:text" gorm:"column:data"`
	CA          string     `json:"ca" yaml:"ca" sql:"null;type:text" gorm:"column:ca"`
	Result      bool       `json:"result" yaml:"result" sql:"null" gorm:"column:result"`
	CreatedAt   time.Time  `json:"create_at" sql:"" gorm:"column:create_at"`
	UpdatedAt   time.Time  `json:"update_at" sql:"" gorm:"column:update_at"`
	DeletedAt   *time.Time `json:"delete_at" sql:"index" gorm:"column:delete_at"`
}

func (d *DeploymentV1) TableName() string {
	return "deployment_v1"
}

var deploymentV1Mutex sync.Mutex

func (d *DeploymentV1) Put(singularV1 int64, tag string) error {
	d.SingularV1, d.Tag = singularV1, tag

	deploymentV1Mutex.Lock()
	defer deploymentV1Mutex.Unlock()

	var count int64

	tx := DB.Begin()
	if err := tx.Model(&DeploymentV1{}).Where("singular_v1 = ? AND tag = ?", singularV1, tag).Count(&count).Error; err != nil {
		tx.Rollback()
		return err
	}

	d.Version = count + 1

	if err := tx.Where("singular_v1 = ? AND tag = ? AND version = ?", singularV1, tag, count+1).FirstOrCreate(&d).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (d *DeploymentV1) Update(id int64, service, log, description string, node int, data, ca string) error {
	tx := DB.Begin()

	if err := tx.Where("id = ?", id).First(&d).Error; err != nil {
		tx.Rollback()
		return err
	}

	deploymentV1Mutex.Lock()
	defer deploymentV1Mutex.Unlock()

	if err := tx.Model(&d).Update(map[string]interface{}{
		"service":     service,
		"log":         log,
		"description": description,
		"node":        node,
		"data":        data,
		"ca":          ca,
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (d *DeploymentV1) UpdateResult(id int64, result bool) error {
	tx := DB.Begin()

	if err := tx.Where("id = ?", id).First(&d).Error; err != nil {
		tx.Rollback()
		return err
	}

	deploymentV1Mutex.Lock()
	defer deploymentV1Mutex.Unlock()

	if err := tx.Model(&d).Update(map[string]interface{}{
		"result": result,
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

type InfraV1 struct {
	ID           int64      `json:"id" yaml:"id"  gorm:"column:id;primary_key"`
	DeploymentV1 int64      `json:"deployment_v1" yaml:"deployment_v1" sql:"not null;default:0" gorm:"column:deployment_v1"`
	Name         string     `json:"name" yaml:"name" sql:"not null;type:varchar(255)" gorm:"column:name"`
	Version      string     `json:"version" yaml:"version" sql:"not null;type:varchar(255)" gorm:"column:version"`
	Master       int        `json:"master" yaml:"master" sql:"not null" gorm:"column:master"`
	Minion       int        `json:"minion" yaml:"minion" sql:"not null" gorm:"column:minion"`
	Log          string     `json:"log" yaml:"log" sql:"null;type:text" gorm:"column:log"`
	CA           string     `json:"ca" yaml:"ca" sql:"null;type:text" gorm:"column:ca"`
	Setting      string     `json:"setting" yaml:"setting" sql:"null;type:text" gorm:"column:setting"`
	Systemd      string     `json:"systemd" yaml:"systemd" sql:"null;type:text" gorm:"column:systemd"`
	CreatedAt    time.Time  `json:"create_at" sql:"" gorm:"column:create_at"`
	UpdatedAt    time.Time  `json:"update_at" sql:"" gorm:"column:update_at"`
	DeletedAt    *time.Time `json:"delete_at" sql:"index" gorm:"column:delete_at"`
}

func (i *InfraV1) TableName() string {
	return "infra_v1"
}

var infraV1Mutex sync.Mutex

func (i *InfraV1) Put(deploymentID int64, name, version string) error {
	i.DeploymentV1, i.Name, i.Version = deploymentID, name, version

	infraV1Mutex.Lock()
	defer infraV1Mutex.Unlock()

	tx := DB.Begin()
	if err := tx.Where("deployment_v1 = ? AND name = ?", deploymentID, name).FirstOrCreate(&i).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (i *InfraV1) Update(id int64, master, minion int, log, systemd, setting, ca string) error {
	tx := DB.Begin()

	infraV1Mutex.Lock()
	defer infraV1Mutex.Unlock()

	if err := tx.Model(&i).Where("id = ?", id).Update(map[string]interface{}{
		"master":  master,
		"minion":  minion,
		"log":     log,
		"systemd": systemd,
		"setting": setting,
		"ca":      ca,
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

type ComponentV1 struct {
	ID        int64      `json:"id" yaml:"id"  gorm:"column:id;primary_key"`
	InfraV1   int64      `json:"infra_v1" yaml:"infra_v1" sql:"not null;default:0" gorm:"column:infra_v1"`
	Name      string     `json:"name" yaml:"name" sql:"not null;type:varchar(255)" gorm:"column:name"`
	URL       string     `json:"url" yaml:"url" sql:"not null;type:text" gorm:"column:url"`
	Package   bool       `json:"package" yaml:"package" sql:"null;default:false" gorm:"column:package"`
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

var componentV1Mutex sync.Mutex

func (c *ComponentV1) Put(infraID int64, binary, url string) error {
	c.InfraV1, c.Name, c.URL = infraID, binary, url

	componentV1Mutex.Lock()
	defer componentV1Mutex.Unlock()

	tx := DB.Begin()
	if err := tx.Where("infra_v1 = ? AND name = ?", infraID, binary).FirstOrCreate(&c).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (c *ComponentV1) Update(id int64, url, before, after string, p bool) error {
	tx := DB.Begin()

	componentV1Mutex.Lock()
	defer componentV1Mutex.Unlock()

	if err := tx.Model(&c).Where("id = ?", id).Update(map[string]interface{}{
		"url":     url,
		"before":  before,
		"after":   after,
		"package": p,
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
