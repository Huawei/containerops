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
	"errors"
	"fmt"
	"time"
)

// AppV1 is the App V1 repository
type AppV1 struct {
	ID          int64      `json:"id" gorm:"primary_key"`
	Namespace   string     `json:"namespace" sql:"not null;type:varchar(255)"`
	Repository  string     `json:"repository" sql:"not null;type:varchar(255)"`
	Short       string     `json:"short" sql:"null;type:text"`
	Description string     `json:"description" sql:"null;type:text"`
	Manifests   string     `json:"manifests" sql:"null;type:text"`
	Keys        string     `json:"keys" sql:"null;type:text"`
	Size        int64      `json:"size" sql:"default:0"`
	Locked      bool       `json:"locked" sql:"default:false"`
	CreatedAt   time.Time  `json:"create_at" sql:""`
	UpdatedAt   time.Time  `json:"update_at" sql:""`
	DeletedAt   *time.Time `json:"delete_at" sql:"index"`
}

// TableName returns the name of AppV1 table in mysql
func (*AppV1) TableName() string {
	return "app_v1"
}

// ArtifactV1 is the Artifcat V1 object
type ArtifactV1 struct {
	ID            int64  `json:"id" gorm:"primary_key"`
	AppV1ID       int64  `json:"app_v1_id" sql:"not null;default:0"`
	OS            string `json:"os" sql:"null;type:varchar(255)"`
	Arch          string `json:"arch" sql:"null;type:varchar(255)"`
	Type          string `json:"type" sql:"null;type:varchar(255)"`
	App           string `json:"app" sql:"not null;varchar(255)" gorm:"unique_index:app_tag"`
	Tag           string `json:"tag" sql:"null;varchar(255)" gorm:"unique_index:app_tag"`
	Manifests     string `json:"manifests" sql:"null;type:text"`
	OSS           string `json:"oss" sql:"null;type:text"`
	EncryptMethod string `json:"encrypt" sql:"null;type:text"`
	// FIXME: Path is both the `URL` of the local storage and the `KEY` of the object storage
	Path      string     `json:"path" sql:"null;type:text"`
	Size      int64      `json:"size" sql:"default:0"`
	Locked    bool       `json:"locked" sql:"default:false"`
	CreatedAt time.Time  `json:"create_at" sql:""`
	UpdatedAt time.Time  `json:"update_at" sql:""`
	DeletedAt *time.Time `json:"delete_at" sql:"index"`
}

// TableName returns the name of ArtifactV1 table in mysql
func (*ArtifactV1) TableName() string {
	return "artifact_v1"
}

// NewAppV1 returns the namespace/repository, it will create the repository if it is not exist
func NewAppV1(namespace, repository string) (AppV1, error) {
	var app AppV1
	app.Namespace = namespace
	app.Repository = repository

	//TODO: create or query in db
	return app, nil
}

// Put adds an artifact to a repository
func (app *AppV1) Put(artifact ArtifactV1) error {
	if app.Locked {
		return fmt.Errorf("AppV1 repository %s/%s is locked, please try it later.", app.Namespace, app.Repository)
	}

	//TODO: here we should both set the lock status and updated status
	//DB.SetLock(app, true)
	//defer DB.SetLock(app, false)
	return nil
}

// Delete removes an artifact from a repository
func (app *AppV1) Delete(artifact ArtifactV1) error {
	if app.Locked {
		return fmt.Errorf("AppV1 repository %s/%s is locked, please try it later.", app.Namespace, app.Repository)
	}

	return nil
}

// Get gets full info by os/arch/app/tag
func (a *ArtifactV1) Get() (ArtifactV1, error) {
	// TODO
	return *a, nil
}

func (a *ArtifactV1) GetName() string {
	if ok, _ := a.isValid(); !ok {
		return ""
	}

	var name string
	if a.Tag == "" {
		name = a.App
	} else {
		name = a.App + ":" + a.Tag
	}

	return fmt.Sprintf("%s/%s/%s", a.OS, a.Arch, name)
}

func (a *ArtifactV1) isValid() (bool, error) {
	if a.OS == "" || a.Arch == "" || a.App == "" {
		return false, errors.New("OS, Arch and App are mandatory fields")
	}

	return true, nil
}
