package gb28181_server

import (
	"encoding/json"
	"net/http"
)

func ApiHookRouters() {
	http.HandleFunc("/index/hook/on_flow_report", ApiHookOnFlowReport)
	http.HandleFunc("/index/hook/on_http_access", ApiHookOnHttpAccess)
	http.HandleFunc("/index/hook/on_play", ApiHookOnplay)
	http.HandleFunc("/index/hook/on_publish", ApiHookOnPublish)
	http.HandleFunc("/index/hook/on_record_mp4", ApiHookOnRecordMp4)
	http.HandleFunc("/index/hook/on_rtsp_auth", ApiHookOnRtspAuth)
	http.HandleFunc("/index/hook/on_rtsp_realm", ApiHookOnRtspRealm)
	http.HandleFunc("/index/hook/on_shell_login", ApiHookOnShellLogin)
	http.HandleFunc("/index/hook/on_stream_changed", ApiHookOnStreamChanged)
	http.HandleFunc("/index/hook/on_stream_none_reader", ApiHookOnStreamNoneReader)
	http.HandleFunc("/index/hook/on_stream_not_found", ApiHookOnStreamNotFound)
	http.HandleFunc("/index/hook/on_server_started", ApiHookOnServerStarted)
	http.HandleFunc("/index/hook/on_server_keepalive", ApiHookOnServerKeepalive)
	http.HandleFunc("/index/hook/on_rtp_server_timeout", ApiHookOnRtpServerTimeout)
}

// ApiHookOnFlowReport 浏览统计事件
func ApiHookOnFlowReport(_ http.ResponseWriter, _ *http.Request) {
	logger.Info("ApiHookOnFlowReport")
}

// ApiHookOnHttpAccess 文件访问鉴权事件: 访问 http 文件服务器上 hls 之外的文件时触发
func ApiHookOnHttpAccess(_ http.ResponseWriter, _ *http.Request) {
	logger.Info("ApiHookOnHttpAccess")
}

// ApiHookOnplay 播放器鉴权事件: 拉流播放时触发
func ApiHookOnplay(w http.ResponseWriter, _ *http.Request) {
	logger.Info("ApiHookOnplay")
	resp := &struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}{
		Code: 0,
		Msg:  "success",
	}

	msg, _ := json.Marshal(resp)
	_, _ = w.Write(msg)
}

// ApiHookOnPublish 推流鉴权事件
func ApiHookOnPublish(w http.ResponseWriter, r *http.Request) {
	logger.Info("ApiHookOnplay")
	var req = &struct {
		App           string `json:"app"`
		Id            string `json:"id"`
		Ip            string `json:"ip"`
		Params        string `json:"params"`
		Port          uint16 `json:"port"`
		Schema        string `json:"schema"`
		Stream        string `json:"stream"`
		Vhost         string `json:"vhost"`
		MediaServerId string `json:"mediaServerId"`
	}{}

	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		logger.Errorf("ApiHookOnPublish 请求解析失败: %v", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	PublishStore.Store(req.Stream, true)
	resp := &struct {
		Code       int    `json:"code"`
		Msg        string `json:"msg"`
		EnableRtmp bool   `json:"enable_rtmp"`
	}{
		Code:       0,
		Msg:        "success",
		EnableRtmp: true,
	}

	msg, _ := json.Marshal(resp)
	_, _ = w.Write(msg)
}

// ApiHookOnRecordMp4 MP4录制通知事件
func ApiHookOnRecordMp4(http.ResponseWriter, *http.Request) {
	logger.Info("ApiHookOnRecordMp4")
}

// ApiHookOnRtspRealm Rtsp是否启用专用鉴权事件处理函数
func ApiHookOnRtspRealm(_ http.ResponseWriter, _ *http.Request) {
	logger.Info("ApiHookOnRtspRealm")
}

// ApiHookOnRtspAuth Rtsp专用鉴权事件
func ApiHookOnRtspAuth(_ http.ResponseWriter, _ *http.Request) {
	logger.Info("ApiHookOnRtspAuth")
}

// ApiHookOnShellLogin Shell登录鉴权事件
func ApiHookOnShellLogin(_ http.ResponseWriter, _ *http.Request) {
	logger.Info("ApiHookOnShellLogin")
}

// ApiHookOnStreamChanged 流注册注销通知事件
func ApiHookOnStreamChanged(_ http.ResponseWriter, _ *http.Request) {
	logger.Info("ApiHookOnStreamChanged")
}

// ApiHookOnStreamNoneReader 流无人观看通知事件
func ApiHookOnStreamNoneReader(_ http.ResponseWriter, _ *http.Request) {
	logger.Info("ApiHookOnStreamNoneReader")
}

// ApiHookOnStreamNotFound 流未找到事件
func ApiHookOnStreamNotFound(w http.ResponseWriter, r *http.Request) {
	logger.Info("ApiHookOnStreamNotFound")
	var req = &struct {
		App           string `json:"app"`
		Id            string `json:"id"`
		Ip            string `json:"ip"`
		Params        string `json:"params"`
		Port          uint16 `json:"port"`
		Schema        string `json:"schema"`
		Stream        string `json:"stream"`
		Vhost         string `json:"vhost"`
		MediaServerId string `json:"mediaServerId"`
	}{}

	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		logger.Errorf("ApiHookOnStreamNotFound 请求解析失败: %v", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	PublishStore.Delete(req.Stream)
	resp := &struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}{
		Code: 0,
		Msg:  "success",
	}

	msg, _ := json.Marshal(resp)
	_, _ = w.Write(msg)
}

// ApiHookOnServerStarted 服务器启动事件
func ApiHookOnServerStarted(_ http.ResponseWriter, _ *http.Request) {
	logger.Info("ApiHookOnServerStarted")
}

// ApiHookOnServerKeepalive 服务器心跳事件，定时上报时间
func ApiHookOnServerKeepalive(w http.ResponseWriter, _ *http.Request) {
	logger.Info("ApiHookOnServerKeepalive")
	resp := &struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}{
		Code: 0,
		Msg:  "success",
	}

	msg, _ := json.Marshal(resp)
	_, _ = w.Write(msg)
}

// ApiHookOnRtpServerTimeout openRtpServer接口流超时事件
func ApiHookOnRtpServerTimeout(w http.ResponseWriter, r *http.Request) {
	logger.Info("ApiHookOnRtpServerTimeout")
	var req = &struct {
		LocalPort     string `json:"local_port"`
		Id            string `json:"id"`
		Ip            string `json:"ip"`
		Params        string `json:"params"`
		Port          uint16 `json:"port"`
		Schema        string `json:"schema"`
		Stream        string `json:"stream"`
		Vhost         string `json:"vhost"`
		MediaServerId string `json:"mediaServerId"`
	}{}

	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		logger.Errorf("ApiHookOnStreamNotFound 请求解析失败: %v", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	PublishStore.Delete(req.Stream)
	resp := &struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}{
		Code: 0,
		Msg:  "success",
	}

	msg, _ := json.Marshal(resp)
	_, _ = w.Write(msg)
}
