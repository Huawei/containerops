package model

import "time"

type JobV1 struct {
	ID           int64
	ActionID     int64
	Endpoint     string
	Timeout      string
	Status       string
	Environments string
	CreateAt     time.Time
	UpdateAt     time.Time
}

type JobDataV1 struct {
	ID            int64
	JobID         int64
	Number        int64
	Result        string
	Outputs       string
	Subscriptions string
	Start         time.Time
	End           time.Time
}
