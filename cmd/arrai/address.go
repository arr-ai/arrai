package main

import (
	"fmt"
	"strings"
)

func arraiAddressWithPort(addr string, defaultPort int) string {
	if !strings.ContainsRune(addr, ':') {
		addr = fmt.Sprintf("%s:%d", addr, defaultPort)
	}
	return addr
}

func arraiAddress(addr string) string {
	return arraiAddressWithPort(addr, 42241)
}
