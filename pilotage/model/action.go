package model

import "time"

type ActionV1 struct {
	ID       int64     `json:"id" gorm:"primary_key" gorm:"column:id"`
	StageID  int64     `json:"flow_id" sql:"not null;type:bigint(20)" gorm:"column:flow_id"`
	Name     string    `json:"name" sql:"type:varchar(255)" gorm:"column:name"`
	Title    string    `json:"title" sql:"type:text" gorm:"column:title"`
	Status   string    `json:"status" sql:"type:varchar(255)" gorm:"column:status"`
	CreateAt time.Time `json:"create_at" sql:"" gorm:"column:create_at"`
	UpdateAt time.Time `json:"update_at" sql:"" gorm:"column:update_at"`
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
