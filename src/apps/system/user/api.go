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

// GetAllUserByPages 搜索用户
//
//	@Summary		搜索用户
//	@Version		0.0.1
//	@Description	搜索用户
//	@Tags			超级用户操作 /system/admin
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string									false	"访问token"
//	@Param			access_token	query		string									false	"访问token"
//	@Param			page			query		integer									false	"分页查询页码，默认值: 1"
//	@Param			page_size		query		integer									false	"每页查询数量，默认值: 15"
//	@Success		200				{object}	model.HttpResponse{data=SysUserPage}	"查询成功"
//	@Failure		500				{object}	model.HttpResponse{data=string}			"查询数据失败"
//	@Router			/system/user/users [GET]
func GetAllUserByPages(c *gin.Context) {
	pageQuery := c.DefaultQuery("page", "1")
	pageSizeQuery := c.DefaultQuery("page_size", "15")
	page, err := strconv.ParseInt(pageQuery, 10, 0)
	if err != nil {
		panic(http.StatusBadRequest)
	}
	pageSize, err := strconv.ParseInt(pageSizeQuery, 10, 0)
	if err != nil {
		panic(http.StatusBadRequest)
	}

	total, data, err := selectUsersByPages(page, pageSize)
	if err != nil {
		global.Logger.Errorf(err.Error())
		model.HttpResponse{}.FailGin(c, "搜索用户失败")
		return
	} else {
		model.HttpResponse{}.OkGin(c, &SysUserPage{
			Total:    total,
			Data:     data,
			Page:     page,
			PageSize: pageSize,
		})
		return
	}
}

// AdminChangePassword 修改用户密码
//
//	@Summary		修改用户密码
//	@Version		0.0.1
//	@Description	搜索用户
//	@Tags			超级用户操作 /system/admin
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string							false	"访问token"
//	@Param			access_token	query		string							false	"访问token"
//	@Param			SysUserLogin	body		SysUserLogin					true	"用户名、密码"
//	@Success		200				{object}	model.HttpResponse{}			"修改用户密码成功"
//	@Failure		500				{object}	model.HttpResponse{data=string}	"修改用户密码失败"
//	@Router			/system/admin/password [POST]
func AdminChangePassword(c *gin.Context) {
	var sysUserLogin SysUserLogin
	if err := c.ShouldBind(&sysUserLogin); err != nil {
		panic(http.StatusBadRequest)
	} else {
		err = updateUserPassword(sysUserLogin.Username, sysUserLogin.Password)
		if err != nil {
			model.HttpResponse{}.FailGin(c, "修改用户密码失败")
		} else {
			model.HttpResponse{}.OkGin(c, nil)
		}
	}
}

// AdminDeleteUser 删除用户
//
//	@Summary		删除用户
//	@Version		0.0.1
//	@Description	删除用户
//	@Tags			超级用户操作 /system/admin
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string							false	"访问token"
//	@Param			access_token	query		string							false	"访问token"
//	@Param			SysUsername		body		SysUsername						true	"用户名"
//	@Success		200				{object}	model.HttpResponse{}			"删除用户成功"
//	@Failure		500				{object}	model.HttpResponse{data=string}	"删除用户失败"
//	@Router			/system/admin/user [DELETE]
func AdminDeleteUser(c *gin.Context) {
	var sysUsername SysUsername
	if err := c.ShouldBind(&sysUsername); err != nil {
		panic(http.StatusBadRequest)
	} else {
		err = deleteUser(sysUsername.Username)
		if err != nil {
			model.HttpResponse{}.FailGin(c, err.Error())
		} else {
			model.HttpResponse{}.OkGin(c, nil)
		}
	}
}

// GetAllRoles 获取所有权限组
//
//	@Summary		获取所有权限组
//	@Version		0.0.1
//	@Description	获取所有权限组
//	@Tags			权限管理 /system/role
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string									false	"访问token"
//	@Param			access_token	query		string									false	"访问token"
//	@Success		200				{object}	model.HttpResponse{data=[]RoleGroup}	"获取权限组列表成功"
//	@Failure		500				{object}	model.HttpResponse{data=string}			"获取权限组列表失败"
//	@Router			/system/role/roles [GET]
func GetAllRoles(c *gin.Context) {
	roles, err := selectAllRoles()
	if err != nil {
		model.HttpResponse{}.FailGin(c, "修改用户密码失败")
	} else {
		model.HttpResponse{}.OkGin(c, roles)
	}
}

// PostRole 更新权限组
//
//	@Summary		更新权限组
//	@Version		0.0.1
//	@Description	更新权限组
//	@Tags			权限管理 /system/role
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string									false	"访问token"
//	@Param			access_token	query		string									false	"访问token"
//	@Param			RoleGroup		body		RoleGroup								true	"用户权限信息"
//	@Success		200				{object}	model.HttpResponse{data=[]RoleGroup}	"获取权限组列表成功"
//	@Failure		500				{object}	model.HttpResponse{data=string}			"获取权限组列表失败"
//	@Router			/system/role [POST]
func PostRole(c *gin.Context) {
	var role RoleGroup
	if err := c.ShouldBind(&role); err != nil {
		panic(http.StatusBadRequest)
	} else {
		err = checkRole(role.RoleId)
		if err != nil {
			model.HttpResponse{}.FailGin(c, err.Error())
			return
		}

		err = updateRole(role.RoleId, role.RoleName)
		if err != nil {
			model.HttpResponse{}.FailGin(c, "更新角色失败")
		} else {
			model.HttpResponse{}.OkGin(c, nil)
		}
	}
}

// PutRole 新增权限组
//
//	@Summary		新增权限组
//	@Version		0.0.1
//	@Description	新增权限组
//	@Tags			权限管理 /system/role
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string									false	"访问token"
//	@Param			access_token	query		string									false	"访问token"
//	@Param			RoleGroup		body		RoleGroup								true	"用户权限信息"
//	@Success		200				{object}	model.HttpResponse{data=[]RoleGroup}	"获取权限组列表成功"
//	@Failure		500				{object}	model.HttpResponse{data=string}			"获取权限组列表失败"
//	@Router			/system/role [PUT]
func PutRole(c *gin.Context) {
	var role RoleGroup
	if err := c.ShouldBind(&role); err != nil {
		panic(http.StatusBadRequest)
	} else {
		if role.RoleName == "admin" || role.RoleName == "guest" {
			model.HttpResponse{}.FailGin(c, "不允许创建admin、guest权限组")
			return
		}

		err = insertRole(role.RoleName)
		if err != nil {
			model.HttpResponse{}.FailGin(c, "新增角色失败")
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
//	@Tags			权限管理 /system/role
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string									false	"访问token"
//	@Param			access_token	query		string									false	"访问token"
//	@Param			RoleGroup		body		RoleGroup								true	"用户权限信息"
//	@Success		200				{object}	model.HttpResponse{data=[]RoleGroup}	"获取权限组列表成功"
//	@Failure		500				{object}	model.HttpResponse{data=string}			"获取权限组列表失败"
//	@Router			/system/role [DELETE]
func DeleteRole(c *gin.Context) {
	var role RoleGroup
	if err := c.ShouldBind(&role); err != nil {
		panic(http.StatusBadRequest)
	} else {
		err = checkRole(role.RoleId)
		if err != nil {
			model.HttpResponse{}.FailGin(c, err.Error())
			return
		}

		err = deleteRole(role.RoleId)
		if err != nil {
			model.HttpResponse{}.FailGin(c, "删除角色失败")
		} else {
			model.HttpResponse{}.OkGin(c, nil)
		}
	}
}

// DeleteRoles 删除权限组
//
//	@Summary		删除权限组
//	@Version		0.0.1
//	@Description	删除权限组
//	@Tags			权限管理 /system/role
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string									false	"访问token"
//	@Param			access_token	query		string									false	"访问token"
//	@Param			RoleIdsEntity	body		RoleIdsEntity							true	"用户权限信息"
//	@Success		200				{object}	model.HttpResponse{data=[]RoleGroup}	"获取权限组列表成功"
//	@Failure		500				{object}	model.HttpResponse{data=string}			"获取权限组列表失败"
//	@Router			/system/role [DELETE]
func DeleteRoles(c *gin.Context) {
	var role RoleIdsEntity
	if err := c.ShouldBind(&role); err != nil {
		panic(http.StatusBadRequest)
	} else {
		for _, id := range role.Ids {
			err = checkRole(id)
			if err != nil {
				model.HttpResponse{}.FailGin(c, err.Error())
				return
			}
		}

		err = deleteRoles(role.Ids)
		if err != nil {
			model.HttpResponse{}.FailGin(c, "删除角色失败")
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
//	@Tags			权限管理 /system/role
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
		rules, err := getCasbinRuleByName(name)
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
//	@Tags			权限管理 /system/role
//	@Accept			json
//	@Produce		json
//	@Param			Authorization			header		string									false	"访问token"
//	@Param			access_token			query		string									false	"访问token"
//	@Param			RoleRuleUpdateEntity	body		RoleRuleUpdateEntity					true	"用户权限组名称"
//	@Success		200						{object}	model.HttpResponse{data=[]RoleGroup}	"获取权限组列表成功"
//	@Failure		500						{object}	model.HttpResponse{data=string}			"获取权限组列表失败"
//	@Router			/system/role/rules [POST]
func UpdateRoleRules(c *gin.Context) {
	var roleRule RoleRuleUpdateEntity
	if err := c.ShouldBind(&roleRule); err != nil {
		panic(http.StatusBadRequest)
	} else {
		role, err := getRoleById(roleRule.RoleId)
		if err != nil {
			model.HttpResponse{}.FailGin(c, "更新用户权限失败")
			return
		} else if role.RoleName == "admin" || role.RoleName == "guest" {
			model.HttpResponse{}.FailGin(c, "不允许修改admin、guest用户权限")
			return
		}

		err = updateRoleRules(role.RoleName, roleRule.Rules)
		if err != nil {
			model.HttpResponse{}.FailGin(c, err.Error())
		} else {
			model.HttpResponse{}.OkGin(c, nil)
		}
	}
}
