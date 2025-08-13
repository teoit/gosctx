package core

import (
	"time"
)

type SQLModelIDDate struct {
	ID        int64      `json:"id" gorm:"column:id;" db:"id"`
	CreatedAt *time.Time `json:"created_at,omitempty" gorm:"column:created_at;"  db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" gorm:"column:updated_at;"  db:"updated_at"`
}

func NewSQLModelIDDate() SQLModelIDDate {
	now := time.Now().UTC()

	return SQLModelIDDate{
		ID:        CreateIdDate(),
		CreatedAt: &now,
		UpdatedAt: &now,
	}
}
