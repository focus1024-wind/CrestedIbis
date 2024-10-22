package main

import (
	"CrestedIbis/gb28181_server"
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
		go gb28181_server.Run(configFilePath)
		// 启动Web服务器
		initizalize.InitHttpServer()
	},
}

// @contact.name	CrestedIbis
// @contact.url	https://github.com/focus1024-wind/CrestedIbis
// @contact.email	focus1024@foxmail.com
// @license.name	Apache 2.0
// @license.url	http://www.apache.org/licenses/LICENSE-2.0.html
func main() {
	if err := rootCmd.Execute(); err != nil {
		rootCmd.PrintErrf("CrestedIbis root cmd execute: %s", err)
		os.Exit(-1)
	}
}
