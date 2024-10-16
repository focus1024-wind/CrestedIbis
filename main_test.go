package main

import (
	"testing"
)

func Test(t *testing.T) {
	snapshot := "store/snapshot/"
	if '/' == snapshot[len(snapshot)-1] {
		snapshot = snapshot[:len(snapshot)-1]
	}
}
