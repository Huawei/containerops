package model

import "time"

type StageV1 struct {
	ID         int64      `json:"id" gorm:"primary_key" gorm:"column:id"`
	FlowID     int64      `json:"flow_id" sql:"not null;type:bigint(20)" gorm:"column:flow_id"`
	StageType  string     `json:"stage_type" sql:"type:varchar(255)" gorm:"column:stage_type"`
	Name       string     `json:"name" sql:"not null;type:varchar(255)" gorm:"column:name"`
	Title      string     `json:"title" sql:"type:text" gorm:"column:title"`
	Sequencing string     `json:"sequencing" sql:"type:varchar(255)" gorm:"column:sequencing"`
	CreatedAt  time.Time  `json:"created_at" sql:"" gorm:"column:created_at"`
	UpdatedAt  time.Time  `json:"updated_at" sql:"" gorm:"column:updated_at"`
	DeletedAt  *time.Time `json:"deleted_at" sql:"index" gorm:"column:deleted_at"`
}

type StageDataV1 struct {
	ID      int64     `json:"id" gorm:"primary_key" gorm:"column:id"`
	StageID int64     `json:"stage_id" sql:"not null;type:bigint(20)" gorm:"column:stage_id"`
	Number  int64     `json:"number" sql:"not null;type:bigint(20)" gorm:"column:number"`
	Result  string    `json:"result" sql:"type:varchar(255)" gorm:"column:result"`
	Start   time.Time `json:"start" sql:"" gorm:"column:start"`
	End     time.Time `json:"end" sql:"" gorm:"column:end"`
}

func (s *StageV1) TableName() string {
	return "stage_v1"
}

func (sd *StageDataV1) TableName() string {
	return "stage_data_v1"
}

func (s *StageV1) Put(flowID int64, stageType, name, title, sequencing string) (stageID int64, err error) {
	if DisableDB {
		return -1, nil
	}

	s.FlowID, s.StageType, s.Name, s.Title, s.Sequencing = flowID, stageType, name, title, sequencing

	tx := DB.Begin()

	if tx.Where("flow_id = ? AND name = ? ", flowID, name).First(&s).RecordNotFound() {
		s.CreatedAt = time.Now()
		if err := tx.Create(&s).Error; err != nil {
			tx.Rollback()
			return 0, err
		}
	} else {
		if err := tx.Model(&s).Updates(StageV1{StageType: stageType, Title: title, Sequencing: sequencing}).Error; err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	tx.Commit()
	stageID = s.ID

	return stageID, nil
}

func (sd *StageDataV1) Put(stageID, number int64, result string, start, end time.Time) error {
	if DisableDB {
		return nil
	}

	sd.StageID, sd.Number, sd.Result, sd.Start, sd.End = stageID, number, result, start, end

	tx := DB.Begin()
	if err := tx.Create(&sd).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}

func (sd *StageDataV1) GetNumbers(stageID int64) (int64, error) {
	if DisableDB {
		return -1, nil
	}

	tmp := DB.Where("stage_id = ?", stageID).Find(&StageDataV1{})
	if tmp.RecordNotFound() {
		return 0, nil
	}
	if tmp.Error != nil {
		return 0, tmp.Error
	}
	return tmp.RowsAffected, nil
}
