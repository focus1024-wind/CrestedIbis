package gb28181_server

import (
	"CrestedIbis/gb28181_server/utils"
	"encoding/xml"
	"github.com/ghettovoice/gosip/sip"
	"net/http"
	"time"
)

func (config *GB28181Config) SipMessageHandler(req sip.Request, tx sip.ServerTransaction) {
	if device, ok := getOnlineGB28181DeviceBySip(req); ok {
		xmlMessageBody := &struct {
			XMLName          xml.Name
			CmdType          string
			SN               int // 请求序列号，一般用于对应 request 和 response
			DeviceID         string
			SumNum           int // 录像结果的总数 SumNum，录像结果会按照多条消息返回，可用于判断是否全部返回
			DeviceName       string
			Manufacturer     string
			Model            string
			AlarmDescription string
			Channels         []GB28181Channel `xml:"DeviceList>Item"`
		}{}

		err := utils.XMLDecoder(xmlMessageBody, []byte(req.Body()))
		if err != nil {
			logger.Error("[SIP SERVER] MESSAGE xml body 解析失败")
			return
		}

		logger.Info("[SIP SERVER] deviceID %s, Method MESSAGE, CmdType %s", device.DeviceID, xmlMessageBody.CmdType)

		var body string
		switch xmlMessageBody.CmdType {
		case "Keepalive":
			AutoInvite(device.DeviceID, &InviteOptions{})
			DeviceRegister.Store(device.DeviceID, time.Now())
		case "DeviceInfo":
			// 更新设备信息
			device.Name = xmlMessageBody.DeviceName
			device.Manufacturer = xmlMessageBody.Manufacturer
			device.Model = xmlMessageBody.Model
			GlobalGB28181DeviceStore.StoreDevice(device)
		case "Catalog":
			// 更新设备通道信息和设备通道ID信息
			var (
				channels   []GB28181Channel
				channelIDs []string
			)
			for _, channel := range xmlMessageBody.Channels {
				channel.ParentID = device.DeviceID
				channels = append(channels, channel)
				channelIDs = append(channelIDs, channel.DeviceID)
			}
			GlobalGB28181DeviceStore.UpdateChannels(channels)
			DeviceChannels.Store(device.DeviceID, channelIDs)
			AutoInvite(device.DeviceID, &InviteOptions{})
		case "Alarm":
			device.snapshot(1, 1)
		}

		_ = tx.Respond(sip.NewResponseFromRequest("", req, http.StatusOK, "OK", body))
	} else {
		// 设备未注册
		_ = tx.Respond(sip.NewResponseFromRequest("", req, http.StatusUnauthorized, "Unauthorized", ""))
		return
	}
}
