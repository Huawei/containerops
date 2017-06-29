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
	"time"
)

// LabelV1 is
type LabelV1 struct {
	ID        int64      `json:"id" gorm:"column:id;primary_key"`
	Type      string     `json:"type" sql:"not null;type:varchar(255)" gorm:"column:type;unique_index:labelv1_value"`
	Label     string     `json:"label" sql:"not null;type:varchar(255)" gorm:"column:label;unique_index:labelv1_value"`
	Value     string     `json:"value" sql:"not null;type:varchar(255)" gorm:"column:value;unique_index:labelv1_value"`
	Object    int64      `json:"object" sql:"not null; default:0" gorm:"column:object;unique_index:labelv1_value"`
	CreatedAt time.Time  `json:"create_at" sql:"" gorm:"column:create_at"`
	UpdatedAt time.Time  `json:"update_at" sql:"" gorm:"column:update_at"`
	DeletedAt *time.Time `json:"delete_at" sql:"index" gorm:"column:delete_at"`
}

// TableName is
func (l *LabelV1) TableName() string {
	return "label_v1"
}
