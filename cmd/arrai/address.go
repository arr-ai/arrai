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

// grpcPassthroughTarget formats addr for grpc.NewClient with Dial-like
// resolution: no DNS lookup, dial the host:port string as-is. NewClient
// defaults to the "dns" resolver; Dial used "passthrough".
func grpcPassthroughTarget(addr string) string {
	return "passthrough:///" + addr
}
