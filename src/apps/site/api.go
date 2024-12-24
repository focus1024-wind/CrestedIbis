package site

import (
	"CrestedIbis/src/global/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// CreateSite 新建区域
//
//	@Summary		新建区域
//	@Version		0.0.1
//	@Description	新建区域
//	@Tags			区域管理 /site
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string							false	"访问token"
//	@Param			access_token	query		string							false	"访问token"
//	@Param			Site			body		Site							true	"区域信息"
//	@Success		200				{object}	model.HttpResponse{data=Site}	"新建成功"
//	@Failure		500				{object}	model.HttpResponse{data=nil}	"新建失败"
//	@Router			/site [PUT]
func CreateSite(c *gin.Context) {
	var site Site
	if err := c.ShouldBindJSON(&site); err != nil {
		panic(http.StatusBadRequest)
	} else {
		err = Site{}.Insert(site)
		if err != nil {
			model.HttpResponse{}.FailGin(c, err.Error())
		} else {
			model.HttpResponse{}.OkGin(c, site)
		}
	}
}

// UpdateSite 修改区域信息
//
//	@Summary		修改区域信息
//	@Version		0.0.1
//	@Description	修改区域信息
//	@Tags			区域管理 /site
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string							false	"访问token"
//	@Param			access_token	query		string							false	"访问token"
//	@Param			Site			body		Site							true	"区域信息"
//	@Success		200				{object}	model.HttpResponse{data=nil}	"新建成功"
//	@Failure		500				{object}	model.HttpResponse{data=string}	"新建失败"
//	@Router			/site [POST]
func UpdateSite(c *gin.Context) {
	var site Site
	if err := c.ShouldBind(&site); err != nil {
		panic(http.StatusBadRequest)
	} else {
		err = Site{}.Update(site)
		if err != nil {
			model.HttpResponse{}.FailGin(c, "修改区域名称失败")
		} else {
			model.HttpResponse{}.OkGin(c, nil)
		}
	}
}

// DeleteSite 删除区域站点
//
//	@Summary		删除区域站点
//	@Version		0.0.1
//	@Description	删除区域站点
//	@Tags			区域管理 /site
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string					false	"访问token"
//	@Param			access_token	query		string					false	"访问token"
//	@Param			model.IDModel	body		model.IDModel			true	"区域站点ID"
//	@Success		200				{object}	model.HttpResponse{}	"删除区域站点成功"
//	@Failure		500				{object}	model.HttpResponse{}	"删除区域站点失败"
//	@Router			/site [DELETE]
func DeleteSite(c *gin.Context) {
	var idModel model.IDModel
	if err := c.ShouldBind(&idModel); err != nil {
		panic(http.StatusBadRequest)
	} else {
		err = Site{}.Delete(idModel)
		if err != nil {
			model.HttpResponse{}.FailGin(c, "删除区域失败")
		} else {
			model.HttpResponse{}.OkGin(c, nil)
		}
	}
}

// GetSites 获取区域列表
//
//	@Summary		获取区域列表
//	@Version		0.0.1
//	@Description	获取区域列表
//	@Tags			区域管理 /site
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string							false	"访问token"
//	@Param			access_token	query		string							false	"访问token"
//	@Param			pid				query		number							false	"父区域ID"
//	@Param			keywords		query		string							false	"模糊区域名称信息"
//	@Success		200				{object}	model.HttpResponse{data=[]Site}	"查询成功"
//	@Failure		500				{object}	model.HttpResponse{}			"查询失败"
//	@Router			/site/sites [GET]
func GetSites(c *gin.Context) {
	var pid *int64
	pidQuery := c.Query("pid")
	keywords := c.DefaultQuery("keywords", "")
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

	sites, err := Site{}.SelectList(pid, keywords)
	if err != nil {
		model.HttpResponse{}.FailGin(c, "查询失败")
	} else {
		model.HttpResponse{}.OkGin(c, sites)
	}
}

// DeleteSites 批量删除区域
//
//	@Summary		批量删除区域
//	@Version		0.0.1
//	@Description	批量删除区域
//	@Tags			区域管理 /site
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string							false	"访问token"
//	@Param			access_token	query		string							false	"访问token"
//	@Param			model.IDsModel	body		model.IDsModel					true	"区域ID列表"
//	@Success		200				{object}	model.HttpResponse{data=nil}	"删除成功"
//	@Failure		500				{object}	model.HttpResponse{data=string}	"删除失败"
//	@Router			/site/sites [DELETE]
func DeleteSites(c *gin.Context) {
	var idsModel model.IDsModel
	if err := c.ShouldBind(&idsModel); err != nil {
		panic(http.StatusBadRequest)
	} else {
		err = Site{}.Deletes(idsModel)
		if err != nil {
			model.HttpResponse{}.FailGin(c, "删除区域失败")
		} else {
			model.HttpResponse{}.OkGin(c, nil)
		}
	}
}
