package ipc

import "CrestedIbis/src/global"

func InitIpcRouter() {
	ipcDeviceRouter := global.HttpEngine.Group("/ipc")
	{
		ipcDeviceRouter.POST("/device/upload_image", IpcUploadImage)
	}
}
