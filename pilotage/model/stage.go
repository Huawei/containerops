package model

import "time"

type StageV1 struct {
	ID         int64      `json:"id" gorm:"primary_key" gorm:"column:id"`
	FlowID     int64      `json:"flow_id" sql:"not null;type:bigint(20)" gorm:"column:flow_id"`
	StageType  string     `json:"stage_type" sql:"type:varchar(255)" gorm:"column:stage_type"`
	Name       string     `json:"name" sql:"not null;type:varchar(255)" gorm:"column:name"`
	Title      string     `json:"title" sql:"type:text" gorm:"column:title"`
	Sequencing string     `json:"sequencing" sql:"type:varchar(255)" gorm:"column:sequencing"`
	CreateAt   time.Time  `json:"create_at" sql:"" gorm:"column:create_at"`
	UpdateAt   time.Time  `json:"update_at" sql:"" gorm:"column:update_at"`
	DeletedAt  *time.Time `json:"delete_at" sql:"index" gorm:"column:delete_at"`
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

func (s *StageV1) Create(flowID int64, stageType, name, title, sequencing string) (stageID int64, err error) {
	s.FlowID, s.StageType, s.Name, s.Title, s.Sequencing = flowID, stageType, name, title, sequencing
	s.CreateAt = time.Now()

	tx := DB.Begin()
	if err := tx.Where("flow_id = ? AND name = ? ", flowID, name).FirstOrCreate(&s).Error; err != nil {
		tx.Rollback()
		return 0, err
	}
	tx.Commit()
	stageID = s.ID

	return stageID, nil
}

func (sd *StageDataV1) Create(stageID, number int64, result string, start, end time.Time) error {
	sd.StageID, sd.Number, sd.Result, sd.Start, sd.End = stageID, number, result, start, end

	tx := DB.Begin()
	if err := tx.Create(&sd).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}
