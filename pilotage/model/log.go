package model

import (
	"sync"
	"time"
)

type LogV1 struct {
	ID    int64  `json:"id" gorm:"primary_key" gorm:"column:id"`
	Level string `json:"level" sql:"not null;type:varchar(255)" gorm:"column:level"`
	//Phase must be one of 'flow','stage','action' or 'job'
	Phase     string    `json:"phase" sql:"type:varchar(255)" gorm:"column:level"`
	PhaseID   int64     `json:"phase_id" sql:"type:bigint(20)" gorm:"column:phase_id"`
	Content   string    `json:"content" sql:"type:text" gorm:"column:content"`
	EventTime time.Time `json:"envent_time" sql:"" gorm:"column:envent_time"`
}

func (l *LogV1) TableName() string {
	return "log_v1"
}

func (l *LogV1) Create(level, phase string, phaseID int64, content string) error {

	l.Level, l.Phase, l.PhaseID, l.Content = level, phase, phaseID, content
	l.EventTime = time.Now()

	tx := DB.Begin()

	mutex := &sync.Mutex{}
	mutex.Lock()
	defer mutex.Unlock()

	if err := tx.Debug().Create(&l).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
