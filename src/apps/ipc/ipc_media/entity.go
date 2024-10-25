package ipc_media

import "time"

type IpcMediaPlayModel struct {
	DeviceId  string    `json:"device_id"`
	ChannelId string    `json:"channel_id"`
	Start     time.Time `json:"start"`
	End       time.Time `json:"end"`
}
