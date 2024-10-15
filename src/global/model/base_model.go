package model

import (
	"gorm.io/gorm"
	"time"
)

type BaseModel struct {
	CreatedAt time.Time      `gorm:"column:created_time" json:"created_time"`
	UpdatedAt time.Time      `gorm:"column:updated_time" json:"updated_time"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
