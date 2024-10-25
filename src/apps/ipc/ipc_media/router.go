package ipc_media

import "CrestedIbis/src/global"

func InitIpcMediaRouter() {
	ipcDeviceRouter := global.HttpEngine.Group("/ipc/media")
	{
		ipcDeviceRouter.POST("/play", IpcMediaPlay)
	}
}
