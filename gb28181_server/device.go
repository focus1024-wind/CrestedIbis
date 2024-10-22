package gb28181_server

import (
	"CrestedIbis/gb28181_server/utils"
	"context"
	"errors"
	"github.com/ghettovoice/gosip/sip"
	"regexp"
	"strings"
	"time"
)

const (
	DeviceOnLineStatus  = "ON"
	DeviceOffLineStatus = "OFF"
)

type GB28181Device struct {
	DeviceID     string    `desc:"GB28181设备ID"`
	Name         string    `desc:"GB28181设备名称"`
	Manufacturer string    `desc:"GB28181设备制作厂商"`
	Model        string    `desc:"GB28181设备Model"`
	SN           int       `desc:"SIP流媒体命令序列号"`
	FromAddress  string    `desc:"GB28181 FromHeader Uri信息"`
	DeviceAddr   string    `desc:"GB28181设备对应网卡IP"`
	RegisterTime time.Time `desc:"GB28181设备最新注册时间"`
	UpdatedTime  time.Time `desc:"GB28181设备最新更新时间"`
	Status       string    `desc:"GB28181设备状态"`
}

// GB28181DeviceStoreInterface 仅负责 GB28181Device 的存储相关操作
type GB28181DeviceStoreInterface interface {
	// LoadDevice 获取设备
	LoadDevice(deviceID string) (GB28181Device, bool)
	// StoreDevice 存储设备信息
	StoreDevice(gb28181Device GB28181Device)
	// DeviceOffline 注销设备：设备下线
	DeviceOffline(deviceId string)
	// LoadChannel 获取指定通道信息
	LoadChannel(deviceId string, channelId string) (GB28181Channel, bool)
	// LoadChannels 获取通道列表
	LoadChannels(deviceId string) ([]GB28181Channel, bool)
	// UpdateChannels 更新通道信息
	UpdateChannels(channels []GB28181Channel)
	// SnapShotUploadUrl 图片抓拍图像上传地址
	SnapShotUploadUrl(deviceId string) string
}

// ###### GB28181Device ######

// storeDevice 新建设备信息，设备上线
func (gb28181Device *GB28181Device) storeDevice(req sip.Request) {
	from, _ := req.From()

	gb28181Device.DeviceID, _ = getGB28181DeviceIdBySip(req)
	gb28181Device.FromAddress = from.Address.String()
	gb28181Device.DeviceAddr = req.Source()
	gb28181Device.RegisterTime = time.Now()
	gb28181Device.UpdatedTime = time.Now()
	gb28181Device.Status = DeviceOnLineStatus

	GlobalGB28181DeviceStore.StoreDevice(*gb28181Device)
}

// getToAddress 根据device.FromAddress信息获取服务端发起请求时ToAddress信息
func (gb28181Device *GB28181Device) getToAddress() (sip.Address, error) {
	var (
		uri               sip.SipUri
		uriRegExpNoUser   = regexp.MustCompile("^([A-Za-z]+):([^\\s;]+)(.*)$")
		uriRegExpWithUser = regexp.MustCompile("^([A-Za-z]+):([^@]+)@([^\\s;]+)(.*)$")
	)

	result := uriRegExpWithUser.FindStringSubmatch(gb28181Device.FromAddress)
	if len(result) != 5 {
		noUserResult := uriRegExpNoUser.FindStringSubmatch(gb28181Device.FromAddress)
		if len(noUserResult) != 4 {
			return sip.Address{}, errors.New("sip: uri format error")
		} else {
			result = []string{noUserResult[0], noUserResult[1], "", noUserResult[2], noUserResult[3]}
		}
	}

	if result[1] == "sips" {
		uri.FIsEncrypted = true
	}
	if result[2] != "" {
		parts := strings.Split(result[2], ":")
		uri.FUser = sip.String{Str: parts[0]}
		if len(parts) > 1 {
			uri.FPassword = sip.String{Str: parts[1]}
		}
	}
	uri.FHost = result[3]
	return sip.Address{Uri: &uri}, nil
}

// createSipRequest 新建SIP请求
func (gb28181Device *GB28181Device) createSipRequest(method sip.RequestMethod) (req sip.Request) {
	gb28181Device.SN++

	callId := sip.CallID(utils.RandNumString(10))
	userAgent := sip.UserAgentHeader("CrestedIbis")
	maxForwards := sip.MaxForwards(70)
	cseq := sip.CSeq{
		SeqNo:      uint32(gb28181Device.SN),
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
	toAddress, _ := gb28181Device.getToAddress()

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
	req.SetDestination(gb28181Device.DeviceAddr)

	return
}

// syncChannels 同步设备信息、下属通道信息，包括主动查询通道信息，订阅通道变化情况
func (gb28181Device *GB28181Device) syncChannels() {
	gb28181Device.syncDeviceInfo()
	gb28181Device.syncCatalog()
}

// syncDeviceInfo 同步IPC设备信息
func (gb28181Device *GB28181Device) syncDeviceInfo() {
	request := gb28181Device.createSipRequest(sip.MESSAGE)

	contentType := sip.ContentType("Application/MANSCDP+xml")
	request.AppendHeader(&contentType)
	request.SetBody(BuildDeviceInfoXML(gb28181Device.SN, gb28181Device.DeviceID), true)

	_, err := globalSipServer.RequestWithContext(context.Background(), request)
	if err != nil {
		logger.Error("[SIP SERVER] DeviceId: %s 同步设备信息失败", gb28181Device.DeviceID)
	}
}

// syncCatalog 同步设备通道信息
func (gb28181Device *GB28181Device) syncCatalog() {
	request := gb28181Device.createSipRequest(sip.MESSAGE)

	expires := sip.Expires(3600)
	contentType := sip.ContentType("Application/MANSCDP+xml")

	request.AppendHeader(&expires)
	request.AppendHeader(&contentType)
	request.SetBody(BuildCatalogXML(gb28181Device.SN, gb28181Device.DeviceID), true)

	_, err := globalSipServer.RequestWithContext(context.Background(), request)
	if err != nil {
		logger.Error("[SIP SERVER] DeviceId: %s 同步通道信息失败", gb28181Device.DeviceID)
	}
}

// snapShot 图片抓拍
func (gb28181Device *GB28181Device) snapshot(snapNum int, interval int) {
	request := gb28181Device.createSipRequest(sip.MESSAGE)

	contentType := sip.ContentType("Application/MANSCDP+xml")
	request.AppendHeader(&contentType)

	request.SetBody(BuildSnapShotXML(gb28181Device.SN, gb28181Device.DeviceID, snapNum, interval), true)

	_, err := globalSipServer.RequestWithContext(context.Background(), request)
	if err != nil {
		logger.Error("[SIP SERVER] DeviceId: %s 图片抓拍失败", gb28181Device.DeviceID)
	}
}
