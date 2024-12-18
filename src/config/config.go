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
	// 目前仅支持 yaml 类型配置文件
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
		panic(fmt.Sprintf("获取配置文件失败: %s", err.Error()))
	}

	// 读取配置文件
	data, err := os.ReadFile(configFilePath)
	if err != nil {
		panic(fmt.Sprintf("读取配置文件失败: %s", err.Error()))
	}

	// yaml解析配置文件
	err = yaml.Unmarshal(data, config)
	if err != nil {
		panic(fmt.Sprintf("解析 yaml 格式失败: %s", err.Error()))
	}

	if config.HttpServer.PublicHost == "" {
		config.HttpServer.PublicHost = fmt.Sprintf("http://%s:%s", config.HttpServer.IP, config.HttpServer.Port)
	}

	config.Store.Init()
	return config
}
