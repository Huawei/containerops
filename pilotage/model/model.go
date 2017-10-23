package model

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/Huawei/containerops/common"
)

var (
	DB        *gorm.DB
	DisableDB bool = false
)

// OpenDatabase is
func OpenDatabase(dbconfig *common.DatabaseConfig) {
	var err error
	driver, host, port, user, password, db := dbconfig.Driver, dbconfig.Host, dbconfig.Port, dbconfig.User, dbconfig.Password, dbconfig.Name
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=True&loc=Local", user, password, host, port, db)

	if driver == "" {
		DisableDB = true
		return
	}

	if DB, err = gorm.Open(driver, dsn); err != nil {
		log.Fatal("Initlization database connection error.", err)
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
	if DisableDB {
		return
	}
	DB.AutoMigrate(&FlowV1{}, &FlowDataV1{})
	DB.AutoMigrate(&StageV1{}, &StageDataV1{})
	DB.AutoMigrate(&ActionV1{}, &ActionDataV1{})
	DB.AutoMigrate(&JobV1{}, &JobDataV1{})
	DB.AutoMigrate(&LogV1{})
}
