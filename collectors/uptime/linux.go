//go:build linux

package uptime

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func (ns *Info) update() {
	contents, err := os.ReadFile("/proc/uptime")
	if err != nil {
		fmt.Println("Error reading /proc/uptime:", err)
		return
	}
	data := strings.Split(string(contents), " ")[0]
	up, _ := strconv.ParseFloat(data, 64)
	ns.up = up
}
