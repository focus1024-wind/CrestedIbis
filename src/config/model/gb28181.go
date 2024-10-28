package model

type GB28181 struct {
	MediaServer struct {
		IP   string `yaml:"ip"`
		Port uint16 `yaml:"port"`
	} `yaml:"media-server"`
}
