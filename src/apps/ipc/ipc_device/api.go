package ipc_device

import (
	"CrestedIbis/src/global"
	"CrestedIbis/src/global/model"
	"CrestedIbis/src/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"path"
	"strconv"
	"strings"
)

type IpcDevicePage struct {
	Total    int64       `json:"total" desc:"设备总数量"`
	Data     []IpcDevice `json:"data" desc:"设备列表"`
	Page     int64       `json:"page" desc:"页码"`
	PageSize int64       `json:"page_size" desc:"每页查询数量"`
}

// GetIpcDevicesByPages 分页查询IpcDevice设备
//
//	@Summary		分页查询IpcDevice设备
//	@Version		0.0.1
//	@Description	分页查询GB28181 IpcDevice设备
//	@Tags			IPC设备 /ipc/device
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string									false	"访问token"
//	@Param			access_token	query		string									false	"访问token"
//	@Param			page			query		integer									false	"分页查询页码，默认值: 1"
//	@Param			page_size		query		integer									false	"每页查询数量，默认值: 15"
//	@Success		200				{object}	model.HttpResponse{data=IpcDevicePage}	"分页查询成功"
//	@Failure		500				{object}	model.HttpResponse{data=string}			"查询数据失败"
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

// GetIpcChannels 获取设备通道信息
//
//	@Summary		获取设备通道信息
//	@Version		0.0.1
//	@Description	查询GB28181 设备对应通道信息
//	@Tags			IPC设备 /ipc/device
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string									false	"访问token"
//	@Param			access_token	query		string									false	"访问token"
//	@Param			device_id		query		string									true	"设备ID"
//	@Success		200				{object}	model.HttpResponse{data=[]IpcChannel}	"查询数据成功"
//	@Failure		500				{object}	model.HttpResponse{data=string}			"查询数据失败"
//	@Router			/ipc/device/channels [GET]
func GetIpcChannels(c *gin.Context) {
	deviceId := c.Query("device_id")
	if deviceId == "" {
		panic(http.StatusBadRequest)
	}

	ipcChannels, err := selectIpcChannels(deviceId)

	if err != nil {
		model.HttpResponse{}.FailGin(c, "查询数据失败")
		return
	} else {
		model.HttpResponse{}.OkGin(c, ipcChannels)
	}
}

// IpcUploadImage IPC图像上传
//
//	@Summary		IPC图像上传
//	@Version		0.0.1
//	@Description	GB28181图像抓拍，图片上传接口
//	@Tags			IPC设备 /ipc/device
//	@Accept			mpfd
//	@Produce		json
//	@Param			access_token	query		string							true	"访问token"
//	@Param			file			formData	file							true	"上传图片"
//	@Success		200				{object}	model.HttpResponse{data=string}	"上传图片成功"
//	@Failure		500				{object}	model.HttpResponse{data=string}	"上传图片失败"
//	@Router			/ipc/device/upload_image [POST]
func IpcUploadImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		model.HttpResponse{}.FailGin(c, "获取图片数据失败")
		return
	}

	fileExt := strings.ToLower(path.Ext(file.Filename))
	if fileExt == ".jpg" || fileExt == ".png" || fileExt == ".jpeg" {
		filePath := global.Conf.Store.Snapshot
		claims, ok := c.Get("claims")
		if ok {
			switch jwtToken := claims.(type) {
			case utils.JwtToken:
				filePath = fmt.Sprintf("%s/%s", filePath, jwtToken.Username)
			}
		}

		err := c.SaveUploadedFile(file, fmt.Sprintf("%s/%s", filePath, file.Filename))
		if err != nil {
			model.HttpResponse{}.FailGin(c, "保存图片失败")
			return
		} else {
			model.HttpResponse{}.OkGin(c, nil)
			return
		}
	} else {
		model.HttpResponse{}.FailGin(c, "图片格式错误，目前仅支持.jpg, .png, .jpeg格式数据上传")
		return
	}
}
