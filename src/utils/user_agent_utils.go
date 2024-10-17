package utils

import "strings"

func GetOsByUserAgent(userAgent string) string {
	if strings.Contains(userAgent, "Windows") {
		return "Windows"
	} else if strings.Contains(userAgent, "Macintosh") {
		return "MacOs"
	} else if strings.Contains(userAgent, "Linux") && strings.Contains(userAgent, "Android") {
		return "Android"
	} else if strings.Contains(userAgent, "Linux") && !strings.Contains(userAgent, "Android") {
		return "Linux"
	} else {
		return "未知"
	}
}
