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

	// 权限组 权限 接口
	systemRoleRouter := global.HttpEngine.Group("/system/role")
	{
		systemRoleRouter.PUT("", CreateRole)
		systemRoleRouter.POST("", UpdateRole)
		systemRoleRouter.DELETE("", DeleteRole)

		systemRoleRouter.GET("/roles", GetRoles)
		systemRoleRouter.DELETE("/roles", DeleteRoles)

		systemRoleRouter.GET("/rules", GetRoleRules)
		systemRoleRouter.POST("/rules", UpdateRoleRules)
	}
}
