package user

import (
	"CrestedIbis/src/apps/audit_log"
	"CrestedIbis/src/global"
	"CrestedIbis/src/global/model"
	"CrestedIbis/src/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// AdminChangePassword 修改用户密码
//
//	@Summary		修改用户密码
//	@Version		0.0.1
//	@Description	修改用户密码
//	@Tags			超级用户操作
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string							false	"访问token"
//	@Param			access_token	query		string							false	"访问token"
//	@Param			SysUserLogin	body		SysUserLogin					true	"用户名、密码"
//	@Success		200				{object}	model.HttpResponse{data=nil}	"修改用户密码成功"
//	@Failure		500				{object}	model.HttpResponse{data=string}	"修改用户密码失败"
//	@Router			/system/admin/password [POST]
func AdminChangePassword(c *gin.Context) {
	var sysUserLogin SysUserLogin
	if err := c.ShouldBind(&sysUserLogin); err != nil {
		panic(http.StatusBadRequest)
	} else {
		err = SysUser{}.UpdatePassword(sysUserLogin)
		if err != nil {
			model.HttpResponse{}.FailGin(c, "修改用户密码失败")
		} else {
			model.HttpResponse{}.OkGin(c, nil)
		}
	}
}

// UpdateUser 更新用户
//
//	@Summary		更新用户
//	@Version		1.0.0
//	@Description	更新用户信息
//	@Tags			用户管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string							false	"访问token"
//	@Param			access_token	query		string							false	"访问token"
//	@Param			SysUser			body		SysUser							true	"新用户信息，必须携带ID信息"
//	@Success		200				{object}	model.HttpResponse{data=nil}	"更新用户成功"
//	@Failure		500				{object}	model.HttpResponse{data=string}	"更新用户失败"
//	@Router			/system/user [POST]
func UpdateUser(c *gin.Context) {
	var sysUser SysUser
	if err := c.ShouldBind(&sysUser); err != nil {
		panic(http.StatusBadRequest)
	} else {
		err = SysUser{}.Update(sysUser)
		if err != nil {
			model.HttpResponse{}.FailGin(c, err.Error())
		} else {
			model.HttpResponse{}.OkGin(c, nil)
		}
	}
}

// DeleteUser 删除用户
//
//	@Summary		删除用户
//	@Version		1.0.0
//	@Description	删除用户
//	@Tags			用户管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string							false	"访问token"
//	@Param			access_token	query		string							false	"访问token"
//	@Param			model.IDModel	body		model.IDModel					true	"用户ID信息"
//	@Success		200				{object}	model.HttpResponse{data=nil}	"删除成功"
//	@Failure		500				{object}	model.HttpResponse{data=string}	"删除失败"
//	@Router			/system/user [DELETE]
func DeleteUser(c *gin.Context) {
	var idModel model.IDModel
	if err := c.ShouldBind(&idModel); err != nil {
		panic(http.StatusBadRequest)
	} else {
		err = SysUser{}.Delete(idModel)
		if err != nil {
			model.HttpResponse{}.FailGin(c, err.Error())
		} else {
			model.HttpResponse{}.OkGin(c, nil)
		}
	}
}

// Login 用户登录
//
//	@Summary		用户登录
//	@Version		1.0.0
//	@Description	用户登录，返回用户信息和JWT信息
//	@Tags			用户管理
//	@Accept			json
//	@Produce		json
//	@Param			SysUserLogin	body		SysUserLogin									true	"用户登录信息，密码采用加盐加密"
//	@Success		200				{object}	model.HttpResponse{data=SysUserLoginResponse}	"登录成功，响应JWT"
//	@Failure		500				{object}	model.HttpResponse{data=string}					"登录失败，响应失败信息"
//	@Router			/system/user/login [POST]
func Login(c *gin.Context) {
	var sysUserLogin SysUserLogin
	if err := c.ShouldBind(&sysUserLogin); err != nil {
		// 参数错误
		panic(http.StatusBadRequest)
	} else {
		sysUser, err := SysUser{}.Login(sysUserLogin)

		if err != nil {
			// 登陆失败
			_ = audit_log.AuditLogLogin{}.Insert(c, sysUserLogin.Username, false, err.Error())
			model.HttpResponse{}.FailGin(c, err.Error())
		} else {
			// 登陆成功，生成Token
			roles := make([]string, len(sysUser.RoleGroups))
			for _, role := range sysUser.RoleGroups {
				roles = append(roles, role.RoleName)
			}

			token, err := utils.JwtToken{}.GenToken(sysUserLogin.Username, roles)
			if err != nil {
				_ = audit_log.AuditLogLogin{}.Insert(c, sysUserLogin.Username, false, "Token生成失败")
				model.HttpResponse{}.FailGin(c, "Token生成失败")
			} else {
				_ = audit_log.AuditLogLogin{}.Insert(c, sysUserLogin.Username, true, "登陆成功")

				model.HttpResponse{}.OkGin(c, &SysUserLoginResponse{
					SysUser:     sysUser,
					AccessToken: token,
				})
			}
		}
	}
}

// Register 注册用户
//
//	@Summary		注册用户，新增用户
//	@Version		1.0.0
//	@Description	注册用户，新增用户
//	@Tags			用户管理
//	@Accept			json
//	@Produce		json
//	@Param			SysUser	body		SysUser							true	"用户注册信息，密码采用加盐加密"
//	@Success		200		{object}	model.HttpResponse{data=nil}	"注册成功"
//	@Failure		500		{object}	model.HttpResponse{data=string}	"注册失败"
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

// GetUsers 根据查询条件搜索用户
//
//	@Summary		根据查询条件搜索用户
//	@Version		0.0.1
//	@Description	根据查询条件搜索用户
//	@Tags			用户管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string															false	"访问token"
//	@Param			access_token	query		string															false	"访问token"
//	@Param			page			query		integer															false	"分页查询页码，默认值: 1"
//	@Param			page_size		query		integer															false	"每页查询数量，默认值: 15"
//	@Param			keywords		query		string															false	"用户名、昵称、邮箱、手机号等模糊信息"
//	@Success		200				{object}	model.HttpResponse{data=model.BasePageResponse{data=[]SysUser}}	"查询成功"
//	@Failure		500				{object}	model.HttpResponse{data=string}									"查询数据失败"
//	@Router			/system/user/users [GET]
func GetUsers(c *gin.Context) {
	pageQuery := c.DefaultQuery("page", "1")
	pageSizeQuery := c.DefaultQuery("page_size", "15")
	keywords := c.DefaultQuery("keywords", "")

	page, err := strconv.ParseInt(pageQuery, 10, 0)
	if err != nil {
		panic(http.StatusBadRequest)
	}
	pageSize, err := strconv.ParseInt(pageSizeQuery, 10, 0)
	if err != nil {
		panic(http.StatusBadRequest)
	}

	total, data, err := SysUser{}.SelectUsers(page, pageSize, keywords)

	if err != nil {
		global.Logger.Errorf(err.Error())
		model.HttpResponse{}.FailGin(c, "搜索用户失败")
		return
	} else {
		model.HttpResponse{}.OkGin(c, &model.BasePageResponse{
			Total:    total,
			Data:     data,
			Page:     page,
			PageSize: pageSize,
		})
		return
	}
}

// DeleteUsers 批量删除用户
//
//	@Summary		批量删除用户
//	@Version		1.0.0
//	@Description	批量删除用户
//	@Tags			用户管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string							false	"访问token"
//	@Param			access_token	query		string							false	"访问token"
//	@Param			model.IDsModel	body		model.IDsModel					true	"ID列表"
//	@Success		200				{object}	model.HttpResponse{data=nil}	"批量删除用户成功"
//	@Failure		500				{object}	model.HttpResponse{data=string}	"批量删除用户失败"
//	@Router			/system/user/users [DELETE]
func DeleteUsers(c *gin.Context) {
	var idsModel model.IDsModel
	if err := c.ShouldBind(&idsModel); err != nil {
		panic(http.StatusBadRequest)
	} else {
		err = SysUser{}.Deletes(idsModel)
		if err != nil {
			model.HttpResponse{}.FailGin(c, err.Error())
		} else {
			model.HttpResponse{}.OkGin(c, nil)
		}
	}
}

// CreateRole 新增权限组
//
//	@Summary		新增权限组
//	@Version		0.0.1
//	@Description	新增权限组
//	@Tags			权限管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string							false	"访问token"
//	@Param			access_token	query		string							false	"访问token"
//	@Param			RoleGroup		body		RoleGroup						true	"用户权限信息"
//	@Success		200				{object}	model.HttpResponse{data=nil}	"新增权限组成功"
//	@Failure		500				{object}	model.HttpResponse{data=string}	"新增权限组失败"
//	@Router			/system/role [PUT]
func CreateRole(c *gin.Context) {
	var role RoleGroup
	if err := c.ShouldBind(&role); err != nil {
		panic(http.StatusBadRequest)
	} else {
		if role.RoleName == "admin" || role.RoleName == "guest" {
			model.HttpResponse{}.FailGin(c, "不允许创建admin、guest权限组")
			return
		}

		err = RoleGroup{}.Insert(role.RoleName)
		if err != nil {
			model.HttpResponse{}.FailGin(c, "新增权限组失败")
		} else {
			model.HttpResponse{}.OkGin(c, nil)
		}
	}
}

// UpdateRole 更新权限组
//
//	@Summary		更新权限组
//	@Version		0.0.1
//	@Description	更新权限组
//	@Tags			权限管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string							false	"访问token"
//	@Param			access_token	query		string							false	"访问token"
//	@Param			RoleGroup		body		RoleGroup						true	"用户权限信息"
//	@Success		200				{object}	model.HttpResponse{data=nil}	"更新权限组成功"
//	@Failure		500				{object}	model.HttpResponse{data=string}	"更新权限组列表失败"
//	@Router			/system/role [POST]
func UpdateRole(c *gin.Context) {
	var role RoleGroup
	if err := c.ShouldBind(&role); err != nil {
		panic(http.StatusBadRequest)
	} else {
		err = RoleGroup{}.Update(role)
		if err != nil {
			model.HttpResponse{}.FailGin(c, err.Error())
		} else {
			model.HttpResponse{}.OkGin(c, nil)
		}
	}
}

// DeleteRole 删除权限组
//
//	@Summary		删除权限组
//	@Version		0.0.1
//	@Description	删除权限组
//	@Tags			权限管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string							false	"访问token"
//	@Param			access_token	query		string							false	"访问token"
//	@Param			model.IDModel	body		model.IDModel					true	"用户权限ID"
//	@Success		200				{object}	model.HttpResponse{data=nil}	"删除权限组成功"
//	@Failure		500				{object}	model.HttpResponse{data=string}	"删除权限组失败"
//	@Router			/system/role [DELETE]
func DeleteRole(c *gin.Context) {
	var idModel model.IDModel
	if err := c.ShouldBind(&idModel); err != nil {
		panic(http.StatusBadRequest)
	} else {
		err = RoleGroup{}.Delete(idModel)
		if err != nil {
			model.HttpResponse{}.FailGin(c, err.Error())
		} else {
			model.HttpResponse{}.OkGin(c, nil)
		}
	}
}

// GetRoles 获取权限组列表
//
//	@Summary		获取权限组列表
//	@Version		0.0.1
//	@Description	获取权限组列表
//	@Tags			权限管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string									false	"访问token"
//	@Param			access_token	query		string									false	"访问token"
//	@Param			keywords		query		string									false	"模糊权限组名称信息"
//	@Success		200				{object}	model.HttpResponse{data=[]RoleGroup}	"获取权限组列表成功"
//	@Failure		500				{object}	model.HttpResponse{data=string}			"获取权限组列表失败"
//	@Router			/system/role/roles [GET]
func GetRoles(c *gin.Context) {
	keywords := c.DefaultQuery("keywords", "")
	roles, err := RoleGroup{}.Select(keywords)
	if err != nil {
		model.HttpResponse{}.FailGin(c, "修改用户密码失败")
	} else {
		model.HttpResponse{}.OkGin(c, roles)
	}
}

// DeleteRoles 批量删除权限组
//
//	@Summary		批量删除权限组
//	@Version		0.0.1
//	@Description	批量删除权限组
//	@Tags			权限管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string							false	"访问token"
//	@Param			access_token	query		string							false	"访问token"
//	@Param			model.IDsModel	body		model.IDsModel					true	"权限组ID列表"
//	@Success		200				{object}	model.HttpResponse{data=nil}	"批量删除权限组列表成功"
//	@Failure		500				{object}	model.HttpResponse{data=string}	"批量删除权限组列表失败"
//	@Router			/system/role/roles [DELETE]
func DeleteRoles(c *gin.Context) {
	var idsModel model.IDsModel
	if err := c.ShouldBind(&idsModel); err != nil {
		panic(http.StatusBadRequest)
	} else {
		err = RoleGroup{}.Deletes(idsModel)

		if err != nil {
			model.HttpResponse{}.FailGin(c, "删除权限组列表失败")
		} else {
			model.HttpResponse{}.OkGin(c, nil)
		}
	}
}

// GetRoleRules 获取权限组对应权限
//
//	@Summary		获取权限组对应权限
//	@Version		0.0.1
//	@Description	获取权限组对应权限
//	@Tags			权限管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string									false	"访问token"
//	@Param			access_token	query		string									false	"访问token"
//	@Param			name			query		string									true	"用户权限组名称"
//	@Success		200				{object}	model.HttpResponse{data=[]RoleGroup}	"获取权限组列表成功"
//	@Failure		500				{object}	model.HttpResponse{data=string}			"获取权限组列表失败"
//	@Router			/system/role/rules [GET]
func GetRoleRules(c *gin.Context) {
	name := c.DefaultQuery("name", "guest")
	if name == "" {
		panic(http.StatusBadRequest)
	} else {
		rules, err := RoleGroup{}.SelectRules(name)
		if err != nil {
			model.HttpResponse{}.FailGin(c, "获取用户组失败")
		} else {
			model.HttpResponse{}.OkGin(c, rules)
		}
	}
}

// UpdateRoleRules 更新权限组对应权限
//
//	@Summary		更新权限组对应权限
//	@Version		0.0.1
//	@Description	更新权限组对应权限
//	@Tags			权限管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization			header		string							false	"访问token"
//	@Param			access_token			query		string							false	"访问token"
//	@Param			RoleRuleUpdateEntity	body		RoleRuleUpdateEntity			true	"用户权限组ID和对应权限"
//	@Success		200						{object}	model.HttpResponse{data=nil}	"更新权限组对应权限成功"
//	@Failure		500						{object}	model.HttpResponse{data=string}	"更新权限组对应权限失败"
//	@Router			/system/role/rules [POST]
func UpdateRoleRules(c *gin.Context) {
	var roleRule RoleRuleUpdateEntity
	if err := c.ShouldBind(&roleRule); err != nil {
		panic(http.StatusBadRequest)
	} else {
		err = RoleGroup{}.UpdateRules(roleRule)
		if err != nil {
			model.HttpResponse{}.FailGin(c, err.Error())
		} else {
			model.HttpResponse{}.OkGin(c, nil)
		}
	}
}
