package ipc_device

import (
	"CrestedIbis/gb28181_server_back"
	"CrestedIbis/src/global/model"
)

type IpcDevice struct {
	gb28181_server_back.GB28181Device
	ID          int64        `gorm:"primary_key;auto_increment" json:"id"`
	DeviceID    string       `gorm:"uniqueIndex;column:device_id" json:"device_id"`
	IpcChannels []IpcChannel `gorm:"foreignKey:ParentID;references:DeviceID"`
	model.BaseModel
}

type IpcChannel struct {
	gb28181_server_back.GB28181Channel
	ID       int64  `gorm:"primary_key;auto_increment" json:"id"`
	ParentID string `gorm:"uniqueIndex:channel_index"`
	DeviceID string `gorm:"uniqueIndex:channel_index"`
	model.BaseModel
}
