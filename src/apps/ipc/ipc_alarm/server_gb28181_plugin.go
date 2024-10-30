package ipc_alarm

import (
	"CrestedIbis/gb28181_server"
	"CrestedIbis/src/global"
	"fmt"
	"strconv"
	"strings"
)

func (IpcAlarm) Handler(alarm gb28181_server.Alarm) {
	var ipcAlarm = IpcAlarm{
		Alarm: alarm,
	}
	global.Db.Create(&ipcAlarm)
	_, _ = gb28181_server.StartRecordMp4(fmt.Sprintf("alarm_%d", ipcAlarm.ID), alarm.DeviceID, alarm.ChannelID, 15)
}

func (IpcAlarm) RecordHandler(record gb28181_server.Record) {
	if record.TimeLen < 3 {
		return
	}

	recordApps := strings.Split(record.App, "alarm_")
	if len(recordApps) == 2 {
		alarmID, err := strconv.ParseInt(recordApps[1], 10, 0)
		if err != nil {
			return
		}
		var ipcRecord = IpcRecord{
			AlarmID: &alarmID,
			Record:  record,
		}
		global.Db.Create(&ipcRecord)
		return
	} else {
		var ipcRecord = IpcRecord{
			AlarmID: nil,
			Record:  record,
		}
		global.Db.Create(&ipcRecord)
		return
	}
}
