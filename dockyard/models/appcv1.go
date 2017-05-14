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

//AppcV1 is
type AppcV1 struct {
	ID          int64      `json:"id" gorm:"primary_key"`
	Namespace   string     `json:"namespace" sql:"not null;type:varchar(255)"  gorm:"unique_index:appcv1_repository"`
	Repository  string     `json:"repository" sql:"not null;type:varchar(255)" gorm:"unique_index:appcv1_repository"`
	Short       string     `json:"short" sql:"null;type:text"`
	Description string     `json:"description" sql:"null;type:text"`
	Keys        string     `json:"keys" sql:"null;type:text"`
	Size        int64      `json:"size" sql:"default:0"`
	CreatedAt   time.Time  `json:"create_at" sql:""`
	UpdatedAt   time.Time  `json:"update_at" sql:""`
	DeletedAt   *time.Time `json:"delete_at" sql:"index"`
}

//TableName is
func (r *AppcV1) TableName() string {
	return "appc_v1"
}

//ACIv1 is
type ACIv1 struct {
	ID        int64      `json:"id" gorm:"primary_key"`
	AppcV1    int64      `json:"appc_v1" sql:"not null;default:0" gorm:"unique_index:aciv1"`
	Name      string     `json:"name" sql:"not null;type:varchar(255)" gorm:"unique_index:aciv1"`
	OS        string     `json:"os" sql:"null;type:varchar(255)"`
	Arch      string     `json:"arch" sql:"null;type:varchar(255)"`
	Version   string     `json:"version" sql:"null;type:varchar(255)"`
	Manifest  string     `json:"manifest" sql:"null;type:text"`
	OSS       string     `json:"name" sql:"null;type:text"`
	Path      string     `json:"pach" sql:"null;type:text"`
	Sign      string     `json:"sign" sql:"null;type:text"`
	Size      int64      `json:"size" sql:"default:0"`
	Locked    bool       `json:"locked" sql:"default:false"`
	CreatedAt time.Time  `json:"create_at" sql:""`
	UpdatedAt time.Time  `json:"update_at" sql:""`
	DeletedAt *time.Time `json:"delete_at" sql:"index"`
}

//TableName is
func (i *ACIv1) TableName() string {
	return "aci_v1"
}

//Put is
func (r *AppcV1) Put(namespace, repository string) error {
	r.Namespace, r.Repository = namespace, repository

	tx := DB.Begin()

	if err := tx.Debug().Where("namespace = ? AND repository = ? ", namespace, repository).FirstOrCreate(&r).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

//Get is
func (r *AppcV1) Get(namespace, repository string) error {
	if err := DB.Debug().Where("namespace = ? AND repository = ?", namespace, repository).First(&r).Error; err != nil {
		return err
	}

	return nil
}

//PutManifest is
func (i *ACIv1) PutManifest(appcv1 int64, version, name, manifest string) error {
	i.AppcV1, i.Version, i.Name = appcv1, version, name

	tx := DB.Begin()

	if err := tx.Debug().Where("appc_v1 = ? AND name = ?", appcv1, name).FirstOrCreate(&i).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Debug().Model(&i).Updates(map[string]interface{}{"manifest": manifest, "locked": true}).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

//PutSign is
func (i *ACIv1) PutSign(appcv1 int64, version, name, sign string) error {
	i.AppcV1, i.Version, i.Name = appcv1, version, name

	tx := DB.Begin()

	if err := tx.Debug().Where("appc_v1 = ? AND name = ?", appcv1, name).FirstOrCreate(&i).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Debug().Model(&i).Updates(map[string]interface{}{"sign": sign, "locked": true}).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

//PutACI is
func (i *ACIv1) PutACI(appcv1, size int64, version, name, aci string) error {
	i.AppcV1, i.Size, i.Version, i.Name = appcv1, size, version, name

	tx := DB.Begin()

	if err := tx.Debug().Where("appc_v1 = ? AND name = ?", appcv1, name).FirstOrCreate(&i).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Debug().Model(&i).Updates(map[string]interface{}{"path": aci, "size": size, "locked": true}).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

//Get is
func (i *ACIv1) Get(appcv1 int64, version, name string) error {
	i.AppcV1, i.Version, i.Name = appcv1, version, name

	if err := DB.Debug().Where("appc_v1 = ? AND name = ? AND version = ?", appcv1, name, version).First(&i).Error; err != nil {
		return err
	}

	return nil
}

//Unlocked is
func (i *ACIv1) Unlocked(appcv1 int64, version, name string) error {
	i.AppcV1, i.Version, i.Name = appcv1, version, name

	tx := DB.Begin()

	if err := tx.Debug().Where("appc_v1 = ? AND name = ?", appcv1, name).FirstOrCreate(&i).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Debug().Model(&i).Updates(map[string]interface{}{"locked": false}).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
