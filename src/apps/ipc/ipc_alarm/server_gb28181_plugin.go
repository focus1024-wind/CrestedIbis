package ipc_alarm

import (
	"CrestedIbis/gb28181_server"
	"fmt"
)

func (Alarm) Handler(alarm gb28181_server.Alarm) {
	fmt.Println(alarm)
}
