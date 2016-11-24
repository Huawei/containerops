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
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/containerops/configure"
)

var (
	db *gorm.DB
)

// OpenDatabase is
func OpenDatabase() {
	var err error

	if db, err = gorm.Open(configure.GetString("database.driver"), configure.GetString("database.uri")); err != nil {
		log.Fatal("Initlization database connection error.")
		os.Exit(1)
	} else {
		db.DB()
		db.DB().Ping()
		db.DB().SetMaxIdleConns(10)
		db.DB().SetMaxOpenConns(100)
		db.SingularTable(true)
	}
}

//Migrate is
func Migrate() {
	OpenDatabase()

	db.AutoMigrate(&ServiceDefinition{}, &Service{}, &Component{})
	db.AutoMigrate(&Pipeline{}, &PipelineLog{}, &PipelineSequence{}, &Stage{}, &StageLog{}, &Action{}, &ActionLog{}, &Outcome{})
	db.AutoMigrate(&EventDefinition{}, &Event{}, &EventJson{})

	log.Info("AutMigrate database structs.")
}

// GetDB is
func GetDB() *gorm.DB {
	if db != nil {
		return db
	}

	OpenDatabase()
	return db
}
