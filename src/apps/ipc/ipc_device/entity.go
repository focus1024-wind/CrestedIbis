package ipc_device

import (
	"CrestedIbis/gb28181_server"
	"CrestedIbis/src/apps/site"
	"CrestedIbis/src/global/model"
	"encoding/json"
)

type simpleSite struct {
	Id   int64  `gorm:"-" json:"id"`
	Name string `gorm:"-" json:"name"`
}

type IpcDevice struct {
	gb28181_server.GB28181Device
	ID            int64                 `gorm:"primary_key;auto_increment" json:"id"`
	DeviceID      string                `gorm:"uniqueIndex;column:device_id" json:"device_id"`
	IpcChannels   []IpcChannel          `gorm:"foreignKey:ParentID;references:DeviceID" json:"ipc_channels"`
	ChannelNum    int64                 `json:"channel_num" desc:"设备通道数"`
	RegisterTime  model.LocalTime       `json:"register_time" desc:"设备最新注册时间"`
	KeepaliveTime model.LocalTime       `json:"keepalive_time" desc:"设备最新心跳时间"`
	SiteId        *int64                `json:"site_id"`
	Site          *site.SiteParentModel `json:"-" desc:"设备所属区域"`
	Site1         simpleSite            `gorm:"-" json:"site1" desc:"一级区域"`
	Site2         simpleSite            `gorm:"-" json:"site2" desc:"二级区域"`
	Site3         simpleSite            `gorm:"-" json:"site3" desc:"三级区域"`
	model.BaseHardDeleteModel
}

func (ipcDevice *IpcDevice) MarshalJSON() ([]byte, error) {
	if ipcDevice.SiteId != nil {
		siteModel := ipcDevice.Site
		for {
			if siteModel.Level == 1 {
				ipcDevice.Site1.Id = siteModel.Id
				ipcDevice.Site1.Name = siteModel.Name
			}
			if siteModel.Level == 2 {
				ipcDevice.Site2.Id = siteModel.Id
				ipcDevice.Site2.Name = siteModel.Name
			}
			if siteModel.Level == 3 {
				ipcDevice.Site3.Id = siteModel.Id
				ipcDevice.Site3.Name = siteModel.Name
			}
			if siteModel.Pid == nil {
				break
			}
			siteModel = siteModel.Parent
			if siteModel == nil {
				break
			}
		}
	}

	return json.Marshal(*ipcDevice)
}

type IpcChannel struct {
	gb28181_server.GB28181Channel
	ID       int64  `gorm:"primary_key;auto_increment" json:"id"`
	ParentID string `gorm:"uniqueIndex:channel_index" json:"parent_id"`
	DeviceID string `gorm:"uniqueIndex:channel_index" json:"device_id"`
	model.BaseHardDeleteModel
}

// IpcDevicePage 分页查询设备响应结构
type IpcDevicePage struct {
	Total    int64       `json:"total" desc:"设备总数量"`
	Data     []IpcDevice `json:"data" desc:"设备列表"`
	Page     int64       `json:"page" desc:"页码"`
	PageSize int64       `json:"page_size" desc:"每页查询数量"`
}

// IpcDeviceID 请求体设备ID包装结构
type IpcDeviceID struct {
	DeviceID string `json:"device_id"`
}
