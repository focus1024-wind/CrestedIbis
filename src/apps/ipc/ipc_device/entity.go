package ipc_device

import (
	"CrestedIbis/gb28181_server"
	"CrestedIbis/src/global/model"
)

type IpcDevice struct {
	gb28181_server.GB28181Device
	ID            int64           `gorm:"primary_key;auto_increment" json:"id"`
	DeviceID      string          `gorm:"uniqueIndex;column:device_id" json:"device_id"`
	IpcChannels   []IpcChannel    `gorm:"foreignKey:ParentID;references:DeviceID"`
	ChannelNum    int64           `json:"channel_num" desc:"设备通道数"`
	RegisterTime  model.LocalTime `json:"register_time" desc:"设备最新注册时间"`
	KeepaliveTime model.LocalTime `json:"keepalive_time" desc:"设备最新心跳时间"`
	model.BaseModel
}

type IpcChannel struct {
	gb28181_server.GB28181Channel
	ID       int64  `gorm:"primary_key;auto_increment" json:"id"`
	ParentID string `gorm:"uniqueIndex:channel_index" json:"parent_id"`
	DeviceID string `gorm:"uniqueIndex:channel_index" json:"device_id"`
	model.BaseModel
}
