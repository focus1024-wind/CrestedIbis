package gb28181_server

type Alarm struct {
	DeviceID         string  `desc:"报警设备" json:"device_id"`
	AlarmPriority    int     `desc:"报警级别" json:"alarm_priority"`
	AlarmMethod      int     `desc:"报警方式" json:"alarm_method"`
	AlarmTime        float64 `desc:"报警时间" json:"alarm_time"`
	AlarmDescription string  `desc:"报警描述" json:"alarm_description"`
	Longitude        float32 `desc:"报警设备经度" json:"longitude"`
	Latitude         float32 `desc:"报警设备纬度" json:"latitude"`
	AlarmType        int     `json:"alarm_type"`
}

type AlarmHandlerInterface interface {
	Handler(alarm Alarm)
}
