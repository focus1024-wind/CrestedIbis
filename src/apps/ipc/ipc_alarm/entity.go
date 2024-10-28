package ipc_alarm

import (
	"CrestedIbis/gb28181_server"
	"CrestedIbis/src/global/model"
)

type IpcAlarm struct {
	gb28181_server.Alarm
	ID         int64       `gorm:"primary_key;auto_increment" json:"id"`
	IpcRecords []IpcRecord `gorm:"foreignKey:AlarmID;references:ID" json:"ipc_records"`
	model.BaseModel
}

type IpcRecord struct {
	gb28181_server.Record
	ID      int64 `gorm:"primary_key;auto_increment" json:"id"`
	AlarmID int64
	model.BaseModel
}
