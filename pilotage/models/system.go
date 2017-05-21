/*
Copyright 2014 - 2017 Huawei Technologies Co., Ltd. All rights reserved.

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

	"github.com/jinzhu/gorm"
)

//UserSetting is user setting definition unit.
type UserSetting struct {
	ID         int64      `json:"id" gorm:"primary_key"`                       //
	Namespace  string     `json:"namespace" sql:"not null;type:varchar(255)"`  //
	Repository string     `json:"repository" sql:"not null;type:varchar(255)"` //
	Setting    string     `json:"setting" sql:"type:text"`                     //
	CreatedAt  time.Time  `json:"created" sql:""`                              //
	UpdatedAt  time.Time  `json:"updated" sql:""`                              //
	DeletedAt  *time.Time `json:"deleted" sql:"index"`                         //
}

//TableName is return the table name of UserSetting in MySQL database.
func (p *UserSetting) TableName() string {
	return "user_setting"
}

// GetUserSetting is
func (p *UserSetting) GetUserSetting() *gorm.DB {
	return conn.Model(&UserSetting{})
}
