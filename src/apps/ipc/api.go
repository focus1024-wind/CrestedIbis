package ipc

import (
	"CrestedIbis/src/global"
	"CrestedIbis/src/global/model"
	"github.com/gin-gonic/gin"
)

type GB28181Info struct {
	Serial   string `json:"serial"`
	Realm    string `json:"realm"`
	SipIp    string `json:"sip_ip"`
	SipPort  string `json:"sip_port"`
	Password string `json:"password"`
	SipMode  string `json:"sip_mode"`
	MediaIp  string `json:"media_ip"`
}

// GB28181InfoApi 获取GB28181连接信息
//
//	@Summary		获取GB28181连接信息
//	@Version		0.0.1
//	@Description	获取GB28181连接信息
//	@Tags			IPC设备
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string									false	"访问token"
//	@Param			access_token	query		string									false	"访问token"
//	@Success		200				{object}	model.HttpResponse{data=GB28181Info}	"查询成功"
//	@Router			/ipc/gb28181_info [GET]
func GB28181InfoApi(c *gin.Context) {
	gb28181Info := GB28181Info{
		Serial:   global.Conf.GB28181.Serial,
		Realm:    global.Conf.GB28181.Realm,
		SipIp:    global.Conf.GB28181.SipServer.PublicIp,
		SipPort:  global.Conf.GB28181.SipServer.Port,
		Password: global.Conf.GB28181.Password,
		SipMode:  global.Conf.GB28181.SipServer.Mode,
		MediaIp:  global.Conf.GB28181.MediaServer.PublicIp,
	}
	model.HttpResponse{}.OkGin(c, gb28181Info)
}
