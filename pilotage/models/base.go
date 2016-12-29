package models

import "time"

type BaseIDField struct {
	ID uint64 `json:"id" gorm:"primary_key"`
}

type BaseModel struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type BaseLogModel struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}
