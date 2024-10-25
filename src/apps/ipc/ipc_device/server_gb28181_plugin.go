package ipc_device

import (
	"CrestedIbis/gb28181_server"
	"CrestedIbis/src/global"
	"CrestedIbis/src/global/model"
	"fmt"
	"time"
)

func (IpcDevice) LoadDevice(deviceId string) (gb28181_server.GB28181Device, bool) {
	var gb28181Device IpcDevice
	err := global.Db.Where(&IpcDevice{DeviceID: deviceId}).First(&gb28181Device).Error
	if err != nil {
		return gb28181_server.GB28181Device{DeviceID: deviceId}, false
	} else {
		gb28181Device.GB28181Device.DeviceID = gb28181Device.DeviceID
		gb28181Device.GB28181Device.RegisterTime = time.Time(gb28181Device.RegisterTime)
		gb28181Device.GB28181Device.KeepaliveTime = time.Time(gb28181Device.KeepaliveTime)
		return gb28181Device.GB28181Device, true
	}
}

func (IpcDevice) StoreDevice(gb28181Device gb28181_server.GB28181Device) {
	err := global.Db.Where(&IpcDevice{
		DeviceID: gb28181Device.DeviceID,
	}).Save(&IpcDevice{
		DeviceID:      gb28181Device.DeviceID,
		GB28181Device: gb28181Device,
		RegisterTime:  model.LocalTime(gb28181Device.RegisterTime),
		KeepaliveTime: model.LocalTime(gb28181Device.KeepaliveTime),
	}).Error
	if err != nil {
		global.Db.Where(&IpcDevice{
			DeviceID: gb28181Device.DeviceID,
		}).Updates(&IpcDevice{
			DeviceID:      gb28181Device.DeviceID,
			GB28181Device: gb28181Device,
			RegisterTime:  model.LocalTime(gb28181Device.RegisterTime),
			KeepaliveTime: model.LocalTime(gb28181Device.KeepaliveTime),
		})
	}
}

func (IpcDevice) DeviceOffline(deviceId string) {
	var gb28181Device IpcDevice
	global.Db.Where(&IpcDevice{DeviceID: deviceId}).First(&gb28181Device)
	gb28181Device.Status = gb28181_server.DeviceOffLineStatus
	global.Db.Save(&gb28181Device)
}

func (IpcDevice) LoadChannel(deviceId string, channelId string) (gb28181_server.GB28181Channel, bool) {
	var channel IpcChannel
	err := global.Db.Where(&IpcChannel{
		ParentID: deviceId,
		DeviceID: channelId,
	}).First(&channel).Error

	if err != nil {
		return gb28181_server.GB28181Channel{}, false
	} else {
		channel.GB28181Channel.ParentID = channel.ParentID
		channel.GB28181Channel.DeviceID = channel.DeviceID
		return channel.GB28181Channel, true
	}
}

func (IpcDevice) LoadChannels(deviceId string) ([]gb28181_server.GB28181Channel, bool) {
	result := make([]IpcChannel, 0)
	err := global.Db.Where(&IpcChannel{ParentID: deviceId}).Take(&result).Error
	if err != nil {
		return []gb28181_server.GB28181Channel{}, false
	} else {
		gb28181Channel := make([]gb28181_server.GB28181Channel, 0)
		for _, channel := range result {
			channel.GB28181Channel.ParentID = channel.ParentID
			channel.GB28181Channel.DeviceID = channel.DeviceID
			gb28181Channel = append(gb28181Channel, channel.GB28181Channel)
		}

		return gb28181Channel, true
	}
}

func (IpcDevice) UpdateChannels(channels []gb28181_server.GB28181Channel) {
	// 更新通道信息
	for _, channel := range channels {
		err := global.Db.Where(&IpcChannel{
			ParentID: channel.ParentID,
			DeviceID: channel.DeviceID,
		}).Save(&IpcChannel{
			ParentID:       channel.ParentID,
			DeviceID:       channel.DeviceID,
			GB28181Channel: channel,
		}).Error
		if err != nil {
			global.Db.Where(&IpcChannel{
				ParentID: channel.ParentID,
				DeviceID: channel.DeviceID,
			}).Updates(&IpcChannel{
				ParentID:       channel.ParentID,
				DeviceID:       channel.DeviceID,
				GB28181Channel: channel,
			})
		}
	}

	// 更新通道数
	for _, channel := range channels {
		var count int64
		if err := global.Db.Debug().Model(&IpcChannel{}).Where(&IpcChannel{ParentID: channel.ParentID}).Count(&count).Error; err != nil {
			return
		}
		global.Db.Debug().Model(&IpcDevice{}).Where(&IpcDevice{
			DeviceID: channel.ParentID,
		}).Updates(&IpcDevice{
			ChannelNum: count,
		})
		return
	}
}

func (IpcDevice) SnapShotUploadUrl(deviceId string) string {
	accessToken := GenUploadImageAccessToken(deviceId)
	return fmt.Sprintf("%s/ipc/device/upload_image?access_token=%s", global.Conf.HttpServer.PublicHost, accessToken)
}
