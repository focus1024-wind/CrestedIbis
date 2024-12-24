package site

import "CrestedIbis/src/global"

func InitSiteRouter() {
	siteRouter := global.HttpEngine.Group("/site")
	{
		siteRouter.PUT("", CreateSite)
		siteRouter.POST("", UpdateSite)
		siteRouter.DELETE("", DeleteSite)

		siteRouter.GET("/sites", GetSites)
		siteRouter.DELETE("/sites", DeleteSites)
	}
}
