package gb28181_server

import (
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
)

func GetMediaInvitePort(streamId string) (port int, err error) {
	client := resty.New()

	response := &struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Port int    `json:"port"`
	}{}
	_, err = client.R().SetQueryParams(map[string]string{
		"secret":    globalGB28181Config.MediaServer.Secret,
		"port":      "0",
		"tcp_mode":  globalGB28181Config.MediaServer.Mode,
		"vhost":     "__defaultVhost__",
		"app":       "rtp",
		"stream_id": streamId,
	}).SetResult(response).SetError(response).Get(fmt.Sprintf("%s/index/api/openRtpServer", globalGB28181Config.MediaServer.Server))

	if err != nil {
		err = errors.New(response.Msg)
	}
	port = response.Port
	return
}
