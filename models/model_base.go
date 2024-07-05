package models

import (
	"encoding/hex"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ModelBase struct {
	ID        string `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdateAt  time.Time
	DeleteAt  gorm.DeletedAt `gorm:"index"`
}

func GenerateShortUUID() string {
	uuid := uuid.New()
	b := make([]byte, 16)
	copy(b[:], uuid[:])

	shortUUID := hex.EncodeToString(b)
	return shortUUID
}

func (model *ModelBase) BeforeCreate(tx *gorm.DB) (err error) {
	uuidHex := GenerateShortUUID()
	if model.ID == "" {
		model.ID = uuidHex
	}
	return
}
