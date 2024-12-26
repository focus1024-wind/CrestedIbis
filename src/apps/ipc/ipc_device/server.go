package ipc_device

import (
	"CrestedIbis/gb28181_server"
	"CrestedIbis/src/apps/site"
	"CrestedIbis/src/global"
	"CrestedIbis/src/utils"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

func (IpcDevice) Select(deviceID string) (device IpcDevice, err error) {
	err = global.Db.Model(&IpcDevice{}).Where(&IpcDevice{DeviceID: deviceID}).First(&device).Error
	return
}

func (IpcDevice) Update(device IpcDevice) (err error) {
	if device.DeviceID == "" {
		return errors.New("device_id 参数为空")
	}
	err = global.Db.Model(&IpcDevice{}).Where(&IpcDevice{
		DeviceID: device.DeviceID,
	}).Updates(&device).Error
	if err != nil {
		global.Logger.Errorf("更新 %s 设备失败：%s", device.DeviceID, err.Error())
		return errors.New(fmt.Sprintf("更新 %s 设备失败：%s", device.DeviceID, err.Error()))
	}
	return
}

func (IpcDevice) Delete(deviceID string) (err error) {
	// 删除对应通道
	err = global.Db.Model(&IpcChannel{}).Where(&IpcChannel{
		ParentID: deviceID,
	}).Delete(&IpcChannel{}).Error
	if err != nil {
		return
	}

	err = global.Db.Model(&IpcDevice{}).Where(&IpcDevice{
		DeviceID: deviceID,
	}).Delete(&IpcDevice{}).Error
	return
}

// SelectIpcDevices 分页搜索IpcDevices
func (IpcDevice) SelectIpcDevices(page int64, pageSize int64, status string, keywords string) (total int64, ipcDevices []IpcDevice, err error) {
	db := global.Db.Model(IpcDevice{}).Preload("IpcChannels").Preload("Site", site.ExpandSitePreload)

	if status == gb28181_server.DeviceOnLineStatus || status == gb28181_server.DeviceOffLineStatus {
		db = db.Where("status LIKE ?", status)
	}
	if keywords != "" {
		db = db.Where("device_id LIKE ?", "%"+keywords+"%")
	}

	db = db.Session(&gorm.Session{})

	if err = db.Count(&total).Error; err != nil {
		return
	} else {
		offset := (page - 1) * pageSize
		err = db.Order("id").Offset(int(offset)).Limit(int(pageSize)).Find(&ipcDevices).Error
		return
	}
}

func (IpcDevice) Deletes(deviceIDs []string) (err error) {
	var errDeleteDeviceIds []string
	for _, deviceID := range deviceIDs {
		err = IpcDevice{}.Delete(deviceID)
		if err != nil {
			errDeleteDeviceIds = append(errDeleteDeviceIds, deviceID)
			err = nil
		}
	}

	if len(errDeleteDeviceIds) > 0 {
		return errors.New(fmt.Sprintf("%s 删除失败", errDeleteDeviceIds))
	}
	return
}

func (IpcDevice) SelectBySiteID(siteID *int64) (ipcDevices []IpcDevice, err error) {
	err = global.Db.Model(&IpcDevice{}).Where(&IpcDevice{
		SiteId: siteID,
	}).Preload("IpcChannels").Preload("Site", site.ExpandSitePreload).
		Order("id").Find(&ipcDevices).Error
	return
}

func (IpcChannel) Update(channel IpcChannel) (err error) {
	if channel.ParentID == "" || channel.DeviceID == "" {
		return errors.New("parent_id 或 device_id 为空")
	}
	err = global.Db.Model(&IpcChannel{}).Where(&IpcChannel{
		ParentID: channel.ParentID,
		DeviceID: channel.DeviceID,
	}).Select("ptz_type").Updates(&channel).Error
	if err != nil {
		global.Logger.Errorf("更新 %s/%s 设备失败：%s", channel.ParentID, channel.DeviceID, err.Error())
		return errors.New("update device error")
	}
	return
}

func (IpcChannel) SelectChannels(deviceID string) (channels []IpcChannel, err error) {
	err = global.Db.Where(&IpcChannel{ParentID: deviceID}).Find(&channels).Error
	return
}

func GenUploadImageAccessToken(deviceId string) string {
	token, err := utils.JwtToken{}.GenTempAccessToken(deviceId, []string{"ipc_device"}, 180)
	if err != nil {
		return ""
	} else {
		return token
	}
}
