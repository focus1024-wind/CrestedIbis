package user

import (
	"CrestedIbis/src/global/model"
	"CrestedIbis/src/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Login 用户登录
//
//	@Title			用户登录
//	@Version		1.0.0
//	@Description	用户登录并生成用户登录日志信息
//	@Tags			用户管理 /system/user
//	@Accept			json
//	@Produce		json
//	@Param			SysLoginUser	body		SysUserLogin				true	"用户登录信息，密码采用加盐加密"
//	@Success		200				{object}	model.HttpResponse{data=string}	"登录成功，响应JWT"
//	@Failure		500				{object}	model.HttpResponse{data=string}	"登录失败，响应失败信息"
//	@Router			/system/user/login [POST]
func Login(c *gin.Context) {
	var sysUserLogin SysUserLogin
	if err := c.ShouldBind(&sysUserLogin); err != nil {
		// 参数错误
		panic(http.StatusBadRequest)
	} else {
		err = SysUser{}.Login(sysUserLogin)
		if err != nil {
			// 登陆失败
			model.HttpResponse{}.FailGin(c, err.Error())
		} else {
			// 登陆成功，生成Token
			token, err := utils.JwtToken{}.GenToken(sysUserLogin.Username)
			if err != nil {
				model.HttpResponse{}.FailGin(c, "Token生成失败")
			} else {
				model.HttpResponse{}.OkGin(c, token)
			}
		}
	}
}

// Register 注册用户
//
//	@Title			注册用户
//	@Version		1.0.0
//	@Description	注册用户
//	@Tags			用户管理 /system/user
//	@Accept			json
//	@Produce		json
//	@Param			SysUser	body		SysUser					true	"用户注册信息，密码采用加盐加密"
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
