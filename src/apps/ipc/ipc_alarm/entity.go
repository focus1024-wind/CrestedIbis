package ipc_alarm

import (
	"CrestedIbis/gb28181_server"
	"CrestedIbis/src/global/model"
)

type IpcAlarm struct {
	model.IDModel
	gb28181_server.Alarm
	IpcRecords []IpcRecord `gorm:"foreignKey:AlarmID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"ipc_records"`
	model.BaseModel
}

type IpcRecord struct {
	model.IDModel
	gb28181_server.Record
	AlarmID *int64 `json:"alarm_id"`
	model.BaseModel
}
