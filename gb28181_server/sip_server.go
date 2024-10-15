package gb28181_server

import (
	"fmt"
	"github.com/ghettovoice/gosip"
	"github.com/ghettovoice/gosip/sip"
	"strings"
)

var (
	globalSipServer gosip.Server
)

// startSipServer 启动SIP服务器
func (config *GB28181Config) startSipServer() {
	sipAddr := fmt.Sprintf("%s:%d", config.SipServer.IP, config.SipServer.Port)
	sipServerConfig := gosip.ServerConfig{
		Host: config.SipServer.IP,
	}

	logger := NewZapLogger(globalGB28181Plugin.Logger, "[SIP SERVER]", nil)
	globalSipServer = gosip.NewServer(sipServerConfig, nil, nil, logger)

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
