package ipc_alarm

import (
	"CrestedIbis/gb28181_server"
	"CrestedIbis/src/global"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

func selectIpcAlarmsByPages(page int64, pageSize int64, deviceID string, channelID string, start string, end string, keywords string) (total int64, ipcDevices []IpcAlarm, err error) {
	db := global.Db.Model(IpcAlarm{})

	if err = db.Where(IpcAlarm{
		Alarm: gb28181_server.Alarm{
			DeviceID:  deviceID,
			ChannelID: channelID,
		},
	}).Where("created_time BETWEEN ? AND ? ", start, end).
		Where("device_id LIKE ? OR channel_id LIKE ?", "%"+keywords+"%", "%"+keywords+"%").
		Count(&total).Error; err != nil {
		return
	}

	offset := (page - 1) * pageSize
	if err = db.Debug().Where(IpcAlarm{
		Alarm: gb28181_server.Alarm{
			DeviceID:  deviceID,
			ChannelID: channelID,
		},
	}).
		Where("device_id LIKE ? OR channel_id LIKE ?", "%"+keywords+"%", "%"+keywords+"%").
		Preload("IpcRecords").Order("id DESC").Offset(int(offset)).Limit(int(pageSize)).Find(&ipcDevices).Error; err != nil {
		return
	}
	return
}

func DeleteIpcAlarmById(id int64) (err error) {
	var ipcAlarm IpcAlarm
	if err = global.Db.Model(&IpcAlarm{}).Preload("IpcRecords").Where("id = ?", id).First(&ipcAlarm).Error; err != nil {
		return
	}

	for _, record := range ipcAlarm.IpcRecords {
		fmt.Println(record.ID)
		// 删除对应历史回放信息
		DeleteRecordServer(record.ID)
	}

	err = global.Db.Delete(&ipcAlarm).Error
	return
}

func DeleteIpcAlarmByIds(ids []int64) (err error) {
	var ipcAlarms []IpcAlarm
	if err = global.Db.Model(&IpcAlarm{}).Preload("IpcRecords").Find(&ipcAlarms, ids).Error; err != nil {
		return
	}

	for _, ipcAlarm := range ipcAlarms {
		if len(ipcAlarm.IpcRecords) > 0 {
			// 存在录像视频的提前删除
			DeleteIpcAlarmById(ipcAlarm.ID)
		}
	}

	// 批量删除
	err = global.Db.Delete(&ipcAlarms).Error

	return
}

func selectIpcRecordsByPages(page int64, pageSize int64, deviceID string, channelID string, start int64, end int64, keywords string) (total int64, ipcDevices []IpcRecord, err error) {
	db := global.Db.Model(IpcRecord{})
	var stream string
	if deviceID != "" && channelID != "" {
		stream = fmt.Sprintf("%s/%s", deviceID, channelID)
	}

	if err = db.Debug().Where(IpcRecord{
		Record: gb28181_server.Record{
			Stream: stream,
		},
	}).Where("start_time BETWEEN ? AND ? ", start, end).Where("stream LIKE ?", "%"+keywords+"%").Count(&total).Error; err != nil {
		return
	}

	offset := (page - 1) * pageSize
	if err = db.Debug().Where(IpcRecord{
		Record: gb28181_server.Record{
			Stream: stream,
		},
	}).Where("start_time BETWEEN ? AND ? ", start, end).Where("stream LIKE ?", "%"+keywords+"%").Order("id DESC").Offset(int(offset)).Limit(int(pageSize)).Find(&ipcDevices).Error; err != nil {
		return
	}
	return
}

func DeleteRecordServer(id int64) (err error) {
	var ipcRecord IpcRecord
	err = global.Db.Model(&IpcRecord{}).Where(&IpcRecord{
		ID: id,
	}).First(&ipcRecord).Error
	if err != nil {
		return errors.New("数据不存在")
	}
	var (
		app    = ipcRecord.App
		stream = ipcRecord.Stream
		period string
		name   string
	)

	// 判断是删除文件还是删除文件夹
	// 以alarm开头：删除文件夹，其他：删除文件
	if strings.HasPrefix(app, "alarm") {
		name = ""
	} else {
		name = ipcRecord.FileName
	}

	// 获取period信息：文件日期
	parts := strings.Split(ipcRecord.Url, "/")
	if len(parts) > 1 {
		period = parts[len(parts)-2]
	} else {
		panic(http.StatusInternalServerError)
	}

	// 删除磁盘上录像文件
	gb28181_server.DelRecord(app, stream, period, name)

	// 删除数据库对应记录
	global.Db.Model(&IpcRecord{}).Where(&IpcRecord{
		ID: ipcRecord.ID,
	}).Delete(&IpcRecord{})

	return
}
