package model

import (
	"gorm.io/gorm"
)

// BaseHardDeleteModel gorm 硬删除模型
type BaseHardDeleteModel struct {
	CreatedAt LocalTime `gorm:"column:created_time" json:"created_time"`
	UpdatedAt LocalTime `gorm:"column:updated_time" json:"updated_time"`
}

// BaseModel gorm 软删除模型
type BaseModel struct {
	CreatedAt LocalTime      `gorm:"column:created_time" json:"created_time"`
	UpdatedAt LocalTime      `gorm:"column:updated_time" json:"updated_time"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// IDModel 主键ID模型
type IDModel struct {
	ID int64 `gorm:"primary_key;auto_increment;comment:主键ID" json:"id"`
}

// IDsModel ID列表
type IDsModel struct {
	IDs []int64 `json:"ids"`
}

type BasePageResponse struct {
	Total    int64       `json:"total" desc:"总数量"`
	Data     interface{} `json:"data" desc:"数据列表"`
	Page     int64       `json:"page" desc:"页码"`
	PageSize int64       `json:"page_size" desc:"每页查询数量"`
}
