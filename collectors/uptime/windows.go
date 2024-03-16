//go:build windows

package uptime

import (
	"syscall"
)

var (
	modkernel32    = syscall.NewLazyDLL("kernel32.dll")
	getTickCount64 = modkernel32.NewProc("GetTickCount64")
)

func (ns *Info) update() {
	r, _, _ := getTickCount64.Call()
	ns.up = float64(r) / 1000
}
