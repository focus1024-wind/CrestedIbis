package ipc

import (
	"CrestedIbis/src/apps/ipc/ipc_alarm"
	"CrestedIbis/src/apps/ipc/ipc_device"
	"CrestedIbis/src/apps/ipc/ipc_media"
	"CrestedIbis/src/global"
)

func InitIpcRouter() {
	ipc_alarm.InitIpcAlarmRouter()
	ipc_device.InitIpcDeviceRouter()
	ipc_media.InitIpcMediaRouter()
	global.HttpEngine.GET("/ipc/gb28181_info", GB28181InfoApi)
}
