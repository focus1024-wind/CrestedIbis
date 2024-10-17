package audit_log

import (
	"CrestedIbis/src/global/model"
)

type AuditLogLogin struct {
	AuditLogLoginId int64  `gorm:"primary_key;AUTO_INCREMENT" json:"audit_log_login_id"`
	Username        string `gorm:"type:varchar(64);comment:用户名" json:"username"`
	IpAddr          string `gorm:"type:varchar(128);comment:IP地址" json:"ip_addr"`
	IpLocation      string `gorm:"type:varchar(128);comment:IP归属地" json:"ip_location"`
	Os              string `gorm:"type:varchar(64);comment:系统" json:"os"`
	Remark          string `gorm:"type:varchar(256);comment:备注" json:"remark"`
	Status          bool   `gorm:"type:tinyint(1);comment:登录状态" json:"status"`
	Msg             string `gorm:"type:varchar(256);comment:登录信息" json:"msg"`
	model.BaseModel
}
