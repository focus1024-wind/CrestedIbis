package user

import (
	"CrestedIbis/src/global/model"
	"CrestedIbis/src/utils"
	"encoding/json"
)

/**
用户相关Entity
*/

type SysUserLogin struct {
	Username string `gorm:"uniqueIndex;type:varchar(64);comment:用户名" json:"username" binding:"required" example:"admin"`
	Password string `gorm:"type:varchar(128);comment:用户密码" json:"password"  example:"CrestedIbis"`
}

type SysUserFields struct {
	Nickname   string      `gorm:"type:varchar(64);comment:用户昵称" json:"nickname"`
	Phone      string      `gorm:"type:varchar(11);comment:用户手机号" json:"phone"`
	Email      string      `gorm:"type:varchar(32);comment:用户邮箱" json:"email"`
	Avatar     string      `gorm:"type:varchar(255);comment:用户头像路径" json:"avatar"`
	Sex        uint8       `gorm:"type:tinyint(1);comment:用户性别(1: 男性; 2: 女性; other: 未知);default:0" json:"sex"`
	RoleGroups []RoleGroup `gorm:"many2many:user_role_groups;" json:"role_groups"`
}

type SysUser struct {
	model.IDModel
	SysUserLogin
	SysUserFields
	model.BaseModel
}

// MarshalJSON 通过 MarshalJSON 序列化用户，避免隐私数据暴露
func (sysUser *SysUser) MarshalJSON() ([]byte, error) {
	sysUser.Password = ""
	if sysUser.Nickname == "" {
		sysUser.Nickname = sysUser.Username
	}
	return json.Marshal(*sysUser)
}

type SysUserLoginResponse struct {
	SysUser     `json:"user"`
	AccessToken string `json:"access_token"`
}

func (sysUserLoginResponse *SysUserLoginResponse) MarshalJSON() ([]byte, error) {
	sysUserLoginResponse.Password = ""
	return json.Marshal(*sysUserLoginResponse)
}

/**
权限组相关Entity
*/

type RoleGroup struct {
	model.IDModel
	RoleName string    `json:"role_name"`
	User     []SysUser `gorm:"many2many:user_role_groups;" json:"user"`
	model.BaseModel
}

type RoleRuleUpdateEntity struct {
	ID    int64              `json:"id"`
	Rules []utils.CasbinRule `json:"rules"`
}
