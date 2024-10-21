package gb28181_server

import (
	"github.com/ghettovoice/gosip/sip"
	"time"
)

// getGB28181DeviceIdBySip 根据sip请求，获取SIP设备ID信息
func getGB28181DeviceIdBySip(req sip.Request) (string, bool) {
	from, ok := req.From()
	if !ok || from.Address == nil || from.Address.User() == nil {
		return "", false
	} else {
		return from.Address.User().String(), true
	}
}

// getGB28181DeviceBySip 根据sip请求，获取SIP设备信息
func getGB28181DeviceBySip(req sip.Request) (GB28181Device, bool) {
	if deviceId, ok := getGB28181DeviceIdBySip(req); ok {
		if device, ok := GlobalGB28181DeviceStore.LoadDevice(deviceId); ok {
			return device, true
		}
	}
	return GB28181Device{}, false
}

// getOnlineGB28181DeviceBySip 根据sip请求，获取在线设备信息
func getOnlineGB28181DeviceBySip(req sip.Request) (GB28181Device, bool) {
	if deviceId, ok := getGB28181DeviceIdBySip(req); ok {
		registerTime, _ := DeviceRegister.Load(deviceId)

		// 设备已注册，且在3次心跳周期内
		if registerTime != nil && time.Now().Sub(registerTime.(time.Time)).Seconds() < 3*60 {
			return getGB28181DeviceBySip(req)
		}
	}
	return GB28181Device{}, false
}

// getGB28181DeviceById 根据设备ID，获取设备信息
func getGB28181DeviceById(deviceId string) (GB28181Device, bool) {
	if device, ok := GlobalGB28181DeviceStore.LoadDevice(deviceId); ok {
		return device, true
	} else {
		return GB28181Device{}, false
	}
}

// getOnlineGB28181DeviceById 根据设备ID，获取在线设备信息
func getOnlineGB28181DeviceById(deviceId string) (GB28181Device, bool) {
	registerTime, _ := DeviceRegister.Load(deviceId)

	// 设备已注册，且在3次心跳周期内
	if registerTime != nil && time.Now().Sub(registerTime.(time.Time)).Seconds() < 3*60 {
		return getGB28181DeviceById(deviceId)
	} else {
		return GB28181Device{}, false
	}
}
