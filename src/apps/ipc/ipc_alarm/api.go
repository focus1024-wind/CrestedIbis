package ipc_alarm

import (
	"CrestedIbis/src/global/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type IpcAlarmPage struct {
	Total    int64      `json:"total" desc:"设备总数量"`
	Data     []IpcAlarm `json:"data" desc:"设备列表"`
	Page     int64      `json:"page" desc:"页码"`
	PageSize int64      `json:"page_size" desc:"每页查询数量"`
}

type IpcRecordPage struct {
	Total    int64       `json:"total" desc:"设备总数量"`
	Data     []IpcRecord `json:"data" desc:"设备列表"`
	Page     int64       `json:"page" desc:"页码"`
	PageSize int64       `json:"page_size" desc:"每页查询数量"`
}

// GetIpcAlarms 分页查询IpcDevice设备报警信息
//
//	@Summary		分页查询IpcDevice设备报警信息
//	@Version		0.0.1
//	@Description	分页查询IpcDevice设备报警信息
//	@Tags			IPC设备 /ipc/device
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string									false	"访问token"
//	@Param			access_token	query		string									false	"访问token"
//	@Param			page			query		integer									false	"分页查询页码，默认值: 1"
//	@Param			page_size		query		integer									false	"每页查询数量，默认值: 15"
//	@Param			device_id		query		string									false	"设备ID"
//	@Param			channel_id		query		string									false	"通道ID"
//	@Param			start			query		string									false	"开始时间，默认值: 2006-01-02 15:04:05"
//	@Param			end				query		string									false	"结束时间，默认值: 当前时间"
//	@Success		200				{object}	model.HttpResponse{data=IpcAlarmPage}	"分页查询成功"
//	@Failure		500				{object}	model.HttpResponse{data=string}			"查询数据失败"
//	@Router			/ipc/device/alarms [GET]
func GetIpcAlarms(c *gin.Context) {
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

	deviceId := c.Query("device_id")
	channelId := c.Query("channel_id")
	start := c.DefaultQuery("start", time.DateTime)
	end := c.DefaultQuery("end", time.Now().Format(time.DateTime))

	total, data, err := selectIpcAlarmsByPages(page, pageSize, deviceId, channelId, start, end)
	if err != nil {
		model.HttpResponse{}.FailGin(c, "查询数据失败")
		return
	} else {
		model.HttpResponse{}.OkGin(c, &IpcAlarmPage{
			Total:    total,
			Data:     data,
			Page:     page,
			PageSize: pageSize,
		})
	}
}

// GetIpcRecords 分页查询IpcDevice设备录像信息
//
//	@Summary		分页查询IpcDevice设备录像信息
//	@Version		0.0.1
//	@Description	分页查询IpcDevice设备录像信息
//	@Tags			IPC设备 /ipc/device
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string									false	"访问token"
//	@Param			access_token	query		string									false	"访问token"
//	@Param			page			query		integer									false	"分页查询页码，默认值: 1"
//	@Param			page_size		query		integer									false	"每页查询数量，默认值: 15"
//	@Param			device_id		query		string									false	"设备ID"
//	@Param			channel_id		query		string									false	"通道ID"
//	@Param			start			query		string									false	"开始时间，默认值: 2006-01-02 15:04:05"
//	@Param			end				query		string									false	"结束时间，默认值: 当前时间"
//	@Success		200				{object}	model.HttpResponse{data=IpcRecordPage}	"分页查询成功"
//	@Failure		500				{object}	model.HttpResponse{data=string}			"查询数据失败"
//	@Router			/ipc/device/records [GET]
func GetIpcRecords(c *gin.Context) {
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

	deviceId := c.Query("device_id")
	channelId := c.Query("channel_id")
	start := c.DefaultQuery("start", time.DateTime)
	end := c.DefaultQuery("end", time.Now().Format(time.DateTime))

	total, data, err := selectIpcRecordsByPages(page, pageSize, deviceId, channelId, start, end)
	if err != nil {
		model.HttpResponse{}.FailGin(c, "查询数据失败")
		return
	} else {
		model.HttpResponse{}.OkGin(c, &IpcRecordPage{
			Total:    total,
			Data:     data,
			Page:     page,
			PageSize: pageSize,
		})
	}
}
