package ipc_device

import (
	"CrestedIbis/src/global"
	"CrestedIbis/src/utils"
	"fmt"
)

func selectIpcDevice(deviceID string) (device IpcDevice, err error) {
	err = global.Db.Model(&IpcDevice{}).Where(&IpcDevice{DeviceID: deviceID}).First(&device).Error
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
	if err = db.Preload("IpcChannels").Order("id DESC").Offset(int(offset)).Limit(int(pageSize)).Find(&ipcDevices).Error; err != nil {
		return
	}
	return
}

func selectIpcChannels(deviceID string) (channels []IpcChannel, err error) {
	err = global.Db.Where(&IpcChannel{ParentID: deviceID}).Take(&channels).Error
	return
}

func GenUploadImageAccessToken(deviceId string) string {
	token, err := utils.JwtToken{}.GenTempAccessToken(deviceId, []string{"ipc_device"}, 180)
	if err != nil {
		fmt.Println(err)
		return ""
	} else {
		return token
	}
}
