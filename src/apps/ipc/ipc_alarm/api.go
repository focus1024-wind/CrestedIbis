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

// DeleteIpcAlarm 删除告警记录
//
//	@Summary		删除告警记录
//	@Version		0.0.1
//	@Description	删除告警记录
//	@Tags			IPC告警 /ipc/alarm
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string							false	"访问token"
//	@Param			access_token	query		string							false	"访问token"
//	@Param			model.IDModel	body		model.IDModel					true	"告警ID"
//	@Success		200				{object}	model.HttpResponse{data=nil}	"删除告警记录成功"
//	@Failure		500				{object}	model.HttpResponse{data=string}	"删除告警记录失败"
//	@Router			/ipc/alarm [DELETE]
func DeleteIpcAlarm(c *gin.Context) {
	var idModel model.IDModel

	if err := c.ShouldBind(&idModel); err != nil {
		panic(http.StatusBadRequest)
	} else {
		err = IpcAlarm{}.Delete(idModel)
		if err != nil {
			model.HttpResponse{}.FailGin(c, err.Error())
		} else {
			model.HttpResponse{}.OkGin(c, "删除告警记录成功")
		}
	}
}

// GetIpcAlarms 获取IPC设备告警列表
//
//	@Summary		获取IPC设备告警列表
//	@Version		0.0.1
//	@Description	获取IPC设备告警列表
//	@Tags			IPC告警 /ipc/alarm
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string																false	"访问token"
//	@Param			access_token	query		string																false	"访问token"
//	@Param			page			query		integer																false	"分页查询页码，默认值: 1"
//	@Param			page_size		query		integer																false	"每页查询数量，默认值: 15"
//	@Param			device_id		query		string																false	"设备ID"
//	@Param			channel_id		query		string																false	"通道ID"
//	@Param			keywords		query		string																false	"设备/通道 模糊值"
//	@Param			start			query		string																false	"开始时间，默认值: 2006-01-02 15:04:05"
//	@Param			end				query		string																false	"结束时间，默认值: 当前时间"
//	@Success		200				{object}	model.HttpResponse{data=model.BasePageResponse{data=[]IpcAlarm}}	"分页查询成功"
//	@Failure		500				{object}	model.HttpResponse{data=string}										"查询数据失败"
//	@Router			/ipc/alarm/alarms [GET]
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

	total, data, err := IpcAlarm{}.SelectAlarms(page, pageSize, deviceId, channelId, start, end, keywords)
	if err != nil {
		model.HttpResponse{}.FailGin(c, "查询数据失败")
		return
	} else {
		model.HttpResponse{}.OkGin(c, &model.BasePageResponse{
			Total:    total,
			Data:     data,
			Page:     page,
			PageSize: pageSize,
		})
	}
}

// DeleteIpcAlarms 批量删除告警记录
//
//	@Summary		批量删除告警记录
//	@Version		0.0.1
//	@Description	批量删除告警记录
//	@Tags			IPC告警 /ipc/alarm
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string							false	"访问token"
//	@Param			access_token	query		string							false	"访问token"
//	@Param			model.IDsModel	body		model.IDsModel					true	"ID列表"
//	@Success		200				{object}	model.HttpResponse{data=nil}	"批量删除告警记录成功"
//	@Failure		500				{object}	model.HttpResponse{data=string}	"批量删除告警记录失败"
//	@Router			/ipc/alarm/alarms [DELETE]
func DeleteIpcAlarms(c *gin.Context) {
	var idsModel model.IDsModel

	if err := c.ShouldBind(&idsModel); err != nil {
		panic(http.StatusBadRequest)
	} else {
		err = IpcAlarm{}.Deletes(idsModel)
		if err != nil {
			model.HttpResponse{}.FailGin(c, err.Error())
		} else {
			model.HttpResponse{}.OkGin(c, "删除文件成功")
		}
	}
}

// DownloadRecord 下载录像文件
//
//	@Summary		下载录像文件
//	@Version		0.0.1
//	@Description	下载录像文件
//	@Tags			IPC录像 /ipc/record
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header	string	false	"访问token"
//	@Param			access_token	query	string	false	"访问token"
//	@Param			url				query	string	true	"视频地址"
//	@Router			/ipc/record [GET]
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
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			global.Logger.Error(err.Error())
		}
	}(resp.Body)

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
//	@Tags			IPC录像 /ipc/record
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string							false	"访问token"
//	@Param			access_token	query		string							false	"访问token"
//	@Param			model.IDModel	body		model.IDModel					true	"录像ID信息"
//	@Success		200				{object}	model.HttpResponse{data=nil}	"删除录像文件成功"
//	@Failure		500				{object}	model.HttpResponse{data=string}	"删除录像文件失败"
//	@Router			/ipc/record [DELETE]
func DeleteRecord(c *gin.Context) {
	var idEntity model.IDModel

	if err := c.ShouldBind(&idEntity); err != nil {
		panic(http.StatusBadRequest)
	} else {
		err = IpcRecord{}.Delete(idEntity)
		if err != nil {
			model.HttpResponse{}.FailGin(c, err.Error())
		} else {
			model.HttpResponse{}.OkGin(c, "删除文件成功")
		}
	}
}

// GetIpcRecords 获取IPC设备录像列表
//
//	@Summary		获取IPC设备录像列表
//	@Version		0.0.1
//	@Description	获取IPC设备录像列表
//	@Tags			IPC录像 /ipc/record
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string																false	"访问token"
//	@Param			access_token	query		string																false	"访问token"
//	@Param			page			query		integer																false	"分页查询页码，默认值: 1"
//	@Param			page_size		query		integer																false	"每页查询数量，默认值: 15"
//	@Param			device_id		query		string																false	"设备ID"
//	@Param			channel_id		query		string																false	"通道ID"
//	@Param			keywords		query		string																false	"设备/通道 模糊值"
//	@Param			start			query		string																false	"开始时间，默认值: 2006-01-02 15:04:05"
//	@Param			end				query		string																false	"结束时间，默认值: 当前时间"
//	@Success		200				{object}	model.HttpResponse{data=model.BasePageResponse{data=[]IpcRecord}}	"IPC设备录像列表"
//	@Failure		500				{object}	model.HttpResponse{data=string}										"查询数据失败"
//	@Router			/ipc/record/records [GET]
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

	total, data, err := IpcRecord{}.SelectRecords(page, pageSize, deviceId, channelId, start, end, keywords)
	if err != nil {
		model.HttpResponse{}.FailGin(c, "查询数据失败")
		return
	} else {
		model.HttpResponse{}.OkGin(c, &model.BasePageResponse{
			Total:    total,
			Data:     data,
			Page:     page,
			PageSize: pageSize,
		})
	}
}

// DeleteRecords 批量删除录像文件
//
//	@Summary		批量删除录像文件
//	@Version		0.0.1
//	@Description	批量删除录像文件
//	@Tags			IPC录像 /ipc/record
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string							false	"访问token"
//	@Param			access_token	query		string							false	"访问token"
//	@Param			model.IDModel	body		model.IDModel					true	"ID列表"
//	@Success		200				{object}	model.HttpResponse{data=nil}	"批量删除录像文件成功"
//	@Failure		500				{object}	model.HttpResponse{data=string}	"批量删除录像文件失败"
//	@Router			/ipc/record/records [DELETE]
func DeleteRecords(c *gin.Context) {
	var idEntity model.IDsModel

	if err := c.ShouldBind(&idEntity); err != nil {
		panic(http.StatusBadRequest)
	} else {
		var errDeleteRecordIds []int64

		for _, id := range idEntity.IDs {
			err = IpcRecord{}.Delete(model.IDModel{ID: id})
			if err != nil {
				errDeleteRecordIds = append(errDeleteRecordIds, id)
				err = nil
			}
		}

		if len(errDeleteRecordIds) > 0 {
			model.HttpResponse{}.FailGin(c, fmt.Sprintf("%v 删除失败", errDeleteRecordIds))
		} else {
			model.HttpResponse{}.OkGin(c, nil)
		}
	}
}
