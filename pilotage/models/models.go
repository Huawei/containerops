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

	if db, err = gorm.Open(configure.GetString("database.driver"), configure.GetString("database.url")); err != nil {
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

	db.AutoMigrate(&ServiceDefinition{}, &Service{}, &Component{}, &ComponentLog{})
	db.AutoMigrate(&Pipeline{}, &Stage{}, &Action{}, &Outcome{}, &PipelineLog{}, &StageLog{}, &ActionLog{})
	db.AutoMigrate(&EventDefinition{}, &Event{}, &Environment{})

	log.Info("AutMigrate database structs.")
}
