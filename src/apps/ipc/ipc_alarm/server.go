package ipc_alarm

import (
	"CrestedIbis/gb28181_server"
	"CrestedIbis/src/global"
)

func selectIpcAlarmsByPages(page int64, pageSize int64, deviceID string, channelID string, start string, end string) (total int64, ipcDevices []IpcAlarm, err error) {
	db := global.Db.Model(IpcAlarm{})

	if err = db.Where(IpcAlarm{
		Alarm: gb28181_server.Alarm{
			DeviceID:  deviceID,
			ChannelID: channelID,
		},
	}).Where("created_time BETWEEN ? AND ? ", start, end).Count(&total).Error; err != nil {
		return
	}

	offset := (page - 1) * pageSize
	if err = db.Debug().Where(IpcAlarm{
		Alarm: gb28181_server.Alarm{
			DeviceID:  deviceID,
			ChannelID: channelID,
		},
	}).Order("id DESC").Offset(int(offset)).Limit(int(pageSize)).Find(&ipcDevices).Error; err != nil {
		return
	}
	return
}
