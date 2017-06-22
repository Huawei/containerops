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

// BinaryV1 is
type BinaryV1 struct {
	ID          int64      `json:"id" gorm:"primary_key"`
	Namespace   string     `json:"namespace" sql:"not null;type:varchar(255)"  gorm:"unique_index:dockerv2_repository"`
	Repository  string     `json:"repository" sql:"not null;type:varchar(255)"  gorm:"unique_index:dockerv2_repository"`
	Short       string     `json:"short" sql:"null;type:text"`
	Description string     `json:"description" sql:"null;type:text"`
	Size        int64      `json:"size" sql:"default:0"`
	Locked      bool       `json:"locked" sql:"default:false"`
	CreatedAt   time.Time  `json:"create_at" sql:""`
	UpdatedAt   time.Time  `json:"update_at" sql:""`
	DeletedAt   *time.Time `json:"delete_at" sql:"index"`
}

// TableName is
func (b *BinaryV1) TableName() string {
	return "binary_v1"
}

// BinaryFileV1 is
type BinaryFileV1 struct {
	ID        int64      `json:"id" gorm:"primary_key"`
	BinaryV1  int64      `json:"binary_v1" sql:"not null;default:0" gorm:"unique_index:binaryfilev1_file"`
	Name      string     `json:"name" sql:"not null;type:varchar(255)" gorm:"unique_index:binaryfilev1_file"`
	Tag       string     `json:"tag" sql:"not null;type:varchar(255)" gorm:"unique_index:binaryfilev1_file"`
	Agent     string     `json:"agent" sql:"null;type:text"`
	SHA512    string     `json:"sha512" sql:"null;type:varchar(255)"`
	Path      string     `json:"path" sql:"null;type:text"`
	OSS       string     `json:"oss" sql:"null;type:text"`
	Size      int64      `json:"size" sql:"default:0"`
	Locked    bool       `json:"locked" sql:"default:false"`
	CreatedAt time.Time  `json:"create_at" sql:""`
	UpdatedAt time.Time  `json:"update_at" sql:""`
	DeletedAt *time.Time `json:"delete_at" sql:"index"`
}

// TableName is
func (b *BinaryFileV1) TableName() string {
	return "binary_file_v1"
}

// Put is
func (b *BinaryV1) Put(namespace, repository string) error {
	b.Namespace, b.Repository = namespace, repository

	tx := DB.Begin()

	mutex := &sync.Mutex{}
	mutex.Lock()
	defer mutex.Unlock()

	if err := tx.Debug().Where("namespace = ? AND repository = ? ", namespace, repository).FirstOrCreate(&b).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

// Get is
func (b *BinaryV1) Get(namespace, repository string) error {
	b.Namespace, b.Repository = namespace, repository

	if err := DB.Debug().Where("namespace = ? AND repository =? ", namespace, repository).First(&b).Error; err != nil {
		return err
	}

	return nil
}

// Get is
func (f *BinaryFileV1) Get(repository int64, name, tag string) error {
	f.BinaryV1, f.Name, f.Tag = repository, name, tag

	if err := DB.Debug().Where("binary_v1 =? AND name = ? AND tag = ?", repository, name, tag).First(&f).Error; err != nil {
		return err
	}

	return nil
}

// Put is
func (f *BinaryFileV1) Put(repository, size int64, name, tag, sha512, path string) error {
	f.BinaryV1, f.Size, f.Name, f.Tag, f.SHA512, f.Path = repository, size, name, tag, sha512, path

	tx := DB.Begin()

	mutex := &sync.Mutex{}
	mutex.Lock()
	defer mutex.Unlock()

	if err := tx.Debug().Where("binary_V1 = ? AND name = ? AND tag = ?", repository, name, tag).FirstOrCreate(&f).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
