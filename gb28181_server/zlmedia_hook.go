package gb28181_server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

// PublishStore 维护已播放流
// Key: 流ID、Value: Set(流注册协议)
var PublishStore sync.Map

// PublishStoreStore PublishStore对应码流ID，处理注册信息
func PublishStoreStore(streamId string, stream string) {
	value, _ := PublishStore.Load(streamId)

	var streamList map[string]bool
	if value != nil {
		// 流已存在
		streamList = value.(map[string]bool)
	} else {
		// 流未存在
		streamList = map[string]bool{}
	}
	streamList[stream] = true
	PublishStore.Store(streamId, streamList)
}

// PublishStoreDelete PublishStore对应码流ID，处理注销信息
func PublishStoreDelete(streamId string, stream string) {
	value, _ := PublishStore.Load(streamId)

	if value != nil {
		streamList := value.(map[string]bool)
		delete(streamList, stream)
		PublishStore.Store(streamId, streamList)
	}
}

// GetMediaPlayUrl 根据StreamId生成对应播放规则URL
func GetMediaPlayUrl(streamId string) map[string]string {
	var mediaPlayUrl map[string]string

	value, _ := PublishStore.Load(streamId)

	if value != nil {
		streamList := value.(map[string]bool)
		if streamList["rtsp"] {
			mediaPlayUrl["rtsp"] = fmt.Sprintf("rtsp://%s/rtp/%s", globalGB28181Config.MediaServer.IP, streamId)
			mediaPlayUrl["rtsps"] = fmt.Sprintf("rtsps://%s/rtp/%s", globalGB28181Config.MediaServer.IP, streamId)
		}
		if streamList["rtmp"] {
			mediaPlayUrl["rtmp"] = fmt.Sprintf("rtmp://%s/rtp/%s", globalGB28181Config.MediaServer.IP, streamId)
			mediaPlayUrl["rtmps"] = fmt.Sprintf("rtmps://%s/rtp/%s", globalGB28181Config.MediaServer.IP, streamId)
			mediaPlayUrl["flv"] = fmt.Sprintf("http://%s/rtp/%s.live.flv", globalGB28181Config.MediaServer.IP, streamId)
			mediaPlayUrl["https_flv"] = fmt.Sprintf("https://%s/rtp/%s.live.flv", globalGB28181Config.MediaServer.IP, streamId)
			mediaPlayUrl["ws_flv"] = fmt.Sprintf("ws://%s/rtp/%s.live.flv", globalGB28181Config.MediaServer.IP, streamId)
			mediaPlayUrl["wss_flv"] = fmt.Sprintf("wss://%s/rtp/%s.live.flv", globalGB28181Config.MediaServer.IP, streamId)
		}
		if streamList["hls"] {
			mediaPlayUrl["hls"] = fmt.Sprintf("http://%s/rtp/%s/hls.m3u8", globalGB28181Config.MediaServer.IP, streamId)
			mediaPlayUrl["https_hls"] = fmt.Sprintf("https://%s/rtp/%s/hls.m3u8", globalGB28181Config.MediaServer.IP, streamId)
			mediaPlayUrl["hls_fmp4"] = fmt.Sprintf("http://%s/rtp/%s/hls.fmp4.m3u8", globalGB28181Config.MediaServer.IP, streamId)
			mediaPlayUrl["https_hls_fmp4"] = fmt.Sprintf("https://%s/rtp/%s/hls.fmp4.m3u8", globalGB28181Config.MediaServer.IP, streamId)
		}
		if streamList["ts"] {
			mediaPlayUrl["ts"] = fmt.Sprintf("http://%s/rtp/%s.live.ts", globalGB28181Config.MediaServer.IP, streamId)
			mediaPlayUrl["https_ts"] = fmt.Sprintf("https://%s/rtp/%s.live.ts", globalGB28181Config.MediaServer.IP, streamId)
			mediaPlayUrl["ws_ts"] = fmt.Sprintf("ws://%s/rtp/%s.live.ts", globalGB28181Config.MediaServer.IP, streamId)
			mediaPlayUrl["wss_ts"] = fmt.Sprintf("wss://%s/rtp/%s.live.ts", globalGB28181Config.MediaServer.IP, streamId)
		}
		if streamList["fmp4"] {
			mediaPlayUrl["fmp4"] = fmt.Sprintf("http://%s/rtp/%s.live.mp4", globalGB28181Config.MediaServer.IP, streamId)
			mediaPlayUrl["https_fmp4"] = fmt.Sprintf("https://%s/rtp/%s.live.mp4", globalGB28181Config.MediaServer.IP, streamId)
			mediaPlayUrl["ws_fmp4"] = fmt.Sprintf("ws://%s/rtp/%s.live.mp4", globalGB28181Config.MediaServer.IP, streamId)
			mediaPlayUrl["wss_fmp4"] = fmt.Sprintf("wss://%s/rtp/%s.live.mp4", globalGB28181Config.MediaServer.IP, streamId)
		}
	}
	return mediaPlayUrl
}

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

// ApiHookOnFlowReport 流量统计事件
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
// 参数说明请查看：https://docs.zlmediakit.com/zh/guide/media_server/web_hook_api.html#_7%E3%80%81on-publish
func ApiHookOnPublish(w http.ResponseWriter, r *http.Request) {
	logger.Info("ApiHookOnPublish 推流鉴权")
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
		logger.Errorf("ApiHookOnPublish 推流鉴权 请求解析失败: %v", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := &struct {
		Code          int    `json:"code"`
		Msg           string `json:"msg"`
		EnableHls     bool   `json:"enable_hls"`
		EnableHlsFmp4 bool   `json:"enable_hls_fmp4"`
		EnableMp4     bool   `json:"enable_mp4"`
		EnableRtsp    bool   `json:"enable_rtsp"`
		EnableRtmp    bool   `json:"enable_rtmp"`
		EnableTs      bool   `json:"enable_ts"`
		EnableFmp4    bool   `json:"enable_fmp4"`
		HlsDemand     bool   `json:"hls_demand"`
		RtspDemand    bool   `json:"rtsp_demand"`
		RtmpDemand    bool   `json:"rtmp_demand"`
		TsDemand      bool   `json:"ts_demand"`
		Fmp4Demand    bool   `json:"fmp4_demand"`
		EnableAudio   bool   `json:"enable_audio"`
		ModifyStamp   int    `json:"modify_stamp"`
	}{
		Code:          0,
		Msg:           "success",
		EnableHls:     true,
		EnableHlsFmp4: true,
		EnableMp4:     true,
		EnableRtsp:    true,
		EnableRtmp:    true,
		EnableTs:      true,
		EnableFmp4:    true,
		HlsDemand:     true,
		RtspDemand:    true,
		RtmpDemand:    true,
		TsDemand:      true,
		Fmp4Demand:    true,
		EnableAudio:   true,
		ModifyStamp:   1,
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
// 注册注销流同步 PublishStore
func ApiHookOnStreamChanged(w http.ResponseWriter, r *http.Request) {
	logger.Info("ApiHookOnStreamChanged")
	var req = &struct {
		App           string `json:"app"`
		Regist        bool   `json:"regist"`
		Schema        string `json:"schema"`
		Stream        string `json:"stream"`
		Vhost         string `json:"vhost"`
		MediaServerId string `json:"mediaServerId"`
	}{}

	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		logger.Errorf("ApiHookOnStreamChanged 请求解析失败: %v", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Regist {
		logger.Infof("%s 流注册", req.Stream)
		PublishStoreStore(req.Stream, req.Schema)
	} else {
		logger.Infof("%s 流注销", req.Stream)
		PublishStoreDelete(req.Stream, req.Schema)
	}

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

// ApiHookOnStreamNoneReader 流无人观看通知事件
func ApiHookOnStreamNoneReader(w http.ResponseWriter, _ *http.Request) {
	logger.Info("ApiHookOnStreamNoneReader")
	resp := &struct {
		Code  int  `json:"code"`
		Close bool `json:"close"`
	}{
		Code:  0,
		Close: false,
	}

	msg, _ := json.Marshal(resp)
	_, _ = w.Write(msg)
}

// ApiHookOnStreamNotFound 流未找到事件
// 未找到流，删除流
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
// 推流失败，删除流
func ApiHookOnRtpServerTimeout(w http.ResponseWriter, r *http.Request) {
	logger.Info("ApiHookOnRtpServerTimeout")
	var req = &struct {
		LocalPort     uint16 `json:"local_port"`
		ReUsePort     bool   `json:"re_use_port"`
		Ssrc          uint32 `json:"ssrc"`
		StreamId      string `json:"stream_id"`
		TcpMode       int    `json:"tcp_mode"`
		MediaServerId string `json:"mediaServerId"`
	}{}

	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		logger.Errorf("ApiHookOnRtpServerTimeout 请求解析失败: %v", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	PublishStore.Delete(req.StreamId)
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
