package gb28181_server

type Alarm struct {
	DeviceID         string `desc:"报警设备设备ID、对应SIP信息中的设备ID" json:"device_id" required:"true" xml:"-"`
	ChannelID        string `desc:"报警设备通道ID、对应XML中的DeviceID信息" json:"channel_id" required:"true" xml:"DeviceID"`
	AlarmPriority    string `desc:"报警级别" json:"alarm_priority"`
	AlarmMethod      string `desc:"报警方式" json:"alarm_method" enum:"1: 电话报警、2: 设备报警、3: 短信报警、4: GPS报警、5: 视频报警、6: 设备故障报警"`
	AlarmType        string `json:"alarm_type" xml:"Info>AlarmType"`
	AlarmTime        string `desc:"报警时间" json:"alarm_time"`
	AlarmDescription string `desc:"报警描述" json:"alarm_description"`
	Longitude        string `desc:"报警设备经度" json:"longitude"`
	Latitude         string `desc:"报警设备纬度" json:"latitude"`
}

type AlarmHandlerInterface interface {
	Handler(alarm Alarm)
	RecordHandler(record Record)
}
