package gb28181_server

var (
	GlobalDeviceStore GB28181DeviceStoreInterface
)

func InstallDevicePlugin(devicePlugin GB28181DeviceStoreInterface) {
	GlobalDeviceStore = devicePlugin
}
