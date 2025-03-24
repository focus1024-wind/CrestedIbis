package ipc_device

import "CrestedIbis/src/global"

func InitIpcDeviceRouter() {
	ipcDeviceRouter := global.HttpEngine.Group("/ipc/device")
	{
		ipcDeviceRouter.GET("", GetIpcDevice)
		ipcDeviceRouter.POST("", PostIpcDevice)
		ipcDeviceRouter.DELETE("", DeleteIpcDevice)
		ipcDeviceRouter.GET("/devices", GetIpcDevices)
		ipcDeviceRouter.DELETE("/devices", DeleteIpcDevices)
		ipcDeviceRouter.GET("/devices/site_id", GetIpcDevicesBySiteID)
	}

	ipcChannelRouter := global.HttpEngine.Group("/ipc/channel")
	{
		ipcChannelRouter.POST("", PostIpcChannel)
		ipcChannelRouter.GET("/channels", GetIpcChannels)
		ipcChannelRouter.POST("/upload_image", IpcUploadImage)
	}
	global.HttpEngine.POST("/ipc/ptz", ControlPTZ)
}
