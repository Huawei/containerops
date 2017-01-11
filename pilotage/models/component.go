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
	log "github.com/Sirupsen/logrus"
	"github.com/jinzhu/gorm"
	"time"
)

type ComponentType string

const (
	ComponentTypeKubernetes ComponentType = "Kubernetes"
	ComponentTypeMesos      ComponentType = "Mesos"
	ComponentTypeSwarm      ComponentType = "Swarm"
)

var ComponentTypes = []ComponentType{ComponentTypeKubernetes, ComponentTypeMesos, ComponentTypeSwarm}

type Component struct {
	ID           int64  `gorm:"primary_key;AUTO_INCREMENT"`
	Name         string `sql:"not null;type:varchar(100);index:idx_component_1"`
	Version      string `sql:"not null;type:varchar(30);index:idx_component_1"`
	Type         int    `sql:"not null;default:0"` //0-kubernetes 1-mesos 2-swarm
	ImageName    string `sql:"not null;varchar(100);index:idx_component_2"`
	ImageTag     string `sql:"null;varchar(30)";index:idx_component_2`
	ImageSetting string `sql:"null;type:text"`
	Timeout      int    `sql:"default:0"`
	UseAdvanced  bool   `sql:"not null;default:false"`
	KubeSetting  string `sql:"null;type:text"`
	Input        string `sql:"null;type:text"`
	Output       string `sql:"null;type:text"`
	Environment  string `sql:"null;type:text"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}

//TableName is return the table name of Component in MySQL database.
func (c *Component) TableName() string {
	return "component"
}

func (c *Component) GetComponent() *gorm.DB {
	return db.Model(&Component{})
}

func (component *Component) Create() error {
	return db.Create(component).Error
}

func (condition *Component) SelectComponent() (component *Component, err error) {
	var result Component
	err = db.Where(condition).First(&result).Error
	component = &result
	return
}

func SelectComponents(name, version string, fuzzy bool, pageNum, versionNum, offset int) (components []Component, err error) {
	var offsetCond, cond string
	values := make([]interface{}, 0)
	if name != "" {
		if fuzzy {
			cond = " where name like ? "
			values = append(values, name+"%")
		} else {
			cond = " where name = ? "
			values = append(values, name)
		}
	}
	if version != "" {
		if cond == "" {
			cond = " where version = ? "
		} else {
			cond = cond + " version = ? "
		}
		values = append(values, version)
	}
	var max int
	if name != "" && !fuzzy {
		offsetCond = " where version_num > ? and version_num <= ?"
		max = offset + versionNum
		values = append(values, offset, max)
	} else {
		offsetCond = " where page_num > ? and page_num <= ? and version_num <= ?"
		max = offset + pageNum
		values = append(values, offset, max, versionNum)
	}

	components = make([]Component, 0)
	tx := db.Begin()
	defer tx.Rollback()
	err = db.Exec("set @page_num = 0").Error
	if err != nil {
		return
	}
	err = db.Exec("set @version_num = 0").Error
	if err != nil {
		return
	}
	err = db.Exec("set @name = ''").Error
	if err != nil {
		return
	}
	raw := "select id, name, version " +
		"from (select id, name, version, " +
		"(case when @name != name then @page_num := @page_num + 1 else @page_num end) as page_num, " +
		"(case when @name != name then @version_num := 1 else @version_num := @version_num + 1 end) as version_num, " +
		"@name := name " +
		"from component " +
		cond +
		"order by name, version) t" +
		offsetCond
	log.Debugf("SelectComponents raw sql string: %s\n", raw)
	err = db.Raw(raw, values...).Find(&components).Error
	return
}

func (component *Component) Save() error {
	return db.Save(component).Error
}

func (component *Component) Delete() error {
	return db.Delete(component).Error
}
