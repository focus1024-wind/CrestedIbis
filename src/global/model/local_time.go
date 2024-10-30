package model

import (
	"CrestedIbis/src/global"
	"database/sql/driver"
	"fmt"
	"time"
)

type LocalTime time.Time

func (t *LocalTime) MarshalJSON() ([]byte, error) {
	tTime := time.Time(*t)
	return []byte(fmt.Sprintf("\"%v\"", tTime.Format(time.DateTime))), nil
}

func (t *LocalTime) UnmarshalJSON(data []byte) error {
	if data[0] == '"' && data[len(data)-1] == '"' {
		data = data[1 : len(data)-1]
	}
	location, err := time.ParseInLocation(time.DateTime, string(data), time.Local)
	if err != nil {
		global.Logger.Errorf("%s 转 LocalTime 格式转换失败: %s", string(data), err.Error())
		return err
	}
	*t = LocalTime(location)
	return nil
}

// Value gorm在底层通过值调用，这里不要修改类型
func (t LocalTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	tlt := time.Time(t)
	//判断给定时间是否和默认零时间的时间戳相同
	if tlt.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return tlt, nil
}
func (t *LocalTime) Scan(v interface{}) error {
	if value, ok := v.(time.Time); ok {
		*t = LocalTime(value)
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}
