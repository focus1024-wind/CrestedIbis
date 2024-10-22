package gb28181_server

import (
	"github.com/ghettovoice/gosip/log"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var (
	GlobalGB28181DeviceStore GB28181DeviceStoreInterface
	logger                   = NewLogrusLogger()
)

func InstallGB28181DevicePlugin(devicePlugin GB28181DeviceStoreInterface) {
	GlobalGB28181DeviceStore = devicePlugin
}

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
