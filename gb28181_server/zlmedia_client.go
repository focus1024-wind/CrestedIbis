package gb28181_server

import (
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
)

// GetMediaInvitePort 调用 openRtpServer 接口，开启GB28181端口推流
func GetMediaInvitePort(streamId string) (port int, err error) {
	client := resty.New()

	response := &struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Port int    `json:"port"`
	}{}
	tcpMode := "0"
	if globalGB28181Config.MediaServer.Mode == "tcp" {
		tcpMode = "1"
	}
	_, err = client.R().SetQueryParams(map[string]string{
		"secret":    globalGB28181Config.MediaServer.Secret,
		"port":      "0",
		"tcp_mode":  tcpMode,
		"vhost":     "__defaultVhost__",
		"app":       "rtp",
		"stream_id": streamId,
	}).SetResult(response).SetError(response).Get(fmt.Sprintf("%s/index/api/openRtpServer", globalGB28181Config.MediaServer.server))

	if err != nil {
		logger.Error(err)
		logger.Error(response.Msg)
		err = errors.New(response.Msg)
		return 0, err
	}
	port = response.Port
	return
}
