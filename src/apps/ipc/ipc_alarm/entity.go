package ipc_alarm

import (
	"CrestedIbis/gb28181_server"
	"CrestedIbis/src/global/model"
)

type IpcAlarm struct {
	gb28181_server.Alarm
	ID int64 `gorm:"primary_key;auto_increment" json:"id"`
	model.BaseModel
}
