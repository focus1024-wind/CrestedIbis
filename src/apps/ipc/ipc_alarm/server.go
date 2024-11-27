package ipc_alarm

import (
	"CrestedIbis/gb28181_server"
	"CrestedIbis/src/global"
	"fmt"
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
	}).Preload("IpcRecords").Order("id DESC").Offset(int(offset)).Limit(int(pageSize)).Find(&ipcDevices).Error; err != nil {
		return
	}
	return
}

func selectIpcRecordsByPages(page int64, pageSize int64, deviceID string, channelID string, start int64, end int64) (total int64, ipcDevices []IpcRecord, err error) {
	db := global.Db.Model(IpcRecord{})
	var stream string
	if deviceID != "" && channelID != "" {
		stream = fmt.Sprintf("%s/%s", deviceID, channelID)
	}

	if err = db.Where(IpcRecord{
		Record: gb28181_server.Record{
			Stream: stream,
		},
	}).Where("start_time BETWEEN ? AND ? ", start, end).Count(&total).Error; err != nil {
		return
	}

	offset := (page - 1) * pageSize
	if err = db.Debug().Where(IpcRecord{
		Record: gb28181_server.Record{
			Stream: stream,
		},
	}).Order("id DESC").Offset(int(offset)).Limit(int(pageSize)).Find(&ipcDevices).Error; err != nil {
		return
	}
	return
}
