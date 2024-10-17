package main

import (
	"CrestedIbis/src/apps/ipc"
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	fmt.Println(ipc.GenUploadImageAccessToken("123"))
}
