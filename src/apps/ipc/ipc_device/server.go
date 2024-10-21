package ipc_device

import (
	"CrestedIbis/gb28181_server_back"
	"CrestedIbis/src/apps/ipc"
	"CrestedIbis/src/global"
	"fmt"
)

func (device *IpcDevice) DeviceOffline(string) {

}

func (device *IpcDevice) LoadDevice(deviceId string) (gb28181_server_back.GB28181Device, bool) {
	var gb28181Device IpcDevice
	err := global.Db.Where(&IpcDevice{DeviceID: deviceId}).First(&gb28181Device).Error
	if err != nil {
		return gb28181_server_back.GB28181Device{}, false
	} else {
		gb28181Device.GB28181Device.DeviceID = gb28181Device.DeviceID
		return gb28181Device.GB28181Device, true
	}
}

func (device *IpcDevice) LoadChannel(deviceId string, channelId string) (gb28181_server_back.GB28181Channel, bool) {
	var channel IpcChannel
	err := global.Db.Debug().Where(&IpcChannel{
		ParentID: deviceId,
		DeviceID: channelId,
	}).First(&channel).Error

	if err != nil {
		return gb28181_server_back.GB28181Channel{}, false
	} else {
		channel.GB28181Channel.ParentID = channel.ParentID
		channel.GB28181Channel.DeviceID = channel.DeviceID
		return channel.GB28181Channel, true
	}
}

func (device *IpcDevice) LoadChannels(deviceId string) ([]gb28181_server_back.GB28181Channel, bool) {
	result := make([]IpcChannel, 0)
	err := global.Db.Debug().Where(&IpcChannel{ParentID: deviceId}).Take(&result).Error
	if err != nil {
		return []gb28181_server_back.GB28181Channel{}, false
	} else {
		gb28181Channel := make([]gb28181_server_back.GB28181Channel, 0)
		for _, channel := range result {
			channel.GB28181Channel.ParentID = channel.ParentID
			channel.GB28181Channel.DeviceID = channel.DeviceID
			gb28181Channel = append(gb28181Channel, channel.GB28181Channel)
		}

		return gb28181Channel, true
	}
}

func (device *IpcDevice) StoreDevice(gb28181Device gb28181_server_back.GB28181Device) {
	global.Db.Debug().Create(&IpcDevice{DeviceID: gb28181Device.DeviceID, GB28181Device: gb28181Device})
}

func (device *IpcDevice) RecoverDevice(gb28181Device gb28181_server_back.GB28181Device) {
	global.Db.Debug().Where(
		&IpcDevice{
			DeviceID: gb28181Device.DeviceID,
		},
	).Updates(
		&IpcDevice{
			GB28181Device: gb28181Device,
			DeviceID:      gb28181Device.DeviceID,
		},
	)
}

func (device *IpcDevice) UpdateChannels(channels []gb28181_server_back.GB28181Channel) {
	for _, channel := range channels {
		err := global.Db.Debug().Where(&IpcChannel{
			ParentID: channel.ParentID,
			DeviceID: channel.DeviceID,
		}).Save(&IpcChannel{
			ParentID:       channel.ParentID,
			DeviceID:       channel.DeviceID,
			GB28181Channel: channel,
		}).Error
		if err != nil {
			global.Db.Debug().Where(&IpcChannel{
				ParentID: channel.ParentID,
				DeviceID: channel.DeviceID,
			}).Updates(&IpcChannel{
				ParentID:       channel.ParentID,
				DeviceID:       channel.DeviceID,
				GB28181Channel: channel,
			})
		}
	}
}

func (device *IpcDevice) SnapShotUploadUrl(deviceId string) string {
	accessToken := ipc.GenUploadImageAccessToken(deviceId)
	return fmt.Sprintf("%s/ipc/device/upload_image?access_token=%s", global.Conf.HttpServer.PublicHost, accessToken)
}
