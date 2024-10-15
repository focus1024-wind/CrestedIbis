package ipc

import (
	"CrestedIbis/src/utils"
	"fmt"
)

func GenUploadImageAccessToken(deviceId string) string {
	token, err := utils.JwtToken{}.GenTempAccessToken(deviceId, []string{"ipc_device"}, 180)
	if err != nil {
		fmt.Println(err)
		return ""
	} else {
		return token
	}
}
