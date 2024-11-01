package site

import (
	"CrestedIbis/src/global/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// GetSites 获取区域列表
//
//	@Summary		获取区域列表
//	@Version		1.0.0
//	@Description	获取区域列表
//	@Tags			区域管理 /site
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string							false	"访问token"
//	@Param			access_token	query		string							false	"访问token"
//	@Param			pid				query		number							false	"父区域ID"
//	@Success		200				{object}	model.HttpResponse{data=[]Site}	"查询成功"
//	@Failure		500				{object}	model.HttpResponse{}			"查询失败"
//	@Router			/site/sites [GET]
func GetSites(c *gin.Context) {
	var pid *int64
	pidQuery := c.Query("pid")
	if pidQuery == "" {
		pid = nil
	} else {
		_pid, err := strconv.ParseInt(pidQuery, 10, 0)
		if err != nil {
			panic(http.StatusBadRequest)
		} else {
			pid = &_pid
		}
	}

	sites, err := selectSites(pid)
	if err != nil {
		model.HttpResponse{}.FailGin(c, "查询失败")
	} else {
		model.HttpResponse{}.OkGin(c, sites)
	}
}

// PutSite 新建区域
//
//	@Summary		新建区域
//	@Version		1.0.0
//	@Description	新建区域
//	@Tags			区域管理 /site
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string					false	"访问token"
//	@Param			access_token	query		string					false	"访问token"
//	@Param			Site			body		Site					true	"区域信息"
//	@Success		200				{object}	model.HttpResponse{}	"新建成功"
//	@Failure		500				{object}	model.HttpResponse{}	"新建失败"
//	@Router			/site [PUT]
func PutSite(c *gin.Context) {
	var site Site
	if err := c.ShouldBindJSON(&site); err != nil {
		panic(http.StatusBadRequest)
	} else {
		err := insertSite(site)
		if err != nil {
			model.HttpResponse{}.FailGin(c, err.Error())
		} else {
			model.HttpResponse{}.OkGin(c, site)
		}
	}
}
