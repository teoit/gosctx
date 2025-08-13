package core

import (
	"time"
)

type SQLModelCreateAt struct {
	ID        int64      `json:"-" gorm:"column:id;" db:"id"`
	FakeId    *UID       `json:"id" gorm:"-"`
	CreatedAt *time.Time `json:"created_at,omitempty" gorm:"column:created_at;"  db:"created_at"`
}

func NewSQLModelCreateAt() SQLModel {
	now := time.Now().UTC()

	return SQLModel{
		ID:        0,
		CreatedAt: &now,
	}
}

func (sqlModel *SQLModelCreateAt) Mask(objectId int) {
	uid := NewUID(uint32(sqlModel.ID), objectId, 1)
	sqlModel.FakeId = &uid
}
