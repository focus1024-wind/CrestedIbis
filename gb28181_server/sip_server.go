package gb28181_server

import (
	"fmt"
	"github.com/ghettovoice/gosip"
	"github.com/ghettovoice/gosip/sip"
	"strings"
	"time"
)

var (
	globalSipServer gosip.Server
)

// startSipServer 启动SIP服务器
func (config *GB28181Config) startSipServer() {
	sipAddr := fmt.Sprintf("%s:%d", config.SipServer.IP, config.SipServer.Port)

	globalSipServer = gosip.NewServer(gosip.ServerConfig{Host: config.SipServer.IP}, nil, nil, logger)

	_ = globalSipServer.OnRequest(sip.REGISTER, config.SipRegisterHandler)
	_ = globalSipServer.OnRequest(sip.MESSAGE, config.SipMessageHandler)
	_ = globalSipServer.OnRequest(sip.NOTIFY, config.SipNotifyHandler)
	_ = globalSipServer.OnRequest(sip.BYE, config.SipByeHandler)

	err := globalSipServer.Listen(strings.ToLower(config.SipServer.Mode), sipAddr)
	if err != nil {
		logger.Error(config.SipServer)
		logger.Error("Start Server Error:", err)
	} else {
		logger.Info(fmt.Sprintf("[SIP SERVER] start success, %s://%s", config.SipServer.Mode, sipAddr))
	}

	go startJob()
}

func startJob() {
	keepaliveTicker := time.NewTicker(3 * time.Minute)
	publishTicker := time.NewTicker(10 * time.Minute)
	for {
		select {
		case <-keepaliveTicker.C:
			deviceOffline()
		case <-publishTicker.C:
			publishOffline()
		}
	}
}

func deviceOffline() {
	DeviceKeepalive.Range(func(key, value interface{}) bool {
		deviceID := key.(string)
		keepaliveTime := value.(time.Time)
		device, ok := getGB28181DeviceById(deviceID)
		if time.Since(keepaliveTime) > 3*time.Minute || !ok {
			device.Status = DeviceOffLineStatus
			DeviceKeepalive.Delete(key)
			GlobalGB28181DeviceStore.DeviceOffline(deviceID)
			logger.Infof("GB28181设备 %s 心跳超时，已下线", deviceID)
		}
		return true
	})
}

func publishOffline() {
	PublishStore.Range(func(key, v interface{}) bool {
		streamId := key.(string)
		if exist, err := ApiClientGetRtpInfo(streamId); err != nil || !exist {
			// 流不存在
			PublishStore.Delete(key)
			logger.Infof("流 %s 已下线", streamId)
		}
		return true
	})
}
