package gb28181_server

import (
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"time"
)

// ApiClientOpenRtpServer 调用 openRtpServer 接口，开启GB28181端口推流
func ApiClientOpenRtpServer(streamId string) (port int, err error) {
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

func ApiClientGetRtpInfo(streamId string) (exist bool, err error) {
	client := resty.New()

	response := &struct {
		Code  int  `json:"code"`
		Exist bool `json:"exist"`
	}{}

	_, err = client.R().SetQueryParams(map[string]string{
		"secret":    globalGB28181Config.MediaServer.Secret,
		"vhost":     "__defaultVhost__",
		"app":       "rtp",
		"stream_id": streamId,
	}).SetResult(response).SetError(response).Get(fmt.Sprintf("%s/index/api/getRtpInfo", globalGB28181Config.MediaServer.server))

	if err != nil {
		logger.Errorf("获取RTP推流信息失败：%s", err.Error())
		return false, err
	}
	exist = response.Exist
	return
}

func record(app string, streamId string, url string, maxSecond int) (ok bool, err error) {
	if ok, _ := AddStreamProxy(app, streamId, url, true); ok {
		logger.Infof("%s 开始录制", streamId)
		time.Sleep(time.Duration(maxSecond) * time.Second)
		_, _ = DelStreamProxy(fmt.Sprintf("__defaultVhost__/%s/%s", app, streamId))
		logger.Infof("%s 录制结束", streamId)
	}
	return ok, err
}

func StartRecordMp4(app string, deviceId string, channelId string, maxSecond int) (ok bool, err error) {
	streamId := fmt.Sprintf("%s/%s", deviceId, channelId)
	stream, _ := PublishStore.Load(streamId)

	if stream != nil {
		// 流已存在，开始录制
		mediaPlayUrl := GetMediaPlayUrl(streamId)
		_, _ = record(app, streamId, mediaPlayUrl["flv"], maxSecond)
	} else {
		// 流不存在，点播并录制
		logger.Infof("%s 流不存在，重新点播", streamId)
		mediaPlayUrl := Play(deviceId, channelId)
		_, _ = record(app, streamId, mediaPlayUrl["flv"], maxSecond)
		PlayStop(deviceId, channelId)
	}
	return
}

func AddStreamProxy(app string, streamId string, url string, record bool) (ok bool, err error) {
	client := resty.New()

	response := &struct {
		Code int `json:"code"`
	}{}

	enableMp4 := "0"
	if record {
		enableMp4 = "1"
	} else {
		enableMp4 = "0"
	}
	_, err = client.R().SetQueryParams(map[string]string{
		"secret":     globalGB28181Config.MediaServer.Secret,
		"vhost":      "__defaultVhost__",
		"app":        app,
		"stream":     streamId,
		"url":        url,
		"enable_mp4": enableMp4,
	}).SetResult(response).SetError(response).Get(fmt.Sprintf("%s/index/api/addStreamProxy", globalGB28181Config.MediaServer.server))

	if response.Code == 0 {
		return true, err
	} else {
		return false, err
	}
}

func DelStreamProxy(key string) (ok bool, err error) {
	client := resty.New()

	response := &struct {
		Code int `json:"code"`
	}{}

	_, err = client.R().SetQueryParams(map[string]string{
		"secret": globalGB28181Config.MediaServer.Secret,
		"key":    key,
	}).SetResult(response).SetError(response).Get(fmt.Sprintf("%s/index/api/delStreamProxy", globalGB28181Config.MediaServer.server))

	if response.Code == 0 {
		return true, err
	} else {
		return false, err
	}
}

// DelRecord 删除回放
func DelRecord(app string, stream string, period string, name string) (ok bool, err error) {
	client := resty.New()

	response := &struct {
		Code int `json:"code"`
	}{}

	if name == "" {
		_, err = client.R().SetQueryParams(map[string]string{
			"secret": globalGB28181Config.MediaServer.Secret,
			"vhost":  "__defaultVhost__",
			"app":    app,
			"stream": stream,
			"period": period,
		}).SetResult(response).SetError(response).Get(fmt.Sprintf("%s/index/api/deleteRecordDirectory", globalGB28181Config.MediaServer.server))
	} else {
		_, err = client.R().SetQueryParams(map[string]string{
			"secret": globalGB28181Config.MediaServer.Secret,
			"vhost":  "__defaultVhost__",
			"app":    app,
			"stream": stream,
			"period": period,
			"name":   name,
		}).SetResult(response).SetError(response).Get(fmt.Sprintf("%s/index/api/deleteRecordDirectory", globalGB28181Config.MediaServer.server))
	}

	if response.Code == 0 {
		return true, err
	} else {
		return false, err
	}
}
