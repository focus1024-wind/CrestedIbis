package gb28181_server

import (
	"github.com/ghettovoice/gosip/log"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var (
	GlobalGB28181DeviceStore GB28181DeviceStoreInterface
	GlobalAlarmHandler       AlarmHandlerInterface
	logger                   = NewLogrusLogger()
)

// InstallGB28181DevicePlugin 注册国标设备保存插件
// ToDo: 添加本地文件接口事件，在不需要依赖外部接口的情况下保存国标数据
func InstallGB28181DevicePlugin(devicePlugin GB28181DeviceStoreInterface) {
	GlobalGB28181DeviceStore = devicePlugin
}

// InstallAlarmHandlerPlugin 注册国标设备保存插件
// ToDo: 添加本地文件接口事件，在不需要依赖外部接口的情况下注册
func InstallAlarmHandlerPlugin(alarmHandler AlarmHandlerInterface) {
	GlobalAlarmHandler = alarmHandler
}

// NewLogrusLogger 新建 logrus logger对象
func NewLogrusLogger() *log.LogrusLogger {
	logger := logrus.New()
	logger.Level = logrus.ErrorLevel
	logger.Formatter = &prefixed.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05.000",
		ForceColors:     true,
		ForceFormatting: true,
	}

	return log.NewLogrusLogger(logger, "[SIP SERVER]", nil)
}
