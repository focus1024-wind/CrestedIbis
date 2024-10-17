package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ipRealLocationModel struct {
	Status  string `json:"status"`
	Country string `json:"country"`
	City    string `json:"city"`
}

func GetRealLocationByIpAddr(ipAddr string) string {
	if ipAddr == "localhost" || ipAddr == "127.0.0.1" || ipAddr == "::1" {
		return ipAddr
	} else {
		url := fmt.Sprintf("http://ip-api.com/json/%s?lang=zh-CN&fields=status,country,city", ipAddr)

		// 发起HTTP GET请求
		resp, err := http.Get(url)
		if err != nil {
			return "未知"
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				return
			}
		}(resp.Body)

		// 读取响应体
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "未知"
		}

		var responseData ipRealLocationModel
		err = json.Unmarshal(body, &responseData)
		if err != nil {
			return "未知"
		}
		if responseData.Status == "success" {
			return fmt.Sprintf("%s %s", responseData.Country, responseData.City)
		} else {
			return "未知"
		}
	}
}
