package site

import "CrestedIbis/src/global/model"

type Site struct {
	Id       int64  `gorm:"primary_key;auto_increment;comment:ID" json:"id"`
	Pid      *int64 `gorm:"comment:Pid" json:"pid"`
	Name     string `gorm:"comment:Name" json:"name"`
	Level    int8   `gorm:"comment:Level" json:"level"`
	Children []Site `gorm:"foreignkey:Pid" json:"children"`
	model.BaseModel
}

type SiteIdQuery struct {
	Id int64 `json:"id"`
}

type PostSiteQuery struct {
	SiteIdQuery
	Name string `json:"name"`
}
