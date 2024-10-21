package main

import (
	"CrestedIbis/gb28181_server_back"
	"CrestedIbis/src/apps/ipc/ipc_device"
	"CrestedIbis/src/config"
	"CrestedIbis/src/global"
	"CrestedIbis/src/global/initizalize"
	"context"
	"github.com/spf13/cobra"
	"m7s.live/engine/v4"
	_ "m7s.live/plugin/hdl/v4"
	_ "m7s.live/plugin/jessica/v4"
	_ "m7s.live/plugin/preview/v4"
	_ "m7s.live/plugin/ps/v4"
	_ "m7s.live/plugin/record/v4"
	_ "m7s.live/plugin/rtmp/v4"
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
		gb28181_server_back.InstallDevicePlugin(new(ipc_device.IpcDevice))
		// 启动GB28181服务器
		go func() {
			err := engine.Run(context.Background(), configFilePath)
			if err != nil {
				global.Logger.Panicf("GB28181 Server Start error: %s", err)
			}
		}()

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
