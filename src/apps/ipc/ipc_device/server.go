package ipc_device

import (
	"CrestedIbis/src/global"
	"CrestedIbis/src/utils"
	"errors"
)

func selectIpcDevice(deviceID string) (device IpcDevice, err error) {
	err = global.Db.Model(&IpcDevice{}).Where(&IpcDevice{DeviceID: deviceID}).First(&device).Error
	return
}

func updateIpcDevice(device IpcDevice) (err error) {
	if device.DeviceID == "" {
		return errors.New("device_id is empty")
	}
	err = global.Db.Debug().Model(&IpcDevice{}).Where(&IpcDevice{
		DeviceID: device.DeviceID,
	}).Updates(&device).Error
	if err != nil {
		global.Logger.Errorf("更新 %s 设备失败：%s", device.DeviceID, err.Error())
		return errors.New("update device error")
	}
	return
}

func deleteIpcDevice(deviceID string) (err error) {
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

// selectIpcDeviceByPages 分页搜索IpcDevices
// page: 页码，pageSize: 每页的数量
func selectIpcDevicesByPages(page int64, pageSize int64) (total int64, ipcDevices []IpcDevice, err error) {
	db := global.Db.Model(IpcDevice{})

	if err = db.Count(&total).Error; err != nil {
		return
	}

	offset := (page - 1) * pageSize
	if err = db.Preload("IpcChannels").Preload("Site").Order("id DESC").Offset(int(offset)).Limit(int(pageSize)).Find(&ipcDevices).Error; err != nil {
		return
	}
	return
}

func updateIpcChannel(channel IpcChannel) (err error) {
	if channel.ParentID == "" || channel.DeviceID == "" {
		return errors.New("parent_id or device_id is empty")
	}
	err = global.Db.Debug().Model(&IpcChannel{}).Where(&IpcChannel{
		ParentID: channel.ParentID,
		DeviceID: channel.DeviceID,
	}).Updates(&channel).Error
	if err != nil {
		global.Logger.Errorf("更新 %s/%s 设备失败：%s", channel.ParentID, channel.DeviceID, err.Error())
		return errors.New("update device error")
	}
	return
}

func selectIpcChannels(deviceID string) (channels []IpcChannel, err error) {
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
