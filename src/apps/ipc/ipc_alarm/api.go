package ipc_alarm

import (
	"CrestedIbis/src/global"
	"CrestedIbis/src/global/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strconv"
	"strings"
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
//	@Param			keywords		query		string									false	"设备/通道 模糊值"
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
	keywords := c.DefaultQuery("keywords", "")
	start := c.DefaultQuery("start", time.DateTime)
	end := c.DefaultQuery("end", time.Now().Format(time.DateTime))

	total, data, err := selectIpcAlarmsByPages(page, pageSize, deviceId, channelId, start, end, keywords)
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
//	@Param			keywords		query		string									false	"设备/通道 模糊值"
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
	keywords := c.DefaultQuery("keywords", "")
	start := c.DefaultQuery("start", time.DateTime)
	end := c.DefaultQuery("end", time.Now().Format(time.DateTime))

	startTime, err := time.ParseInLocation(time.DateTime, start, time.Local)
	if err != nil {
		panic(http.StatusBadRequest)
	}
	endTime, err := time.ParseInLocation(time.DateTime, end, time.Local)
	if err != nil {
		panic(http.StatusBadRequest)
	}

	total, data, err := selectIpcRecordsByPages(page, pageSize, deviceId, channelId, startTime.Unix(), endTime.Unix(), keywords)
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

// DownloadRecord 下载录像文件
//
//	@Summary		下载录像文件
//	@Version		0.0.1
//	@Description	下载录像文件
//	@Tags			IPC设备 /ipc/device
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string									false	"访问token"
//	@Param			access_token	query		string									false	"访问token"
//	@Param			url				query		string									true	"视频地址"
//	@Success		200				{object}	model.HttpResponse{data=IpcRecordPage}	"分页查询成功"
//	@Failure		500				{object}	model.HttpResponse{data=string}			"查询数据失败"
//	@Router			/ipc/device/record [GET]
func DownloadRecord(c *gin.Context) {
	remoteURL := c.Query("url")
	if remoteURL == "" {
		panic(http.StatusBadRequest)
	}

	remoteURL = fmt.Sprintf("http://%s:%d/%s", global.Conf.GB28181.MediaServer.IP, global.Conf.GB28181.MediaServer.Port, remoteURL)
	// 发起 HTTP 请求获取远程文件
	resp, err := http.Get(remoteURL)
	if err != nil {
		model.HttpResponse{}.FailGin(c, "无法获取远程文件")
		return
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		model.HttpResponse{}.FailGin(c, "远程服务器返回错误")
		return
	}

	parts := strings.Split(remoteURL, "/")
	var filename string
	if len(parts) <= 1 {
		filename = fmt.Sprintf("%s", parts[len(parts)-1])
	} else {
		filename = fmt.Sprintf("%s %s", parts[len(parts)-2], parts[len(parts)-1])
	}

	// 设置响应头
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", resp.Header.Get("Content-Type"))
	c.Header("Content-Length", resp.Header.Get("Content-Length"))

	// 流式传输文件
	if _, err := io.Copy(c.Writer, resp.Body); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}

// DeleteRecord 删除录像文件
//
//	@Summary		删除录像文件
//	@Version		0.0.1
//	@Description	删除录像文件
//	@Tags			IPC设备 /ipc/device
//	@Accept			json
//	@Produce		json
//	@Param			Authorization		header		string									false	"访问token"
//	@Param			access_token		query		string									false	"访问token"
//	@Param			IpcRecordIdEntity	body		IpcRecordIdEntity						true	"视频地址"
//	@Success		200					{object}	model.HttpResponse{data=IpcRecordPage}	"分页查询成功"
//	@Failure		500					{object}	model.HttpResponse{data=string}			"查询数据失败"
//	@Router			/ipc/device/record [DELETE]
func DeleteRecord(c *gin.Context) {
	var idEntity IpcRecordIdEntity

	if err := c.ShouldBind(&idEntity); err != nil {
		panic(http.StatusBadRequest)
	} else {
		err = DeleteRecordServer(idEntity.ID)
		if err != nil {
			model.HttpResponse{}.FailGin(c, err.Error())
		} else {
			model.HttpResponse{}.OkGin(c, "删除文件成功")
		}
	}
}

// DeleteRecords 删除录像文件
//
//	@Summary		删除录像文件
//	@Version		0.0.1
//	@Description	删除录像文件
//	@Tags			IPC设备 /ipc/device
//	@Accept			json
//	@Produce		json
//	@Param			Authorization		header		string									false	"访问token"
//	@Param			access_token		query		string									false	"访问token"
//	@Param			IpcRecordIdEntity	body		IpcRecordIdEntity						true	"视频地址"
//	@Success		200					{object}	model.HttpResponse{data=IpcRecordPage}	"分页查询成功"
//	@Failure		500					{object}	model.HttpResponse{data=string}			"查询数据失败"
//	@Router			/ipc/device/records [DELETE]
func DeleteRecords(c *gin.Context) {
	var idEntity IpcRecordIdEntity

	if err := c.ShouldBind(&idEntity); err != nil {
		panic(http.StatusBadRequest)
	} else {
		for _, id := range idEntity.Ids {
			DeleteRecordServer(id)
		}

		model.HttpResponse{}.OkGin(c, "删除文件成功")
	}
}
