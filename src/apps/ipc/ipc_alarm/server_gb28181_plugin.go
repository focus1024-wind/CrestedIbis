package ipc_alarm

import (
	"CrestedIbis/gb28181_server"
	"CrestedIbis/src/global"
)

func (IpcAlarm) Handler(alarm gb28181_server.Alarm) {
	global.Db.Save(&IpcAlarm{
		Alarm: alarm,
	})
	_, _ = gb28181_server.StartRecordMp4(alarm.DeviceID, alarm.ChannelID, 15)
}
