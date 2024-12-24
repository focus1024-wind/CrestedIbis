package site

import "CrestedIbis/src/global/model"

type Site struct {
	model.IDModel
	Pid      *int64 `gorm:"comment:Pid" json:"pid"`
	Name     string `gorm:"comment:Name" json:"name"`
	Level    int8   `gorm:"comment:Level" json:"level"`
	Children []Site `gorm:"foreignkey:Pid" json:"children"`
	model.BaseModel
}

type SiteParentModel struct {
	model.IDModel
	Pid    *int64           `gorm:"comment:Pid" json:"pid"`
	Name   string           `gorm:"comment:Name" json:"name"`
	Level  int8             `gorm:"comment:Level" json:"level"`
	Parent *SiteParentModel `gorm:"foreignkey:Pid" json:"parent"`
	model.BaseModel
}

func (SiteParentModel) TableName() string {
	return "sites"
}
