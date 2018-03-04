package main

import (
	"strings"
)

func arraiAddress(addr string) string {
	if !strings.ContainsRune(addr, ':') {
		addr += ":42241"
	}
	return addr
}
