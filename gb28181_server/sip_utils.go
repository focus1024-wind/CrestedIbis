package gb28181_server

import (
	"github.com/ghettovoice/gosip/sip"
	"time"
)

// GetSipDeviceId 根据sip请求，获取SIP设备ID信息
func GetSipDeviceId(req sip.Request) (string, bool) {
	from, ok := req.From()
	if !ok || from.Address == nil || from.Address.User() == nil {
		return "", false
	} else {
		return from.Address.User().String(), true
	}
}

// GetSipDevice 根据sip请求，获取SIP设备信息
func GetSipDevice(req sip.Request) (GB28181Device, bool) {
	if deviceId, ok := GetSipDeviceId(req); ok {
		if device, ok := GlobalDeviceStore.LoadDevice(deviceId); ok {
			return device, true
		}
	}
	return GB28181Device{}, false
}

// GetSipOnlineDevice 根据sip请求，获取在线设备信息
func GetSipOnlineDevice(req sip.Request) (GB28181Device, bool) {
	if deviceId, ok := GetSipDeviceId(req); ok {
		registerTime, _ := DeviceRegister.Load(deviceId)

		// 设备已注册，且在3次心跳周期内
		if registerTime != nil && time.Now().Sub(registerTime.(time.Time)).Seconds() < 3*60 {
			return GetSipDevice(req)
		}
	}
	return GB28181Device{}, false
}

// GetDevice 根据设备ID，获取设备信息
func GetDevice(deviceId string) (GB28181Device, bool) {
	if device, ok := GlobalDeviceStore.LoadDevice(deviceId); ok {
		return device, true
	} else {
		return GB28181Device{}, false
	}
}

// GetOnlineDevice 根据设备ID，获取在线设备信息
func GetOnlineDevice(deviceId string) (GB28181Device, bool) {
	registerTime, _ := DeviceRegister.Load(deviceId)

	// 设备已注册，且在3次心跳周期内
	if registerTime != nil && time.Now().Sub(registerTime.(time.Time)).Seconds() < 3*60 {
		return GetDevice(deviceId)
	} else {
		return GB28181Device{}, false
	}
}
