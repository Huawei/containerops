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

import "time"

//ServiceDefinition is the list of DevOps service already integration the Pilotage.
type ServiceDefinition struct {
	ID            int64      `json:"id" gorm:"primary_key"`
	Service       string     `json:"service" sql:"unique;not null;type:varchar(255)"` //Service's name for query
	Title         string     `json:"title" sql:"null;type:varchar(255)"`              //Service's name for display
	Gravatar      string     `json:"gravatar" sql:"null;type:text"`                   //Service's logo.
	Endpoints     string     `json:"endpoints" sql:"null;type:text"`                  //Service's endpoints include multi endpoints url with JSON format.
	Status        string     `json:"status" sql:"null;type:text"`                     //Service's service status endpoint.
	Environment   string     `json:"environment" sql:"null;type:text"`                //
	Authorization string     `json:"authorization" sql:"null;type:text"`              //
	Configuration string     `json:"configuration" sql:"null;type:text"`              //
	Description   string     `json:"description" sql:"null;type:text"`                //
	CreatedAt     time.Time  `json:"created" sql:""`                                  //
	UpdatedAt     time.Time  `json:"updated" sql:""`                                  //
	DeletedAt     *time.Time `json:"deleted" sql:"index"`                             //
}

//TableName is return the table name of Component in MySQL database.
func (sd *ServiceDefinition) TableName() string {
	return "serivce_definition"
}

//Service is third DevOps service outside the system, is must be the one of Service List.
type Service struct {
	ID            int64      `json:"id" gorm:"primary_key"`                                                            //
	Namespace     string     `json:"namespace" sql:"not null;type:varchar(255)" gorm:"unique_index:namespace_service"` //User name or organization name
	Service       string     `json:"service" sql:"not null;type:varchar(255)" gorm:"unique_index:namespace_service"`   //Service's name for query.
	Title         string     `json:"title" sql:"null;type:varchar(255)"`                                               //
	Gravatar      string     `json:"gravatar" sql:"null;type:text"`                                                    //
	Endpoints     string     `json:"endpoints" sql:"null;type:text"`                                                   //
	Environment   string     `json:"environment" sql:"null;type:text"`                                                 //
	Authorization string     `json:"authorization" sql:"null;type:text"`                                               //
	Configuration string     `json:"configuration" sql:"null;type:text"`                                               //
	Description   string     `json:"description" sql:"null;type:text"`                                                 //
	CreatedAt     time.Time  `json:"created" sql:""`                                                                   //
	UpdatedAt     time.Time  `json:"updated" sql:""`                                                                   //
	DeletedAt     *time.Time `json:"deleted" sql:"index"`                                                              //
}

//TableName is return the table name of Component in MySQL database.
func (s *Service) TableName() string {
	return "service"
}

//Create is
func (sd *ServiceDefinition) Create(service, title, gravatar, endpoints, environments, authorizations, configurations, descriptin string) (int64, error) {
	tx := conn.Begin()

	if err := tx.Debug().Where("service = ?", service).FirstOrCreate(&sd).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()
	return sd.ID, nil
}
