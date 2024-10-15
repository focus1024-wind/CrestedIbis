package model

type Datasource struct {
	Type     string `yaml:"type"`
	Host     string `yaml:"host"`
	Port     uint16 `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DbName   string `yaml:"db-name"`
	MaxIdle  int    `default:"10" yaml:"max-idle"`
	MaxOpen  int    `default:"10" yaml:"max-open"`
}
