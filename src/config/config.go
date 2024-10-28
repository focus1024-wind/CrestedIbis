package config

import (
	"CrestedIbis/src/config/model"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

// Config 配置文件结构订阅
type Config struct {
	HttpServer *model.HttpServer `yaml:"http-server"`
	Datasource *model.Datasource `yaml:"datasource"`
	Jwt        *model.Jwt        `yaml:"jwt"`
	Log        *model.Log        `yaml:"log"`
	Store      *model.Store      `yaml:"store"`
	GB28181    *model.GB28181    `yaml:"gb28181"`
}

// InitConfig 读取配置文件，生成配置
func InitConfig(configFilePath string) *Config {
	switch filepath.Ext(configFilePath) {
	case ".yml", ".yaml":
		return initYamlConfig(configFilePath)
	default:
		panic(fmt.Sprintf("config file type not support: %s", filepath.Ext(configFilePath)))
	}
}

// initYamlConfig 以yaml格式读取配置文件
func initYamlConfig(configFilePath string) *Config {
	config := &Config{}

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

	if config.HttpServer.PublicHost == "" {
		config.HttpServer.PublicHost = fmt.Sprintf("http://%s:%s", config.HttpServer.IP, config.HttpServer.Port)
	}
	config.Store.Init()
	return config
}
