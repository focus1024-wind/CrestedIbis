package ipc_alarm

import "CrestedIbis/src/global"

func InitIpcAlarmRouter() {
	global.HttpEngine.GET("/ipc/device/alarms", GetIpcAlarms)
}
