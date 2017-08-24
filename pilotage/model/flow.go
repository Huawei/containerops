package model

import "time"

type FlowV1 struct {
	ID         string
	Namespace  string
	Repository string
	Tag        string
	Version    int64
	Number     int64
	Title      string
	Timeout    int64
	Content    string
	CreateAt   time.Time
	UpdateAt   time.Time
}


type FlowRecordV1 struct {
	ID        string
	FlowID    string
	Number    int64
	Result    string
	StartTime time.Time
	EndTime   time.Time
}
