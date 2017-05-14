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
	"encoding/json"
	"time"
)

//DockerV1 is Docker Repository V1 repository.
type DockerV1 struct {
	ID          int64      `json:"id" gorm:"primary_key"`
	Namespace   string     `json:"namespace" sql:"not null;type:varchar(255)" gorm:"unique_index:dockerv1_repository"`
	Repository  string     `json:"repository" sql:"not null;type:varchar(255)" gorm:"unique_index:dockerv1_repository"`
	JSON        string     `json:"json" sql:"null;type:text"`
	Manifests   string     `json:"manifests" sql:"null;type:text"`
	Agent       string     `json:"agent" sql:"null;type:text"`
	Short       string     `json:"short" sql:"null;type:text"`
	Description string     `json:"description" sql:"null;type:text"`
	Size        int64      `json:"size" sql:"default:0"`
	Locked      bool       `json:"locked" sql:"default:false"` //When create/update the repository, the locked will be true.
	CreatedAt   time.Time  `json:"create_at" sql:""`
	UpdatedAt   time.Time  `json:"update_at" sql:""`
	DeletedAt   *time.Time `json:"delete_at" sql:"index"`
}

//TableName in mysql is "docker_v1".
func (r *DockerV1) TableName() string {
	return "docker_v1"
}

//DockerImageV1 is
type DockerImageV1 struct {
	ID         int64      `json:"id" gorm:"primary_key"`
	ImageID    string     `json:"image_id" sql:"not null;unique;varchar(255)"`
	JSON       string     `json:"json" sql:"null;type:text"`
	Ancestry   string     `json:"ancestry" sql:"null;type:text"`
	Checksum   string     `json:"checksum" sql:"null;type:varchar(255)"`
	Payload    string     `json:"payload" sql:"null;type:varchar(255)"`
	Path       string     `json:"path" sql:"null;type:text"`
	OSS        string     `json:"oss" sql:"null;type:text"`
	Size       int64      `json:"size" sql:"default:0"`
	Uploaded   bool       `json:"uploaded" sql:"default:false"`
	Checksumed bool       `json:"checksumed" sql:"default:false"`
	Locked     bool       `json:"locked" sql:"default:false"`
	CreatedAt  time.Time  `json:"create_at" sql:""`
	UpdatedAt  time.Time  `json:"update_at" sql:""`
	DeletedAt  *time.Time `json:"delete_at" sql:"index"`
}

//TableName in mysql is "docker_image_v1".
func (i *DockerImageV1) TableName() string {
	return "docker_image_v1"
}

//DockerTagV1 is
type DockerTagV1 struct {
	ID        int64      `json:"id" gorm:"primary_key"`
	DockerV1  int64      `json:"docker_v1" sql:"not null;default:0"`
	Tag       string     `json:"tag" sql:"not null;varchar(255)"`
	ImageID   string     `json:"image_id" sql:"not null;varchar(255)"`
	CreatedAt time.Time  `json:"create_at" sql:""`
	UpdatedAt time.Time  `json:"update_at" sql:""`
	DeletedAt *time.Time `json:"delete_at" sql:"index"`
}

//TableName in mysql is "docker_tag_v1".
func (t *DockerTagV1) TableName() string {
	return "docker_tag_v1"
}

//Put function will create or update repository.
func (r *DockerV1) Put(namespace, repository, json, agent string) error {
	r.Namespace, r.Repository, r.JSON, r.Agent, r.Locked = namespace, repository, json, agent, true

	tx := DB.Begin()

	if err := tx.Debug().Where("namespace = ? AND repository = ? ", namespace, repository).FirstOrCreate(&r).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Debug().Model(&r).Updates(map[string]interface{}{"json": json, "agent": agent, "locked": true}).Error; err != nil {
		tx.Rollback()
		return err
	} else if err == nil {
		tx.Commit()
		return nil
	}

	tx.Commit()
	return nil
}

//Unlocked is Unlocked repository data so could pull.
func (r *DockerV1) Unlocked(namespace, repository string) error {
	tx := DB.Begin()

	if err := tx.Debug().Where("namespace = ? AND repository = ?", namespace, repository).First(&r).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Debug().Model(&r).Updates(map[string]interface{}{"locked": false}).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

//Get return Docker V1 repository data.
func (r *DockerV1) Get(namespace, repository string) (DockerV1, error) {
	if err := DB.Debug().Where("namespace = ? AND repository = ?", namespace, repository).First(&r).Error; err != nil {
		return *new(DockerV1), err
	} else {
		return *r, nil
	}
}

//GetTags return tas data of repository.
func (r *DockerV1) GetTags(namespace, repository string) (map[string]string, error) {
	if err := DB.Debug().Where("namespace = ? AND repository = ?", namespace, repository).First(&r).Error; err != nil {
		return map[string]string{}, err
	} else {
		var tags []DockerTagV1
		result := map[string]string{}

		if err := DB.Debug().Where("docker_v1 = ?", r.ID).Find(&tags).Error; err != nil {
			return map[string]string{}, err
		}

		for _, tag := range tags {
			result[tag.Tag] = tag.ImageID
		}

		return result, nil
	}
}

//Get is search image by ImageID.
func (i *DockerImageV1) Get(imageID string) (DockerImageV1, error) {
	if err := DB.Debug().Where("image_id = ?", imageID).First(&i).Error; err != nil {
		return *i, err
	} else {
		return *i, nil
	}
}

//PutJSON is put image json by ImageID.
func (i *DockerImageV1) PutJSON(imageID, json string) error {
	i.ImageID = imageID

	tx := DB.Begin()

	if err := tx.Debug().Where("image_id = ?", imageID).FirstOrCreate(&i).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Debug().Model(&i).Updates(map[string]interface{}{"json": json, "uploaded": false, "checksumed": false}).Error; err != nil {
		tx.Rollback()
		return err
	} else if err == nil {
		tx.Commit()
		return nil
	}

	tx.Commit()
	return nil
}

//PutLayer is put image layer, path, uploaded and size.
func (i *DockerImageV1) PutLayer(imageID, path string, size int64) error {
	tx := DB.Begin()

	if err := tx.Debug().Where("image_id = ?", imageID).First(&i).Updates(map[string]interface{}{"path": path, "uploaded": true, "size": size}).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

//PutChecksum is put image's checksum, payload and ancestry.
func (i *DockerImageV1) PutChecksum(imageID, checksum, payload string) error {
	tx := DB.Begin()

	var data map[string]interface{}
	var ancestries []string
	var parentAnestries []string

	if err := tx.Debug().Where("image_id = ?", imageID).First(&i).Error; err != nil {
		tx.Rollback()
		return err
	}

	ancestries = append(ancestries, imageID)

	if err := json.Unmarshal([]byte(i.JSON), &data); err != nil {
		tx.Rollback()
		return err
	}

	if value, has := data["parent"]; has == true {
		image := new(DockerImageV1)

		if err := tx.Debug().Where("image_id = ?", value.(string)).First(&image).Error; err != nil {
			tx.Rollback()
			return err
		}

		if err := json.Unmarshal([]byte(image.Ancestry), &parentAnestries); err != nil {
			tx.Rollback()
			return err
		}

		ancestries = append(ancestries, parentAnestries...)
	}

	ancestry, _ := json.Marshal(ancestries)

	if err := tx.Debug().Model(&i).Updates(map[string]interface{}{"checksum": checksum, "payload": payload, "ancestry": string(ancestry)}).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

//Put is set tag in the database.
func (t *DockerTagV1) Put(imageID, tag, namespace, repository string) error {
	tx := DB.Begin()

	r := new(DockerV1)
	if err := tx.Debug().Where("namespace = ? AND repository = ?", namespace, repository).First(&r).Error; err != nil {
		tx.Rollback()
		return err
	}

	t.DockerV1 = r.ID
	t.ImageID = imageID
	t.Tag = tag
	if err := tx.Debug().Where("docker_v1 = ? AND tag = ?", r.ID, tag).FirstOrCreate(&t).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Debug().Model(&t).Updates(map[string]interface{}{"image_id": imageID}).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
