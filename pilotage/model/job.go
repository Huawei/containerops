package model

import "time"

type JobV1 struct {
	ID            int64      `json:"id" gorm:"primary_key" gorm:"column:id"`
	ActionID      int64      `json:"action_id" sql:"not null;type:bigint(20)" gorm:"column:action_id"`
	Name          string     `json:"name" sql:"type:varchar(255)" gorm:"column:name"`
	JobType       string     `json:"type" sql:"type:varchar(255)" gorm:"column:type"`
	Endpoint      string     `json:"endpoint" sql:"type:varchar(255)" gorm:"column:endpoint"`
	Timeout       int64      `json:"timeout" sql:"type:bigint(20);default:0" gorm:"column:timeout"`
	Resources     string     `json:"resources" sql:"type:text" gorm:"column:resources"`
	Environments  string     `json:"environments" sql:"type:text" gorm:"column:environments"`
	Outputs       string     `json:"outputs" sql:"type:text" gorm:"column:outputs"`
	Subscriptions string     `json:"subscriptions" sql:"type:text" gorm:"column:subscriptions"`
	CreatedAt     time.Time  `json:"created_at" sql:"" gorm:"column:created_at"`
	UpdatedAt     time.Time  `json:"updatde_at" sql:"" gorm:"column:updated_at"`
	DeletedAt     *time.Time `json:"deleted_at" sql:"index" gorm:"column:deleted_at"`
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

func (j *JobV1) Put(actionID, timeout int64, name, jobType, endpoint, resources, environments, outputs, subscriptions string) (jobID int64, err error) {
	if DisableDB {
		return -1, nil
	}
	j.ActionID, j.Name, j.JobType, j.Endpoint, j.Timeout, j.Environments = actionID, name, jobType, endpoint, timeout, environments
	j.Resources, j.Outputs, j.Subscriptions = resources, outputs, subscriptions

	tx := DB.Begin()
	if tx.Where("action_id = ? AND name = ? ", actionID, name).First(&j).RecordNotFound() {
		j.CreatedAt = time.Now()
		if err := tx.Create(&j).Error; err != nil {
			tx.Rollback()
			return 0, err
		}
	} else {
		if err := tx.Model(&j).Updates(JobV1{JobType: jobType, Endpoint: endpoint, Timeout: timeout, Resources: resources,
			Environments: environments, Outputs: outputs, Subscriptions: subscriptions}).Error; err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	tx.Commit()
	jobID = j.ID

	return jobID, nil
}

func (jd *JobDataV1) Put(jobID, number int64, result string, start, end time.Time) error {
	if DisableDB {
		return nil
	}

	jd.JobID, jd.Number, jd.Result, jd.Start, jd.End = jobID, number, result, start, end

	tx := DB.Begin()
	if err := tx.Create(&jd).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}

func (jd *JobDataV1) GetNumbers(jobID int64) (int64, error) {
	if DisableDB {
		return -1, nil
	}

	tmp := DB.Where("job_id = ?", jobID).Find(&JobDataV1{})
	if tmp.RecordNotFound() {
		return 0, nil
	}
	if tmp.Error != nil {
		return 0, tmp.Error
	}
	return tmp.RowsAffected, nil
}
