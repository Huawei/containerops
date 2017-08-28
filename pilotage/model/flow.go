package model

import (
	"github.com/Huawei/containerops/pilotage/module"
	"time"
)

type FlowV1 struct {
	ID         int64     `json:"id" gorm:"primary_key" gorm:"column:id"`
	Namespace  string    `json:"namespace" sql:"not null;type:varchar(255)" gorm:"column:namespace"`
	Repository string    `json:"repository" sql:"not null;type:varchar(255)" gorm:"column:tag"`
	Tag        string    `json:"tag" sql:"not null;type:varchar(255)" gorm:"column:tag"`
	Version    int64     `json:"version" sql:"type:bigint(20)" gorm:"column:version"`
	Title      string    `json:"title" sql:"type:text" gorm:"column:title"`
	Timeout    int64     `json:"timeout" sql:"default:0" gorm:"column:timeout"`
	Content    string    `json:"content" sql:"type:text" gorm:"column:content"`
	CreateAt   time.Time `json:"create_at" sql:"" gorm:"column:create_at"`
	UpdateAt   time.Time `json:"update_at" sql:"" gorm:"column:update_at"`
}

type FlowDataV1 struct {
	ID     int64     `json:"id" gorm:"primary_key" gorm:"column:id"`
	FlowID int64     `json:"flow_id" sql:"not null;type:bigint(20)" gorm:"column:flow_id"`
	Number int64     `json:"number" sql:"not null;type:bigint(20)" gorm:"column:number"`
	Result string    `json:"result" sql:"type:varchar(255)" gorm:"column:result"`
	Start  time.Time `json:"start" sql:"" gorm:"column:start"`
	End    time.Time `json:"end" sql:"" gorm:"column:end"`
}


func (f *FlowV1) TableName() string {
	return "flow_v1"
}

func (fd *FlowDataV1) TableName() string {
	return "flow_data_v1"
}