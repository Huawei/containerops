package model

import "time"

type JobV1 struct {
	ID           string
	Endpoint     string
	Timeout      string
	Status       string
	Environments string
	CreateAt     time.Time
	UpdateAt     time.Time
}

type JobRecordV1 struct {
	ID            string
	JobID         string
	Number        int64
	Result        string
	Outputs       string
	Subscriptions string
	StartTime     time.Time
	EndTime       time.Time
}
