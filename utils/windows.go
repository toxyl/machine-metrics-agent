//go:build windows

package utils

func IsRunningInLXC() bool {
	return false
}
