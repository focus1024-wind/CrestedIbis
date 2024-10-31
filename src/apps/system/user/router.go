package user

import "CrestedIbis/src/global"

func InitSystemUserRouter() {
	systemUserRouter := global.HttpEngine.Group("/system/user")
	{
		systemUserRouter.POST("/login", Login)
		systemUserRouter.POST("/register", Register)
		systemUserRouter.GET("/users", AdminGetAllUserByPages)
	}
}
