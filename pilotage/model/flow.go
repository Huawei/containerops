package model

import (
	"github.com/Huawei/containerops/pilotage/module"
	"time"
)

type FlowV1 struct {
	ID         int64
	Namespace  string
	Repository string
	Tag        string
	Version    int64
	Title      string
	Timeout    int64
	Content    string
	CreateAt   time.Time
	UpdateAt   time.Time
}

type FlowDataV1 struct {
	ID     int64
	FlowID int64
	Number int64
	Result string
	Start  time.Time
	End    time.Time
}

func (f *FlowV1) Save(flow *module.Flow) error {

	return nil
}
