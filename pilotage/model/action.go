package model

import "time"

type ActionV1 struct {
	ID       string
	StageID  string
	Name     string
	Title    string
	Status   string
	CreateAt time.Time
	UpdateAt time.Time
}

type ActionRecordV1 struct {
	ID        string
	ActionID  string
	Number    int64
	Result    string
	StartTime time.Time
	EndTime   time.Time
}
