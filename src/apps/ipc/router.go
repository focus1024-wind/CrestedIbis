package ipc

import (
	"CrestedIbis/src/apps/ipc/ipc_device"
	"CrestedIbis/src/global"
)

func InitIpcRouter() {
	ipc_device.InitIpcDeviceRouter()
	ipcRouter := global.HttpEngine.Group("/ipc")
	{
		ipcRouter.POST("/device/upload_image", IpcUploadImage)
	}
}
