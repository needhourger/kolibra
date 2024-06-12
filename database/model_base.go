package database

import (
	"encoding/hex"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ModelBase struct {
	ID        string `gorm:"primary_key"`
	CreatedAt time.Time
	UpdateAt  time.Time
	DeleteAt  gorm.DeletedAt `gorm:"index"`
}

func generateShortUUID() (string, error) {
	u, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	b := make([]byte, 16)
	copy(b[:], u[:])

	shortUUID := hex.EncodeToString(b)
	return shortUUID, nil
}

func (model *ModelBase) BeforeCreate(tx *gorm.DB) (err error) {
	uuidHex, err := generateShortUUID()
	if err != nil {
		return err
	}
	model.ID = uuidHex
	return
}
