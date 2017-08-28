package model

import "time"

type JobV1 struct {
	ID           int64     `json:"id" gorm:"primary_key" gorm:"column:id"`
	ActionID     int64     `json:"action_id" sql:"not null;type:bigint(20)" gorm:"column:action_id"`
	Name         string    `json:"name" sql:"type:varchar(255)" gorm:"column:name"`
	Endpoint     string    `json:"endpoint" sql:"type:varchar(255)" gorm:"column:endpoint"`
	Timeout      string    `json:"timeout" sql:"type:varchar(255)" gorm:"column:timeout"`
	Environments string    `json:"environments" sql:"type:text" gorm:"column:environments"`
	CreateAt     time.Time `json:"create_at" sql:"" gorm:"column:create_at"`
	UpdateAt     time.Time `json:"update_at" sql:"" gorm:"column:update_at"`
}

type JobDataV1 struct {
	ID            int64     `json:"id" gorm:"primary_key" gorm:"column:id"`
	JobID         int64     `json:"job_id" sql:"not null;type:bigint(20)" gorm:"column:job_id"`
	Number        int64     `json:"number" sql:"not null;type:bigint(20)" gorm:"column:number"`
	Result        string    `json:"result" sql:"type:varchar(255)" gorm:"column:result"`
	Outputs       string    `json:"outputs" sql:"type:text" gorm:"column:outputs"`
	Subscriptions string    `json:"subscriptions" sql:"type:text" gorm:"column:subscriptions"`
	Start         time.Time `json:"start" sql:"" gorm:"column:start"`
	End           time.Time `json:"end" sql:"" gorm:"column:end"`
}

func (j *JobV1) TableName() string {
	return "job_v1"
}

func (jd *JobDataV1) TableName() string {
	return "job_data_v1"
}
