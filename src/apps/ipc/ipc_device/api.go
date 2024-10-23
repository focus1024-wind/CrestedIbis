package ipc_device

import (
	"CrestedIbis/src/global/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type IpcDevicePage struct {
	Total    int64       `json:"total" desc:"设备总数量"`
	Data     []IpcDevice `json:"data" desc:"设备列表"`
	Page     int64       `json:"page" desc:"页码"`
	PageSize int64       `json:"page_size" desc:"每页查询数量"`
}

// GetIpcDevicesByPages 分页查询IpcDevice设备
//
//	@Title			分页查询IpcDevice设备
//	@Version		1.0.0
//	@Description	分页查询GB28181 IpcDevice设备
//	@Tags			IPC设备 /ipc/device
//	@Accept			json
//	@Produce		json
//	@Param			access_token	query		string									false	"访问token"
//	@Param			page			query		integer									false	"分页查询页码，默认值: 1"
//	@Param			page_size		query		integer									false	"每页查询数量，默认值: 15"
//	@Success		200				{object}	model.HttpResponse{data=IpcDevicePage}	"分页查询成功"
//	@Failure		200				{object}	model.HttpResponse{data=string}			"查询数据失败"
//	@Router			/ipc/device/devices [GET]
func GetIpcDevicesByPages(c *gin.Context) {
	pageQuery := c.DefaultQuery("page", "1")
	pageSizeQuery := c.DefaultQuery("page_size", "15")
	page, err := strconv.ParseInt(pageQuery, 10, 0)
	if err != nil {
		panic(http.StatusBadRequest)
	}
	pageSize, err := strconv.ParseInt(pageSizeQuery, 10, 0)
	if err != nil {
		panic(http.StatusBadRequest)
	}

	total, data, err := selectIpcDevicesByPages(page, pageSize)
	if err != nil {
		model.HttpResponse{}.FailGin(c, "查询数据失败")
		return
	} else {
		model.HttpResponse{}.OkGin(c, &IpcDevicePage{
			Total:    total,
			Data:     data,
			Page:     page,
			PageSize: pageSize,
		})
	}
}
