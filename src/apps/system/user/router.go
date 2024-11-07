package user

import "CrestedIbis/src/global"

func InitSystemUserRouter() {
	systemUserRouter := global.HttpEngine.Group("/system/user")
	{
		systemUserRouter.POST("/login", Login)
		systemUserRouter.POST("/register", Register)
	}
	systemAdminRouter := global.HttpEngine.Group("/system/admin")
	{
		systemAdminRouter.GET("/users", AdminGetAllUserByPages)
		systemAdminRouter.POST("/password", AdminChangePassword)
		systemAdminRouter.DELETE("/user", AdminDeleteUser)
	}
}
