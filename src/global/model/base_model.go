package model

import (
	"gorm.io/gorm"
)

type BaseModel struct {
	CreatedAt LocalTime      `gorm:"column:created_time" json:"created_time"`
	UpdatedAt LocalTime      `gorm:"column:updated_time" json:"updated_time"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
