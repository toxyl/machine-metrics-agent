//go:build windows

package net

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func (ns *Info) update() {
	cmd := exec.Command("wmic", "path", "Win32_PerfRawData_Tcpip_NetworkInterface", "get", "BytesReceivedPersec,BytesSentPersec")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error executing wmic command:", err)
		return
	}
	bIn := uint64(0)
	bOut := uint64(0)
	lines := strings.Split(string(output), "\r\n")
	for _, line := range lines[1:] { // line 0 is the header
		fields := strings.Fields(line)
		if len(fields) >= 2 {
			bytesIn, _ := strconv.ParseUint(fields[0], 10, 64)
			bytesOut, _ := strconv.ParseUint(fields[1], 10, 64)
			bIn += bytesIn
			bOut += bytesOut
		}
	}
	ns.in = bIn
	ns.out = bOut
}
