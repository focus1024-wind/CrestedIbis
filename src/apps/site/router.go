package site

import "CrestedIbis/src/global"

func InitSiteRouter() {
	siteRouter := global.HttpEngine.Group("/site")
	{
		siteRouter.POST("", PostSite)
		siteRouter.PUT("", PutSite)
		siteRouter.DELETE("", DeleteSite)
		siteRouter.GET("/sites", GetSites)
	}
}
