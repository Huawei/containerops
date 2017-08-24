package model

import "time"

type StageV1 struct {
	ID         string
	FlowID     string
	T          string
	Name       string
	Title      string
	Sequencing string
	CreateAt time.Time
	UpdateAt time.Time
}

type StageRecord struct {
	ID      string
	StageID string
	Number    int64
	Result    string
	Status  string
}
