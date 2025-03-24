package gb28181_server

import (
	"errors"
	"fmt"
)

// PTZ控制命令常量
const (
	PTZ_CMD_STOP       = "stop"       // 停止
	PTZ_CMD_RIGHT      = "right"      // 右转
	PTZ_CMD_LEFT       = "left"       // 左转
	PTZ_CMD_DOWN       = "down"       // 下转
	PTZ_CMD_DOWN_RIGHT = "down_right" // 右下
	PTZ_CMD_DOWN_LEFT  = "down_left"  // 左下
	PTZ_CMD_UP         = "up"         // 上转
	PTZ_CMD_UP_RIGHT   = "up_right"   // 右上
	PTZ_CMD_UP_LEFT    = "up_left"    // 左上
	PTZ_CMD_ZOOM_IN    = "zoom_in"    // 放大
	PTZ_CMD_ZOOM_OUT   = "zoom_out"   // 缩小
	// 扩展PTZ命令
	PTZ_CMD_FOCUS_NEAR = "focus_near" // 焦距近
	PTZ_CMD_FOCUS_FAR  = "focus_far"  // 焦距远
	PTZ_CMD_IRIS_OPEN  = "iris_open"  // 光圈开
	PTZ_CMD_IRIS_CLOSE = "iris_close" // 光圈关
)

// 默认速度值
const (
	DEFAULT_HORIZONTAL_SPEED = 50 // 默认水平速度
	DEFAULT_VERTICAL_SPEED   = 50 // 默认垂直速度
	DEFAULT_ZOOM_SPEED       = 50 // 默认缩放速度
)

var (
	cmdName2Code = map[string]byte{
		PTZ_CMD_STOP:       0,
		PTZ_CMD_RIGHT:      1,
		PTZ_CMD_LEFT:       2,
		PTZ_CMD_DOWN:       4,
		PTZ_CMD_DOWN_RIGHT: 5,
		PTZ_CMD_DOWN_LEFT:  6,
		PTZ_CMD_UP:         8,
		PTZ_CMD_UP_RIGHT:   9,
		PTZ_CMD_UP_LEFT:    10,
		PTZ_CMD_ZOOM_IN:    16,
		PTZ_CMD_ZOOM_OUT:   32,
		PTZ_CMD_FOCUS_NEAR: 64,
		PTZ_CMD_FOCUS_FAR:  65,
		PTZ_CMD_IRIS_OPEN:  66,
		PTZ_CMD_IRIS_CLOSE: 67,
	}
)

// PTZControlOptions PTZ控制选项
type PTZControlOptions struct {
	Command         string `json:"command" desc:"PTZ命令名称"`
	HorizontalSpeed uint8  `json:"horizontal_speed" desc:"水平速度, 0-255"`
	VerticalSpeed   uint8  `json:"vertical_speed" desc:"垂直速度, 0-255"`
	ZoomSpeed       uint8  `json:"zoom_speed" desc:"缩放速度, 0-255"`
}

// NewPTZControlOptions 创建默认PTZ控制选项
func NewPTZControlOptions(command string) *PTZControlOptions {
	return &PTZControlOptions{
		Command:         command,
		HorizontalSpeed: DEFAULT_HORIZONTAL_SPEED,
		VerticalSpeed:   DEFAULT_VERTICAL_SPEED,
		ZoomSpeed:       DEFAULT_ZOOM_SPEED,
	}
}

// toPtzStrByCmdName 根据命令名称生成PTZ控制字符串
func toPtzStrByCmdName(cmdName string, horizontalSpeed, verticalSpeed, zoomSpeed uint8) (string, error) {
	if cmdCode, ok := cmdName2Code[cmdName]; ok {
		return toPtzStr(cmdCode, horizontalSpeed, verticalSpeed, zoomSpeed), nil
	} else {
		return "", errors.New("unknown command")
	}
}

// toPtzStr 根据命令代码生成PTZ控制字符串
func toPtzStr(cmdCode, horizontalSpeed, verticalSpeed, zoomSpeed uint8) string {
	checkCode := uint16(0xA5+0x0F+0x01+cmdCode+horizontalSpeed+verticalSpeed+(zoomSpeed&0xF0)) % 0x100

	return fmt.Sprintf("A50F01%02X%02X%02X%01X0%02X",
		cmdCode,
		horizontalSpeed,
		verticalSpeed,
		zoomSpeed>>4,
		checkCode,
	)
}

// ControlPTZ 控制通道PTZ
func ControlPTZ(channelID string, options *PTZControlOptions) (int, error) {
	// 查找通道
	var channel GB28181Channel
	var found bool

	// 从所有在线设备中查找通道
	DeviceChannels.Range(func(deviceID, channelIDs interface{}) bool {
		for _, id := range channelIDs.([]string) {
			if id == channelID {
				if ch, ok := GlobalGB28181DeviceStore.LoadChannel(deviceID.(string), channelID); ok {
					channel = ch
					found = true
					return false // 停止遍历
				}
			}
		}
		return true // 继续遍历
	})

	if !found {
		return 404, errors.New("未找到通道信息")
	}

	// 生成PTZ控制字符串
	ptzStr, err := toPtzStrByCmdName(options.Command, options.HorizontalSpeed, options.VerticalSpeed, options.ZoomSpeed)
	if err != nil {
		return 400, err
	}

	// 发送PTZ控制命令
	return channel.PtzControl(ptzStr), nil
}

// ControlChannelPTZ 直接通过通道对象控制PTZ
func (channel *GB28181Channel) ControlPTZ(options *PTZControlOptions) (int, error) {
	// 生成PTZ控制字符串
	ptzStr, err := toPtzStrByCmdName(options.Command, options.HorizontalSpeed, options.VerticalSpeed, options.ZoomSpeed)
	if err != nil {
		return 400, err
	}

	// 发送PTZ控制命令
	return channel.PtzControl(ptzStr), nil
}
