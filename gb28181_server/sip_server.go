package gb28181_server

import (
	"fmt"
	"github.com/ghettovoice/gosip"
	"github.com/ghettovoice/gosip/log"
	"github.com/ghettovoice/gosip/sip"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"strings"
)

var (
	globalSipServer gosip.Server
)

// startSipServer 启动SIP服务器
func (config *GB28181Config) startSipServer() {
	sipAddr := fmt.Sprintf("%s:%d", config.SipServer.IP, config.SipServer.Port)

	globalSipServer = gosip.NewServer(gosip.ServerConfig{
		Host: config.SipServer.IP,
	}, nil, nil, NewLogrusLogger())

	_ = globalSipServer.OnRequest(sip.REGISTER, config.SipRegisterHandler)
	_ = globalSipServer.OnRequest(sip.MESSAGE, config.SipMessageHandler)
	_ = globalSipServer.OnRequest(sip.NOTIFY, config.SipNotifyHandler)
	_ = globalSipServer.OnRequest(sip.BYE, config.SipByeHandler)

	err := globalSipServer.Listen(strings.ToLower(config.SipServer.Mode), sipAddr)
	if err != nil {
		fmt.Println("[SIP SERVER] Start Server Error:", err)
	} else {
		fmt.Println(fmt.Sprintf("[SIP SERVER] start success, %s://%s", config.SipServer.Mode, sipAddr))
	}
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

	return log.NewLogrusLogger(logger, "main", nil)
}
