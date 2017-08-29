package model

import "time"

type JobV1 struct {
	ID            int64     `json:"id" gorm:"primary_key" gorm:"column:id"`
	ActionID      int64     `json:"action_id" sql:"not null;type:bigint(20)" gorm:"column:action_id"`
	Name          string    `json:"name" sql:"type:varchar(255)" gorm:"column:name"`
	Endpoint      string    `json:"endpoint" sql:"type:varchar(255)" gorm:"column:endpoint"`
	Timeout       string    `json:"timeout" sql:"type:varchar(255)" gorm:"column:timeout"`
	Environments  string    `json:"environments" sql:"type:text" gorm:"column:environments"`
	Outputs       string    `json:"outputs" sql:"type:text" gorm:"column:outputs"`
	Subscriptions string    `json:"subscriptions" sql:"type:text" gorm:"column:subscriptions"`
	CreateAt      time.Time `json:"create_at" sql:"" gorm:"column:create_at"`
	UpdateAt      time.Time `json:"update_at" sql:"" gorm:"column:update_at"`
	DeletedAt   *time.Time `json:"delete_at" sql:"index" gorm:"column:delete_at"`
}

type JobDataV1 struct {
	ID     int64     `json:"id" gorm:"primary_key" gorm:"column:id"`
	JobID  int64     `json:"job_id" sql:"not null;type:bigint(20)" gorm:"column:job_id"`
	Number int64     `json:"number" sql:"not null;type:bigint(20)" gorm:"column:number"`
	Result string    `json:"result" sql:"type:varchar(255)" gorm:"column:result"`
	Start  time.Time `json:"start" sql:"" gorm:"column:start"`
	End    time.Time `json:"end" sql:"" gorm:"column:end"`
}

func (j *JobV1) TableName() string {
	return "job_v1"
}

func (jd *JobDataV1) TableName() string {
	return "job_data_v1"
}

func (j *JobV1) Create(actionID int64, name, endpoint, timeout, environments, outputs, subscriptions string) (jobID int64, err error) {
	j.ActionID, j.Name, j.Endpoint, j.Timeout, j.Environments = actionID, name, endpoint, timeout, environments
	j.Outputs, j.Subscriptions = outputs, subscriptions

	tx := DB.Begin()
	if err := tx.Where("action_id = ? AND name = ? ", actionID, name).FirstOrCreate(&j).Error; err != nil {
		tx.Rollback()
		return 0, err
	}
	tx.Commit()
	jobID = j.ID

	return jobID, nil
}

func (jd *JobDataV1) Create(jobID, number int64, result string, start, end time.Time) error {
	jd.JobID, jd.Number, jd.Result, jd.Start, jd.End = jobID, number, result, start, end

	tx := DB.Begin()
	if err := tx.Create(&jd).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}
