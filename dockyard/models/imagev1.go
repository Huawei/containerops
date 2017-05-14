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
	"time"
)

//
type ImageV1 struct {
	ID          int64      `json:"id" gorm:"primary_key"`
	Namespace   string     `json:"namespace" sql:"not null;type:varchar(255)"`
	Repository  string     `json:"repository" sql:"not null;type:varchar(255)"`
	Short       string     `json:"short" sql:"null;type:text"`
	Description string     `json:"description" sql:"null;type:text"`
	Manifests   string     `json:"manifests" sql:"null;type:text"`
	Type        string     `json:"type" sql:"not null;type:varchar(255)"`
	Keys        string     `json:"keys" sql:"null;type:text"`
	Size        int64      `json:"size" sql:"default:0"`
	Locked      bool       `json:"locked" sql:"default:false"`
	CreatedAt   time.Time  `json:"create_at" sql:""`
	UpdatedAt   time.Time  `json:"update_at" sql:""`
	DeletedAt   *time.Time `json:"delete_at" sql:"index"`
}

//
func (*ImageV1) TableName() string {
	return "image_v1"
}

//
type VirtualV1 struct {
	ID        int64      `json:"id" gorm:"primary_key"`
	ImageV1   int64      `json:"image_v1" sql:"not null"`
	OS        string     `json:"os" sql:"null;type:varchar(255)"`
	Arch      string     `json:"arch" sql:"null;type:varchar(255)"`
	Image     string     `json:"image" sql:"not null;varchar(255)" gorm:"unique_index:image_tag"`
	Tag       string     `json:"tag" sql:"null;varchar(255)" gorm:"unique_index:image_tag"`
	Manifests string     `json:"manifests" sql:"null;type:text"`
	OSS       string     `json:"oss" sql:"null;type:text"`
	Path      string     `json:"arch" sql:"null;type:text"`
	Size      int64      `json:"size" sql:"default:0"`
	Locked    bool       `json:"locked" sql:"default:false"`
	CreatedAt time.Time  `json:"create_at" sql:""`
	UpdatedAt time.Time  `json:"update_at" sql:""`
	DeletedAt *time.Time `json:"delete_at" sql:"index"`
}

func (*VirtualV1) TableName() string {
	return "virtual_v1"
}
