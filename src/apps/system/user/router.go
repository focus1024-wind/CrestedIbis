package user

import "CrestedIbis/src/global"

func InitSystemUserRouter() {
	// 超级管理员接口
	systemAdminRouter := global.HttpEngine.Group("/system/admin")
	{
		systemAdminRouter.POST("/password", AdminChangePassword)
	}

	// 用户接口
	systemUserRouter := global.HttpEngine.Group("/system/user")
	{
		systemUserRouter.POST("", UpdateUser)
		systemUserRouter.DELETE("", DeleteUser)
		systemUserRouter.POST("/login", Login)
		systemUserRouter.POST("/register", Register)
		systemUserRouter.GET("/users", GetUsers)
		systemUserRouter.DELETE("/users", DeleteUsers)
	}

	systemRoleRouter := global.HttpEngine.Group("/system/role")
	{
		systemRoleRouter.GET("/roles", GetAllRoles)
		systemRoleRouter.DELETE("/roles", DeleteRoles)
		systemRoleRouter.POST("", PostRole)
		systemRoleRouter.PUT("", PutRole)
		systemRoleRouter.DELETE("/", DeleteRole)
	}

	systemRoleRuleRouter := global.HttpEngine.Group("/system/role/rules")
	{
		systemRoleRuleRouter.GET("", GetRoleRules)
		systemRoleRuleRouter.POST("", UpdateRoleRules)
	}
}
