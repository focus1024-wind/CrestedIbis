package user

import (
	"CrestedIbis/src/global/model"
)

type SysUserLogin struct {
	Username string `gorm:"uniqueIndex;type:varchar(64);comment:用户名" json:"username" binding:"required" example:"admin"`
	Password string `gorm:"type:varchar(128);comment:用户密码" json:"password"  binding:"required" example:"CrestedIbis"`
}

type SysUserId struct {
	UserId int64 `gorm:"primary_key;auto_increment;comment:用户ID" json:"user_id"`
}

type SysUserFields struct {
	Nickname   string      `gorm:"type:varchar(64);comment:用户昵称" json:"nickname"`
	Phone      string      `gorm:"type:varchar(11);comment:用户手机号" json:"phone"`
	Email      string      `gorm:"type:varchar(32);comment:用户邮箱" json:"email"`
	Avatar     string      `gorm:"type:varchar(255);comment:用户头像路径" json:"avatar"`
	Sex        uint8       `gorm:"type:tinyint(1);comment:用户性别(1: 男性; 0: 女性; other: 未知);default:9" json:"sex"`
	RoleGroups []RoleGroup `gorm:"many2many:user_role_groups;"`
	model.BaseModel
}

type SysUser struct {
	SysUserId
	SysUserLogin
	SysUserFields
}

type RoleGroup struct {
	RoleId   int64     `gorm:"primary_key;AUTO_INCREMENT;comment:权限ID" json:"role_id"`
	RoleName string    `json:"role_name"`
	User     []SysUser `gorm:"many2many:user_role_groups;" json:"user"`
	model.BaseModel
}
