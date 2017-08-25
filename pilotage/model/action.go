package model

import "time"

type ActionV1 struct {
	ID       int64
	StageID  int64
	Name     string
	Title    string
	Status   string
	CreateAt time.Time
	UpdateAt time.Time
}

type ActionDataV1 struct {
	ID       int64
	ActionID int64
	Number   int64
	Result   string
	Start    time.Time
	End      time.Time
}
