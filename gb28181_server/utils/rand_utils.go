package utils

import (
	"math/rand"
	"time"
)

func RandNum(n int) string {
	letterBytes := "0123456789"
	return randStringBySource(letterBytes, n)
}

// RandNumString 生成 n 位字符串数据
func RandNumString(n int) string {
	letterBytes := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	return randStringBySource(letterBytes, n)
}

// randStringBySource 根据 src 生成长度为 n 的随机字符串数据
func randStringBySource(src string, n int) string {
	randSeed := make([]byte, n)

	rand.New(rand.NewSource(time.Now().UnixNano())).Read(randSeed)

	srcLen := len(src)

	// fill output
	output := make([]byte, n)
	for pos := range output {
		random := randSeed[pos]
		randomPos := random % uint8(srcLen)
		output[pos] = src[randomPos]
	}

	return string(output)
}
