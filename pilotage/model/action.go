package model

import "time"

type ActionV1 struct {
	ID       int64     `json:"id" gorm:"primary_key" gorm:"column:id"`
	StageID  int64     `json:"stage_id" sql:"not null;type:bigint(20)" gorm:"column:stage_id"`
	Name     string    `json:"name" sql:"type:varchar(255)" gorm:"column:name"`
	Title    string    `json:"title" sql:"type:text" gorm:"column:title"`
	CreateAt time.Time `json:"create_at" sql:"" gorm:"column:create_at"`
	UpdateAt time.Time `json:"update_at" sql:"" gorm:"column:update_at"`
	DeletedAt   *time.Time `json:"delete_at" sql:"index" gorm:"column:delete_at"`
}

type ActionDataV1 struct {
	ID       int64     `json:"id" gorm:"primary_key" gorm:"column:id"`
	ActionID int64     `json:"action_id" sql:"not null;type:bigint(20)" gorm:"column:action_id"`
	Number   int64     `json:"number" sql:"not null;type:bigint(20)" gorm:"column:number"`
	Result   string    `json:"result" sql:"type:varchar(255)" gorm:"column:result"`
	Start    time.Time `json:"start" sql:"" gorm:"column:start"`
	End      time.Time `json:"end" sql:"" gorm:"column:end"`
}

func (a *ActionV1) TableName() string {
	return "action_v1"
}

func (ad *ActionDataV1) TableName() string {
	return "action_data_v1"
}

func (a *ActionV1) Create(stageID int64, name, title string) (actionID int64, err error) {
	a.StageID, a.Name, a.Title = stageID, name, title
	a.CreateAt = time.Now()

	tx := DB.Begin()
	if err := tx.Where("stage_id = ? AND name = ? ", stageID, name).FirstOrCreate(&a).Error; err != nil {
		tx.Rollback()
		return 0, err
	}
	tx.Commit()
	actionID = a.ID

	return actionID, nil
}

func (ad *ActionDataV1) Create(actionID, number int64, result string, start, end time.Time) error {
	ad.ActionID, ad.Number, ad.Result, ad.Start, ad.End = actionID, number, result, start, end

	tx := DB.Begin()
	if err := tx.Create(&ad).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}
