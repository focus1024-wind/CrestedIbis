package gb28181_server

import (
	"fmt"
	"time"
)

// Play 实时音视频点播
func Play(deviceID string, channelID string) map[string]string {
	var streamPath = fmt.Sprintf("%s/%s", deviceID, channelID)
	stream, _ := PublishStore.Load(streamPath)

	if stream != nil {
		// 流已存在，不重复拉流
		logger.Info("[Stream] 已存在码流, streamPath", streamPath)
	} else {
		channel, exist := GlobalGB28181DeviceStore.LoadChannel(deviceID, channelID)
		if exist {
			_ = channel.Invite(&InviteOptions{})
			time.Sleep(3 * time.Second)
		}
	}
	return GetMediaPlayUrl(streamPath)
}

// PlayStop 停止音视频点播
func PlayStop(deviceID string, channelID string) {
	if channel, ok := GlobalGB28181DeviceStore.LoadChannel(deviceID, channelID); ok {
		channel.Bye()
	}
}
