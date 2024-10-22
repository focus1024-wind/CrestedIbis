package gb28181_server

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"net/http"
	"os"
	"path/filepath"
)

var globalGB28181Config GB28181Config

type GB28181Config struct {
	SipServer struct {
		IP   string `default:"0.0.0.0" yaml:"ip" desc:"SIP服务器绑定网卡IP"`
		Port uint16 `default:"5060" yaml:"port" desc:"SIP服务器绑定端口"`
		Mode string `default:"udp" yaml:"mode" desc:"SIP服务器模式" enum:"udp: UDP模式，tcp: TCP模式"`
	} `yaml:"sip-server"`
	MediaServer struct {
		Secret string `yaml:"secret"`
		Server string `default:"http://localhost:80" yaml:"server" desc:"媒体服务器服务" yaml:"server"`
		Mode   string `default:"udp" yaml:"mode" desc:"SIP服务器模式" enum:"udp: UDP模式，tcp: TCP模式"`
	} `yaml:"media-server"`
	Serial     string `default:"61040200492007000001" yaml:"serial" desc:"GB28181服务器统一编码，参考GB28181-2022-附录D、GB/T2260-2007(https://openstd.samr.gov.cn/bzgk/gb/newGbInfo?hcno=C9C488FD717AFDCD52157F41C3302C6D)"`
	Realm      string `default:"6104020049" desc:"sip 服务域" yaml:"realm"`
	Username   string `desc:"sip 服务账号" yaml:"username"`
	Password   string `desc:"SIP服务器密码" yaml:"password"`
	AutoInvite bool   `default:"true" desc:"拉流模式" enum:"true: 自动拉流，false：手动拉流" yaml:"auto-invite"`
	LogLevel   string `default:"Info" yaml:"log-level"`
}

func Run(configFilePath string) {
	config := &struct {
		GB28181 GB28181Config
	}{}
	// 获取配置文件绝对路径
	configFilePath, err := filepath.Abs(configFilePath)
	if err != nil {
		panic(fmt.Sprintf("get config file abs path error: %s", err))
	}

	// 读取配置文件
	data, err := os.ReadFile(configFilePath)
	if err != nil {
		panic(fmt.Sprintf("read config file error: %s", err))
	}

	// yaml解析配置文件
	err = yaml.Unmarshal(data, config)
	if err != nil {
		panic(fmt.Sprintf("read yaml data from config file error: %s", err))
	}

	globalGB28181Config = config.GB28181

	level, err := logrus.ParseLevel(globalGB28181Config.LogLevel)
	if err != nil {
		panic(fmt.Sprintf("Invalid log level: %s", globalGB28181Config.LogLevel))
	}
	logger.SetLevel(uint32(level))

	if '/' == globalGB28181Config.MediaServer.Server[len(globalGB28181Config.MediaServer.Server)-1] {
		globalGB28181Config.MediaServer.Server = globalGB28181Config.MediaServer.Server[:len(globalGB28181Config.MediaServer.Server)-1]
	}

	globalGB28181Config.startSipServer()

	ApiHookRouters()
	err = http.ListenAndServe(":18080", nil)
	if err != nil {
		logger.Panicf("Web Server Start Error: %s", err)
		return
	}
}
