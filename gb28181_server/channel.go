package gb28181_server

import (
	"CrestedIbis/gb28181_server/utils"
	"context"
	"fmt"
	"github.com/ghettovoice/gosip/sip"
	"net/http"
	"strconv"
	"strings"
)

type GB28181Channel struct {
	ParentID     string `desc:"GB28181父设备ID" json:"parent_id"`
	DeviceID     string `desc:"GB28181通道ID" json:"device_id"`
	Name         string `desc:"GB28181设备名称" json:"name"`
	Manufacturer string `desc:"GB28181设备制作厂商" json:"manufacturer"`
	Model        string `desc:"GB28181设备Model" json:"model"`
	Status       string `json:"status"`
	State        int    `desc:"通道状态" enum:"0: 空闲，1：Invite，2：正在播放/对讲" json:"state"`
}

// CreateSipRequest 创建通用SIP请求
func (channel *GB28181Channel) CreateSipRequest(method sip.RequestMethod) (req sip.Request) {
	if device, ok := getOnlineGB28181DeviceById(channel.ParentID); ok {
		device.SN++

		callId := sip.CallID(utils.RandNumString(10))
		userAgent := sip.UserAgentHeader("CrestedIbis")
		maxForwards := sip.MaxForwards(70)
		cseq := sip.CSeq{
			SeqNo:      uint32(device.SN),
			MethodName: method,
		}
		port := sip.Port(globalGB28181Config.SipServer.Port)

		fromAddress := sip.Address{
			Uri: &sip.SipUri{
				FUser: sip.String{Str: globalGB28181Config.Serial},
				FHost: globalGB28181Config.SipServer.IP,
				FPort: &port,
			},
			Params: sip.NewParams().Add("tag", sip.String{Str: utils.RandNumString(9)}),
		}
		toAddress, _ := device.getToAddress()

		req = sip.NewRequest(
			"",
			method,
			toAddress.Uri,
			"SIP/2.0",
			[]sip.Header{
				fromAddress.AsFromHeader(),
				toAddress.AsToHeader(),
				&callId,
				&userAgent,
				&cseq,
				&maxForwards,
				fromAddress.AsContactHeader(),
			},
			"",
			nil,
		)

		req.SetTransport(globalGB28181Config.SipServer.Mode)
		req.SetDestination(device.DeviceAddr)

		return req
	} else {
		return req
	}
}

// Invite 点播
func (channel *GB28181Channel) Invite(opt *InviteOptions) (err error) {
	logger.Infof("[STREAM] 开始拉流，流ID %s/%s", channel.ParentID, channel.DeviceID)
	var (
		streamPath string
		streamMode string
	)

	opt.CreateSSRC()
	if opt.IsRecord() {
		// 回放流
		streamPath = fmt.Sprintf("%s/%s/%d-%d", channel.ParentID, channel.DeviceID, opt.Start, opt.End)
		streamMode = "Playback"
	} else {
		// 直播流
		streamPath = fmt.Sprintf("%s/%s", channel.ParentID, channel.DeviceID)
		streamMode = "Play"
	}

	port, err := GetMediaInvitePort(streamPath)

	if err != nil {
		logger.Errorf("拉流失败: %v", err.Error())
		return err
	} else {
		opt.MediaPort = uint16(port)
	}

	protocol := ""
	if globalGB28181Config.MediaServer.Mode == "tcp" {
		protocol = "tcp"
	}
	// sdp信息
	sdpInfo := []string{
		"v=0",
		fmt.Sprintf("o=%s 0 0 IN IP4 %s", channel.DeviceID, globalGB28181Config.MediaServer.IP),
		"s=" + streamMode,
		"c=IN IP4 " + globalGB28181Config.MediaServer.IP,
		opt.String(),
		fmt.Sprintf("m=video %d %sRTP/AVP 96 97 98", opt.MediaPort, protocol),
		"a=recvonly",
		"a=rtpmap:96 PS/90000",
		"a=rtpmap:97 MPEG4/90000",
		"a=rtpmap:98 H264/90000",
		"y=" + opt.ssrc,
	}

	if globalGB28181Config.MediaServer.Mode == "tcp" {
		sdpInfo = append(sdpInfo, "a=setup:passive", "a=connection:new")
	}

	invite := channel.CreateSipRequest(sip.INVITE)

	contentType := sip.ContentType("APPLICATION/SDP")
	invite.AppendHeader(&contentType)

	invite.SetBody(strings.Join(sdpInfo, "\r\n")+"\r\n", true)

	subject := sip.GenericHeader{
		HeaderName: "Subject",
		Contents:   fmt.Sprintf("%s:%s,%s:0", channel.DeviceID, opt.ssrc, globalGB28181Config.Serial),
	}
	invite.AppendHeader(&subject)

	inviteRes, err := globalSipServer.RequestWithContext(context.Background(), invite)
	if err != nil {
		return err
	}

	if int(inviteRes.StatusCode()) == http.StatusOK {
		inviteResBodyLines := strings.Split(inviteRes.Body(), "\r\n")

		for _, line := range inviteResBodyLines {
			if ls := strings.Split(line, "="); len(ls) > 1 {
				if ls[0] == "y" && len(ls[1]) > 0 {
					if _ssrc, err := strconv.ParseInt(ls[1], 10, 0); err == nil {
						opt.SSRC = uint32(_ssrc)
					}
				}
				if ls[0] == "m" && len(ls[1]) > 0 {
					networkInfo := strings.Split(ls[1], " ")
					if strings.ToUpper(networkInfo[2]) != "TCP/RTP/AVP" {
						logger.Debug("[Stream] ipc_device not support tcp, streamPath: ", streamPath)
					}
				}
			}
		}
		err = globalSipServer.Send(sip.NewAckRequest("", invite, inviteRes, "", nil))
	}

	return nil
}

// AutoInvite 自动拉流
func AutoInvite(deviceID string, opt *InviteOptions) {
	if globalGB28181Config.AutoInvite {
		value, _ := DeviceChannels.Load(deviceID)

		if value != nil {
			channelIDs := value.([]string)

			for _, channelID := range channelIDs {
				streamPath := fmt.Sprintf("%s/%s", deviceID, channelID)
				stream, _ := PublishStore.Load(streamPath)

				if stream != nil {
					// 流已存在，不重复拉流
					logger.Info("[Stream] 已存在码流, streamPath", streamPath)
				} else {
					channel, exist := GlobalGB28181DeviceStore.LoadChannel(deviceID, channelID)
					if exist {
						_ = channel.Invite(opt)
					}
				}
			}
		}
	}
}

// Bye 停止点播
func (channel *GB28181Channel) Bye() int {
	logger.Infof("码流 %s/%s 停止点播", channel.ParentID, channel.DeviceID)
	request := channel.CreateSipRequest(sip.BYE)

	resp, err := globalSipServer.RequestWithContext(context.Background(), request)
	if err != nil {
		return http.StatusInternalServerError
	}

	return int(resp.StatusCode())
}
