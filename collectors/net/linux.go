//go:build linux

package net

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func (ns *Info) update() {
	contents, err := os.ReadFile("/proc/net/dev")
	if err != nil {
		fmt.Println("Error reading /proc/net/dev:", err)
		return
	}

	bIn := uint64(0)
	bOut := uint64(0)
	for _, line := range strings.Split(string(contents), "\n") {
		fields := strings.Fields(line)
		if len(fields) >= 17 {
			interfaceName := strings.TrimSuffix(strings.TrimPrefix(fields[0], ":"), ":")
			if interfaceName != "lo" &&
				!strings.HasPrefix(interfaceName, "vmbr") &&
				!strings.HasPrefix(interfaceName, "veth") &&
				!strings.HasPrefix(interfaceName, "tap") &&
				!strings.HasPrefix(interfaceName, "fwln") &&
				!strings.HasPrefix(interfaceName, "fwbr") &&
				!strings.HasPrefix(interfaceName, "fwpr") {
				bytesIn, _ := strconv.ParseUint(fields[1], 10, 64)
				bytesOut, _ := strconv.ParseUint(fields[9], 10, 64)
				bIn += bytesIn
				bOut += bytesOut
			}
		}
	}
	ns.in = bIn
	ns.out = bOut
}
