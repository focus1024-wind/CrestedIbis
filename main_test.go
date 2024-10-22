package main

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"testing"
)

type DefaultMediaResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func GetMediaInvitePort() {
	client := resty.New()

	response := &struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Port int    `json:"port"`
	}{}
	resp, err := client.R().SetQueryParams(map[string]string{
		"secret":    "Uva5PKPA1aGOoBiQHglYExH2kMi2cX2S",
		"port":      "0",
		"tcp_mode":  "1",
		"vhost":     "__defaultVhost__",
		"app":       "rtp",
		"stream_id": "test",
	}).SetResult(response).SetError(response).Get("http://localhost:80/index/api/openRtpServer")
	fmt.Println(resp)
	fmt.Println(err)
	fmt.Println(response)
}

func Test(t *testing.T) {
	GetMediaInvitePort()
}
