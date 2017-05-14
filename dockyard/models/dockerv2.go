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

//DockerV2
type DockerV2 struct {
	ID            int64      `json:"id" gorm:"primary_key"`
	Namespace     string     `json:"namespace" sql:"not null;type:varchar(255)"  gorm:"unique_index:dockerv2_repository"`
	Repository    string     `json:"repository" sql:"not null;type:varchar(255)"  gorm:"unique_index:dockerv2_repository"`
	SchemaVersion string     `json:"schema_version" sql:"not null;type:varchar(255)"`
	Manifests     string     `json:"manifests" sql:"null;type:text"`
	Agent         string     `json:"agent" sql:"null;type:text"`
	Short         string     `json:"short" sql:"null;type:text"`
	Description   string     `json:"description" sql:"null;type:text"`
	Size          int64      `json:"size" sql:"default:0"`
	Locked        bool       `json:"locked" sql:"default:false"`
	CreatedAt     time.Time  `json:"create_at" sql:""`
	UpdatedAt     time.Time  `json:"update_at" sql:""`
	DeletedAt     *time.Time `json:"delete_at" sql:"index"`
}

//TableName is
func (r *DockerV2) TableName() string {
	return "docker_v2"
}

//DockerImageV2 is
type DockerImageV2 struct {
	ID              int64      `json:"id" gorm:"primary_key"`
	ImageID         string     `json:"image_id" sql:"null;type:varchar(255)"`
	BlobSum         string     `json:"blob_sum" sql:"null;type:varchar(255)"`
	V1Compatibility string     `json:"v1_compatibility" sql:"null;type:text"`
	Path            string     `json:"path" sql:"null;type:text"`
	OSS             string     `json:"oss" sql:"null;type:text"`
	Size            int64      `json:"size" sql:"default:0"`
	Locked          bool       `json:"locked" sql:"default:false"`
	CreatedAt       time.Time  `json:"create_at" sql:""`
	UpdatedAt       time.Time  `json:"update_at" sql:""`
	DeletedAt       *time.Time `json:"delete_at" sql:"index"`
}

//TableName is
func (i *DockerImageV2) TableName() string {
	return "docker_image_v2"
}

//DockerTagV2 is
type DockerTagV2 struct {
	ID            int64      `json:"id" gorm:"primary_key"`
	DockerV2      int64      `json:"docker_v2" sql:"not null;default:0"`
	Tag           string     `json:"tag" sql:"not null;type:varchar(255)"`
	ImageID       string     `json:"image_id" sql:"not null;type:varchar(255)"`
	Manifest      string     `json:"manifest" sql:"null;type:text"`
	SchemaVersion string     `json:"schema_version" sql:"not null;type:varchar(255)"`
	CreatedAt     time.Time  `json:"create_at" sql:""`
	UpdatedAt     time.Time  `json:"update_at" sql:""`
	DeletedAt     *time.Time `json:"delete_at" sql:"index"`
}

//TableName is
func (t *DockerTagV2) TableName() string {
	return "docker_tag_v2"
}

//GetTags return tas data of repository.
func (r *DockerV2) GetTags(namespace, repository string) ([]string, error) {
	if err := DB.Debug().Where("namespace = ? AND repository = ?", namespace, repository).First(&r).Error; err != nil {
		return []string{}, err
	} else {
		var tags []DockerTagV2
		result := []string{}

		if err := DB.Debug().Where("docker_v2 = ?", r.ID).Find(&tags).Error; err != nil {
			return []string{}, err
		}

		for _, tag := range tags {
			result = append(result, tag.Tag)
		}

		return result, nil
	}
}

//Get is
func (r *DockerV2) Get(namespace, repository string) error {
	if err := DB.Debug().Where("namespace = ? AND repository =? ", namespace, repository).First(&r).Error; err != nil {
		return err
	}

	return nil
}

//Put is
func (r *DockerV2) Put(namespace, repository string) error {
	r.Namespace, r.Repository = namespace, repository

	tx := DB.Begin()

	if err := tx.Debug().Where("namespace = ? AND repository = ? ", namespace, repository).FirstOrCreate(&r).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

//PutAgent is
func (r *DockerV2) PutAgent(namespace, repository, agent, version string) error {
	r.Namespace, r.Repository = namespace, repository

	tx := DB.Begin()

	if err := tx.Debug().Where("namespace = ? AND repository = ? ", namespace, repository).FirstOrCreate(&r).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Debug().Model(&r).Updates(map[string]string{"agent": agent, "schema_version": version}).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

//Get is
func (i *DockerImageV2) Get(blobsum string) error {
	i.BlobSum = blobsum

	if err := DB.Debug().Where("blob_sum = ?", blobsum).First(&i).Error; err != nil {
		return err
	}

	return nil
}

//Put is
func (i *DockerImageV2) Put(tarsum, path string, size int64) error {
	i.BlobSum, i.Path, i.Size = tarsum, path, size

	tx := DB.Begin()

	if err := tx.Debug().Where("blob_sum = ? ", tarsum).FirstOrCreate(&i).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

//Get is get DockerTagV2 data
func (t *DockerTagV2) Get(namespace, repository, tag string) (*DockerTagV2, error) {
	r := new(DockerV2)

	if err := DB.Debug().Where("namespace = ? AND repository = ?", namespace, repository).First(&r).Error; err != nil {
		return t, err
	}

	if err := DB.Debug().Where("docker_v2 = ? AND tag = ?", r.ID, tag).First(&t).Error; err != nil {
		return t, err
	}

	return t, nil
}

//Put is
func (t *DockerTagV2) Put(namespace, repository, tag, imageID, manifest, schema string) error {
	r := new(DockerV2)

	if err := DB.Debug().Where("namespace = ? AND repository = ? ", namespace, repository).First(&r).Error; err != nil {
		return err
	}

	tx := DB.Begin()
	t.DockerV2, t.Tag, t.ImageID, t.Manifest, t.SchemaVersion = r.ID, tag, imageID, manifest, schema

	if err := tx.Debug().Where("docker_v2 = ? AND tag = ?", r.ID, tag).FirstOrCreate(&t).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Debug().Model(&t).Updates(map[string]interface{}{"image_id": imageID, "manifest": manifest, "schema_version": string(schema)}).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
