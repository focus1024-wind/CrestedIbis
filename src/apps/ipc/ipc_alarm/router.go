package ipc_alarm

import (
	"CrestedIbis/src/global"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http/httputil"
	"net/url"
)

func InitIpcAlarmRouter() {
	ipcAlarmRouter := global.HttpEngine.Group("/ipc/alarm")
	{
		ipcAlarmRouter.DELETE("", DeleteIpcAlarm)
		ipcAlarmRouter.GET("/alarms", GetIpcAlarms)
		ipcAlarmRouter.DELETE("/alarms", DeleteIpcAlarms)
	}

	ipcRecordRouter := global.HttpEngine.Group("/ipc/record")
	{
		ipcRecordRouter.GET("", DownloadRecord)
		ipcRecordRouter.DELETE("", DeleteRecord)
		ipcRecordRouter.GET("/records", GetIpcRecords)
		ipcRecordRouter.DELETE("/records", DeleteRecords)
	}

	global.HttpEngine.Any("/record/*name", proxyRecordHandler)
}

// proxyRecordHandler /record 路径：代理媒体服务器
func proxyRecordHandler(c *gin.Context) {
	var target = fmt.Sprintf("http://%s:%d/", global.Conf.GB28181.MediaServer.IP, global.Conf.GB28181.MediaServer.Port)
	proxyUrl, _ := url.Parse(target)
	proxy := httputil.NewSingleHostReverseProxy(proxyUrl)
	proxy.ServeHTTP(c.Writer, c.Request)
}
