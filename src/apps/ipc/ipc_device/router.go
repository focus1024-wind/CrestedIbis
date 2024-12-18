package ipc_device

import "CrestedIbis/src/global"

func InitIpcDeviceRouter() {
	ipcDeviceRouter := global.HttpEngine.Group("/ipc/device")
	{
		ipcDeviceRouter.GET("", GetIpcDevice)
		ipcDeviceRouter.POST("", PostIpcDevice)
		ipcDeviceRouter.DELETE("", DeleteIpcDevice)
		ipcDeviceRouter.DELETE("/devices", DeleteIpcDevices)
		ipcDeviceRouter.GET("/status", GetIpcDevicesStatus)
		ipcDeviceRouter.GET("/devices", GetIpcDevicesByPages)
		ipcDeviceRouter.GET("/devices/site_id", GetIpcDevicesBySite)
		ipcDeviceRouter.POST("/channel", PostIpcChannel)
		ipcDeviceRouter.GET("/channels", GetIpcChannels)
		ipcDeviceRouter.POST("/upload_image", IpcUploadImage)
	}
}
