package gb28181_server

import (
	"CrestedIbis/gb28181_server/utils"
	"fmt"
)

var (
	// DeviceInfoXML 查询设备详情xml样式
	DeviceInfoXML = `<?xml version="1.0"?>
<Query>
<CmdType>DeviceInfo</CmdType>
<SN>%d</SN>
<DeviceID>%s</DeviceID>
</Query>`

	// CatalogXML 获取设备子通道列表
	CatalogXML = `<?xml version="1.0"?>
<Query>
<CmdType>Catalog</CmdType>
<SN>%d</SN>
<DeviceID>%s</DeviceID>
</Query>`

	SnapShotXML = `<?xml version="1.0" ?>
<Control>
    <CmdType>DeviceConfig</CmdType>
    <SN>%d</SN>
    <DeviceID>%s</DeviceID>
    <SnapShotConfig>
        <SnapNum>%d</SnapNum>
        <Interval>%d</Interval>
        <UploadURL>%s</UploadURL>
        <SessionID>%s</SessionID>
    </SnapShotConfig>
</Control>`
)

// BuildDeviceInfoXML 获取设备详情指令
func BuildDeviceInfoXML(sn int, id string) string {
	return fmt.Sprintf(DeviceInfoXML, sn, id)
}

// BuildCatalogXML 获取NVR下设备列表指令
func BuildCatalogXML(sn int, id string) string {
	return fmt.Sprintf(CatalogXML, sn, id)
}

// BuildSnapShotXML 图片抓拍
func BuildSnapShotXML(sn int, id string, snapNum int, interval int) string {
	uploadUrl := GlobalDeviceStore.SnapShotUploadUrl(id)
	return fmt.Sprintf(SnapShotXML, sn, id, snapNum, interval, uploadUrl, utils.RandNumString(32))
}
