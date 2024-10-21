package gb28181_server

var globalGB28181Config GB28181Config

type GB28181Config struct {
	SipServer struct {
		IP   string `default:"0.0.0.0" desc:"SIP服务器绑定网卡IP"`
		Port uint16 `default:"5060" desc:"SIP服务器绑定端口"`
		Mode string `default:"udp" desc:"SIP服务器模式" enum:"udp: UDP模式，tcp: TCP模式"`
	} `yaml:"sip-server"`
	Serial   string `default:"61040200492007000001" desc:"GB28181服务器统一编码，参考GB28181-2022-附录D、GB/T2260-2007(https://openstd.samr.gov.cn/bzgk/gb/newGbInfo?hcno=C9C488FD717AFDCD52157F41C3302C6D)"`
	Realm    string `default:"6104020049" desc:"sip 服务域"`
	Username string `desc:"sip 服务账号"`
	Password string `desc:"SIP服务器密码"`
}
