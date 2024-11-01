package site

import "CrestedIbis/src/global"

func InitSiteRouter() {
	siteRouter := global.HttpEngine.Group("/site")
	{
		siteRouter.GET("/sites", GetSites)
		siteRouter.PUT("", PutSite)
	}
}
