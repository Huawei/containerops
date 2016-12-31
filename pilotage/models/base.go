package models

import "time"

type BaseIDField struct {
	ID uint64 `json:"id;omitempty" gorm:"primary_key;AUTO_INCREMENT"`
}

type BaseModel struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type BaseLogModel struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
