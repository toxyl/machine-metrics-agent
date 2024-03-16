//go:build windows

package cpu

import (
	"math"
	"runtime"
	"syscall"
	"time"
	"unsafe"
)

type FILETIME struct {
	dwLowDateTime  uint64
	dwHighDateTime uint64
}

type SystemTimes struct {
	IdleTime   FILETIME
	KernelTime FILETIME
	UserTime   FILETIME
}

var (
	modKernel32        = syscall.NewLazyDLL("kernel32.dll")
	procGetSystemTimes = modKernel32.NewProc("GetSystemTimes")
)

func (c *Info) update() {
	var idleTime, kernelTime, userTime FILETIME
	procGetSystemTimes.Call(
		uintptr(unsafe.Pointer(&idleTime)),
		uintptr(unsafe.Pointer(&kernelTime)),
		uintptr(unsafe.Pointer(&userTime)),
	)
	t := SystemTimes{IdleTime: idleTime, KernelTime: kernelTime, UserTime: userTime}
	idleFirst := t.IdleTime.dwLowDateTime | (t.IdleTime.dwHighDateTime << 32)
	kernelFirst := t.KernelTime.dwLowDateTime | (t.KernelTime.dwHighDateTime << 32)
	userFirst := t.UserTime.dwLowDateTime | (t.UserTime.dwHighDateTime << 32)

	time.Sleep(time.Second)

	procGetSystemTimes.Call(
		uintptr(unsafe.Pointer(&idleTime)),
		uintptr(unsafe.Pointer(&kernelTime)),
		uintptr(unsafe.Pointer(&userTime)),
	)
	t2 := SystemTimes{IdleTime: idleTime, KernelTime: kernelTime, UserTime: userTime}
	idleSecond := t2.IdleTime.dwLowDateTime | (t2.IdleTime.dwHighDateTime << 32)
	kernelSecond := t2.KernelTime.dwLowDateTime | (t2.KernelTime.dwHighDateTime << 32)
	userSecond := t2.UserTime.dwLowDateTime | (t2.UserTime.dwHighDateTime << 32)

	totalSys := (kernelSecond - kernelFirst) + (userSecond - userFirst)
	usedPct := math.Max(0.0, float64(totalSys-(idleSecond-idleFirst))/float64(totalSys))

	c.UsedPct = usedPct
	c.Cores = uint(runtime.NumCPU())
}
