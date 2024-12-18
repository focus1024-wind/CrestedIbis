package ipc_device

import (
	"CrestedIbis/gb28181_server"
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

// GetIpcDevice 根据device_id获取IPC设备
//
//	@Summary		根据device_id获取IPC设备
//	@Version		0.0.1
//	@Description	根据device_id获取IPC设备
//	@Tags			IPC设备 /ipc/device
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string								false	"访问token"
//	@Param			access_token	query		string								false	"访问token"
//	@Param			device_id		query		string								true	"设备ID"
//	@Success		200				{object}	model.HttpResponse{data=IpcDevice}	"查询成功"
//	@Failure		500				{object}	model.HttpResponse{data=string}		"查询数据失败"
//	@Router			/ipc/device [GET]
func GetIpcDevice(c *gin.Context) {
	deviceID := c.Query("device_id")
	if deviceID == "" {
		panic(http.StatusBadRequest)
	} else {
		device, err := selectIpcDevice(deviceID)
		if err != nil {
			model.HttpResponse{}.FailGin(c, "查询设备失败")
		} else {
			model.HttpResponse{}.OkGin(c, device)
		}
	}
}

// PostIpcDevice 更新IPC设备
//
//	@Summary		更新IPC设备
//	@Version		0.0.1
//	@Description	更新IPC设备
//	@Tags			IPC设备 /ipc/device
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string								false	"访问token"
//	@Param			access_token	query		string								false	"访问token"
//	@Param			IpcDevice		body		IpcDevice							true	"设备信息"
//	@Success		200				{object}	model.HttpResponse{data=IpcDevice}	"更新成功"
//	@Failure		500				{object}	model.HttpResponse{data=string}		"查询数据失败"
//	@Router			/ipc/device [POST]
func PostIpcDevice(c *gin.Context) {
	var ipcDevice IpcDevice
	if err := c.ShouldBind(&ipcDevice); err != nil {
		panic(http.StatusBadRequest)
	}
	err := updateIpcDevice(ipcDevice)
	if err != nil {
		model.HttpResponse{}.FailGin(c, err.Error())
	} else {
		model.HttpResponse{}.OkGin(c, ipcDevice)
	}
}

// DeleteIpcDevice 删除IPC设备
//
//	@Summary		删除IPC设备
//	@Version		0.0.1
//	@Description	删除IPC设备及对应通道，该删除仅为删除数据库记录，不影响IPC设备的重新注册
//	@Tags			IPC设备 /ipc/device
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string							false	"访问token"
//	@Param			access_token	query		string							false	"访问token"
//	@Param			IpcDeviceID		body		IpcDeviceID						true	"设备ID"
//	@Success		200				{object}	model.HttpResponse{data=string}	"查询数据成功"
//	@Failure		500				{object}	model.HttpResponse{data=string}	"查询数据失败"
//	@Router			/ipc/device [DELETE]
func DeleteIpcDevice(c *gin.Context) {
	var deviceID IpcDeviceID
	if err := c.ShouldBind(&deviceID); err != nil {
		panic(http.StatusBadRequest)
	}

	err := deleteIpcDevice(deviceID.DeviceID)
	if err != nil {
		model.HttpResponse{}.FailGin(c, "删除设备失败")
	} else {
		model.HttpResponse{}.OkGin(c, "删除设备成功")
	}
}

// DeleteIpcDevices 删除IPC设备
//
//	@Summary		删除IPC设备
//	@Version		0.0.1
//	@Description	删除IPC设备及对应通道，该删除仅为删除数据库记录，不影响IPC设备的重新注册
//	@Tags			IPC设备 /ipc/device
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string							false	"访问token"
//	@Param			access_token	query		string							false	"访问token"
//	@Param			IpcDeviceID		body		IpcDeviceID						true	"设备ID"
//	@Success		200				{object}	model.HttpResponse{data=string}	"查询数据成功"
//	@Failure		500				{object}	model.HttpResponse{data=string}	"查询数据失败"
//	@Router			/ipc/device/devices [DELETE]
func DeleteIpcDevices(c *gin.Context) {
	var idsModel IpcDeviceIDs
	if err := c.ShouldBind(&idsModel); err != nil {
		panic(http.StatusBadRequest)
	}

	err := deleteIpcDevices(idsModel.DeviceIDs)
	if err != nil {
		model.HttpResponse{}.FailGin(c, "删除设备失败")
	} else {
		model.HttpResponse{}.OkGin(c, "删除设备成功")
	}
}

// GetIpcDevicesStatus 获取设备状态信息
//
//	@Summary		获取设备状态信息
//	@Version		0.0.1
//	@Description	获取设备状态信息，总设备量，在线设备量，离线设备量
//	@Tags			IPC设备 /ipc/device
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string							false	"访问token"
//	@Param			access_token	query		string							false	"访问token"
//	@Success		200				{object}	model.HttpResponse{data=string}	"查询数据成功"
//	@Failure		500				{object}	model.HttpResponse{data=string}	"查询数据失败"
//	@Router			/ipc/device/status [GET]
func GetIpcDevicesStatus(c *gin.Context) {
	var (
		total     int64
		online    int64
		offline   int64
		ipcDevice IpcDevice
	)

	global.Db.Model(IpcDevice{}).Count(&total)
	ipcDevice.Status = gb28181_server.DeviceOnLineStatus
	global.Db.Model(IpcDevice{}).Where(IpcDevice{
		GB28181Device: gb28181_server.GB28181Device{
			Status: gb28181_server.DeviceOnLineStatus,
		},
	}).Count(&online)
	global.Db.Model(IpcDevice{}).Where(IpcDevice{
		GB28181Device: gb28181_server.GB28181Device{
			Status: gb28181_server.DeviceOffLineStatus,
		},
	}).Count(&offline)

	model.HttpResponse{}.OkGin(c, &struct {
		Total   int64 `json:"total"`
		Online  int64 `json:"online"`
		Offline int64 `json:"offline"`
	}{
		Total:   total,
		Online:  online,
		Offline: offline,
	})
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

// GetIpcDevicesBySite 根据区域ID查询IpcDevice设备
//
//	@Summary		根据区域ID查询IpcDevice设备
//	@Version		0.0.1
//	@Description	分页查询GB28181 IpcDevice设备
//	@Tags			IPC设备 /ipc/device
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string									false	"访问token"
//	@Param			access_token	query		string									false	"访问token"
//	@Param			site_id			query		integer									true	"区域ID"
//	@Success		200				{object}	model.HttpResponse{data=IpcDevicePage}	"查询成功"
//	@Failure		500				{object}	model.HttpResponse{data=string}			"查询数据失败"
//	@Router			/ipc/device/devices/site_id [GET]
func GetIpcDevicesBySite(c *gin.Context) {
	siteIdQuery := c.Query("site_id")
	if siteIdQuery == "" {
		panic(http.StatusBadRequest)
	}

	siteId, err := strconv.ParseInt(siteIdQuery, 10, 0)
	if err != nil {
		panic(http.StatusBadRequest)
	}

	data, err := selectIpcDevicesBySiteId(&siteId)
	if err != nil {
		model.HttpResponse{}.FailGin(c, "查询数据失败")
		return
	} else {
		model.HttpResponse{}.OkGin(c, data)
	}
}

// PostIpcChannel 更新IPC通道
//
//	@Summary		更新IPC通道
//	@Version		0.0.1
//	@Description	更新IPC通道
//	@Tags			IPC设备 /ipc/device
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string								false	"访问token"
//	@Param			access_token	query		string								false	"访问token"
//	@Param			IpcDevice		body		IpcDevice							true	"设备信息"
//	@Success		200				{object}	model.HttpResponse{data=IpcChannel}	"更新成功"
//	@Failure		500				{object}	model.HttpResponse{data=string}		"查询数据失败"
//	@Router			/ipc/device/channel [POST]
func PostIpcChannel(c *gin.Context) {
	var ipcChannel IpcChannel
	if err := c.ShouldBind(&ipcChannel); err != nil {
		panic(http.StatusBadRequest)
	}
	err := updateIpcChannel(ipcChannel)
	if err != nil {
		model.HttpResponse{}.FailGin(c, err.Error())
	} else {
		model.HttpResponse{}.OkGin(c, ipcChannel)
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
		return
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
