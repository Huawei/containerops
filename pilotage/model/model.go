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
	DB *gorm.DB
)

// OpenDatabase is
func OpenDatabase(dbconfig *common.DatabaseConfig) {
	var err error
	driver, host, port, user, password, db := dbconfig.Driver, dbconfig.Host, dbconfig.Port, dbconfig.User, dbconfig.Password, dbconfig.Name
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=True&loc=Local", user, password, host, port, db)
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


