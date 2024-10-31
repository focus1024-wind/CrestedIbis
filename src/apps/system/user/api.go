package user

import (
	"CrestedIbis/src/apps/audit_log"
	"CrestedIbis/src/global/model"
	"CrestedIbis/src/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Login 用户登录
//
//	@Summary		用户登录
//	@Version		1.0.0
//	@Description	用户登录并生成用户登录日志信息
//	@Tags			用户管理 /system/user
//	@Accept			json
//	@Produce		json
//	@Param			SysLoginUser	body		SysUserLogin					true	"用户登录信息，密码采用加盐加密"
//	@Success		200				{object}	model.HttpResponse{data=string}	"登录成功，响应JWT"
//	@Failure		500				{object}	model.HttpResponse{data=string}	"登录失败，响应失败信息"
//	@Router			/system/user/login [POST]
func Login(c *gin.Context) {
	var sysUserLogin SysUserLogin
	if err := c.ShouldBind(&sysUserLogin); err != nil {
		// 参数错误
		panic(http.StatusBadRequest)
	} else {
		roles, err := SysUser{}.Login(sysUserLogin)

		if err != nil {
			// 登陆失败
			_ = audit_log.AuditLogLogin{}.Insert(c, sysUserLogin.Username, false, err.Error())
			model.HttpResponse{}.FailGin(c, err.Error())
		} else {
			// 登陆成功，生成Token
			token, err := utils.JwtToken{}.GenToken(sysUserLogin.Username, roles)
			if err != nil {
				_ = audit_log.AuditLogLogin{}.Insert(c, sysUserLogin.Username, false, "Token生成失败")
				model.HttpResponse{}.FailGin(c, "Token生成失败")
			} else {
				_ = audit_log.AuditLogLogin{}.Insert(c, sysUserLogin.Username, true, "登陆成功")
				model.HttpResponse{}.OkGin(c, token)
			}
		}
	}
}

// Register 注册用户
//
//	@Summary		注册用户
//	@Version		1.0.0
//	@Description	注册用户
//	@Tags			用户管理 /system/user
//	@Accept			json
//	@Produce		json
//	@Param			SysUser	body		SysUser							true	"用户注册信息，密码采用加盐加密"
//	@Success		200		{object}	model.HttpResponse{}			"注册成功"
//	@Failure		500		{object}	model.HttpResponse{data=string}	"注册失败，响应失败信息"
//	@Router			/system/user/register [POST]
func Register(c *gin.Context) {
	var sysUser SysUser
	if err := c.ShouldBind(&sysUser); err != nil {
		panic(http.StatusBadRequest)
	} else {
		err = SysUser{}.Insert(sysUser)
		if err != nil {
			model.HttpResponse{}.FailGin(c, err.Error())
		} else {
			model.HttpResponse{}.OkGin(c, nil)
		}
	}
}

// AdminGetAllUserByPages 搜索用户
//
//	@Summary		搜索用户
//	@Version		0.0.1
//	@Description	搜索用户
//	@Tags			用户管理 /system/user
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string								false	"访问token"
//	@Param			access_token	query		string								false	"访问token"
//	@Success		200				{object}	model.HttpResponse{data=[]SysUser}	"查询成功"
//	@Failure		500				{object}	model.HttpResponse{data=string}		"查询数据失败"
//	@Router			/system/user/users [GET]
func AdminGetAllUserByPages(c *gin.Context) {
	if claims, exists := c.Get("claims"); exists {
		claims := claims.(*utils.JwtToken)
		isAdmin := false
		for _, role := range claims.Roles {
			if role == "admin" {
				isAdmin = true
				break
			}
		}
		if isAdmin {
			users, err := selectUsers()
			if err != nil {
				model.HttpResponse{}.FailGin(c, "搜索用户失败")
			}
			model.HttpResponse{}.OkGin(c, users)
			return
		}
	}
	model.HttpResponse{}.FailGin(c, "无权限")
}
