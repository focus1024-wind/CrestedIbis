package main

import (
	"CrestedIbis/gb28181_server"
	"CrestedIbis/src/apps/ipc/ipc_alarm"
	"CrestedIbis/src/apps/ipc/ipc_device"
	"CrestedIbis/src/config"
	"CrestedIbis/src/global"
	"CrestedIbis/src/global/initizalize"
	"github.com/spf13/cobra"
	"os"
)

var configFilePath string

func init() {
	rootCmd.Flags().StringVarP(&configFilePath, "config", "c", "./config.yaml", "config file path, currently only support yaml")
}

var rootCmd = &cobra.Command{
	Use:   "CrestedIbis",
	Short: "CrestedIbis Web video platform",
	PreRun: func(cmd *cobra.Command, args []string) {
		global.Conf = config.InitConfig(configFilePath)
		global.Logger = initizalize.InitLogger(global.Conf.Log)
		global.Db = initizalize.InitDatabase(global.Conf.Datasource)
		initizalize.InitDbTable()
	},
	Run: func(cmd *cobra.Command, args []string) {
		// 启动GB28181服务器
		gb28181_server.InstallGB28181DevicePlugin(new(ipc_device.IpcDevice))
		gb28181_server.InstallAlarmHandlerPlugin(new(ipc_alarm.IpcAlarm))
		go gb28181_server.Run(configFilePath)
		// 启动Web服务器
		initizalize.InitHttpServer()
	},
}

// @title			CrestedIbis
// @version		0.0.1
// @description	CrestedIbis目前是一个基于GB28181标准实现的音视频云平台，负责实现GB28181信令和设备管理，未来将会是一个支持物联网设备接入，算法训练和部署的综合物联网平台。
// @contact.name	北溪入江流(focus1024)
// @contact.url	http://focus1024.com(https://github.com/focus1024-wind/CrestedIbis)
// @contact.email	focus1024@foxmail.com
// @license.name	Apache 2.0
// @license.url	http://www.apache.org/licenses/LICENSE-2.0.html
func main() {
	if err := rootCmd.Execute(); err != nil {
		rootCmd.PrintErrf("CrestedIbis root cmd execute: %s", err)
		os.Exit(-1)
	}
}
