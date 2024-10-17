package audit_log

import (
	"CrestedIbis/src/global"
	"CrestedIbis/src/utils"
	"github.com/gin-gonic/gin"
)

func (AuditLogLogin) Insert(c *gin.Context, username string, status bool, msg string) (err error) {
	login := AuditLogLogin{
		Username:   username,
		IpAddr:     c.ClientIP(),
		IpLocation: utils.GetRealLocationByIpAddr(c.ClientIP()),
		Os:         utils.GetOsByUserAgent(c.Request.UserAgent()),
		Remark:     c.Request.UserAgent(),
		Status:     status,
		Msg:        msg,
	}
	return global.Db.Create(&login).Error
}
