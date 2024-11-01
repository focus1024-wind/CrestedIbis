package gb28181_server

import (
	"fmt"
	"time"
)

// Play 实时音视频点播
func Play(deviceID string, channelID string) map[string]string {
	var streamPath = fmt.Sprintf("%s/%s", deviceID, channelID)

	if exist, err := ApiClientGetRtpInfo(streamPath); err == nil && exist {
		logger.Info("[Stream] 已存在码流, streamPath", streamPath)
	} else {
		channel, exist := GlobalGB28181DeviceStore.LoadChannel(deviceID, channelID)
		if exist {
			_ = channel.Invite(&InviteOptions{})
			// 等待流注册完毕或流注册超时后返回
			ticker := time.NewTicker(time.Second)
			timeout := time.After(3 * time.Second)
			stream, _ := PublishStore.Load(streamPath)
			for stream == nil {
				select {
				case <-ticker.C:
					stream, _ = PublishStore.Load(streamPath)
				case <-timeout:
					// break在for-select中无法跳出，采用goto跳出多层循环
					goto Loop
				}
			}
		}
	}
Loop:
	return GetMediaPlayUrl(streamPath)
}

// PlayStop 停止音视频点播
func PlayStop(deviceID string, channelID string) {
	if channel, ok := GlobalGB28181DeviceStore.LoadChannel(deviceID, channelID); ok {
		channel.Bye()
	}
}
