//go:build linux

package utils

import (
	"bufio"
	"os"
	"strings"
)

func IsRunningInLXC() bool {
	file, err := os.Open("/proc/self/mounts")
	if err != nil {
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "lxcfs /proc") { // LXC containers on Proxmox always have several lxcfs mounted in /proc
			return true
		}
	}
	return false
}
