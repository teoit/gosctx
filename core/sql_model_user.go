package core

import "time"

type SQLModelUser struct {
	ID        int64      `gorm:"type:int(11);primaryKey;autoIncrement" json:"id"`
	CreatedBy *int       `json:"created_by,omitempty" gorm:"column:created_by;"  db:"created_by"`
	UpdatedBy *int       `json:"updated_by,omitempty" gorm:"column:updated_by;"  db:"updated_by"`
	CreatedAt *time.Time `json:"created_at,omitempty" gorm:"column:created_at;"  db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" gorm:"column:updated_at;"  db:"updated_at"`
}

func NewSQLModelUser() SQLModelUser {
	now := time.Now().UTC()

	return SQLModelUser{
		ID:        0,
		CreatedAt: &now,
		UpdatedAt: &now,		
	}
}
