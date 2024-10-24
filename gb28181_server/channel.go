package gb28181_server

import (
	"CrestedIbis/gb28181_server/utils"
	"context"
	"fmt"
	"github.com/ghettovoice/gosip/sip"
	"go.uber.org/zap"
	"m7s.live/plugin/ps/v4"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

var PublishStore sync.Map

type GB28181Channel struct {
	ParentID     string `desc:"GB28181父设备ID"`
	DeviceID     string `desc:"GB28181通道ID"`
	Name         string `desc:"GB28181设备名称"`
	Manufacturer string `desc:"GB28181设备制作厂商"`
	Model        string `desc:"GB28181设备Model"`
	Status       string
	State        int `desc:"通道状态" enum:"0: 空闲，1：Invite，2：正在播放/对讲"`
}

func (channel *GB28181Channel) CreateSipRequest(method sip.RequestMethod) (req sip.Request) {
	if device, ok := GetOnlineDevice(channel.ParentID); ok {
		device.SN++

		callId := sip.CallID(utils.RandNumString(10))
		userAgent := sip.UserAgentHeader(UserAgent)
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
		toAddress, _ := device.GetToAddress()

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

func (channel *GB28181Channel) Invite(opt *InviteOptions) (err error) {
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

	if globalGB28181Config.portsManager.Valid {
		opt.MediaPort, err = globalGB28181Config.portsManager.GetPort()
		opt.recyclePort = globalGB28181Config.portsManager.Recycle
	}
	if err != nil {
		return err
	}

	protocol := ""
	if globalGB28181Config.MediaServer.Mode == "tcp" {
		protocol = "tcp"
	}
	// sdp信息
	sdpInfo := []string{
		"v=0",
		fmt.Sprintf("o=%s 0 0 IN IP4 %s", channel.DeviceID, "192.168.1.11"),
		"s=" + streamMode,
		"c=IN IP4 " + "192.168.1.11",
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
		if opt.recyclePort != nil {
			_ = opt.recyclePort(opt.MediaPort)
		}
		return err
	}

	if int(inviteRes.StatusCode()) == http.StatusOK {
		var networkType = globalGB28181Config.MediaServer.Mode
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
						globalGB28181Plugin.Debug("[Stream] ipc_device not support tcp, ", zap.String("streamPath", streamPath))
						networkType = "udp"
					}
				}
			}
		}
		var psPublisher ps.PSPublisher

		err = psPublisher.Receive(streamPath, opt.dump, fmt.Sprintf("%s:%d", networkType, opt.MediaPort), opt.SSRC, false)
		if err != nil {
			if opt.recyclePort != nil {
				_ = opt.recyclePort(opt.MediaPort)
			}
			return err
		}

		if !opt.IsLive() {
			// 超时关闭无数据关闭
			if psPublisher.Stream.DelayCloseTimeout == 0 {
				psPublisher.Stream.DelayCloseTimeout = time.Second * time.Duration(globalGB28181Config.MediaServer.timeout)
			}
			if psPublisher.Stream.IdleTimeout == 0 {
				psPublisher.Stream.IdleTimeout = time.Duration(globalGB28181Config.MediaServer.timeout)
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

				if stream != nil && stream.(bool) {
					// 流已存在，不重复拉流
					globalGB28181Plugin.Info("[Stream] 已存在码流", zap.String("streamPath", streamPath))
				} else {
					channel, exist := GlobalDeviceStore.LoadChannel(deviceID, channelID)
					if exist {
						_ = channel.Invite(opt)
					}
				}
			}
		}
	}
}
