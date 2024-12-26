package ipc_alarm

import (
	"CrestedIbis/gb28181_server"
	"CrestedIbis/src/global"
	"CrestedIbis/src/global/model"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

func (IpcAlarm) SelectAlarms(page int64, pageSize int64, deviceID string, channelID string, start string, end string, keywords string) (total int64, ipcAlarms []IpcAlarm, err error) {
	db := global.Db.Model(&IpcAlarm{})

	if deviceID != "" {
		db = db.Where("device_id=?", deviceID)
	}
	if channelID != "" {
		db = db.Where("channel_id=?", channelID)
	}
	if start != "" && end != "" {
		db = db.Where("created_time BETWEEN ? AND ? ", start, end)
	}
	if keywords != "" {
		db = db.Where("device_id LIKE ? OR channel_id LIKE ?", "%"+keywords+"%", "%"+keywords+"%")
	}

	db = db.Session(&gorm.Session{})

	if err = db.Count(&total).Error; err != nil {
		return
	} else {
		offset := (page - 1) * pageSize
		err = db.Preload("IpcRecords").Order("id DESC").Offset(int(offset)).Limit(int(pageSize)).Find(&ipcAlarms).Error
		return
	}
}

func (IpcAlarm) Delete(idModel model.IDModel) (err error) {
	var ipcAlarm IpcAlarm
	if err = global.Db.Model(&IpcAlarm{}).Preload("IpcRecords").Where("id = ?", idModel.ID).First(&ipcAlarm).Error; err != nil {
		return
	}

	for _, record := range ipcAlarm.IpcRecords {
		// 删除对应历史回放信息
		err = IpcRecord{}.Delete(model.IDModel{ID: record.ID})
		if err != nil {
			global.Logger.Errorf("删除告警信息对应告警视频失败: %s", err.Error())
			return
		}
	}

	err = global.Db.Delete(&ipcAlarm).Error
	return
}

func (IpcAlarm) Deletes(ids model.IDsModel) (err error) {
	var errDeleteAlarmIds []int64

	var ipcAlarms []IpcAlarm
	if err = global.Db.Model(&IpcAlarm{}).Preload("IpcRecords").Find(&ipcAlarms, ids).Error; err != nil {
		return
	}

	for _, ipcAlarm := range ipcAlarms {
		if len(ipcAlarm.IpcRecords) > 0 {
			// 存在录像视频的提前删除
			err = IpcAlarm{}.Delete(model.IDModel{ID: ipcAlarm.ID})
			if err != nil {
				errDeleteAlarmIds = append(errDeleteAlarmIds, ipcAlarm.ID)
				err = nil
			}
		}
	}

	if len(errDeleteAlarmIds) > 0 {
		return errors.New(fmt.Sprintf("%v 删除失败", errDeleteAlarmIds))
	}

	// 批量删除
	err = global.Db.Delete(&ipcAlarms).Error

	return
}

func (IpcRecord) Delete(idModel model.IDModel) (err error) {
	var ipcRecord IpcRecord
	err = global.Db.Model(&IpcRecord{}).Where(&IpcRecord{
		IDModel: idModel,
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
	_, err = gb28181_server.DelRecord(app, stream, period, name)
	if err != nil {
		global.Logger.Error("删除录像视频失败")
		return
	}

	// 删除数据库对应记录
	return global.Db.Model(&IpcRecord{}).Where(&IpcRecord{
		IDModel: idModel,
	}).Delete(&IpcRecord{}).Error
}

func (IpcRecord) SelectRecords(page int64, pageSize int64, deviceID string, channelID string, start string, end string, keywords string) (total int64, ipcRecords []IpcRecord, err error) {
	db := global.Db.Model(&IpcRecord{})

	if deviceID != "" && channelID != "" {
		db = db.Where("stream=?", fmt.Sprintf("%s/%s", deviceID, channelID))
	}
	if start != "" && end != "" {
		db = db.Where("created_time BETWEEN ? AND ? ", start, end)
	}
	if keywords != "" {
		db = db.Where("stream LIKE ?", "%"+keywords+"%")
	}

	db = db.Session(&gorm.Session{})

	if err = db.Count(&total).Error; err != nil {
		return
	} else {
		offset := (page - 1) * pageSize
		err = db.Order("id DESC").Offset(int(offset)).Limit(int(pageSize)).Find(&ipcRecords).Error
		return
	}
}
