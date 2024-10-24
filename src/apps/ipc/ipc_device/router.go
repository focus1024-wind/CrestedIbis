package ipc_device

import "CrestedIbis/src/global"

func InitIpcDeviceRouter() {
	ipcDeviceRouter := global.HttpEngine.Group("/ipc/device")
	{
		ipcDeviceRouter.GET("/devices", GetIpcDevicesByPages)
		ipcDeviceRouter.GET("/channels", GetIpcChannels)
		ipcDeviceRouter.POST("/upload_image", IpcUploadImage)
	}
}
