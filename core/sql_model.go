package core

import (
	"time"
)

type SQLModel struct {
	ID        int64      `json:"-" gorm:"column:id;" db:"id"`
	FakeId    *UID       `json:"id" gorm:"-"`
	CreatedAt *time.Time `json:"created_at,omitempty" gorm:"column:created_at;"  db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" gorm:"column:updated_at;"  db:"updated_at"`
}

func NewSQLModel() SQLModel {
	now := time.Now().UTC()

	return SQLModel{
		ID:        0,
		CreatedAt: &now,
		UpdatedAt: &now,
	}
}

func (m *SQLModel) FullFill() {
	t := time.Now()

	if m.UpdatedAt == nil {
		m.UpdatedAt = &t
	}
}

func NewUpsertSQLModel(id int64) *SQLModel {
	t := time.Now()

	return &SQLModel{
		ID:        id,
		CreatedAt: &t,
		UpdatedAt: &t,
	}
}

func NewUpsertWithoutIdSQLModel() *SQLModel {
	t := time.Now()

	return &SQLModel{
		CreatedAt: &t,
		UpdatedAt: &t,
	}
}

func (sqlModel *SQLModel) Mask(objectId int) {
	uid := NewUID(uint32(sqlModel.ID), objectId, 1)
	sqlModel.FakeId = &uid
}

func (sqlModel *SQLModel) MaskField(id int64, objectId int) {
	uid := NewUID(uint32(id), objectId, 1)
	sqlModel.FakeId = &uid
}