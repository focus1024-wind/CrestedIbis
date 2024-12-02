package gb28181_server

import (
	"errors"
	"fmt"
)

var (
	cmdName2Code = map[string]byte{
		"stop":       0,
		"right":      1,
		"left":       2,
		"down":       4,
		"down_right": 5,
		"down_left":  6,
		"up":         8,
		"up_right":   9,
		"up_left":    10,
		"zoom_in":    16,
		"zoom_out":   32,
	}
)

func toPtzStrByCmdName(cmdName string, horizontalSpeed, verticalSpeed, zoomSpeed uint8) (string, error) {
	if cmdCode, ok := cmdName2Code[cmdName]; ok {
		return toPtzStr(cmdCode, horizontalSpeed, verticalSpeed, zoomSpeed), nil
	} else {
		return "", errors.New("unknown command")
	}
}

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
