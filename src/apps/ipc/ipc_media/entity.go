package ipc_media

import "time"

type IpcMediaPlayModel struct {
	DeviceId  string    `json:"device_id" desc:"设备ID"`
	ChannelId string    `json:"channel_id" desc:"通道ID"`
	Start     time.Time `json:"start"`
	End       time.Time `json:"end"`
}
