package ipc

import (
	"CrestedIbis/src/apps/ipc/ipc_device"
	"CrestedIbis/src/apps/ipc/ipc_media"
)

func InitIpcRouter() {
	ipc_device.InitIpcDeviceRouter()
	ipc_media.InitIpcMediaRouter()
}
