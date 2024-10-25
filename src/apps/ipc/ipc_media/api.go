package ipc_media

import (
	"CrestedIbis/gb28181_server"
	"CrestedIbis/src/global/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

// IpcMediaPlay Ipc设备点播
//
//	@Summary		Ipc设备点播
//	@Version		1.0.0
//	@Description	Ipc设备点播
//	@Tags			设备点播 /ipc/media
//	@Accept			json
//	@Produce		json
//	@Param			Authorization		header		string							false	"访问token"
//	@Param			access_token		query		string							false	"访问token"
//	@Param			IpcMediaPlayModel	body		IpcMediaPlayModel				true	"点播参数"
//	@Success		200					{object}	model.HttpResponse{data=string}	"注册成功，响应点播地址"
//	@Failure		500					{object}	model.HttpResponse{data=string}	"注册失败，响应失败信息"
//	@Router			/ipc/media/play [POST]
func IpcMediaPlay(c *gin.Context) {
	var ipcMediaPlay IpcMediaPlayModel
	if err := c.ShouldBind(&ipcMediaPlay); err != nil {
		panic(http.StatusBadRequest)
	} else {
		playUrl := gb28181_server.Play(ipcMediaPlay.DeviceId, ipcMediaPlay.ChannelId)
		model.HttpResponse{}.OkGin(c, playUrl)
		return
	}
}
