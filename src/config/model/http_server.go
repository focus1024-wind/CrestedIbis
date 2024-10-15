package model

type HttpServer struct {
	IP         string `yaml:"ip"`
	Port       uint16 `yaml:"port"`
	PublicHost string `yaml:"public-host"`
}
