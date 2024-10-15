package model

type Jwt struct {
	Key        string `json:"key"`
	ExpireTime int    `yaml:"expire-time"`
}
