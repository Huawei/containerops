package model

import "time"

type StageV1 struct {
	ID         int64
	FlowID     int64
	StageType  string
	Name       string
	Title      string
	Sequencing string
	CreateAt   time.Time
	UpdateAt   time.Time
}

type StageDataV1 struct {
	ID      int64
	StageID int64
	Number  int64
	Result  string
	Status  string
	Start   time.Time
	End     time.Time
}
