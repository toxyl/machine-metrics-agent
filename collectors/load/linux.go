//go:build linux

package load

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func (la *Info) update() {
	contents, err := os.ReadFile("/proc/loadavg")
	if err != nil {
		fmt.Println("Error reading /proc/loadavg:", err)
		return
	}

	fields := strings.Fields(string(contents))
	if len(fields) >= 3 {
		la.avg1, _ = strconv.ParseFloat(fields[0], 64)
		la.avg5, _ = strconv.ParseFloat(fields[1], 64)
		la.avg15, _ = strconv.ParseFloat(fields[2], 64)
	}
}
