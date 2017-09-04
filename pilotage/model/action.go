package model

import "time"

type ActionV1 struct {
	ID        int64      `json:"id" gorm:"primary_key" gorm:"column:id"`
	StageID   int64      `json:"stage_id" sql:"not null;type:bigint(20)" gorm:"column:stage_id"`
	Name      string     `json:"name" sql:"type:varchar(255)" gorm:"column:name"`
	Title     string     `json:"title" sql:"type:text" gorm:"column:title"`
	CreatedAt time.Time  `json:"created_at" sql:"" gorm:"column:created_at"`
	UpdatedAt time.Time  `json:"updated_at" sql:"" gorm:"column:updated_at"`
	DeletedAt *time.Time `json:"deleted_at" sql:"index" gorm:"column:deleted_at"`
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

func (a *ActionV1) Put(stageID int64, name, title string) (actionID int64, err error) {
	a.StageID, a.Name, a.Title = stageID, name, title

	tx := DB.Begin()
	if tx.Where("stage_id = ? AND name = ? ", stageID, name).First(&a).RecordNotFound() {
		a.CreatedAt = time.Now()
		if err := tx.Create(&a).Error; err != nil {
			tx.Rollback()
			return 0, err
		}
	} else {
		if err := tx.Model(&a).Updates(ActionV1{Title: title}).Error; err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	tx.Commit()
	actionID = a.ID

	return actionID, nil
}

func (ad *ActionDataV1) Put(actionID, number int64, result string, start, end time.Time) error {
	ad.ActionID, ad.Number, ad.Result, ad.Start, ad.End = actionID, number, result, start, end

	tx := DB.Begin()
	if err := tx.Create(&ad).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}

func (ad *ActionDataV1) GetNumbers(actionID int64) (int64, error) {
	tmp := DB.Where("action_id = ?", actionID).Find(&ActionDataV1{})
	if tmp.RecordNotFound() {
		return 0, nil
	}
	if tmp.Error != nil {
		return 0, tmp.Error
	}
	return tmp.RowsAffected, nil
}
