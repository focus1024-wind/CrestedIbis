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
//	@Version		0.0.1
//	@Description	Ipc设备点播
//	@Tags			IPC设备点播 /ipc/media
//	@Accept			json
//	@Produce		json
//	@Param			Authorization		header		string										false	"访问token"
//	@Param			access_token		query		string										false	"访问token"
//	@Param			IpcMediaPlayModel	body		IpcMediaPlayModel							true	"点播参数"
//	@Success		200					{object}	model.HttpResponse{data=map[string]string}	"点播成功，响应点播地址"
//	@Failure		500					{object}	model.HttpResponse{data=nil}				"点播失败"
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

// IpcMediaStop Ipc设备停止点播
//
//	@Summary		Ipc设备停止点播
//	@Version		0.0.1
//	@Description	Ipc设备停止点播
//	@Tags			IPC设备点播 /ipc/media
//	@Accept			json
//	@Produce		json
//	@Param			Authorization		header		string							false	"访问token"
//	@Param			access_token		query		string							false	"访问token"
//	@Param			IpcMediaPlayModel	body		IpcMediaPlayModel				true	"点播参数"
//	@Success		200					{object}	model.HttpResponse{data=nil}	"停止点播成功"
//	@Failure		500					{object}	model.HttpResponse{data=nil}	"停止点播失败"
//	@Router			/ipc/media/stop [POST]
func IpcMediaStop(c *gin.Context) {
	var ipcMediaPlay IpcMediaPlayModel
	if err := c.ShouldBind(&ipcMediaPlay); err != nil {
		panic(http.StatusBadRequest)
	} else {
		gb28181_server.PlayStop(ipcMediaPlay.DeviceId, ipcMediaPlay.ChannelId)
		model.HttpResponse{}.OkGin(c, nil)
		return
	}
}
