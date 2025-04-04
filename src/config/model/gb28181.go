package model

type GB28181 struct {
	Serial     string `yaml:"serial"`
	Realm      string `yaml:"realm"`
	Password   string `yaml:"password"`
	AutoInvite bool   `yaml:"auto-invite"`
	LogInfo    string `yaml:"log-info"`
	SipServer  struct {
		IP       string `yaml:"ip"`
		PublicIp string `yaml:"public-ip"`
		Port     string `yaml:"port"`
		Mode     string `yaml:"mode"`
	} `yaml:"sip-server"`
	MediaServer struct {
		IP       string `yaml:"ip"`
		PublicIp string `yaml:"public-ip"`
		Port     uint16 `yaml:"port"`
		Mode     string `yaml:"mode"`
	} `yaml:"media-server"`
	HttpServer struct {
		Port uint16 `yaml:"port"`
	} `yaml:"http-server"`
}

func (gb28181 *GB28181) Init() {
	if gb28181.SipServer.PublicIp == "" {
		gb28181.SipServer.PublicIp = gb28181.SipServer.IP
	}

	if gb28181.MediaServer.PublicIp == "" {
		gb28181.MediaServer.PublicIp = gb28181.MediaServer.IP
	}
}
