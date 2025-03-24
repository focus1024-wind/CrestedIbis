package ipc_device

import (
	"CrestedIbis/gb28181_server"
	"CrestedIbis/src/global"
	"CrestedIbis/src/global/model"
	"CrestedIbis/src/utils"
	"fmt"
	"net/http"
	"path"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// GetIpcDevice 获取IPC设备
//
//	@Summary		获取IPC设备
//	@Version		0.0.1
//	@Description	获取IPC设备
//	@Tags			IPC设备
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
		device, err := IpcDevice{}.Select(deviceID)
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
//	@Tags			IPC设备
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string							false	"访问token"
//	@Param			access_token	query		string							false	"访问token"
//	@Param			IpcDevice		body		IpcDevice						true	"设备信息"
//	@Success		200				{object}	model.HttpResponse{data=nil}	"更新成功"
//	@Failure		500				{object}	model.HttpResponse{data=string}	"查询数据失败"
//	@Router			/ipc/device [POST]
func PostIpcDevice(c *gin.Context) {
	var ipcDevice IpcDevice
	if err := c.ShouldBind(&ipcDevice); err != nil {
		panic(http.StatusBadRequest)
	}
	err := IpcDevice{}.Update(ipcDevice)
	if err != nil {
		model.HttpResponse{}.FailGin(c, err.Error())
	} else {
		model.HttpResponse{}.OkGin(c, nil)
	}
}

// DeleteIpcDevice 删除IPC设备
//
//	@Summary		删除IPC设备
//	@Version		0.0.1
//	@Description	删除IPC设备及对应通道，该删除仅为删除数据库记录，不影响IPC设备的重新注册
//	@Tags			IPC设备
//	@Accept			json
//	@Produce		json
//	@Param			Authorization		header		string							false	"访问token"
//	@Param			access_token		query		string							false	"访问token"
//	@Param			IpcDeviceIDModel	body		IpcDeviceIDModel				true	"设备ID"
//	@Success		200					{object}	model.HttpResponse{data=nil}	"删除IPC设备成功"
//	@Failure		500					{object}	model.HttpResponse{data=string}	"删除IPC设备失败"
//	@Router			/ipc/device [DELETE]
func DeleteIpcDevice(c *gin.Context) {
	var deviceID IpcDeviceIDModel
	if err := c.ShouldBind(&deviceID); err != nil {
		panic(http.StatusBadRequest)
	}

	err := IpcDevice{}.Delete(deviceID.DeviceID)
	if err != nil {
		model.HttpResponse{}.FailGin(c, "删除设备失败")
	} else {
		model.HttpResponse{}.OkGin(c, nil)
	}
}

// GetIpcDevices 获取IPC设备列表
//
//	@Summary		获取IPC设备列表
//	@Version		0.0.1
//	@Description	获取IPC设备列表
//	@Tags			IPC设备
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string																false	"访问token"
//	@Param			access_token	query		string																false	"访问token"
//	@Param			page			query		integer																false	"分页查询页码，默认值: 1"
//	@Param			page_size		query		integer																false	"每页查询数量，默认值: 15"
//	@Param			status			query		string																false	"设备状态，支持: ALl、ON、OFF，默认值: ALL"
//	@Param			keywords		query		string																false	"设备模型查询信息"
//	@Success		200				{object}	model.HttpResponse{data=model.BasePageResponse{data=[]IpcDevice}}	"获取IPC设备列表成功"
//	@Failure		500				{object}	model.HttpResponse{data=string}										"获取IPC设备列表失败"
//	@Router			/ipc/device/devices [GET]
func GetIpcDevices(c *gin.Context) {
	pageQuery := c.DefaultQuery("page", "1")
	pageSizeQuery := c.DefaultQuery("page_size", "15")
	statusQuery := c.DefaultQuery("status", "ALL")
	keywords := c.DefaultQuery("keywords", "")

	page, err := strconv.ParseInt(pageQuery, 10, 0)
	if err != nil {
		panic(http.StatusBadRequest)
	}
	pageSize, err := strconv.ParseInt(pageSizeQuery, 10, 0)
	if err != nil {
		panic(http.StatusBadRequest)
	}

	total, data, err := IpcDevice{}.SelectIpcDevices(page, pageSize, statusQuery, keywords)
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

// DeleteIpcDevices 批量删除IPC设备
//
//	@Summary		批量删除IPC设备
//	@Version		0.0.1
//	@Description	批量删除IPC设备及对应通道，该删除仅为删除数据库记录，不影响IPC设备的重新注册
//	@Tags			IPC设备
//	@Accept			json
//	@Produce		json
//	@Param			Authorization		header		string							false	"访问token"
//	@Param			access_token		query		string							false	"访问token"
//	@Param			IpcDeviceIDsModel	body		IpcDeviceIDsModel				true	"设备ID列表"
//	@Success		200					{object}	model.HttpResponse{data=nil}	"批量删除IPC设备成功"
//	@Failure		500					{object}	model.HttpResponse{data=string}	"批量删除IPC设备失败"
//	@Router			/ipc/device/devices [DELETE]
func DeleteIpcDevices(c *gin.Context) {
	var idsModel IpcDeviceIDsModel
	if err := c.ShouldBind(&idsModel); err != nil {
		panic(http.StatusBadRequest)
	}

	err := IpcDevice{}.Deletes(idsModel.DeviceIDs)
	if err != nil {
		model.HttpResponse{}.FailGin(c, err.Error())
	} else {
		model.HttpResponse{}.OkGin(c, nil)
	}
}

// GetIpcDevicesBySiteID 根据区域ID查询Ipc设备
//
//	@Summary		根据区域ID查询Ipc设备
//	@Version		0.0.1
//	@Description	根据区域ID查询Ipc设备
//	@Tags			IPC设备
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string									false	"访问token"
//	@Param			access_token	query		string									false	"访问token"
//	@Param			site_id			query		integer									true	"区域ID"
//	@Success		200				{object}	model.HttpResponse{data=[]IpcDevice}	"查询成功"
//	@Failure		500				{object}	model.HttpResponse{data=string}			"查询数据失败"
//	@Router			/ipc/device/devices/site_id [GET]
func GetIpcDevicesBySiteID(c *gin.Context) {
	siteIdQuery := c.Query("site_id")
	if siteIdQuery == "" {
		panic(http.StatusBadRequest)
	}

	siteId, err := strconv.ParseInt(siteIdQuery, 10, 0)
	if err != nil {
		panic(http.StatusBadRequest)
	}

	data, err := IpcDevice{}.SelectBySiteID(&siteId)
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
//	@Tags			IPC设备通道
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string							false	"访问token"
//	@Param			access_token	query		string							false	"访问token"
//	@Param			IpcChannel		body		IpcChannel						true	"设备通道信息"
//	@Success		200				{object}	model.HttpResponse{data=nil}	"更新IPC通道成功"
//	@Failure		500				{object}	model.HttpResponse{data=string}	"更新IPC通道失败"
//	@Router			/ipc/channel [POST]
func PostIpcChannel(c *gin.Context) {
	var ipcChannel IpcChannel
	if err := c.ShouldBind(&ipcChannel); err != nil {
		panic(http.StatusBadRequest)
	}
	err := IpcChannel{}.Update(ipcChannel)
	if err != nil {
		model.HttpResponse{}.FailGin(c, err.Error())
	} else {
		model.HttpResponse{}.OkGin(c, nil)
	}
}

// GetIpcChannels 获取设备通道列表
//
//	@Summary		获取设备通道列表
//	@Version		0.0.1
//	@Description	获取设备通道列表
//	@Tags			IPC设备通道
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string									false	"访问token"
//	@Param			access_token	query		string									false	"访问token"
//	@Param			device_id		query		string									true	"设备ID"
//	@Success		200				{object}	model.HttpResponse{data=[]IpcChannel}	"查询数据成功"
//	@Failure		500				{object}	model.HttpResponse{data=string}			"查询数据失败"
//	@Router			/ipc/channel/channels [GET]
func GetIpcChannels(c *gin.Context) {
	deviceId := c.Query("device_id")
	if deviceId == "" {
		panic(http.StatusBadRequest)
	}

	ipcChannels, err := IpcChannel{}.SelectChannels(deviceId)

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
//	@Tags			IPC设备通道
//	@Accept			mpfd
//	@Produce		json
//	@Param			access_token	query		string							true	"访问token"
//	@Param			file			formData	file							true	"上传图片"
//	@Success		200				{object}	model.HttpResponse{data=string}	"上传图片成功"
//	@Failure		500				{object}	model.HttpResponse{data=string}	"上传图片失败"
//	@Router			/ipc/channel/upload_image [POST]
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

// ControlPTZ PTZ控制
//
//	@Summary		PTZ控制
//	@Version		0.0.1
//	@Description	PTZ控制接口
//	@Tags			IPC控制
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string								false	"访问token"
//	@Param			access_token	query		string								false	"访问token"
//	@Param			channel_id		query		string								true	"通道ID"
//	@Param			options			body		gb28181_server.PTZControlOptions	true	"PTZ控制参数"
//	@Success		200				{object}	model.HttpResponse{data=string}		"PTZ控制成功"
//	@Failure		500				{object}	model.HttpResponse{data=string}		"PTZ控制失败"
//	@Router			/ipc/ptz [POST]
func ControlPTZ(c *gin.Context) {
	channelID := c.Query("channel_id")
	if channelID == "" {
		model.HttpResponse{}.FailGin(c, "通道ID不能为空")
		return
	}
	var options gb28181_server.PTZControlOptions
	if err := c.ShouldBind(&options); err != nil {
		panic(http.StatusBadRequest)
	}

	_, err := gb28181_server.ControlPTZ(channelID, &options)
	if err != nil {
		model.HttpResponse{}.FailGin(c, err.Error())
	} else {
		model.HttpResponse{}.OkGin(c, nil)
	}
}
