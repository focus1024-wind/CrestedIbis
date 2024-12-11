package user

import "CrestedIbis/src/global"

func InitSystemUserRouter() {
	systemUserRouter := global.HttpEngine.Group("/system/user")
	{
		systemUserRouter.POST("/login", Login)
		systemUserRouter.POST("/register", Register)
		systemUserRouter.GET("/users", GetAllUserByPages)
	}
	systemAdminRouter := global.HttpEngine.Group("/system/admin")
	{
		systemAdminRouter.GET("/users", GetAllUserByPages)
		systemAdminRouter.POST("/password", AdminChangePassword)
		systemAdminRouter.DELETE("/user", AdminDeleteUser)
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
