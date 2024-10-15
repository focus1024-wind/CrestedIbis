package gb28181_server

import (
	"fmt"
	. "m7s.live/engine/v4"
	"strconv"
	"strings"
)

const UserAgent = "CrestedIbis"

// GB28181Config GB28181-Server插件配置项
type GB28181Config struct {
	SipServer struct {
		IP   string `default:"0.0.0.0" desc:"SIP服务器绑定网卡IP"`
		Port uint16 `default:"5060" desc:"SIP服务器绑定端口"`
		Mode string `default:"udp" desc:"SIP服务器模式" enum:"udp: UDP模式，tcp: TCP模式"`
	} `yaml:"sip-server"`
	MediaServer struct {
		IP      string `default:"0.0.0.0" desc:"媒体服务器IP"`
		Port    string `default:"58200-59200" desc:"媒体服务端口"`
		Mode    string `default:"udp" desc:"媒体服务模式" enum:"udp: UDP模式，tcp: TCP模式"`
		timeout uint8  `default:"10" desc:"媒体流超时时间，单位s"`
	} `yaml:"media-server"`
	Serial     string `default:"61040200492007000001" desc:"GB28181服务器统一编码，参考GB28181-2022-附录D、GB/T2260-2007(https://openstd.samr.gov.cn/bzgk/gb/newGbInfo?hcno=C9C488FD717AFDCD52157F41C3302C6D)"`
	Realm      string `default:"6104020049" desc:"sip 服务域"`
	Username   string `desc:"sip 服务账号"`
	Password   string `desc:"SIP服务器密码"`
	AutoInvite bool   `default:"true" desc:"拉流模式" enum:"true: 自动拉流，false：手动拉流"`

	portsManager PortManager
}

// OnEvent 实现 Monibuca 事件
func (config *GB28181Config) OnEvent(event any) {
	switch value := event.(type) {
	// 插件初始化逻辑
	case FirstConfig:
		// 初始化端口管理
		portRange := strings.Split(config.MediaServer.Port, "-")
		minPort, err := strconv.ParseInt(portRange[0], 10, 0)
		if err != nil {
			fmt.Println("[GB28181] config: gb28181.media-server.port error")
		}
		maxPort, err := strconv.ParseInt(portRange[0], 10, 0)
		if err != nil {
			fmt.Println("[GB28181] config: gb28181.media-server.port error")
		}
		config.portsManager = PortManager{}
		config.portsManager.Init(uint16(minPort), uint16(maxPort))

		config.startSipServer()
	// 插件热更新逻辑
	case UpdateConfig:
		fmt.Println("Monibuca### UpdateConfig")
	// 按需拉流逻辑
	case InvitePublish:
		fmt.Println("Monibuca### InvitePublish")
	// 由于发布者掉线等待发布者
	case SEwaitPublish:
		fmt.Println("Monibuca### SEwaitPublish")
	// 首次进入发布状态
	case SEpublish:
		fmt.Println("Monibuca### SEpublish")
		PublishStore.Store(value.Target.Path, true)
	// 再次进入发布状态
	case SErepublish:
		fmt.Println("Monibuca### SErepublish")
	// 由于最后一个订阅者离开等待关闭流
	case SEwaitClose:
		fmt.Println("Monibuca### SEwaitClose")
	// 关闭流
	case SEclose:
		fmt.Println("Monibuca### SEclose")
		PublishStore.Delete(value.Target.Path)
	// 订阅者离开
	case UnsubscribeEvent:
		fmt.Println("Monibuca### UnsubscribeEvent")
	}
}

var (
	globalGB28181Config GB28181Config
	// 注册 Monibuca 插件
	globalGB28181Plugin = InstallPlugin(&globalGB28181Config)
)
