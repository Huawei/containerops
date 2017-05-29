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
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/jinzhu/gorm"
)

var (
	DB *gorm.DB
)

// init()
func init() {

}

// OpenDatabase is
func OpenDatabase() {
	var err error
	if DB, err = gorm.Open("database.driver", "database.uri"); err != nil {
		log.Fatal("Initlization database connection error.")
		os.Exit(1)
	} else {
		DB.DB()
		DB.DB().Ping()
		DB.DB().SetMaxIdleConns(10)
		DB.DB().SetMaxOpenConns(100)
		DB.SingularTable(true)
	}
}

// Migrate is
func Migrate() {
	OpenDatabase()

	DB.AutoMigrate(&DockerV2{}, &DockerImageV2{}, &DockerTagV2{})

	log.Info("Auto Migrate Dockyard Database Structs Done.")
}
