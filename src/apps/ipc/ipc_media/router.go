package ipc_media

import "CrestedIbis/src/global"

func InitIpcMediaRouter() {
	ipcMediaRouter := global.HttpEngine.Group("/ipc/media")
	{
		ipcMediaRouter.POST("/play", IpcMediaPlay)
		ipcMediaRouter.POST("/stop", IpcMediaStop)
	}
}
